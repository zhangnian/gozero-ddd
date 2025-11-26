package event

import (
	"time"

	"github.com/google/uuid"

	"gozero-ddd/internal/domain/valueobject"
)

// DomainEvent 领域事件接口
// 领域事件用于记录领域中发生的重要事件
// 领域事件的特点：
// 1. 不可变：一旦创建就不能修改
// 2. 时间戳：记录事件发生的时间
// 3. 唯一标识：每个事件都有唯一 ID
type DomainEvent interface {
	EventID() string    // 事件唯一标识
	EventName() string  // 事件名称
	OccurredAt() time.Time // 发生时间
	AggregateID() string   // 聚合根ID
}

// BaseEvent 基础事件
// 所有领域事件的基类，提供通用字段
type BaseEvent struct {
	eventID     string
	occurredAt  time.Time
	aggregateID string
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(aggregateID string) BaseEvent {
	return BaseEvent{
		eventID:     uuid.New().String(),
		occurredAt:  time.Now(),
		aggregateID: aggregateID,
	}
}

func (e BaseEvent) EventID() string {
	return e.eventID
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

// ==================== 知识库相关事件 ====================

// KnowledgeBaseCreatedEvent 知识库创建事件
// 当新的知识库被创建时触发
type KnowledgeBaseCreatedEvent struct {
	BaseEvent
	KnowledgeBaseID valueobject.KnowledgeBaseID
	Name            string
	Description     string
}

func NewKnowledgeBaseCreatedEvent(id valueobject.KnowledgeBaseID, name string) *KnowledgeBaseCreatedEvent {
	return &KnowledgeBaseCreatedEvent{
		BaseEvent:       NewBaseEvent(id.String()),
		KnowledgeBaseID: id,
		Name:            name,
	}
}

func (e *KnowledgeBaseCreatedEvent) EventName() string {
	return "knowledge_base.created"
}

// KnowledgeBaseUpdatedEvent 知识库更新事件
type KnowledgeBaseUpdatedEvent struct {
	BaseEvent
	KnowledgeBaseID valueobject.KnowledgeBaseID
	OldName         string
	NewName         string
	OldDescription  string
	NewDescription  string
}

func NewKnowledgeBaseUpdatedEvent(
	id valueobject.KnowledgeBaseID,
	oldName, newName string,
	oldDesc, newDesc string,
) *KnowledgeBaseUpdatedEvent {
	return &KnowledgeBaseUpdatedEvent{
		BaseEvent:       NewBaseEvent(id.String()),
		KnowledgeBaseID: id,
		OldName:         oldName,
		NewName:         newName,
		OldDescription:  oldDesc,
		NewDescription:  newDesc,
	}
}

func (e *KnowledgeBaseUpdatedEvent) EventName() string {
	return "knowledge_base.updated"
}

// KnowledgeBaseDeletedEvent 知识库删除事件
type KnowledgeBaseDeletedEvent struct {
	BaseEvent
	KnowledgeBaseID valueobject.KnowledgeBaseID
	Name            string
}

func NewKnowledgeBaseDeletedEvent(id valueobject.KnowledgeBaseID, name string) *KnowledgeBaseDeletedEvent {
	return &KnowledgeBaseDeletedEvent{
		BaseEvent:       NewBaseEvent(id.String()),
		KnowledgeBaseID: id,
		Name:            name,
	}
}

func (e *KnowledgeBaseDeletedEvent) EventName() string {
	return "knowledge_base.deleted"
}

// ==================== 文档相关事件 ====================

// DocumentAddedEvent 文档添加事件
// 当文档被添加到知识库时触发
type DocumentAddedEvent struct {
	BaseEvent
	DocumentID      valueobject.DocumentID
	KnowledgeBaseID valueobject.KnowledgeBaseID
	Title           string
	Tags            []string
}

func NewDocumentAddedEvent(docID valueobject.DocumentID, kbID valueobject.KnowledgeBaseID, title string) *DocumentAddedEvent {
	return &DocumentAddedEvent{
		BaseEvent:       NewBaseEvent(kbID.String()), // 聚合根是知识库
		DocumentID:      docID,
		KnowledgeBaseID: kbID,
		Title:           title,
	}
}

func (e *DocumentAddedEvent) EventName() string {
	return "document.added"
}

// DocumentRemovedEvent 文档删除事件
// 当文档从知识库中移除时触发
type DocumentRemovedEvent struct {
	BaseEvent
	DocumentID      valueobject.DocumentID
	KnowledgeBaseID valueobject.KnowledgeBaseID
	Title           string // 删除前的标题，用于日志记录
}

func NewDocumentRemovedEvent(docID valueobject.DocumentID, kbID valueobject.KnowledgeBaseID) *DocumentRemovedEvent {
	return &DocumentRemovedEvent{
		BaseEvent:       NewBaseEvent(kbID.String()),
		DocumentID:      docID,
		KnowledgeBaseID: kbID,
	}
}

func (e *DocumentRemovedEvent) EventName() string {
	return "document.removed"
}

// DocumentUpdatedEvent 文档更新事件
type DocumentUpdatedEvent struct {
	BaseEvent
	DocumentID      valueobject.DocumentID
	KnowledgeBaseID valueobject.KnowledgeBaseID
	OldTitle        string
	NewTitle        string
}

func NewDocumentUpdatedEvent(
	docID valueobject.DocumentID,
	kbID valueobject.KnowledgeBaseID,
	oldTitle, newTitle string,
) *DocumentUpdatedEvent {
	return &DocumentUpdatedEvent{
		BaseEvent:       NewBaseEvent(kbID.String()),
		DocumentID:      docID,
		KnowledgeBaseID: kbID,
		OldTitle:        oldTitle,
		NewTitle:        newTitle,
	}
}

func (e *DocumentUpdatedEvent) EventName() string {
	return "document.updated"
}

