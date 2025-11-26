package event

import "context"

// EventPublisher 事件发布器接口
// 定义在领域层，实现在基础设施层
// 这体现了依赖倒置原则
type EventPublisher interface {
	// Publish 发布单个事件
	Publish(ctx context.Context, event DomainEvent) error

	// PublishAll 发布多个事件
	PublishAll(ctx context.Context, events []DomainEvent) error
}

// EventHandler 事件处理器接口
// 用于订阅和处理领域事件
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event DomainEvent) error

	// EventName 返回处理器关注的事件名称
	// 返回空字符串表示处理所有事件
	EventName() string
}

// EventSubscriber 事件订阅器接口
// 用于注册事件处理器
type EventSubscriber interface {
	// Subscribe 订阅特定事件
	Subscribe(eventName string, handler EventHandler)

	// SubscribeAll 订阅所有事件
	SubscribeAll(handler EventHandler)
}

// EventBus 事件总线接口
// 组合了发布和订阅功能
type EventBus interface {
	EventPublisher
	EventSubscriber
}

