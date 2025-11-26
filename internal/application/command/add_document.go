package command

import (
	"context"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain"
	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/event"
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
// 演示：在事务操作中如何正确处理领域事件
type AddDocumentHandler struct {
	unitOfWork     repository.UnitOfWork
	kbRepo         repository.KnowledgeBaseRepository
	docRepo        repository.DocumentRepository
	eventPublisher event.EventPublisher // 事件发布器
}

// NewAddDocumentHandler 创建处理器
func NewAddDocumentHandler(
	uow repository.UnitOfWork,
	kbRepo repository.KnowledgeBaseRepository,
	docRepo repository.DocumentRepository,
	ep event.EventPublisher,
) *AddDocumentHandler {
	return &AddDocumentHandler{
		unitOfWork:     uow,
		kbRepo:         kbRepo,
		docRepo:        docRepo,
		eventPublisher: ep,
	}
}

// Handle 处理添加文档命令
// 使用事务确保数据一致性
// 关键点：事件在事务提交成功后才发布
func (h *AddDocumentHandler) Handle(ctx context.Context, cmd *AddDocumentCommand) (*dto.DocumentDTO, error) {
	// 验证 ID 格式
	kbID, err := valueobject.KnowledgeBaseIDFromString(cmd.KnowledgeBaseID)
	if err != nil {
		return nil, err
	}

	var result *dto.DocumentDTO
	var kb *entity.KnowledgeBase // 保存聚合根引用，用于事务后获取事件

	// 使用事务包裹所有数据库操作
	err = h.unitOfWork.Transaction(ctx, func(txCtx context.Context) error {
		// 查找知识库
		var err error
		kb, err = h.kbRepo.FindByID(txCtx, kbID)
		if err != nil {
			return err
		}
		if kb == nil {
			return domain.ErrKnowledgeBaseNotFound
		}

		// 通过聚合根添加文档（此时会收集 DocumentAddedEvent）
		doc, err := kb.AddDocument(cmd.Title, cmd.Content, cmd.Tags)
		if err != nil {
			return err
		}

		// 保存文档
		if err := h.docRepo.Save(txCtx, doc); err != nil {
			return err
		}

		// 更新知识库
		if err := h.kbRepo.Save(txCtx, kb); err != nil {
			return err
		}

		result = dto.DocumentFromEntity(doc)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 事务成功后发布事件
	// 注意：这里在事务外发布，确保只有成功持久化的操作才触发事件
	if kb != nil && h.eventPublisher != nil {
		events := kb.PullEvents()
		if len(events) > 0 {
			_ = h.eventPublisher.PublishAll(ctx, events)
		}
	}

	return result, nil
}
