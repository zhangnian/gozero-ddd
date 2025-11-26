package entity

import (
	"time"

	"gozero-ddd/internal/domain"
	"gozero-ddd/internal/domain/event"
	"gozero-ddd/internal/domain/valueobject"
)

// KnowledgeBase 知识库实体（聚合根）
// 作为聚合根，KnowledgeBase 负责管理其下所有的 Document 实体
// 外部不能直接操作 Document，必须通过 KnowledgeBase 进行
//
// 领域事件：聚合根负责收集领域事件，在应用层持久化成功后发布
// 这确保了事件与状态变更的一致性
type KnowledgeBase struct {
	id          valueobject.KnowledgeBaseID // 唯一标识
	name        string                      // 知识库名称
	description string                      // 描述
	documents   []*Document                 // 文档集合
	createdAt   time.Time                   // 创建时间
	updatedAt   time.Time                   // 更新时间

	// 领域事件收集器
	// 聚合根在业务操作时收集事件，由应用层负责发布
	events []event.DomainEvent
}

// NewKnowledgeBase 创建新的知识库
// 创建时会收集 KnowledgeBaseCreatedEvent 事件
func NewKnowledgeBase(name, description string) (*KnowledgeBase, error) {
	if name == "" {
		return nil, domain.ErrKnowledgeBaseNameEmpty
	}

	now := time.Now()
	id := valueobject.NewKnowledgeBaseID()

	kb := &KnowledgeBase{
		id:          id,
		name:        name,
		description: description,
		documents:   make([]*Document, 0),
		createdAt:   now,
		updatedAt:   now,
		events:      make([]event.DomainEvent, 0),
	}

	// 收集创建事件
	kb.addEvent(event.NewKnowledgeBaseCreatedEvent(id, name))

	return kb, nil
}

// ReconstructKnowledgeBase 从持久化数据重建知识库实体
// 用于仓储层从数据库加载数据时使用
// 注意：重建不会产生领域事件（因为这不是新的业务操作）
func ReconstructKnowledgeBase(
	id valueobject.KnowledgeBaseID,
	name, description string,
	documents []*Document,
	createdAt, updatedAt time.Time,
) *KnowledgeBase {
	return &KnowledgeBase{
		id:          id,
		name:        name,
		description: description,
		documents:   documents,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		events:      make([]event.DomainEvent, 0), // 重建不产生事件
	}
}

// ID 获取知识库ID
func (kb *KnowledgeBase) ID() valueobject.KnowledgeBaseID {
	return kb.id
}

// Name 获取知识库名称
func (kb *KnowledgeBase) Name() string {
	return kb.name
}

// Description 获取知识库描述
func (kb *KnowledgeBase) Description() string {
	return kb.description
}

// Documents 获取文档列表（返回副本，保护内部状态）
func (kb *KnowledgeBase) Documents() []*Document {
	result := make([]*Document, len(kb.documents))
	copy(result, kb.documents)
	return result
}

// CreatedAt 获取创建时间
func (kb *KnowledgeBase) CreatedAt() time.Time {
	return kb.createdAt
}

// UpdatedAt 获取更新时间
func (kb *KnowledgeBase) UpdatedAt() time.Time {
	return kb.updatedAt
}

// UpdateInfo 更新知识库信息
// 会收集 KnowledgeBaseUpdatedEvent 事件
func (kb *KnowledgeBase) UpdateInfo(name, description string) error {
	if name == "" {
		return domain.ErrKnowledgeBaseNameEmpty
	}

	// 保存旧值用于事件
	oldName := kb.name
	oldDesc := kb.description

	// 更新值
	kb.name = name
	kb.description = description
	kb.updatedAt = time.Now()

	// 收集更新事件
	kb.addEvent(event.NewKnowledgeBaseUpdatedEvent(kb.id, oldName, name, oldDesc, description))

	return nil
}

// AddDocument 添加文档到知识库
// 通过聚合根添加文档，确保业务规则的一致性
// 会收集 DocumentAddedEvent 事件
func (kb *KnowledgeBase) AddDocument(title, content string, tags []string) (*Document, error) {
	doc, err := NewDocument(kb.id, title, content, tags)
	if err != nil {
		return nil, err
	}
	kb.documents = append(kb.documents, doc)
	kb.updatedAt = time.Now()

	// 收集文档添加事件
	kb.addEvent(event.NewDocumentAddedEvent(doc.ID(), kb.id, title))

	return doc, nil
}

// RemoveDocument 从知识库移除文档
// 会收集 DocumentRemovedEvent 事件
func (kb *KnowledgeBase) RemoveDocument(docID valueobject.DocumentID) error {
	for i, doc := range kb.documents {
		if doc.ID() == docID {
			kb.documents = append(kb.documents[:i], kb.documents[i+1:]...)
			kb.updatedAt = time.Now()

			// 收集文档删除事件
			kb.addEvent(event.NewDocumentRemovedEvent(docID, kb.id))

			return nil
		}
	}
	return domain.ErrDocumentNotFound
}

// GetDocument 获取指定文档
func (kb *KnowledgeBase) GetDocument(docID valueobject.DocumentID) (*Document, error) {
	for _, doc := range kb.documents {
		if doc.ID() == docID {
			return doc, nil
		}
	}
	return nil, domain.ErrDocumentNotFound
}

// DocumentCount 获取文档数量
func (kb *KnowledgeBase) DocumentCount() int {
	return len(kb.documents)
}

// ==================== 领域事件相关方法 ====================

// addEvent 添加领域事件（内部方法）
func (kb *KnowledgeBase) addEvent(e event.DomainEvent) {
	kb.events = append(kb.events, e)
}

// PullEvents 拉取并清空所有领域事件
// 应用层在持久化成功后调用此方法获取事件并发布
// 这确保了"先持久化，后发布事件"的顺序，保证一致性
func (kb *KnowledgeBase) PullEvents() []event.DomainEvent {
	events := kb.events
	kb.events = make([]event.DomainEvent, 0) // 清空事件
	return events
}

// HasEvents 检查是否有未发布的事件
func (kb *KnowledgeBase) HasEvents() bool {
	return len(kb.events) > 0
}
