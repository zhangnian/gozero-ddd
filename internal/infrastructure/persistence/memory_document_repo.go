package persistence

import (
	"context"
	"sync"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
)

// MemoryDocumentRepository 内存文档仓储实现
type MemoryDocumentRepository struct {
	mu   sync.RWMutex
	data map[valueobject.DocumentID]*entity.Document
}

// NewMemoryDocumentRepository 创建内存文档仓储
func NewMemoryDocumentRepository() *MemoryDocumentRepository {
	return &MemoryDocumentRepository{
		data: make(map[valueobject.DocumentID]*entity.Document),
	}
}

// 确保实现了接口
var _ repository.DocumentRepository = (*MemoryDocumentRepository)(nil)

// Save 保存文档
func (r *MemoryDocumentRepository) Save(ctx context.Context, doc *entity.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[doc.ID()] = doc
	return nil
}

// FindByID 根据ID查找文档
func (r *MemoryDocumentRepository) FindByID(ctx context.Context, id valueobject.DocumentID) (*entity.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	doc, exists := r.data[id]
	if !exists {
		return nil, nil
	}
	return doc, nil
}

// FindByKnowledgeBaseID 根据知识库ID查找所有文档
func (r *MemoryDocumentRepository) FindByKnowledgeBaseID(ctx context.Context, kbID valueobject.KnowledgeBaseID) ([]*entity.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.Document, 0)
	for _, doc := range r.data {
		if doc.KnowledgeBaseID() == kbID {
			result = append(result, doc)
		}
	}
	return result, nil
}

// Delete 删除文档
func (r *MemoryDocumentRepository) Delete(ctx context.Context, id valueobject.DocumentID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

// DeleteByKnowledgeBaseID 删除知识库下所有文档
func (r *MemoryDocumentRepository) DeleteByKnowledgeBaseID(ctx context.Context, kbID valueobject.KnowledgeBaseID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, doc := range r.data {
		if doc.KnowledgeBaseID() == kbID {
			delete(r.data, id)
		}
	}
	return nil
}

// SearchByTags 根据标签搜索文档
func (r *MemoryDocumentRepository) SearchByTags(ctx context.Context, tags []string) ([]*entity.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.Document, 0)
	tagSet := make(map[string]struct{})
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}

	for _, doc := range r.data {
		for _, docTag := range doc.Tags() {
			if _, exists := tagSet[docTag]; exists {
				result = append(result, doc)
				break
			}
		}
	}
	return result, nil
}

