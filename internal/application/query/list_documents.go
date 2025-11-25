package query

import (
	"context"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

// ListDocumentsQuery 列出文档查询
type ListDocumentsQuery struct {
	KnowledgeBaseID string
}

// ListDocumentsHandler 列出文档查询处理器
type ListDocumentsHandler struct {
	docRepo repository.DocumentRepository
}

// NewListDocumentsHandler 创建处理器
func NewListDocumentsHandler(docRepo repository.DocumentRepository) *ListDocumentsHandler {
	return &ListDocumentsHandler{
		docRepo: docRepo,
	}
}

// Handle 处理列出文档查询
func (h *ListDocumentsHandler) Handle(ctx context.Context, query *ListDocumentsQuery) (*dto.DocumentListDTO, error) {
	docs, err := h.docRepo.FindByKnowledgeBaseID(ctx, valueobject.KnowledgeBaseIDFromString(query.KnowledgeBaseID))
	if err != nil {
		return nil, err
	}

	items := make([]*dto.DocumentDTO, len(docs))
	for i, doc := range docs {
		items[i] = dto.DocumentFromEntity(doc)
	}

	return &dto.DocumentListDTO{
		Items: items,
		Total: len(items),
	}, nil
}

