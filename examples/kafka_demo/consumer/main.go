package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gozero-ddd/internal/domain/event"
	"gozero-ddd/internal/infrastructure/eventbus"
)

/*
	Kafka 领域事件消费者示例

	演示如何使用 Kafka 消费 DDD 领域事件

	运行步骤：
	1. 启动 Kafka: docker-compose up -d
	2. 运行消费者: go run examples/kafka_demo/consumer/main.go
	3. 运行生产者: go run examples/kafka_demo/producer/main.go (另一个终端)
*/

func main() {
	log.Println("🚀 启动 Kafka 领域事件消费者示例...")

	// 创建 Kafka 配置
	config := eventbus.KafkaConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "domain-events",
		GroupID:     "knowledge-consumer",
		ReadTimeout: 10 * time.Second,
	}

	// 创建 Kafka 事件消费器
	consumer := eventbus.NewKafkaEventConsumer(config)

	// 注册事件处理器
	registerHandlers(consumer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动消费者
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("❌ 启动消费者失败: %v", err)
	}

	log.Println("✅ 消费者已启动，等待事件...")
	log.Println("   按 Ctrl+C 退出")
	log.Println("---")

	// 监听退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("\n📥 收到退出信号，停止消费...")
	if err := consumer.Stop(); err != nil {
		log.Printf("❌ 停止消费者失败: %v", err)
	}

	log.Println("✅ 消费者示例结束")
}

// registerHandlers 注册事件处理器
func registerHandlers(consumer *eventbus.KafkaEventConsumer) {
	// 注册知识库创建事件处理器
	consumer.Subscribe("knowledge_base.created", &KnowledgeBaseCreatedHandler{})
	
	// 注册知识库更新事件处理器
	consumer.Subscribe("knowledge_base.updated", &KnowledgeBaseUpdatedHandler{})
	
	// 注册文档添加事件处理器
	consumer.Subscribe("document.added", &DocumentAddedHandler{})
	
	// 注册文档删除事件处理器
	consumer.Subscribe("document.removed", &DocumentRemovedHandler{})
	
	// 注册全局审计日志处理器
	consumer.SubscribeAll(&AuditLogHandler{})
	
	// 注册搜索索引更新处理器（演示同一事件多个处理器）
	consumer.Subscribe("document.added", &SearchIndexHandler{})

	log.Println("📫 已注册事件处理器:")
	log.Println("   - KnowledgeBaseCreatedHandler")
	log.Println("   - KnowledgeBaseUpdatedHandler")
	log.Println("   - DocumentAddedHandler")
	log.Println("   - DocumentRemovedHandler")
	log.Println("   - AuditLogHandler (全局)")
	log.Println("   - SearchIndexHandler")
}

// ==================== 事件处理器实现 ====================

// KnowledgeBaseCreatedHandler 知识库创建事件处理器
type KnowledgeBaseCreatedHandler struct{}

func (h *KnowledgeBaseCreatedHandler) EventName() string {
	return "knowledge_base.created"
}

func (h *KnowledgeBaseCreatedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	log.Printf("🎯 [KnowledgeBaseCreatedHandler] 处理事件")
	
	// 从包装事件中提取原始数据
	if wrapped, ok := evt.(*eventbus.WrappedDomainEvent); ok {
		var data struct {
			KnowledgeBaseID string `json:"KnowledgeBaseID"`
			Name            string `json:"Name"`
		}
		if err := json.Unmarshal(wrapped.Payload(), &data); err == nil {
			log.Printf("   📦 知识库ID: %s", data.KnowledgeBaseID)
			log.Printf("   📝 名称: %s", data.Name)
		}
	}
	
	// 模拟业务处理
	log.Println("   ✅ 发送欢迎邮件通知...")
	log.Println("   ✅ 初始化默认设置...")
	
	return nil
}

// KnowledgeBaseUpdatedHandler 知识库更新事件处理器
type KnowledgeBaseUpdatedHandler struct{}

func (h *KnowledgeBaseUpdatedHandler) EventName() string {
	return "knowledge_base.updated"
}

func (h *KnowledgeBaseUpdatedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	log.Printf("🎯 [KnowledgeBaseUpdatedHandler] 处理事件")
	
	if wrapped, ok := evt.(*eventbus.WrappedDomainEvent); ok {
		var data struct {
			KnowledgeBaseID string `json:"KnowledgeBaseID"`
			OldName         string `json:"OldName"`
			NewName         string `json:"NewName"`
		}
		if err := json.Unmarshal(wrapped.Payload(), &data); err == nil {
			log.Printf("   📦 知识库ID: %s", data.KnowledgeBaseID)
			log.Printf("   📝 名称变更: %s -> %s", data.OldName, data.NewName)
		}
	}
	
	log.Println("   ✅ 清除相关缓存...")
	log.Println("   ✅ 发送变更通知...")
	
	return nil
}

// DocumentAddedHandler 文档添加事件处理器
type DocumentAddedHandler struct{}

func (h *DocumentAddedHandler) EventName() string {
	return "document.added"
}

func (h *DocumentAddedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	log.Printf("🎯 [DocumentAddedHandler] 处理事件")
	
	if wrapped, ok := evt.(*eventbus.WrappedDomainEvent); ok {
		var data struct {
			DocumentID      string `json:"DocumentID"`
			KnowledgeBaseID string `json:"KnowledgeBaseID"`
			Title           string `json:"Title"`
		}
		if err := json.Unmarshal(wrapped.Payload(), &data); err == nil {
			log.Printf("   📄 文档ID: %s", data.DocumentID)
			log.Printf("   📦 所属知识库: %s", data.KnowledgeBaseID)
			log.Printf("   📝 标题: %s", data.Title)
		}
	}
	
	log.Println("   ✅ 触发 AI 内容分析...")
	log.Println("   ✅ 生成文档摘要...")
	
	return nil
}

// DocumentRemovedHandler 文档删除事件处理器
type DocumentRemovedHandler struct{}

func (h *DocumentRemovedHandler) EventName() string {
	return "document.removed"
}

func (h *DocumentRemovedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	log.Printf("🎯 [DocumentRemovedHandler] 处理事件")
	
	if wrapped, ok := evt.(*eventbus.WrappedDomainEvent); ok {
		var data struct {
			DocumentID      string `json:"DocumentID"`
			KnowledgeBaseID string `json:"KnowledgeBaseID"`
		}
		if err := json.Unmarshal(wrapped.Payload(), &data); err == nil {
			log.Printf("   📄 文档ID: %s", data.DocumentID)
			log.Printf("   📦 所属知识库: %s", data.KnowledgeBaseID)
		}
	}
	
	log.Println("   ✅ 清理相关资源...")
	log.Println("   ✅ 更新统计信息...")
	
	return nil
}

// AuditLogHandler 审计日志处理器（全局）
type AuditLogHandler struct{}

func (h *AuditLogHandler) EventName() string {
	return "" // 空字符串表示处理所有事件
}

func (h *AuditLogHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	log.Printf("📋 [AuditLog] EventID=%s EventName=%s AggregateID=%s OccurredAt=%s",
		evt.EventID(),
		evt.EventName(),
		evt.AggregateID(),
		evt.OccurredAt().Format("2006-01-02 15:04:05"),
	)
	
	// 模拟写入审计日志数据库
	// 实际项目中可以写入 audit_logs 表
	
	return nil
}

// SearchIndexHandler 搜索索引处理器
// 演示同一事件可以有多个处理器
type SearchIndexHandler struct{}

func (h *SearchIndexHandler) EventName() string {
	return "document.added"
}

func (h *SearchIndexHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	log.Printf("🔍 [SearchIndexHandler] 更新搜索索引")
	
	if wrapped, ok := evt.(*eventbus.WrappedDomainEvent); ok {
		var data struct {
			DocumentID string `json:"DocumentID"`
			Title      string `json:"Title"`
		}
		if err := json.Unmarshal(wrapped.Payload(), &data); err == nil {
			log.Printf("   📄 索引文档: %s - %s", data.DocumentID, data.Title)
		}
	}
	
	// 模拟更新 Elasticsearch 索引
	log.Println("   ✅ 已添加到 Elasticsearch 索引")
	
	return nil
}

