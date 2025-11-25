package dto

import (
	"time"

	"gozero-ddd/internal/domain/entity"
)

// KnowledgeBaseDTO 知识库数据传输对象
// DTO 用于在层之间传递数据，解耦领域层和接口层
type KnowledgeBaseDTO struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	DocumentCount int           `json:"document_count"`
	Documents     []DocumentDTO `json:"documents,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// KnowledgeBaseFromEntity 从实体转换为DTO
func KnowledgeBaseFromEntity(kb *entity.KnowledgeBase, includeDocuments bool) *KnowledgeBaseDTO {
	dto := &KnowledgeBaseDTO{
		ID:            kb.ID().String(),
		Name:          kb.Name(),
		Description:   kb.Description(),
		DocumentCount: kb.DocumentCount(),
		CreatedAt:     kb.CreatedAt(),
		UpdatedAt:     kb.UpdatedAt(),
	}

	if includeDocuments {
		docs := kb.Documents()
		dto.Documents = make([]DocumentDTO, len(docs))
		for i, doc := range docs {
			dto.Documents[i] = *DocumentFromEntity(doc)
		}
	}

	return dto
}

// KnowledgeBaseListDTO 知识库列表DTO
type KnowledgeBaseListDTO struct {
	Items []*KnowledgeBaseDTO `json:"items"`
	Total int                 `json:"total"`
}

// DocumentDTO 文档数据传输对象
type DocumentDTO struct {
	ID              string    `json:"id"`
	KnowledgeBaseID string    `json:"knowledge_base_id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Tags            []string  `json:"tags"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DocumentFromEntity 从实体转换为DTO
func DocumentFromEntity(doc *entity.Document) *DocumentDTO {
	return &DocumentDTO{
		ID:              doc.ID().String(),
		KnowledgeBaseID: doc.KnowledgeBaseID().String(),
		Title:           doc.Title(),
		Content:         doc.Content(),
		Tags:            doc.Tags(),
		CreatedAt:       doc.CreatedAt(),
		UpdatedAt:       doc.UpdatedAt(),
	}
}

// DocumentListDTO 文档列表DTO
type DocumentListDTO struct {
	Items []*DocumentDTO `json:"items"`
	Total int            `json:"total"`
}

