package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/valueobject"
)

// StringSlice 自定义类型，用于 GORM JSON 序列化
type StringSlice []string

// Scan 实现 sql.Scanner 接口
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = make([]string, 0)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to scan StringSlice")
	}

	if len(bytes) == 0 {
		*s = make([]string, 0)
		return nil
	}

	return json.Unmarshal(bytes, s)
}

// Value 实现 driver.Valuer 接口
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// DocumentModel 文档数据库模型
type DocumentModel struct {
	ID              string      `gorm:"column:id;type:varchar(36);primaryKey"`
	KnowledgeBaseID string      `gorm:"column:knowledge_base_id;type:varchar(36);index;not null"`
	Title           string      `gorm:"column:title;type:varchar(500);not null"`
	Content         string      `gorm:"column:content;type:longtext;not null"`
	Tags            StringSlice `gorm:"column:tags;type:json"`
	CreatedAt       time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time   `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (DocumentModel) TableName() string {
	return "documents"
}

// ToEntity 将数据库模型转换为领域实体
func (m *DocumentModel) ToEntity() *entity.Document {
	tags := make([]string, len(m.Tags))
	copy(tags, m.Tags)

	return entity.ReconstructDocument(
		valueobject.DocumentIDFromString(m.ID),
		valueobject.KnowledgeBaseIDFromString(m.KnowledgeBaseID),
		m.Title,
		m.Content,
		tags,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// FromEntity 从领域实体创建数据库模型
func DocumentModelFromEntity(doc *entity.Document) *DocumentModel {
	return &DocumentModel{
		ID:              doc.ID().String(),
		KnowledgeBaseID: doc.KnowledgeBaseID().String(),
		Title:           doc.Title(),
		Content:         doc.Content(),
		Tags:            StringSlice(doc.Tags()),
		CreatedAt:       doc.CreatedAt(),
		UpdatedAt:       doc.UpdatedAt(),
	}
}
