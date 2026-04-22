// Package event 提供基於 RabbitMQ 的事件發布客戶端實現
package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"self_go_gin/infra/env"

	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient RabbitMQ 客戶端，實現 Publisher 接口
type RabbitMQClient struct {
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	exchange  string
	durable   bool
	autoAck   bool
	exclusive bool
	noWait    bool
	args      amqp091.Table
}

var rabbitmqClient *RabbitMQClient

// InitRabbitMQClient 初始化 RabbitMQ 客戶端
func InitRabbitMQClient(serverConfig *env.ServerConfig) (*RabbitMQClient, error) {
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

	rabbitmqClient = &RabbitMQClient{
		conn:      conn,
		channel:   channel,
		exchange:  exchangeName,
		durable:   true,
		autoAck:   false,
		exclusive: false,
		noWait:    false,
		args:      nil,
	}

	fmt.Println("RabbitMQ client initialized successfully")
	return rabbitmqClient, nil
}

// GetRabbitMQClient 獲取 RabbitMQ 客戶端實例
func GetRabbitMQClient() *RabbitMQClient {
	return rabbitmqClient
}

// Publish 發布事件到默認隊列
func (c *RabbitMQClient) Publish(ctx context.Context, event *Event) error {
	opts := &PublishOptions{
		Queue:      DefaultQueue,
		MaxRetry:   3,
		Priority:   5,
		RoutingKey: event.Type, // 使用事件類型作為路由鍵
	}
	return c.PublishWithOptions(ctx, event, opts)
}

// PublishWithOptions 使用自定義選項發布事件（實現 Publisher 接口）
func (c *RabbitMQClient) PublishWithOptions(ctx context.Context, event *Event, opts *PublishOptions) error {
	if opts == nil {
		opts = &PublishOptions{
			Queue:      DefaultQueue,
			MaxRetry:   3,
			Priority:   5,
			RoutingKey: DefaultQueue, // 默認使用隊列名稱作為路由鍵
		}
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 設置發布參數
	headers := amqp091.Table{
		"x-max-priority": opts.Priority,  // 設置消息優先級
		"x-max-length":   int32(1000000), // 隊列最大消息數
	}

	// 如果有延遲，添加延遲頭信息
	publishing := amqp091.Publishing{
		ContentType:  "application/json",
		Body:         payload,
		Headers:      headers,
		DeliveryMode: amqp091.Persistent, // 持久化消息
		Timestamp:    time.Now(),
	}

	// 發布消息
	err = c.channel.PublishWithContext(
		ctx,
		c.exchange,      // exchange
		opts.RoutingKey, // routing key
		false,           // mandatory
		false,           // immediate
		publishing,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	fmt.Printf("Event published to RabbitMQ: type=%s, queue=%s, routing_key=%s\n", event.Type, opts.Queue, opts.RoutingKey)
	return nil
}

// PublishDeferred 延遲發布事件
func (c *RabbitMQClient) PublishDeferred(ctx context.Context, event *Event, delaySeconds int) error {
	opts := &PublishOptions{
		Queue:    DefaultQueue,
		MaxRetry: 3,
		Delay:    time.Duration(delaySeconds) * time.Second,
		Priority: 5,
	}
	return c.PublishWithOptions(ctx, event, opts)
}

// PublishToHighPriorityQueue 發布到高優先級隊列
func (c *RabbitMQClient) PublishToHighPriorityQueue(ctx context.Context, event *Event) error {
	opts := &PublishOptions{
		Queue:    HighPriorityQueue,
		MaxRetry: 5,
		Priority: 10,
	}
	return c.PublishWithOptions(ctx, event, opts)
}

// PublishToLowPriorityQueue 發布到低優先級隊列
func (c *RabbitMQClient) PublishToLowPriorityQueue(ctx context.Context, event *Event) error {
	opts := &PublishOptions{
		Queue:    LowPriorityQueue,
		MaxRetry: 1,
		Priority: 1,
	}
	return c.PublishWithOptions(ctx, event, opts)
}

// Close 關閉客戶端連接
func (c *RabbitMQClient) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}
	fmt.Println("RabbitMQ client closed")
	return nil
}
