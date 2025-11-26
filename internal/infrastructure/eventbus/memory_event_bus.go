package eventbus

import (
	"context"
	"log"
	"sync"

	"gozero-ddd/internal/domain/event"
)

// MemoryEventBus 内存事件总线实现
// 适用于单体应用，事件在进程内同步处理
// 生产环境可以替换为基于消息队列的实现（如 Kafka、RabbitMQ）
type MemoryEventBus struct {
	mu       sync.RWMutex
	handlers map[string][]event.EventHandler // 特定事件的处理器
	allHandlers []event.EventHandler         // 处理所有事件的处理器
}

// NewMemoryEventBus 创建内存事件总线
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers:    make(map[string][]event.EventHandler),
		allHandlers: make([]event.EventHandler, 0),
	}
}

// 确保实现了接口
var _ event.EventBus = (*MemoryEventBus)(nil)

// Subscribe 订阅特定事件
// eventName 为事件名称，如 "knowledge_base.created"
func (b *MemoryEventBus) Subscribe(eventName string, handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventName] = append(b.handlers[eventName], handler)
	log.Printf("📫 [EventBus] 注册事件处理器: %s", eventName)
}

// SubscribeAll 订阅所有事件
// 用于日志记录、审计等通用处理
func (b *MemoryEventBus) SubscribeAll(handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.allHandlers = append(b.allHandlers, handler)
	log.Printf("📫 [EventBus] 注册全局事件处理器")
}

// Publish 发布单个事件
func (b *MemoryEventBus) Publish(ctx context.Context, evt event.DomainEvent) error {
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
func (b *MemoryEventBus) PublishAll(ctx context.Context, events []event.DomainEvent) error {
	for _, evt := range events {
		if err := b.Publish(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}

// invokeHandler 调用处理器（带错误处理）
func (b *MemoryEventBus) invokeHandler(ctx context.Context, handler event.EventHandler, evt event.DomainEvent) error {
	// 在生产环境中，可以考虑：
	// 1. 异步处理（使用 goroutine）
	// 2. 重试机制
	// 3. 死信队列
	// 4. 熔断器
	if err := handler.Handle(ctx, evt); err != nil {
		log.Printf("❌ [EventBus] 事件处理失败: %s, 错误: %v", evt.EventName(), err)
		// 根据业务需求决定是否继续处理其他事件
		// 这里选择继续，不中断后续处理
		return nil
	}
	return nil
}

// AsyncMemoryEventBus 异步内存事件总线
// 事件处理在单独的 goroutine 中执行，不阻塞发布者
type AsyncMemoryEventBus struct {
	*MemoryEventBus
	workerCount int
	eventChan   chan eventWrapper
	wg          sync.WaitGroup
}

type eventWrapper struct {
	ctx   context.Context
	event event.DomainEvent
}

// NewAsyncMemoryEventBus 创建异步内存事件总线
func NewAsyncMemoryEventBus(workerCount int) *AsyncMemoryEventBus {
	if workerCount <= 0 {
		workerCount = 1
	}

	bus := &AsyncMemoryEventBus{
		MemoryEventBus: NewMemoryEventBus(),
		workerCount:    workerCount,
		eventChan:      make(chan eventWrapper, 1000), // 缓冲区大小可配置
	}

	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		bus.wg.Add(1)
		go bus.worker(i)
	}

	return bus
}

// worker 工作协程
func (b *AsyncMemoryEventBus) worker(id int) {
	defer b.wg.Done()
	log.Printf("🚀 [EventBus] 工作协程 #%d 启动", id)

	for wrapper := range b.eventChan {
		b.processEvent(wrapper.ctx, wrapper.event)
	}

	log.Printf("🛑 [EventBus] 工作协程 #%d 停止", id)
}

// processEvent 处理事件
func (b *AsyncMemoryEventBus) processEvent(ctx context.Context, evt event.DomainEvent) {
	// 调用父类的同步发布方法
	if err := b.MemoryEventBus.Publish(ctx, evt); err != nil {
		log.Printf("❌ [EventBus] 异步事件处理失败: %v", err)
	}
}

// Publish 异步发布事件
func (b *AsyncMemoryEventBus) Publish(ctx context.Context, evt event.DomainEvent) error {
	select {
	case b.eventChan <- eventWrapper{ctx: ctx, event: evt}:
		log.Printf("📤 [EventBus] 事件已入队: %s", evt.EventName())
		return nil
	default:
		log.Printf("⚠️ [EventBus] 事件队列已满，同步处理: %s", evt.EventName())
		// 队列满时降级为同步处理
		return b.MemoryEventBus.Publish(ctx, evt)
	}
}

// Close 关闭事件总线
func (b *AsyncMemoryEventBus) Close() {
	close(b.eventChan)
	b.wg.Wait()
	log.Println("🛑 [EventBus] 异步事件总线已关闭")
}

