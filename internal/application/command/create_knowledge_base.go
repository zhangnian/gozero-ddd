package command

import (
	"context"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/domain/event"
	"gozero-ddd/internal/domain/service"
)

// CreateKnowledgeBaseCommand 创建知识库命令
type CreateKnowledgeBaseCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateKnowledgeBaseHandler 创建知识库命令处理器
// 职责：
// 1. 接收命令参数
// 2. 调用领域服务执行业务逻辑
// 3. 在持久化成功后发布领域事件
// 4. 返回 DTO
type CreateKnowledgeBaseHandler struct {
	knowledgeService *service.KnowledgeService
	eventPublisher   event.EventPublisher // 事件发布器
}

// NewCreateKnowledgeBaseHandler 创建处理器
func NewCreateKnowledgeBaseHandler(
	ks *service.KnowledgeService,
	ep event.EventPublisher,
) *CreateKnowledgeBaseHandler {
	return &CreateKnowledgeBaseHandler{
		knowledgeService: ks,
		eventPublisher:   ep,
	}
}

// Handle 处理创建知识库命令
// 关键点：先持久化，后发布事件
// 这确保了只有成功持久化的操作才会触发事件
func (h *CreateKnowledgeBaseHandler) Handle(ctx context.Context, cmd *CreateKnowledgeBaseCommand) (*dto.KnowledgeBaseDTO, error) {
	// 1. 调用领域服务创建知识库（包含持久化）
	kb, err := h.knowledgeService.CreateKnowledgeBase(ctx, cmd.Name, cmd.Description)
	if err != nil {
		return nil, err
	}

	// 2. 从聚合根拉取领域事件并发布
	// 注意：事件在持久化成功后才发布
	events := kb.PullEvents()
	if len(events) > 0 && h.eventPublisher != nil {
		// 发布事件（如果发布失败，可以选择记录日志但不影响主流程）
		if err := h.eventPublisher.PublishAll(ctx, events); err != nil {
			// 记录错误但不阻断主流程
			// 在生产环境中，可以将失败的事件存入"待发送队列"后续重试
			// log.Printf("⚠️ 事件发布失败: %v", err)
			_ = err
		}
	}

	// 3. 返回 DTO
	return dto.KnowledgeBaseFromEntity(kb, false), nil
}

