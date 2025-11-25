package command

import (
	"context"
	"errors"

	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
)

// RemoveDocumentCommand 删除文档命令
type RemoveDocumentCommand struct {
	KnowledgeBaseID string `json:"knowledge_base_id"`
	DocumentID      string `json:"document_id"`
}

// RemoveDocumentHandler 删除文档命令处理器
type RemoveDocumentHandler struct {
	kbRepo  repository.KnowledgeBaseRepository
	docRepo repository.DocumentRepository
}

// NewRemoveDocumentHandler 创建处理器
func NewRemoveDocumentHandler(
	kbRepo repository.KnowledgeBaseRepository,
	docRepo repository.DocumentRepository,
) *RemoveDocumentHandler {
	return &RemoveDocumentHandler{
		kbRepo:  kbRepo,
		docRepo: docRepo,
	}
}

// Handle 处理删除文档命令
func (h *RemoveDocumentHandler) Handle(ctx context.Context, cmd *RemoveDocumentCommand) error {
	// 查找知识库
	kb, err := h.kbRepo.FindByID(ctx, valueobject.KnowledgeBaseIDFromString(cmd.KnowledgeBaseID))
	if err != nil {
		return err
	}
	if kb == nil {
		return ErrKnowledgeBaseNotFound
	}

	docID := valueobject.DocumentIDFromString(cmd.DocumentID)

	// 通过聚合根删除文档
	if err := kb.RemoveDocument(docID); err != nil {
		return err
	}

	// 删除文档持久化数据
	if err := h.docRepo.Delete(ctx, docID); err != nil {
		return err
	}

	// 更新知识库
	return h.kbRepo.Save(ctx, kb)
}

