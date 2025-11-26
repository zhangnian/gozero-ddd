package command

import (
	"context"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

// UpdateKnowledgeBaseCommand 更新知识库命令
type UpdateKnowledgeBaseCommand struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateKnowledgeBaseHandler 更新知识库命令处理器
type UpdateKnowledgeBaseHandler struct {
	kbRepo repository.KnowledgeBaseRepository
}

// NewUpdateKnowledgeBaseHandler 创建处理器
func NewUpdateKnowledgeBaseHandler(kbRepo repository.KnowledgeBaseRepository) *UpdateKnowledgeBaseHandler {
	return &UpdateKnowledgeBaseHandler{
		kbRepo: kbRepo,
	}
}

// Handle 处理更新知识库命令
func (h *UpdateKnowledgeBaseHandler) Handle(ctx context.Context, cmd *UpdateKnowledgeBaseCommand) (*dto.KnowledgeBaseDTO, error) {
	// 验证 ID 格式
	kbID, err := valueobject.KnowledgeBaseIDFromString(cmd.ID)
	if err != nil {
		return nil, err
	}

	// 查找知识库
	kb, err := h.kbRepo.FindByID(ctx, kbID)
	if err != nil {
		return nil, err
	}
	if kb == nil {
		return nil, domain.ErrKnowledgeBaseNotFound
	}

	// 更新信息
	if err := kb.UpdateInfo(cmd.Name, cmd.Description); err != nil {
		return nil, err
	}

	// 保存
	if err := h.kbRepo.Save(ctx, kb); err != nil {
		return nil, err
	}

	return dto.KnowledgeBaseFromEntity(kb, false), nil
}
