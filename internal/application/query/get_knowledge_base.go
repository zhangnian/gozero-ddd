package query

import (
	"context"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
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
	// 验证 ID 格式
	kbID, err := valueobject.KnowledgeBaseIDFromString(query.ID)
	if err != nil {
		return nil, err
	}

	kb, err := h.kbRepo.FindByID(ctx, kbID)
	if err != nil {
		return nil, err
	}
	if kb == nil {
		return nil, domain.ErrKnowledgeBaseNotFound
	}

	return dto.KnowledgeBaseFromEntity(kb, query.IncludeDocuments), nil
}
