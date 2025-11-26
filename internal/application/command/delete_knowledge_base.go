package command

import (
	"context"

	"gozero-ddd/internal/domain"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/service"
	"gozero-ddd/internal/domain/valueobject"
)

// DeleteKnowledgeBaseCommand 删除知识库命令
type DeleteKnowledgeBaseCommand struct {
	ID string `json:"id"`
}

// DeleteKnowledgeBaseHandler 删除知识库命令处理器
type DeleteKnowledgeBaseHandler struct {
	kbRepo           repository.KnowledgeBaseRepository
	knowledgeService *service.KnowledgeService
}

// NewDeleteKnowledgeBaseHandler 创建处理器
func NewDeleteKnowledgeBaseHandler(
	kbRepo repository.KnowledgeBaseRepository,
	ks *service.KnowledgeService,
) *DeleteKnowledgeBaseHandler {
	return &DeleteKnowledgeBaseHandler{
		kbRepo:           kbRepo,
		knowledgeService: ks,
	}
}

// Handle 处理删除知识库命令
func (h *DeleteKnowledgeBaseHandler) Handle(ctx context.Context, cmd *DeleteKnowledgeBaseCommand) error {
	// 验证 ID 格式
	kbID, err := valueobject.KnowledgeBaseIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// 查找知识库
	kb, err := h.kbRepo.FindByID(ctx, kbID)
	if err != nil {
		return err
	}
	if kb == nil {
		return domain.ErrKnowledgeBaseNotFound
	}

	// 使用领域服务删除（包含删除关联文档的逻辑）
	return h.knowledgeService.DeleteKnowledgeBase(ctx, kb)
}
