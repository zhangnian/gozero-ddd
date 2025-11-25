package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
	"gozero-ddd/internal/infrastructure/persistence/model"
)

// MysqlDocumentRepository MySQL 文档仓储实现
type MysqlDocumentRepository struct {
	conn sqlx.SqlConn
}

// NewMysqlDocumentRepository 创建 MySQL 文档仓储
func NewMysqlDocumentRepository(conn sqlx.SqlConn) *MysqlDocumentRepository {
	return &MysqlDocumentRepository{conn: conn}
}

// 确保实现了接口
var _ repository.DocumentRepository = (*MysqlDocumentRepository)(nil)

// Save 保存文档
func (r *MysqlDocumentRepository) Save(ctx context.Context, doc *entity.Document) error {
	m := model.DocumentModelFromEntity(doc)

	query := `
		REPLACE INTO documents (id, knowledge_base_id, title, content, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.conn.ExecCtx(ctx, query,
		m.ID, m.KnowledgeBaseID, m.Title, m.Content, m.Tags, m.CreatedAt, m.UpdatedAt)
	return err
}

// FindByID 根据ID查找文档
func (r *MysqlDocumentRepository) FindByID(ctx context.Context, id valueobject.DocumentID) (*entity.Document, error) {
	query := `
		SELECT id, knowledge_base_id, title, content, tags, created_at, updated_at 
		FROM documents 
		WHERE id = ?
	`

	var m model.DocumentModel
	err := r.conn.QueryRowCtx(ctx, &m, query, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return m.ToEntity(), nil
}

// FindByKnowledgeBaseID 根据知识库ID查找所有文档
func (r *MysqlDocumentRepository) FindByKnowledgeBaseID(ctx context.Context, kbID valueobject.KnowledgeBaseID) ([]*entity.Document, error) {
	query := `
		SELECT id, knowledge_base_id, title, content, tags, created_at, updated_at 
		FROM documents 
		WHERE knowledge_base_id = ?
		ORDER BY created_at DESC
	`

	var models []model.DocumentModel
	err := r.conn.QueryRowsCtx(ctx, &models, query, kbID.String())
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
func (r *MysqlDocumentRepository) Delete(ctx context.Context, id valueobject.DocumentID) error {
	query := `DELETE FROM documents WHERE id = ?`
	_, err := r.conn.ExecCtx(ctx, query, id.String())
	return err
}

// DeleteByKnowledgeBaseID 删除知识库下所有文档
func (r *MysqlDocumentRepository) DeleteByKnowledgeBaseID(ctx context.Context, kbID valueobject.KnowledgeBaseID) error {
	query := `DELETE FROM documents WHERE knowledge_base_id = ?`
	_, err := r.conn.ExecCtx(ctx, query, kbID.String())
	return err
}

// SearchByTags 根据标签搜索文档
func (r *MysqlDocumentRepository) SearchByTags(ctx context.Context, tags []string) ([]*entity.Document, error) {
	if len(tags) == 0 {
		return make([]*entity.Document, 0), nil
	}

	// 使用 JSON_CONTAINS 进行标签匹配
	// 这里简化处理，只要包含任一标签即返回
	query := `
		SELECT id, knowledge_base_id, title, content, tags, created_at, updated_at 
		FROM documents 
		WHERE `

	// 构建 OR 条件
	conditions := ""
	args := make([]interface{}, 0, len(tags))
	for i, tag := range tags {
		if i > 0 {
			conditions += " OR "
		}
		tagJSON, _ := json.Marshal(tag)
		conditions += "JSON_CONTAINS(tags, ?)"
		args = append(args, string(tagJSON))
	}
	query += conditions + " ORDER BY created_at DESC"

	var models []model.DocumentModel
	err := r.conn.QueryRowsCtx(ctx, &models, query, args...)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Document, len(models))
	for i, m := range models {
		result[i] = m.ToEntity()
	}

	return result, nil
}

