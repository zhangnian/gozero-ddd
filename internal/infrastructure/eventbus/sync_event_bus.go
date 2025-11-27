package eventbus

import (
	"context"
	"log"
	"sync"

	"gozero-ddd/internal/domain/event"
)

// SyncEventBus 同步事件总线实现
// 适用于单体应用，事件在进程内同步处理
// 这是一个简单的事件总线，不依赖外部消息队列
type SyncEventBus struct {
	mu          sync.RWMutex
	handlers    map[string][]event.EventHandler // 特定事件的处理器
	allHandlers []event.EventHandler            // 处理所有事件的处理器
}

// NewSyncEventBus 创建同步事件总线
func NewSyncEventBus() *SyncEventBus {
	return &SyncEventBus{
		handlers:    make(map[string][]event.EventHandler),
		allHandlers: make([]event.EventHandler, 0),
	}
}

// 确保实现了接口
var _ event.EventBus = (*SyncEventBus)(nil)

// Subscribe 订阅特定事件
// eventName 为事件名称，如 "knowledge_base.created"
func (b *SyncEventBus) Subscribe(eventName string, handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventName] = append(b.handlers[eventName], handler)
	log.Printf("📫 [EventBus] 注册事件处理器: %s", eventName)
}

// SubscribeAll 订阅所有事件
// 用于日志记录、审计等通用处理
func (b *SyncEventBus) SubscribeAll(handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.allHandlers = append(b.allHandlers, handler)
	log.Printf("📫 [EventBus] 注册全局事件处理器")
}

// Publish 发布单个事件
func (b *SyncEventBus) Publish(ctx context.Context, evt event.DomainEvent) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	eventName := evt.EventName()
	log.Printf("📤 [EventBus] 发布事件: %s", eventName)

	// 调用特定事件的处理器
	if handlers, ok := b.handlers[eventName]; ok {
		for _, handler := range handlers {
			if err := b.invokeHandler(ctx, handler, evt); err != nil {
				return err
			}
		}
	}

	// 调用全局处理器
	for _, handler := range b.allHandlers {
		if err := b.invokeHandler(ctx, handler, evt); err != nil {
			return err
		}
	}

	return nil
}

// PublishAll 发布多个事件
func (b *SyncEventBus) PublishAll(ctx context.Context, events []event.DomainEvent) error {
	for _, evt := range events {
		if err := b.Publish(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}

// invokeHandler 调用处理器（带错误处理）
func (b *SyncEventBus) invokeHandler(ctx context.Context, handler event.EventHandler, evt event.DomainEvent) error {
	if err := handler.Handle(ctx, evt); err != nil {
		log.Printf("❌ [EventBus] 事件处理失败: %s, 错误: %v", evt.EventName(), err)
		// 记录错误但不中断后续处理
		return nil
	}
	return nil
}

