package persistence

import (
	"context"
	"sync"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

// MemoryKnowledgeBaseRepository 内存知识库仓储实现
// 这是一个简单的内存实现，用于演示目的
// 生产环境应该使用数据库实现
type MemoryKnowledgeBaseRepository struct {
	mu   sync.RWMutex
	data map[valueobject.KnowledgeBaseID]*entity.KnowledgeBase
}

// NewMemoryKnowledgeBaseRepository 创建内存知识库仓储
func NewMemoryKnowledgeBaseRepository() *MemoryKnowledgeBaseRepository {
	return &MemoryKnowledgeBaseRepository{
		data: make(map[valueobject.KnowledgeBaseID]*entity.KnowledgeBase),
	}
}

// 确保实现了接口
var _ repository.KnowledgeBaseRepository = (*MemoryKnowledgeBaseRepository)(nil)

// Save 保存知识库
func (r *MemoryKnowledgeBaseRepository) Save(ctx context.Context, kb *entity.KnowledgeBase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[kb.ID()] = kb
	return nil
}

// FindByID 根据ID查找知识库
func (r *MemoryKnowledgeBaseRepository) FindByID(ctx context.Context, id valueobject.KnowledgeBaseID) (*entity.KnowledgeBase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	kb, exists := r.data[id]
	if !exists {
		return nil, nil
	}
	return kb, nil
}

// FindAll 查找所有知识库
func (r *MemoryKnowledgeBaseRepository) FindAll(ctx context.Context) ([]*entity.KnowledgeBase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.KnowledgeBase, 0, len(r.data))
	for _, kb := range r.data {
		result = append(result, kb)
	}
	return result, nil
}

// Delete 删除知识库
func (r *MemoryKnowledgeBaseRepository) Delete(ctx context.Context, id valueobject.KnowledgeBaseID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

// ExistsByName 检查名称是否已存在
func (r *MemoryKnowledgeBaseRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, kb := range r.data {
		if kb.Name() == name {
			return true, nil
		}
	}
	return false, nil
}

