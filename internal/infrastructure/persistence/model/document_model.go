package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/valueobject"
)

// DocumentModel 文档数据库模型
type DocumentModel struct {
	ID              string         `db:"id"`
	KnowledgeBaseID string         `db:"knowledge_base_id"`
	Title           string         `db:"title"`
	Content         string         `db:"content"`
	Tags            sql.NullString `db:"tags"` // JSON 格式存储标签
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

// TableName 返回表名
func (m *DocumentModel) TableName() string {
	return "documents"
}

// ToEntity 将数据库模型转换为领域实体
func (m *DocumentModel) ToEntity() *entity.Document {
	var tags []string
	if m.Tags.Valid && m.Tags.String != "" {
		_ = json.Unmarshal([]byte(m.Tags.String), &tags)
	}
	if tags == nil {
		tags = make([]string, 0)
	}

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
	tagsJSON, _ := json.Marshal(doc.Tags())

	return &DocumentModel{
		ID:              doc.ID().String(),
		KnowledgeBaseID: doc.KnowledgeBaseID().String(),
		Title:           doc.Title(),
		Content:         doc.Content(),
		Tags:            sql.NullString{String: string(tagsJSON), Valid: true},
		CreatedAt:       doc.CreatedAt(),
		UpdatedAt:       doc.UpdatedAt(),
	}
}

