package entity

import (
	"errors"
	"time"

	"gozero-ddd/internal/domain/valueobject"
)

var (
	ErrKnowledgeBaseNameEmpty = errors.New("knowledge base name cannot be empty")
	ErrDocumentNotFound       = errors.New("document not found in knowledge base")
)

// KnowledgeBase 知识库实体（聚合根）
// 作为聚合根，KnowledgeBase 负责管理其下所有的 Document 实体
// 外部不能直接操作 Document，必须通过 KnowledgeBase 进行
type KnowledgeBase struct {
	id          valueobject.KnowledgeBaseID // 唯一标识
	name        string                      // 知识库名称
	description string                      // 描述
	documents   []*Document                 // 文档集合
	createdAt   time.Time                   // 创建时间
	updatedAt   time.Time                   // 更新时间
}

// NewKnowledgeBase 创建新的知识库
func NewKnowledgeBase(name, description string) (*KnowledgeBase, error) {
	if name == "" {
		return nil, ErrKnowledgeBaseNameEmpty
	}

	now := time.Now()
	return &KnowledgeBase{
		id:          valueobject.NewKnowledgeBaseID(),
		name:        name,
		description: description,
		documents:   make([]*Document, 0),
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ReconstructKnowledgeBase 从持久化数据重建知识库实体
// 用于仓储层从数据库加载数据时使用
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
func (kb *KnowledgeBase) UpdateInfo(name, description string) error {
	if name == "" {
		return ErrKnowledgeBaseNameEmpty
	}
	kb.name = name
	kb.description = description
	kb.updatedAt = time.Now()
	return nil
}

// AddDocument 添加文档到知识库
// 通过聚合根添加文档，确保业务规则的一致性
func (kb *KnowledgeBase) AddDocument(title, content string, tags []string) (*Document, error) {
	doc, err := NewDocument(kb.id, title, content, tags)
	if err != nil {
		return nil, err
	}
	kb.documents = append(kb.documents, doc)
	kb.updatedAt = time.Now()
	return doc, nil
}

// RemoveDocument 从知识库移除文档
func (kb *KnowledgeBase) RemoveDocument(docID valueobject.DocumentID) error {
	for i, doc := range kb.documents {
		if doc.ID() == docID {
			kb.documents = append(kb.documents[:i], kb.documents[i+1:]...)
			kb.updatedAt = time.Now()
			return nil
		}
	}
	return ErrDocumentNotFound
}

// GetDocument 获取指定文档
func (kb *KnowledgeBase) GetDocument(docID valueobject.DocumentID) (*Document, error) {
	for _, doc := range kb.documents {
		if doc.ID() == docID {
			return doc, nil
		}
	}
	return nil, ErrDocumentNotFound
}

// DocumentCount 获取文档数量
func (kb *KnowledgeBase) DocumentCount() int {
	return len(kb.documents)
}
