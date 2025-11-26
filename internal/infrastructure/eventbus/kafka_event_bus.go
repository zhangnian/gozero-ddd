package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"

	"gozero-ddd/internal/domain/event"
)

// KafkaConfig Kafka 配置
type KafkaConfig struct {
	Brokers       []string      // Kafka broker 地址列表
	Topic         string        // 事件主题
	GroupID       string        // 消费者组ID
	WriteTimeout  time.Duration // 写入超时
	ReadTimeout   time.Duration // 读取超时
	BatchSize     int           // 批量发送大小
	BatchTimeout  time.Duration // 批量发送超时
	RequiredAcks  int           // 确认模式: -1=all, 0=none, 1=leader
	Async         bool          // 是否异步发送
	AutoCreateTopic bool        // 是否自动创建主题
}

// DefaultKafkaConfig 默认配置
func DefaultKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers:       []string{"localhost:9092"},
		Topic:         "domain-events",
		GroupID:       "knowledge-service",
		WriteTimeout:  10 * time.Second,
		ReadTimeout:   10 * time.Second,
		BatchSize:     100,
		BatchTimeout:  time.Second,
		RequiredAcks:  -1, // 等待所有副本确认
		Async:         false,
		AutoCreateTopic: true,
	}
}

// ==================== 事件消息结构 ====================

// EventMessage Kafka 消息结构
// 用于序列化/反序列化领域事件
type EventMessage struct {
	EventID     string          `json:"event_id"`
	EventName   string          `json:"event_name"`
	AggregateID string          `json:"aggregate_id"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Payload     json.RawMessage `json:"payload"` // 事件具体数据
	Metadata    EventMetadata   `json:"metadata"`
}

// EventMetadata 事件元数据
type EventMetadata struct {
	TraceID     string `json:"trace_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	Version     string `json:"version,omitempty"`
}

// ==================== Kafka 事件发布器 ====================

// KafkaEventPublisher Kafka 事件发布器
// 负责将领域事件发布到 Kafka
type KafkaEventPublisher struct {
	writer *kafka.Writer
	config KafkaConfig
	mu     sync.Mutex
}

// NewKafkaEventPublisher 创建 Kafka 事件发布器
func NewKafkaEventPublisher(config KafkaConfig) (*KafkaEventPublisher, error) {
	// 自动创建 topic（如果启用）
	if config.AutoCreateTopic {
		if err := createTopicIfNotExists(config); err != nil {
			log.Printf("⚠️ [Kafka] 自动创建主题失败: %v (可能主题已存在)", err)
		}
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{}, // 负载均衡策略
		WriteTimeout: config.WriteTimeout,
		BatchSize:    config.BatchSize,
		BatchTimeout: config.BatchTimeout,
		RequiredAcks: kafka.RequiredAcks(config.RequiredAcks),
		Async:        config.Async,
	}

	log.Printf("📤 [Kafka] 事件发布器已创建: brokers=%v, topic=%s", config.Brokers, config.Topic)

	return &KafkaEventPublisher{
		writer: writer,
		config: config,
	}, nil
}

// createTopicIfNotExists 如果主题不存在则创建
func createTopicIfNotExists(config KafkaConfig) error {
	conn, err := kafka.Dial("tcp", config.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", controller.Host+":"+string(rune(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             config.Topic,
			NumPartitions:     3,  // 分区数
			ReplicationFactor: 1,  // 副本因子（单机环境设为1）
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}

// 确保实现了接口
var _ event.EventPublisher = (*KafkaEventPublisher)(nil)

// Publish 发布单个事件到 Kafka
func (p *KafkaEventPublisher) Publish(ctx context.Context, evt event.DomainEvent) error {
	return p.PublishAll(ctx, []event.DomainEvent{evt})
}

// PublishAll 批量发布事件到 Kafka
func (p *KafkaEventPublisher) PublishAll(ctx context.Context, events []event.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	messages := make([]kafka.Message, 0, len(events))
	
	for _, evt := range events {
		msg, err := p.buildMessage(ctx, evt)
		if err != nil {
			log.Printf("❌ [Kafka] 构建消息失败: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		log.Printf("❌ [Kafka] 发布事件失败: %v", err)
		return err
	}

	log.Printf("📤 [Kafka] 成功发布 %d 个事件到主题 %s", len(messages), p.config.Topic)
	return nil
}

// buildMessage 构建 Kafka 消息
func (p *KafkaEventPublisher) buildMessage(ctx context.Context, evt event.DomainEvent) (kafka.Message, error) {
	// 序列化事件数据
	payload, err := json.Marshal(evt)
	if err != nil {
		return kafka.Message{}, err
	}

	// 构建事件消息
	eventMsg := EventMessage{
		EventID:     evt.EventID(),
		EventName:   evt.EventName(),
		AggregateID: evt.AggregateID(),
		OccurredAt:  evt.OccurredAt(),
		Payload:     payload,
		Metadata: EventMetadata{
			ServiceName: "knowledge-service",
			Version:     "1.0",
		},
	}

	// 从 context 提取 trace_id（如果有）
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		eventMsg.Metadata.TraceID = traceID
	}

	value, err := json.Marshal(eventMsg)
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{
		Key:   []byte(evt.AggregateID()), // 使用聚合根ID作为分区键，保证同一聚合的事件有序
		Value: value,
		Headers: []kafka.Header{
			{Key: "event_name", Value: []byte(evt.EventName())},
			{Key: "event_id", Value: []byte(evt.EventID())},
		},
	}, nil
}

// Close 关闭发布器
func (p *KafkaEventPublisher) Close() error {
	log.Println("🛑 [Kafka] 关闭事件发布器")
	return p.writer.Close()
}

// ==================== Kafka 事件消费器 ====================

// KafkaEventConsumer Kafka 事件消费器
// 负责从 Kafka 消费领域事件并分发给处理器
type KafkaEventConsumer struct {
	reader   *kafka.Reader
	config   KafkaConfig
	handlers map[string][]event.EventHandler // 事件处理器映射
	allHandlers []event.EventHandler          // 全局处理器
	mu       sync.RWMutex
	running  bool
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewKafkaEventConsumer 创建 Kafka 事件消费器
func NewKafkaEventConsumer(config KafkaConfig) *KafkaEventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.Topic,
		GroupID:        config.GroupID,
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		MaxWait:        time.Second, // 最大等待时间
		StartOffset:    kafka.FirstOffset,
		CommitInterval: time.Second, // 自动提交间隔
	})

	log.Printf("📥 [Kafka] 事件消费器已创建: brokers=%v, topic=%s, group=%s",
		config.Brokers, config.Topic, config.GroupID)

	return &KafkaEventConsumer{
		reader:      reader,
		config:      config,
		handlers:    make(map[string][]event.EventHandler),
		allHandlers: make([]event.EventHandler, 0),
		stopCh:      make(chan struct{}),
	}
}

// 确保实现了接口
var _ event.EventSubscriber = (*KafkaEventConsumer)(nil)

// Subscribe 订阅特定事件
func (c *KafkaEventConsumer) Subscribe(eventName string, handler event.EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers[eventName] = append(c.handlers[eventName], handler)
	log.Printf("📫 [Kafka] 注册事件处理器: %s", eventName)
}

// SubscribeAll 订阅所有事件
func (c *KafkaEventConsumer) SubscribeAll(handler event.EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.allHandlers = append(c.allHandlers, handler)
	log.Printf("📫 [Kafka] 注册全局事件处理器")
}

// Start 启动消费者
func (c *KafkaEventConsumer) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return errors.New("消费者已在运行")
	}
	c.running = true
	c.mu.Unlock()

	log.Println("🚀 [Kafka] 启动事件消费者...")

	c.wg.Add(1)
	go c.consumeLoop(ctx)

	return nil
}

// consumeLoop 消费循环
func (c *KafkaEventConsumer) consumeLoop(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 [Kafka] 收到停止信号，退出消费循环")
			return
		case <-c.stopCh:
			log.Println("🛑 [Kafka] 收到停止信号，退出消费循环")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Printf("❌ [Kafka] 读取消息失败: %v", err)
				continue
			}

			c.handleMessage(ctx, msg)
		}
	}
}

// handleMessage 处理 Kafka 消息
func (c *KafkaEventConsumer) handleMessage(ctx context.Context, msg kafka.Message) {
	var eventMsg EventMessage
	if err := json.Unmarshal(msg.Value, &eventMsg); err != nil {
		log.Printf("❌ [Kafka] 解析消息失败: %v", err)
		return
	}

	log.Printf("📥 [Kafka] 收到事件: %s, EventID=%s, AggregateID=%s",
		eventMsg.EventName, eventMsg.EventID, eventMsg.AggregateID)

	// 创建包装的事件对象
	wrappedEvent := &WrappedDomainEvent{
		eventMsg: eventMsg,
	}

	// 调用处理器
	c.dispatchEvent(ctx, eventMsg.EventName, wrappedEvent)
}

// dispatchEvent 分发事件给处理器
func (c *KafkaEventConsumer) dispatchEvent(ctx context.Context, eventName string, evt event.DomainEvent) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 调用特定事件处理器
	if handlers, ok := c.handlers[eventName]; ok {
		for _, handler := range handlers {
			if err := handler.Handle(ctx, evt); err != nil {
				log.Printf("❌ [Kafka] 事件处理失败: %s, 错误: %v", eventName, err)
			}
		}
	}

	// 调用全局处理器
	for _, handler := range c.allHandlers {
		if err := handler.Handle(ctx, evt); err != nil {
			log.Printf("❌ [Kafka] 全局处理器执行失败: %v", err)
		}
	}
}

// Stop 停止消费者
func (c *KafkaEventConsumer) Stop() error {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = false
	c.mu.Unlock()

	close(c.stopCh)
	c.wg.Wait()

	log.Println("🛑 [Kafka] 关闭事件消费器")
	return c.reader.Close()
}

// ==================== 包装事件 ====================

// WrappedDomainEvent 包装的领域事件
// 用于从 Kafka 消息反序列化后的事件
type WrappedDomainEvent struct {
	eventMsg EventMessage
}

func (e *WrappedDomainEvent) EventID() string {
	return e.eventMsg.EventID
}

func (e *WrappedDomainEvent) EventName() string {
	return e.eventMsg.EventName
}

func (e *WrappedDomainEvent) OccurredAt() time.Time {
	return e.eventMsg.OccurredAt
}

func (e *WrappedDomainEvent) AggregateID() string {
	return e.eventMsg.AggregateID
}

// Payload 获取原始事件数据
func (e *WrappedDomainEvent) Payload() json.RawMessage {
	return e.eventMsg.Payload
}

// Metadata 获取事件元数据
func (e *WrappedDomainEvent) Metadata() EventMetadata {
	return e.eventMsg.Metadata
}

// ==================== Kafka 事件总线 ====================

// KafkaEventBus Kafka 事件总线
// 组合了发布器和消费器，实现完整的事件总线功能
type KafkaEventBus struct {
	publisher *KafkaEventPublisher
	consumer  *KafkaEventConsumer
}

// NewKafkaEventBus 创建 Kafka 事件总线
func NewKafkaEventBus(config KafkaConfig) (*KafkaEventBus, error) {
	publisher, err := NewKafkaEventPublisher(config)
	if err != nil {
		return nil, err
	}

	consumer := NewKafkaEventConsumer(config)

	return &KafkaEventBus{
		publisher: publisher,
		consumer:  consumer,
	}, nil
}

// 确保实现了接口
var _ event.EventBus = (*KafkaEventBus)(nil)

// Publish 发布事件
func (b *KafkaEventBus) Publish(ctx context.Context, evt event.DomainEvent) error {
	return b.publisher.Publish(ctx, evt)
}

// PublishAll 批量发布事件
func (b *KafkaEventBus) PublishAll(ctx context.Context, events []event.DomainEvent) error {
	return b.publisher.PublishAll(ctx, events)
}

// Subscribe 订阅特定事件
func (b *KafkaEventBus) Subscribe(eventName string, handler event.EventHandler) {
	b.consumer.Subscribe(eventName, handler)
}

// SubscribeAll 订阅所有事件
func (b *KafkaEventBus) SubscribeAll(handler event.EventHandler) {
	b.consumer.SubscribeAll(handler)
}

// Start 启动事件总线（开始消费事件）
func (b *KafkaEventBus) Start(ctx context.Context) error {
	return b.consumer.Start(ctx)
}

// Close 关闭事件总线
func (b *KafkaEventBus) Close() error {
	if err := b.consumer.Stop(); err != nil {
		log.Printf("❌ [Kafka] 关闭消费器失败: %v", err)
	}
	return b.publisher.Close()
}

