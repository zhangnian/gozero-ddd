package model

import (
	"time"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/valueobject"
)

// KnowledgeBaseModel 知识库数据库模型
// 用于数据库表映射，与领域实体分离
type KnowledgeBaseModel struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TableName 返回表名
func (m *KnowledgeBaseModel) TableName() string {
	return "knowledge_bases"
}

// ToEntity 将数据库模型转换为领域实体
func (m *KnowledgeBaseModel) ToEntity(documents []*entity.Document) *entity.KnowledgeBase {
	return entity.ReconstructKnowledgeBase(
		valueobject.KnowledgeBaseIDFromString(m.ID),
		m.Name,
		m.Description,
		documents,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// FromEntity 从领域实体创建数据库模型
func KnowledgeBaseModelFromEntity(kb *entity.KnowledgeBase) *KnowledgeBaseModel {
	return &KnowledgeBaseModel{
		ID:          kb.ID().String(),
		Name:        kb.Name(),
		Description: kb.Description(),
		CreatedAt:   kb.CreatedAt(),
		UpdatedAt:   kb.UpdatedAt(),
	}
}

