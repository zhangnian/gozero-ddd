package command

import (
	"context"
	"errors"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

var (
	ErrSourceKnowledgeBaseNotFound = errors.New("source knowledge base not found")
	ErrTargetKnowledgeBaseNotFound = errors.New("target knowledge base not found")
	ErrCannotMergeSameKnowledgeBase = errors.New("cannot merge knowledge base with itself")
)

// MergeKnowledgeBasesCommand 合并知识库命令
// 将源知识库的所有文档移动到目标知识库，然后删除源知识库
// 这是一个需要事务保证的操作
type MergeKnowledgeBasesCommand struct {
	SourceID string `json:"source_id"` // 源知识库ID（将被删除）
	TargetID string `json:"target_id"` // 目标知识库ID（保留）
}

// MergeKnowledgeBasesHandler 合并知识库命令处理器
// 演示如何在应用层正确使用事务
type MergeKnowledgeBasesHandler struct {
	unitOfWork repository.UnitOfWork
	kbRepo     repository.KnowledgeBaseRepository
	docRepo    repository.DocumentRepository
}

// NewMergeKnowledgeBasesHandler 创建处理器
func NewMergeKnowledgeBasesHandler(
	uow repository.UnitOfWork,
	kbRepo repository.KnowledgeBaseRepository,
	docRepo repository.DocumentRepository,
) *MergeKnowledgeBasesHandler {
	return &MergeKnowledgeBasesHandler{
		unitOfWork: uow,
		kbRepo:     kbRepo,
		docRepo:    docRepo,
	}
}

// Handle 处理合并知识库命令
// 使用事务确保操作的原子性：要么全部成功，要么全部失败
func (h *MergeKnowledgeBasesHandler) Handle(ctx context.Context, cmd *MergeKnowledgeBasesCommand) (*dto.MergeResultDTO, error) {
	sourceID := valueobject.KnowledgeBaseIDFromString(cmd.SourceID)
	targetID := valueobject.KnowledgeBaseIDFromString(cmd.TargetID)

	// 检查是否合并自己
	if cmd.SourceID == cmd.TargetID {
		return nil, ErrCannotMergeSameKnowledgeBase
	}

	var result *dto.MergeResultDTO

	// 使用工作单元执行事务
	// Transaction 方法会自动处理提交和回滚
	err := h.unitOfWork.Transaction(ctx, func(txCtx context.Context) error {
		// ========== 以下所有操作都在同一个事务中 ==========

		// 1. 查找源知识库
		sourceKB, err := h.kbRepo.FindByID(txCtx, sourceID)
		if err != nil {
			return err
		}
		if sourceKB == nil {
			return ErrSourceKnowledgeBaseNotFound
		}

		// 2. 查找目标知识库
		targetKB, err := h.kbRepo.FindByID(txCtx, targetID)
		if err != nil {
			return err
		}
		if targetKB == nil {
			return ErrTargetKnowledgeBaseNotFound
		}

		// 3. 获取源知识库的所有文档
		sourceDocs, err := h.docRepo.FindByKnowledgeBaseID(txCtx, sourceID)
		if err != nil {
			return err
		}

		// 4. 将每个文档添加到目标知识库
		movedCount := 0
		for _, doc := range sourceDocs {
			// 通过聚合根添加文档到目标知识库
			newDoc, err := targetKB.AddDocument(doc.Title(), doc.Content(), doc.Tags())
			if err != nil {
				return err
			}

			// 保存新文档
			if err := h.docRepo.Save(txCtx, newDoc); err != nil {
				return err
			}

			// 删除原文档
			if err := h.docRepo.Delete(txCtx, doc.ID()); err != nil {
				return err
			}

			movedCount++
		}

		// 5. 更新目标知识库
		if err := h.kbRepo.Save(txCtx, targetKB); err != nil {
			return err
		}

		// 6. 删除源知识库
		if err := h.kbRepo.Delete(txCtx, sourceID); err != nil {
			return err
		}

		// 构建结果
		result = &dto.MergeResultDTO{
			SourceID:        cmd.SourceID,
			SourceName:      sourceKB.Name(),
			TargetID:        cmd.TargetID,
			TargetName:      targetKB.Name(),
			DocumentsMoved:  movedCount,
			SourceDeleted:   true,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

