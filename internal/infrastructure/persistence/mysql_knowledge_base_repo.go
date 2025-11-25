package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/valueobject"
	"gozero-ddd/internal/infrastructure/persistence/model"
)

// MysqlKnowledgeBaseRepository MySQL 知识库仓储实现
type MysqlKnowledgeBaseRepository struct {
	conn    sqlx.SqlConn
	docRepo repository.DocumentRepository
}

// NewMysqlKnowledgeBaseRepository 创建 MySQL 知识库仓储
func NewMysqlKnowledgeBaseRepository(conn sqlx.SqlConn, docRepo repository.DocumentRepository) *MysqlKnowledgeBaseRepository {
	return &MysqlKnowledgeBaseRepository{
		conn:    conn,
		docRepo: docRepo,
	}
}

// 确保实现了接口
var _ repository.KnowledgeBaseRepository = (*MysqlKnowledgeBaseRepository)(nil)

// Save 保存知识库（创建或更新）
func (r *MysqlKnowledgeBaseRepository) Save(ctx context.Context, kb *entity.KnowledgeBase) error {
	m := model.KnowledgeBaseModelFromEntity(kb)

	// 使用 REPLACE INTO 实现 upsert
	query := `
		REPLACE INTO knowledge_bases (id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.conn.ExecCtx(ctx, query, m.ID, m.Name, m.Description, m.CreatedAt, m.UpdatedAt)
	return err
}

// FindByID 根据ID查找知识库
func (r *MysqlKnowledgeBaseRepository) FindByID(ctx context.Context, id valueobject.KnowledgeBaseID) (*entity.KnowledgeBase, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM knowledge_bases WHERE id = ?`

	var m model.KnowledgeBaseModel
	err := r.conn.QueryRowCtx(ctx, &m, query, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
func (r *MysqlKnowledgeBaseRepository) FindAll(ctx context.Context) ([]*entity.KnowledgeBase, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM knowledge_bases ORDER BY created_at DESC`

	var models []model.KnowledgeBaseModel
	err := r.conn.QueryRowsCtx(ctx, &models, query)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.KnowledgeBase, len(models))
	for i, m := range models {
		// 为了性能考虑，列表查询不加载文档
		result[i] = m.ToEntity(nil)
	}

	return result, nil
}

// Delete 删除知识库
func (r *MysqlKnowledgeBaseRepository) Delete(ctx context.Context, id valueobject.KnowledgeBaseID) error {
	query := `DELETE FROM knowledge_bases WHERE id = ?`
	_, err := r.conn.ExecCtx(ctx, query, id.String())
	return err
}

// ExistsByName 检查名称是否已存在
func (r *MysqlKnowledgeBaseRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT COUNT(1) FROM knowledge_bases WHERE name = ?`

	var count int
	err := r.conn.QueryRowCtx(ctx, &count, query, name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

