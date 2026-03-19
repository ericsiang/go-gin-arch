// Package event 提供事件代理的工廠和配置
package event

import (
	"context"
	"fmt"
	"sync"

	"self_go_gin/infra/env"
)

// BrokerType 事件代理類型
type BrokerType string

const (
	// BrokerTypeAsynq Asynq (Redis) 實現
	BrokerTypeAsynq BrokerType = "asynq"
	// BrokerTypeRabbitMQ RabbitMQ 實現
	BrokerTypeRabbitMQ BrokerType = "rabbitmq"
	// BrokerTypeKafka Kafka 實現
	BrokerTypeKafka BrokerType = "kafka"
)

// Broker 事件代理，封裝 Publisher 和 Subscriber
type Broker struct {
	brokerType BrokerType
	publisher  Publisher
	subscriber Subscriber
	mu         sync.RWMutex
}

// NewBroker 創建新的事件代理實例
func NewBroker(brokerType BrokerType, config *env.ServerConfig) (*Broker, error) {
	var publisher Publisher
	var subscriber Subscriber

	switch brokerType {
	case BrokerTypeAsynq:
		publisher = InitAsynqClient(config)
		subscriber = InitAsynqServer(config)

	case BrokerTypeRabbitMQ:
		// TODO: 實現 RabbitMQ
		return nil, fmt.Errorf("rabbitMQ broker not implemented yet")

	case BrokerTypeKafka:
		// TODO: 實現 Kafka
		return nil, fmt.Errorf("kafka broker not implemented yet")

	default:
		return nil, fmt.Errorf("unsupported broker type: %s", brokerType)
	}

	broker := &Broker{
		brokerType: brokerType,
		publisher:  publisher,
		subscriber: subscriber,
	}

	fmt.Printf("Event broker initialized: type=%s\n", brokerType)
	return broker, nil
}

// Publisher 獲取事件發布者
func (b *Broker) Publisher() Publisher {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.publisher
}

// Subscriber 獲取事件訂閱者
func (b *Broker) Subscriber() Subscriber {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.subscriber
}

// BrokerType 獲取代理類型
func (b *Broker) BrokerType() BrokerType {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.brokerType
}

// Close 關閉事件代理
func (b *Broker) Close(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var errs []error

	if b.publisher != nil {
		if err := b.publisher.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close publisher: %w", err))
		}
		b.publisher = nil
	}

	if b.subscriber != nil {
		if err := b.subscriber.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown subscriber: %w", err))
		}
		b.subscriber = nil
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing broker: %v", errs)
	}

	fmt.Println("Event broker closed successfully")
	return nil
}
