package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
	"gozero-ddd/internal/infrastructure/persistence/model"
)

// GormDocumentRepository GORM 文档仓储实现
type GormDocumentRepository struct {
	db *gorm.DB
}

// NewGormDocumentRepository 创建 GORM 文档仓储
func NewGormDocumentRepository(db *gorm.DB) *GormDocumentRepository {
	return &GormDocumentRepository{db: db}
}

// 确保实现了接口
var _ repository.DocumentRepository = (*GormDocumentRepository)(nil)

// getDB 获取数据库连接（支持事务）
func (r *GormDocumentRepository) getDB(ctx context.Context) *gorm.DB {
	return GetDBFromContext(ctx, r.db)
}

// Save 保存文档
func (r *GormDocumentRepository) Save(ctx context.Context, doc *entity.Document) error {
	m := model.DocumentModelFromEntity(doc)
	return r.getDB(ctx).WithContext(ctx).Save(m).Error
}

// FindByID 根据ID查找文档
func (r *GormDocumentRepository) FindByID(ctx context.Context, id valueobject.DocumentID) (*entity.Document, error) {
	var m model.DocumentModel

	err := r.getDB(ctx).WithContext(ctx).Where("id = ?", id.String()).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return m.ToEntity(), nil
}

// FindByKnowledgeBaseID 根据知识库ID查找所有文档
func (r *GormDocumentRepository) FindByKnowledgeBaseID(ctx context.Context, kbID valueobject.KnowledgeBaseID) ([]*entity.Document, error) {
	var models []model.DocumentModel

	err := r.getDB(ctx).WithContext(ctx).
		Where("knowledge_base_id = ?", kbID.String()).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Document, len(models))
	for i, m := range models {
		result[i] = m.ToEntity()
	}

	return result, nil
}

// Delete 删除文档
func (r *GormDocumentRepository) Delete(ctx context.Context, id valueobject.DocumentID) error {
	return r.getDB(ctx).WithContext(ctx).Where("id = ?", id.String()).Delete(&model.DocumentModel{}).Error
}

// DeleteByKnowledgeBaseID 删除知识库下所有文档
func (r *GormDocumentRepository) DeleteByKnowledgeBaseID(ctx context.Context, kbID valueobject.KnowledgeBaseID) error {
	return r.getDB(ctx).WithContext(ctx).Where("knowledge_base_id = ?", kbID.String()).Delete(&model.DocumentModel{}).Error
}

// SearchByTags 根据标签搜索文档
func (r *GormDocumentRepository) SearchByTags(ctx context.Context, tags []string) ([]*entity.Document, error) {
	if len(tags) == 0 {
		return make([]*entity.Document, 0), nil
	}

	query := r.getDB(ctx).WithContext(ctx).Model(&model.DocumentModel{})

	for i, tag := range tags {
		if i == 0 {
			query = query.Where("JSON_CONTAINS(tags, ?)", `"`+tag+`"`)
		} else {
			query = query.Or("JSON_CONTAINS(tags, ?)", `"`+tag+`"`)
		}
	}

	var models []model.DocumentModel
	err := query.Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Document, len(models))
	for i, m := range models {
		result[i] = m.ToEntity()
	}

	return result, nil
}
