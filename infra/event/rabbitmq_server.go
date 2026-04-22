// Package event 提供基於 RabbitMQ 的事件處理服務器實現
package event

import (
	"context"
	"encoding/json"
	"fmt"

	"self_go_gin/infra/env"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// RabbitMQServer RabbitMQ 服務器，實現 Subscriber 接口
type RabbitMQServer struct {
	conn        *amqp091.Connection
	channel     *amqp091.Channel
	exchange    string
	handlers    map[string]Handler
	deliveries  map[string]<-chan amqp091.Delivery
	closeSignal chan bool
	isRunning   bool
}

var rabbitmqServer *RabbitMQServer

// InitRabbitMQServer 初始化 RabbitMQ 服務器
func InitRabbitMQServer(serverConfig *env.ServerConfig) (*RabbitMQServer, error) {
	// 構建 RabbitMQ 連接字符串
	rabbitmqURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		serverConfig.RabbitMQ.User,
		serverConfig.RabbitMQ.Password,
		serverConfig.RabbitMQ.Host,
		serverConfig.RabbitMQ.Port,
		serverConfig.RabbitMQ.VHost,
	)

	// 建立連接
	conn, err := amqp091.Dial(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// 建立通道
	channel, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close connection after channel error: %w", err)
		}
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 設置通道的 QoS（Quality of Service）
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		err := channel.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close channel after QoS error: %w", err)
		}
		err = conn.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close connection after QoS error: %w", err)
		}
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	// 聲明交換機
	exchangeName := serverConfig.RabbitMQ.Exchange
	err = channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		err := channel.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close channel after exchange declare error: %w", err)
		}
		err = conn.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close connection after exchange declare error: %w", err)
		}
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	rabbitmqServer = &RabbitMQServer{
		conn:        conn,
		channel:     channel,
		exchange:    exchangeName,
		handlers:    make(map[string]Handler),
		deliveries:  make(map[string]<-chan amqp091.Delivery),
		closeSignal: make(chan bool),
		isRunning:   false,
	}

	fmt.Println("RabbitMQ server initialized successfully")
	return rabbitmqServer, nil
}

// GetRabbitMQServer 獲取 RabbitMQ 服務器實例
func GetRabbitMQServer() *RabbitMQServer {
	return rabbitmqServer
}

// Subscribe 訂閱事件，註冊事件處理器和隊列綁定
func (s *RabbitMQServer) Subscribe(handler Handler) error {
	eventType := handler.EventType()

	// 檢查是否已經註冊
	if _, exists := s.handlers[eventType]; exists {
		return fmt.Errorf("handler for event type %s already registered", eventType)
	}

	// 保存處理器
	s.handlers[eventType] = handler

	// 聲明隊列
	queue, err := s.channel.QueueDeclare(
		eventType, // 队列名称
		true,      // 是否持久化
		false,     // 是否自动删除
		false,     // 是否为排他性队列
		true,      // 是否等待服务器返回响应
		amqp091.Table{ // 额外参数
			"x-max-priority": int32(10), // 啟用優先級隊列
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue for event type %s: %w", eventType, err)
	}

	// 綁定隊列到交換機
	err = s.channel.QueueBind(
		queue.Name, // queue name
		eventType,  // routing key
		s.exchange, // exchange name
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue for event type %s: %w", eventType, err)
	}

	// 開始消費消息
	deliveries, err := s.channel.Consume(
		queue.Name, // queue
		eventType,  // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages for event type %s: %w", eventType, err)
	}

	// 保存 deliveries 通道
	s.deliveries[eventType] = deliveries

	// 啟動協程處理消息
	go s.handleMessages(eventType, deliveries, handler)

	fmt.Printf("Event handler registered and subscribed: %s\n", eventType)
	return nil
}

// handleMessages 處理來自隊列的消息
func (s *RabbitMQServer) handleMessages(eventType string, deliveries <-chan amqp091.Delivery, handler Handler) {
	for delivery := range deliveries {
		var event Event
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			zap.S().Error("Failed to unmarshal event", zap.String("event_type", eventType), zap.Error(err))
			// 拒絕消息，不重新入隊
			err := delivery.Nack(false, false)
			if err != nil {
				zap.S().Error("Failed to nack message", zap.String("event_type", eventType), zap.Error(err))
			}
			continue
		}

		fmt.Printf("Processing event from RabbitMQ: type=%s, source=%s\n", event.Type, event.Source)

		// 調用處理器
		ctx := context.Background()
		if err := handler.Handle(ctx, &event); err != nil {
			zap.S().Error("Error processing event", zap.String("event_type", eventType), zap.Error(err))
			// 拒絕消息，重新入隊
			err = delivery.Nack(false, true)
			if err != nil {
				zap.S().Error("Failed to nack message", zap.String("event_type", eventType), zap.Error(err))
			}
		} else {
			// 確認消息
			err := delivery.Ack(false)
			if err != nil {
				zap.S().Error("Failed to ack message", zap.String("event_type", eventType), zap.Error(err))
			}
		}
	}

	fmt.Printf("Message handler closed for event type: %s\n", eventType)
}

// Start 啟動事件處理服務器（非阻塞）
func (s *RabbitMQServer) Start() error {
	if s.isRunning {
		return fmt.Errorf("RabbitMQ server is already running")
	}

	s.isRunning = true
	fmt.Println("Starting RabbitMQ event server...")

	// 監聽連接關閉信號
	go func() {
		<-s.conn.NotifyClose(make(chan *amqp091.Error))
		s.isRunning = false
		fmt.Println("RabbitMQ connection closed")
	}()

	return nil
}

// Run 運行服務器（阻塞方法）
func (s *RabbitMQServer) Run() error {
	if s.isRunning {
		return fmt.Errorf("RabbitMQ server is already running")
	}

	s.isRunning = true
	fmt.Println("Running RabbitMQ event server...")

	// 等待關閉信號
	<-s.closeSignal

	return nil
}

// Shutdown 優雅關閉服務器
func (s *RabbitMQServer) Shutdown(ctx context.Context) error {
	if !s.isRunning {
		return fmt.Errorf("RabbitMQ server is not running")
	}

	s.isRunning = false

	// 關閉所有消費者
	for eventType := range s.handlers {
		err := s.channel.Cancel(eventType, false)
		if err != nil {
			zap.S().Error("Failed to cancel consumer", zap.String("event_type", eventType), zap.Error(err))
		}
	}

	// 關閉通道
	if err := s.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	// 關閉連接
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	// 發送關閉信號
	select {
	case s.closeSignal <- true:
	case <-ctx.Done():
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	default:
	}

	fmt.Println("RabbitMQ server shutdown")
	return nil
}

// IsRunning 檢查服務器是否在運行
func (s *RabbitMQServer) IsRunning() bool {
	return s.isRunning
}
