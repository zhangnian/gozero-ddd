package model

import (
	"time"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/valueobject"
)

// KnowledgeBaseModel 知识库数据库模型
// GORM 模型，用于数据库表映射
type KnowledgeBaseModel struct {
	ID          string    `gorm:"column:id;type:varchar(36);primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(255);uniqueIndex;not null"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (KnowledgeBaseModel) TableName() string {
	return "knowledge_bases"
}

// ToEntity 将数据库模型转换为领域实体
// 使用 MustKnowledgeBaseIDFromString 因为数据来自数据库，是可信的
func (m *KnowledgeBaseModel) ToEntity(documents []*entity.Document) *entity.KnowledgeBase {
	return entity.ReconstructKnowledgeBase(
		valueobject.MustKnowledgeBaseIDFromString(m.ID),
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
