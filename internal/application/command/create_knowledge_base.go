package command

import (
	"context"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain/service"
)

// CreateKnowledgeBaseCommand 创建知识库命令
type CreateKnowledgeBaseCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateKnowledgeBaseHandler 创建知识库命令处理器
type CreateKnowledgeBaseHandler struct {
	knowledgeService *service.KnowledgeService
}

// NewCreateKnowledgeBaseHandler 创建处理器
func NewCreateKnowledgeBaseHandler(ks *service.KnowledgeService) *CreateKnowledgeBaseHandler {
	return &CreateKnowledgeBaseHandler{
		knowledgeService: ks,
	}
}

// Handle 处理创建知识库命令
func (h *CreateKnowledgeBaseHandler) Handle(ctx context.Context, cmd *CreateKnowledgeBaseCommand) (*dto.KnowledgeBaseDTO, error) {
	kb, err := h.knowledgeService.CreateKnowledgeBase(ctx, cmd.Name, cmd.Description)
	if err != nil {
		return nil, err
	}

	return dto.KnowledgeBaseFromEntity(kb, false), nil
}

