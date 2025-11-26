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

// GormKnowledgeBaseRepository GORM 知识库仓储实现
type GormKnowledgeBaseRepository struct {
	db      *gorm.DB
	docRepo repository.DocumentRepository
}

// NewGormKnowledgeBaseRepository 创建 GORM 知识库仓储
func NewGormKnowledgeBaseRepository(db *gorm.DB, docRepo repository.DocumentRepository) *GormKnowledgeBaseRepository {
	return &GormKnowledgeBaseRepository{
		db:      db,
		docRepo: docRepo,
	}
}

// 确保实现了接口
var _ repository.KnowledgeBaseRepository = (*GormKnowledgeBaseRepository)(nil)

// getDB 获取数据库连接（支持事务）
func (r *GormKnowledgeBaseRepository) getDB(ctx context.Context) *gorm.DB {
	return GetDBFromContext(ctx, r.db)
}

// Save 保存知识库（创建或更新）
func (r *GormKnowledgeBaseRepository) Save(ctx context.Context, kb *entity.KnowledgeBase) error {
	m := model.KnowledgeBaseModelFromEntity(kb)
	return r.getDB(ctx).WithContext(ctx).Save(m).Error
}

// FindByID 根据ID查找知识库
func (r *GormKnowledgeBaseRepository) FindByID(ctx context.Context, id valueobject.KnowledgeBaseID) (*entity.KnowledgeBase, error) {
	var m model.KnowledgeBaseModel

	err := r.getDB(ctx).WithContext(ctx).Where("id = ?", id.String()).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// 加载关联的文档
	docs, err := r.docRepo.FindByKnowledgeBaseID(ctx, id)
	if err != nil {
		return nil, err
	}

	return m.ToEntity(docs), nil
}

// FindAll 查找所有知识库
func (r *GormKnowledgeBaseRepository) FindAll(ctx context.Context) ([]*entity.KnowledgeBase, error) {
	var models []model.KnowledgeBaseModel

	err := r.getDB(ctx).WithContext(ctx).Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.KnowledgeBase, len(models))
	for i, m := range models {
		result[i] = m.ToEntity(nil)
	}

	return result, nil
}

// Delete 删除知识库
func (r *GormKnowledgeBaseRepository) Delete(ctx context.Context, id valueobject.KnowledgeBaseID) error {
	return r.getDB(ctx).WithContext(ctx).Where("id = ?", id.String()).Delete(&model.KnowledgeBaseModel{}).Error
}

// ExistsByName 检查名称是否已存在
func (r *GormKnowledgeBaseRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64

	err := r.getDB(ctx).WithContext(ctx).Model(&model.KnowledgeBaseModel{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
