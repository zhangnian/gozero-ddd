package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"gozero-ddd/internal/domain/event"
	"gozero-ddd/internal/domain/valueobject"
	"gozero-ddd/internal/infrastructure/eventbus"
)

/*
	Kafka 领域事件生产者示例

	演示如何使用 Kafka 发布 DDD 领域事件

	运行步骤：
	1. 启动 Kafka: docker-compose up -d
	2. 运行生产者: go run examples/kafka_demo/producer/main.go
	3. 运行消费者: go run examples/kafka_demo/consumer/main.go (另一个终端)
*/

func main() {
	log.Println("🚀 启动 Kafka 领域事件生产者示例...")

	// 创建 Kafka 配置
	config := eventbus.KafkaConfig{
		Brokers:         []string{"localhost:9092"},
		Topic:           "domain-events",
		GroupID:         "knowledge-producer",
		WriteTimeout:    10 * time.Second,
		BatchSize:       1,              // 示例中设为1，立即发送
		BatchTimeout:    time.Millisecond * 100,
		RequiredAcks:    -1,             // 等待所有副本确认
		Async:           false,          // 同步发送，便于观察
		AutoCreateTopic: true,
	}

	// 创建 Kafka 事件发布器
	publisher, err := eventbus.NewKafkaEventPublisher(config)
	if err != nil {
		log.Fatalf("❌ 创建 Kafka 发布器失败: %v", err)
	}
	defer publisher.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("📤 收到退出信号，停止发送...")
		cancel()
	}()

	// 模拟发布领域事件
	simulateEvents(ctx, publisher)

	log.Println("✅ 生产者示例结束")
}

// simulateEvents 模拟发布领域事件
func simulateEvents(ctx context.Context, publisher *eventbus.KafkaEventPublisher) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	eventCount := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			eventCount++
			// 轮流发布不同类型的事件
			switch eventCount % 4 {
			case 1:
				publishKnowledgeBaseCreated(ctx, publisher)
			case 2:
				publishKnowledgeBaseUpdated(ctx, publisher)
			case 3:
				publishDocumentAdded(ctx, publisher)
			case 0:
				publishDocumentRemoved(ctx, publisher)
			}
		}
	}
}

// publishKnowledgeBaseCreated 发布知识库创建事件
func publishKnowledgeBaseCreated(ctx context.Context, publisher *eventbus.KafkaEventPublisher) {
	kbID := valueobject.NewKnowledgeBaseID()
	name := "技术文档库-" + uuid.New().String()[:8]

	evt := event.NewKnowledgeBaseCreatedEvent(kbID, name)
	
	log.Printf("📤 发布事件: %s", evt.EventName())
	log.Printf("   KnowledgeBaseID: %s", kbID.String())
	log.Printf("   Name: %s", name)

	if err := publisher.Publish(ctx, evt); err != nil {
		log.Printf("❌ 发布失败: %v", err)
	} else {
		log.Println("   ✅ 发布成功!")
	}
	log.Println("---")
}

// publishKnowledgeBaseUpdated 发布知识库更新事件
func publishKnowledgeBaseUpdated(ctx context.Context, publisher *eventbus.KafkaEventPublisher) {
	kbID := valueobject.NewKnowledgeBaseID()
	oldName := "旧名称-" + uuid.New().String()[:8]
	newName := "新名称-" + uuid.New().String()[:8]

	evt := event.NewKnowledgeBaseUpdatedEvent(kbID, oldName, newName, "", "")
	
	log.Printf("📤 发布事件: %s", evt.EventName())
	log.Printf("   KnowledgeBaseID: %s", kbID.String())
	log.Printf("   OldName: %s -> NewName: %s", oldName, newName)

	if err := publisher.Publish(ctx, evt); err != nil {
		log.Printf("❌ 发布失败: %v", err)
	} else {
		log.Println("   ✅ 发布成功!")
	}
	log.Println("---")
}

// publishDocumentAdded 发布文档添加事件
func publishDocumentAdded(ctx context.Context, publisher *eventbus.KafkaEventPublisher) {
	docID := valueobject.NewDocumentID()
	kbID := valueobject.NewKnowledgeBaseID()
	title := "Go语言最佳实践-" + uuid.New().String()[:8]

	evt := event.NewDocumentAddedEvent(docID, kbID, title)
	
	log.Printf("📤 发布事件: %s", evt.EventName())
	log.Printf("   DocumentID: %s", docID.String())
	log.Printf("   KnowledgeBaseID: %s", kbID.String())
	log.Printf("   Title: %s", title)

	if err := publisher.Publish(ctx, evt); err != nil {
		log.Printf("❌ 发布失败: %v", err)
	} else {
		log.Println("   ✅ 发布成功!")
	}
	log.Println("---")
}

// publishDocumentRemoved 发布文档删除事件
func publishDocumentRemoved(ctx context.Context, publisher *eventbus.KafkaEventPublisher) {
	docID := valueobject.NewDocumentID()
	kbID := valueobject.NewKnowledgeBaseID()

	evt := event.NewDocumentRemovedEvent(docID, kbID)
	
	log.Printf("📤 发布事件: %s", evt.EventName())
	log.Printf("   DocumentID: %s", docID.String())
	log.Printf("   KnowledgeBaseID: %s", kbID.String())

	if err := publisher.Publish(ctx, evt); err != nil {
		log.Printf("❌ 发布失败: %v", err)
	} else {
		log.Println("   ✅ 发布成功!")
	}
	log.Println("---")
}

