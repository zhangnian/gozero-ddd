package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gozero-ddd/internal/domain/repository"
)

// 定义上下文键，用于存储事务连接
type txKey struct{}

// GormUnitOfWork GORM 工作单元实现
type GormUnitOfWork struct {
	db *gorm.DB
}

// NewGormUnitOfWork 创建 GORM 工作单元
func NewGormUnitOfWork(db *gorm.DB) *GormUnitOfWork {
	return &GormUnitOfWork{db: db}
}

// 确保实现了接口
var _ repository.UnitOfWork = (*GormUnitOfWork)(nil)

// Begin 开始事务
func (u *GormUnitOfWork) Begin(ctx context.Context) (context.Context, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}
	// 将事务连接存入上下文
	return context.WithValue(ctx, txKey{}, tx), nil
}

// Commit 提交事务
func (u *GormUnitOfWork) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if !ok {
		return errors.New("no transaction in context")
	}
	return tx.Commit().Error
}

// Rollback 回滚事务
func (u *GormUnitOfWork) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if !ok {
		return errors.New("no transaction in context")
	}
	return tx.Rollback().Error
}

// Transaction 在事务中执行函数（推荐方式）
func (u *GormUnitOfWork) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// 开始事务
	txCtx, err := u.Begin(ctx)
	if err != nil {
		return err
	}

	// 执行业务逻辑
	if err := fn(txCtx); err != nil {
		// 出错时回滚
		_ = u.Rollback(txCtx)
		return err
	}

	// 成功时提交
	return u.Commit(txCtx)
}

// GetTxFromContext 从上下文获取事务连接
// 供仓储实现使用
func GetTxFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
}

// GetDBFromContext 从上下文获取数据库连接（优先事务连接）
func GetDBFromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := GetTxFromContext(ctx); ok {
		return tx
	}
	return defaultDB
}

