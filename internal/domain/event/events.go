package event

import (
	"time"

	"gozero-ddd/internal/domain/valueobject"
)

// DomainEvent 领域事件接口
// 领域事件用于记录领域中发生的重要事件
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

// BaseEvent 基础事件
type BaseEvent struct {
	occurredAt time.Time
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// KnowledgeBaseCreatedEvent 知识库创建事件
type KnowledgeBaseCreatedEvent struct {
	BaseEvent
	KnowledgeBaseID valueobject.KnowledgeBaseID
	Name            string
}

func NewKnowledgeBaseCreatedEvent(id valueobject.KnowledgeBaseID, name string) *KnowledgeBaseCreatedEvent {
	return &KnowledgeBaseCreatedEvent{
		BaseEvent:       BaseEvent{occurredAt: time.Now()},
		KnowledgeBaseID: id,
		Name:            name,
	}
}

func (e *KnowledgeBaseCreatedEvent) EventName() string {
	return "knowledge_base.created"
}

// DocumentAddedEvent 文档添加事件
type DocumentAddedEvent struct {
	BaseEvent
	DocumentID      valueobject.DocumentID
	KnowledgeBaseID valueobject.KnowledgeBaseID
	Title           string
}

func NewDocumentAddedEvent(docID valueobject.DocumentID, kbID valueobject.KnowledgeBaseID, title string) *DocumentAddedEvent {
	return &DocumentAddedEvent{
		BaseEvent:       BaseEvent{occurredAt: time.Now()},
		DocumentID:      docID,
		KnowledgeBaseID: kbID,
		Title:           title,
	}
}

func (e *DocumentAddedEvent) EventName() string {
	return "document.added"
}

// DocumentRemovedEvent 文档删除事件
type DocumentRemovedEvent struct {
	BaseEvent
	DocumentID      valueobject.DocumentID
	KnowledgeBaseID valueobject.KnowledgeBaseID
}

func NewDocumentRemovedEvent(docID valueobject.DocumentID, kbID valueobject.KnowledgeBaseID) *DocumentRemovedEvent {
	return &DocumentRemovedEvent{
		BaseEvent:       BaseEvent{occurredAt: time.Now()},
		DocumentID:      docID,
		KnowledgeBaseID: kbID,
	}
}

func (e *DocumentRemovedEvent) EventName() string {
	return "document.removed"
}

