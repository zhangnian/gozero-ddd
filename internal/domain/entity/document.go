package entity

import (
	"errors"
	"time"

	"gozero-ddd/internal/domain/valueobject"
)

var (
	ErrDocumentTitleEmpty   = errors.New("document title cannot be empty")
	ErrDocumentContentEmpty = errors.New("document content cannot be empty")
)

// Document 文档实体
// Document 属于 KnowledgeBase 聚合，不是聚合根
// 外部操作 Document 必须通过 KnowledgeBase 进行
type Document struct {
	id              valueobject.DocumentID      // 唯一标识
	knowledgeBaseID valueobject.KnowledgeBaseID // 所属知识库ID
	title           string                      // 文档标题
	content         string                      // 文档内容
	tags            []string                    // 标签
	createdAt       time.Time                   // 创建时间
	updatedAt       time.Time                   // 更新时间
}

// NewDocument 创建新文档
func NewDocument(kbID valueobject.KnowledgeBaseID, title, content string, tags []string) (*Document, error) {
	if title == "" {
		return nil, ErrDocumentTitleEmpty
	}
	if content == "" {
		return nil, ErrDocumentContentEmpty
	}

	if tags == nil {
		tags = make([]string, 0)
	}

	now := time.Now()
	return &Document{
		id:              valueobject.NewDocumentID(),
		knowledgeBaseID: kbID,
		title:           title,
		content:         content,
		tags:            tags,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

// ReconstructDocument 从持久化数据重建文档实体
func ReconstructDocument(
	id valueobject.DocumentID,
	kbID valueobject.KnowledgeBaseID,
	title, content string,
	tags []string,
	createdAt, updatedAt time.Time,
) *Document {
	return &Document{
		id:              id,
		knowledgeBaseID: kbID,
		title:           title,
		content:         content,
		tags:            tags,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

// ID 获取文档ID
func (d *Document) ID() valueobject.DocumentID {
	return d.id
}

// KnowledgeBaseID 获取所属知识库ID
func (d *Document) KnowledgeBaseID() valueobject.KnowledgeBaseID {
	return d.knowledgeBaseID
}

// Title 获取文档标题
func (d *Document) Title() string {
	return d.title
}

// Content 获取文档内容
func (d *Document) Content() string {
	return d.content
}

// Tags 获取标签列表
func (d *Document) Tags() []string {
	result := make([]string, len(d.tags))
	copy(result, d.tags)
	return result
}

// CreatedAt 获取创建时间
func (d *Document) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt 获取更新时间
func (d *Document) UpdatedAt() time.Time {
	return d.updatedAt
}

// UpdateContent 更新文档内容
func (d *Document) UpdateContent(title, content string) error {
	if title == "" {
		return ErrDocumentTitleEmpty
	}
	if content == "" {
		return ErrDocumentContentEmpty
	}
	d.title = title
	d.content = content
	d.updatedAt = time.Now()
	return nil
}

// UpdateTags 更新标签
func (d *Document) UpdateTags(tags []string) {
	if tags == nil {
		tags = make([]string, 0)
	}
	d.tags = tags
	d.updatedAt = time.Now()
}

