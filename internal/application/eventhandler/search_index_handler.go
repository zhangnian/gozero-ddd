package eventhandler

import (
	"context"
	"log"

	"gozero-ddd/internal/domain/event"
)

// SearchIndexHandler 搜索索引事件处理器
// 示例：当文档发生变化时，更新搜索索引（如 Elasticsearch）
// 这是领域事件的典型应用场景
type SearchIndexHandler struct {
	// 在实际项目中，这里会注入 Elasticsearch 客户端
	// esClient *elasticsearch.Client
}

// NewSearchIndexHandler 创建搜索索引处理器
func NewSearchIndexHandler() *SearchIndexHandler {
	return &SearchIndexHandler{}
}

// 确保实现了接口
var _ event.EventHandler = (*SearchIndexHandler)(nil)

// EventName 返回空字符串，表示处理所有事件
// 这个处理器会处理所有文档相关的事件
func (h *SearchIndexHandler) EventName() string {
	return "" // 处理多个事件类型
}

// Handle 处理事件，更新搜索索引
func (h *SearchIndexHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	switch e := evt.(type) {
	case *event.DocumentAddedEvent:
		return h.handleDocumentAdded(ctx, e)
	case *event.DocumentRemovedEvent:
		return h.handleDocumentRemoved(ctx, e)
	case *event.DocumentUpdatedEvent:
		return h.handleDocumentUpdated(ctx, e)
	case *event.KnowledgeBaseDeletedEvent:
		return h.handleKnowledgeBaseDeleted(ctx, e)
	default:
		// 其他事件不处理
		return nil
	}
}

// handleDocumentAdded 处理文档添加事件
func (h *SearchIndexHandler) handleDocumentAdded(ctx context.Context, e *event.DocumentAddedEvent) error {
	log.Printf("🔍 [SearchIndex] 索引新文档: DocID=%s, Title=%s", e.DocumentID, e.Title)

	// 在实际项目中，这里会：
	// 1. 从数据库加载文档完整内容
	// 2. 对内容进行分词处理
	// 3. 将文档索引到 Elasticsearch
	//
	// 示例代码：
	// doc, _ := h.docRepo.FindByID(ctx, e.DocumentID)
	// h.esClient.Index(
	//     h.indexName,
	//     doc.ID().String(),
	//     map[string]interface{}{
	//         "title":   doc.Title(),
	//         "content": doc.Content(),
	//         "tags":    doc.Tags(),
	//         "kb_id":   doc.KnowledgeBaseID().String(),
	//     },
	// )

	return nil
}

// handleDocumentRemoved 处理文档删除事件
func (h *SearchIndexHandler) handleDocumentRemoved(ctx context.Context, e *event.DocumentRemovedEvent) error {
	log.Printf("🔍 [SearchIndex] 从索引删除文档: DocID=%s", e.DocumentID)

	// 在实际项目中，这里会：
	// h.esClient.Delete(h.indexName, e.DocumentID.String())

	return nil
}

// handleDocumentUpdated 处理文档更新事件
func (h *SearchIndexHandler) handleDocumentUpdated(ctx context.Context, e *event.DocumentUpdatedEvent) error {
	log.Printf("🔍 [SearchIndex] 更新文档索引: DocID=%s, OldTitle=%s -> NewTitle=%s",
		e.DocumentID, e.OldTitle, e.NewTitle)

	// 在实际项目中，这里会：
	// 1. 从数据库加载更新后的文档
	// 2. 更新 Elasticsearch 中的索引

	return nil
}

// handleKnowledgeBaseDeleted 处理知识库删除事件
func (h *SearchIndexHandler) handleKnowledgeBaseDeleted(ctx context.Context, e *event.KnowledgeBaseDeletedEvent) error {
	log.Printf("🔍 [SearchIndex] 删除知识库下所有文档索引: KnowledgeBaseID=%s", e.KnowledgeBaseID)

	// 在实际项目中，这里会：
	// h.esClient.DeleteByQuery(h.indexName, map[string]interface{}{
	//     "query": map[string]interface{}{
	//         "term": map[string]interface{}{
	//             "kb_id": e.KnowledgeBaseID.String(),
	//         },
	//     },
	// })

	return nil
}

