package command

import (
	"context"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

// AddDocumentCommand 添加文档命令
type AddDocumentCommand struct {
	KnowledgeBaseID string   `json:"knowledge_base_id"`
	Title           string   `json:"title"`
	Content         string   `json:"content"`
	Tags            []string `json:"tags"`
}

// AddDocumentHandler 添加文档命令处理器
type AddDocumentHandler struct {
	kbRepo  repository.KnowledgeBaseRepository
	docRepo repository.DocumentRepository
}

// NewAddDocumentHandler 创建处理器
func NewAddDocumentHandler(
	kbRepo repository.KnowledgeBaseRepository,
	docRepo repository.DocumentRepository,
) *AddDocumentHandler {
	return &AddDocumentHandler{
		kbRepo:  kbRepo,
		docRepo: docRepo,
	}
}

// Handle 处理添加文档命令
func (h *AddDocumentHandler) Handle(ctx context.Context, cmd *AddDocumentCommand) (*dto.DocumentDTO, error) {
	// 查找知识库
	kb, err := h.kbRepo.FindByID(ctx, valueobject.KnowledgeBaseIDFromString(cmd.KnowledgeBaseID))
	if err != nil {
		return nil, err
	}
	if kb == nil {
		return nil, ErrKnowledgeBaseNotFound
	}

	// 通过聚合根添加文档
	doc, err := kb.AddDocument(cmd.Title, cmd.Content, cmd.Tags)
	if err != nil {
		return nil, err
	}

	// 保存文档
	if err := h.docRepo.Save(ctx, doc); err != nil {
		return nil, err
	}

	// 更新知识库
	if err := h.kbRepo.Save(ctx, kb); err != nil {
		return nil, err
	}

	return dto.DocumentFromEntity(doc), nil
}

