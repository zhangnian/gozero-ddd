package repository

import (
	"context"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/valueobject"
)

// KnowledgeBaseRepository 知识库仓储接口
// 仓储接口定义在领域层，实现在基础设施层
// 这体现了依赖倒置原则：领域层不依赖基础设施层
type KnowledgeBaseRepository interface {
	// Save 保存知识库（创建或更新）
	Save(ctx context.Context, kb *entity.KnowledgeBase) error

	// FindByID 根据ID查找知识库
	FindByID(ctx context.Context, id valueobject.KnowledgeBaseID) (*entity.KnowledgeBase, error)

	// FindAll 查找所有知识库
	FindAll(ctx context.Context) ([]*entity.KnowledgeBase, error)

	// Delete 删除知识库
	Delete(ctx context.Context, id valueobject.KnowledgeBaseID) error

	// ExistsByName 检查名称是否已存在
	ExistsByName(ctx context.Context, name string) (bool, error)
}

