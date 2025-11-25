package repository

import (
	"context"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/valueobject"
)

// DocumentRepository 文档仓储接口
// 虽然 Document 不是聚合根，但为了查询方便，提供独立的仓储接口
type DocumentRepository interface {
	// Save 保存文档
	Save(ctx context.Context, doc *entity.Document) error

	// FindByID 根据ID查找文档
	FindByID(ctx context.Context, id valueobject.DocumentID) (*entity.Document, error)

	// FindByKnowledgeBaseID 根据知识库ID查找所有文档
	FindByKnowledgeBaseID(ctx context.Context, kbID valueobject.KnowledgeBaseID) ([]*entity.Document, error)

	// Delete 删除文档
	Delete(ctx context.Context, id valueobject.DocumentID) error

	// DeleteByKnowledgeBaseID 删除知识库下所有文档
	DeleteByKnowledgeBaseID(ctx context.Context, kbID valueobject.KnowledgeBaseID) error

	// SearchByTags 根据标签搜索文档
	SearchByTags(ctx context.Context, tags []string) ([]*entity.Document, error)
}

