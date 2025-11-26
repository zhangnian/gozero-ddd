package command

import (
	"context"

	"gozero-ddd/internal/domain"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

// RemoveDocumentCommand 删除文档命令
type RemoveDocumentCommand struct {
	KnowledgeBaseID string `json:"knowledge_base_id"`
	DocumentID      string `json:"document_id"`
}

// RemoveDocumentHandler 删除文档命令处理器
type RemoveDocumentHandler struct {
	unitOfWork repository.UnitOfWork
	kbRepo     repository.KnowledgeBaseRepository
	docRepo    repository.DocumentRepository
}

// NewRemoveDocumentHandler 创建处理器
func NewRemoveDocumentHandler(
	uow repository.UnitOfWork,
	kbRepo repository.KnowledgeBaseRepository,
	docRepo repository.DocumentRepository,
) *RemoveDocumentHandler {
	return &RemoveDocumentHandler{
		unitOfWork: uow,
		kbRepo:     kbRepo,
		docRepo:    docRepo,
	}
}

// Handle 处理删除文档命令
// 使用事务确保数据一致性
func (h *RemoveDocumentHandler) Handle(ctx context.Context, cmd *RemoveDocumentCommand) error {
	// 验证 ID 格式
	kbID, err := valueobject.KnowledgeBaseIDFromString(cmd.KnowledgeBaseID)
	if err != nil {
		return err
	}

	docID, err := valueobject.DocumentIDFromString(cmd.DocumentID)
	if err != nil {
		return err
	}

	// 使用事务包裹所有数据库操作
	return h.unitOfWork.Transaction(ctx, func(txCtx context.Context) error {
		// 查找知识库
		kb, err := h.kbRepo.FindByID(txCtx, kbID)
		if err != nil {
			return err
		}
		if kb == nil {
			return domain.ErrKnowledgeBaseNotFound
		}

		// 通过聚合根删除文档
		if err := kb.RemoveDocument(docID); err != nil {
			return err
		}

		// 删除文档持久化数据
		if err := h.docRepo.Delete(txCtx, docID); err != nil {
			return err
		}

		// 更新知识库
		return h.kbRepo.Save(txCtx, kb)
	})
}
