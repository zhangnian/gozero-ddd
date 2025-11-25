package query

import (
	"context"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain/repository"
)

// ListKnowledgeBasesQuery 列出知识库查询
type ListKnowledgeBasesQuery struct {
	// 可以添加分页、过滤等参数
}

// ListKnowledgeBasesHandler 列出知识库查询处理器
type ListKnowledgeBasesHandler struct {
	kbRepo repository.KnowledgeBaseRepository
}

// NewListKnowledgeBasesHandler 创建处理器
func NewListKnowledgeBasesHandler(kbRepo repository.KnowledgeBaseRepository) *ListKnowledgeBasesHandler {
	return &ListKnowledgeBasesHandler{
		kbRepo: kbRepo,
	}
}

// Handle 处理列出知识库查询
func (h *ListKnowledgeBasesHandler) Handle(ctx context.Context, query *ListKnowledgeBasesQuery) (*dto.KnowledgeBaseListDTO, error) {
	kbs, err := h.kbRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*dto.KnowledgeBaseDTO, len(kbs))
	for i, kb := range kbs {
		items[i] = dto.KnowledgeBaseFromEntity(kb, false)
	}

	return &dto.KnowledgeBaseListDTO{
		Items: items,
		Total: len(items),
	}, nil
}

