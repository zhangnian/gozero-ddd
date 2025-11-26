package eventhandler

import (
	"context"
	"log"

	"gozero-ddd/internal/domain/event"
)

// KnowledgeBaseCreatedHandler 知识库创建事件处理器
// 示例：当知识库创建后，可以进行一些后续操作
// 例如：发送通知、初始化缓存、记录审计日志等
type KnowledgeBaseCreatedHandler struct{}

// NewKnowledgeBaseCreatedHandler 创建处理器
func NewKnowledgeBaseCreatedHandler() *KnowledgeBaseCreatedHandler {
	return &KnowledgeBaseCreatedHandler{}
}

// 确保实现了接口
var _ event.EventHandler = (*KnowledgeBaseCreatedHandler)(nil)

// EventName 返回处理的事件名称
func (h *KnowledgeBaseCreatedHandler) EventName() string {
	return "knowledge_base.created"
}

// Handle 处理知识库创建事件
func (h *KnowledgeBaseCreatedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	// 类型断言，获取具体事件
	e, ok := evt.(*event.KnowledgeBaseCreatedEvent)
	if !ok {
		return nil
	}

	log.Printf("📝 [EventHandler] 处理知识库创建事件: EventID=%s, KnowledgeBaseID=%s, Name=%s",
		e.EventID(), e.KnowledgeBaseID, e.Name)

	// 这里可以执行后续操作：
	// 1. 发送欢迎通知
	// 2. 初始化默认文档
	// 3. 记录审计日志
	// 4. 更新搜索索引
	// 5. 发送 Webhook

	return nil
}

// KnowledgeBaseUpdatedHandler 知识库更新事件处理器
type KnowledgeBaseUpdatedHandler struct{}

// NewKnowledgeBaseUpdatedHandler 创建处理器
func NewKnowledgeBaseUpdatedHandler() *KnowledgeBaseUpdatedHandler {
	return &KnowledgeBaseUpdatedHandler{}
}

// 确保实现了接口
var _ event.EventHandler = (*KnowledgeBaseUpdatedHandler)(nil)

// EventName 返回处理的事件名称
func (h *KnowledgeBaseUpdatedHandler) EventName() string {
	return "knowledge_base.updated"
}

// Handle 处理知识库更新事件
func (h *KnowledgeBaseUpdatedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(*event.KnowledgeBaseUpdatedEvent)
	if !ok {
		return nil
	}

	log.Printf("📝 [EventHandler] 处理知识库更新事件: KnowledgeBaseID=%s, OldName=%s -> NewName=%s",
		e.KnowledgeBaseID, e.OldName, e.NewName)

	// 这里可以执行后续操作：
	// 1. 清除缓存
	// 2. 更新搜索索引
	// 3. 发送变更通知

	return nil
}

// DocumentAddedHandler 文档添加事件处理器
// 示例：当文档添加后，更新搜索索引
type DocumentAddedHandler struct{}

// NewDocumentAddedHandler 创建处理器
func NewDocumentAddedHandler() *DocumentAddedHandler {
	return &DocumentAddedHandler{}
}

// 确保实现了接口
var _ event.EventHandler = (*DocumentAddedHandler)(nil)

// EventName 返回处理的事件名称
func (h *DocumentAddedHandler) EventName() string {
	return "document.added"
}

// Handle 处理文档添加事件
func (h *DocumentAddedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(*event.DocumentAddedEvent)
	if !ok {
		return nil
	}

	log.Printf("📝 [EventHandler] 处理文档添加事件: DocID=%s, KnowledgeBaseID=%s, Title=%s",
		e.DocumentID, e.KnowledgeBaseID, e.Title)

	// 这里可以执行后续操作：
	// 1. 更新全文搜索索引（如 Elasticsearch）
	// 2. 生成文档摘要
	// 3. 触发 AI 分析
	// 4. 发送通知

	return nil
}

// DocumentRemovedHandler 文档删除事件处理器
type DocumentRemovedHandler struct{}

// NewDocumentRemovedHandler 创建处理器
func NewDocumentRemovedHandler() *DocumentRemovedHandler {
	return &DocumentRemovedHandler{}
}

// 确保实现了接口
var _ event.EventHandler = (*DocumentRemovedHandler)(nil)

// EventName 返回处理的事件名称
func (h *DocumentRemovedHandler) EventName() string {
	return "document.removed"
}

// Handle 处理文档删除事件
func (h *DocumentRemovedHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(*event.DocumentRemovedEvent)
	if !ok {
		return nil
	}

	log.Printf("📝 [EventHandler] 处理文档删除事件: DocID=%s, KnowledgeBaseID=%s",
		e.DocumentID, e.KnowledgeBaseID)

	// 这里可以执行后续操作：
	// 1. 从搜索索引中删除
	// 2. 清理相关缓存
	// 3. 记录审计日志

	return nil
}

// AuditLogHandler 审计日志处理器
// 示例：记录所有领域事件到审计日志
// 这是一个"全局处理器"，处理所有事件
type AuditLogHandler struct{}

// NewAuditLogHandler 创建审计日志处理器
func NewAuditLogHandler() *AuditLogHandler {
	return &AuditLogHandler{}
}

// 确保实现了接口
var _ event.EventHandler = (*AuditLogHandler)(nil)

// EventName 返回空字符串，表示处理所有事件
func (h *AuditLogHandler) EventName() string {
	return "" // 空字符串表示处理所有事件
}

// Handle 记录审计日志
func (h *AuditLogHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	log.Printf("📋 [AuditLog] EventID=%s, EventName=%s, AggregateID=%s, OccurredAt=%s",
		evt.EventID(), evt.EventName(), evt.AggregateID(), evt.OccurredAt())

	// 在实际项目中，可以将审计日志：
	// 1. 写入数据库（创建 audit_logs 表）
	// 2. 发送到日志服务（如 ELK、Loki）
	// 3. 写入消息队列供后续分析
	// 4. 用于事件溯源（Event Sourcing）

	return nil
}

