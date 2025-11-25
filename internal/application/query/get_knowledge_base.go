package query

import (
	"context"
	"errors"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

var (
	ErrKnowledgeBaseNotFound = errors.New("knowledge base not found")
)

// GetKnowledgeBaseQuery 获取知识库查询
type GetKnowledgeBaseQuery struct {
	ID               string
	IncludeDocuments bool
}

// GetKnowledgeBaseHandler 获取知识库查询处理器
type GetKnowledgeBaseHandler struct {
	kbRepo  repository.KnowledgeBaseRepository
	docRepo repository.DocumentRepository
}

// NewGetKnowledgeBaseHandler 创建处理器
func NewGetKnowledgeBaseHandler(
	kbRepo repository.KnowledgeBaseRepository,
	docRepo repository.DocumentRepository,
) *GetKnowledgeBaseHandler {
	return &GetKnowledgeBaseHandler{
		kbRepo:  kbRepo,
		docRepo: docRepo,
	}
}

// Handle 处理获取知识库查询
func (h *GetKnowledgeBaseHandler) Handle(ctx context.Context, query *GetKnowledgeBaseQuery) (*dto.KnowledgeBaseDTO, error) {
	kb, err := h.kbRepo.FindByID(ctx, valueobject.KnowledgeBaseIDFromString(query.ID))
	if err != nil {
		return nil, err
	}
	if kb == nil {
		return nil, ErrKnowledgeBaseNotFound
	}

	return dto.KnowledgeBaseFromEntity(kb, query.IncludeDocuments), nil
}

