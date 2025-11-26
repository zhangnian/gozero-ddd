package persistence

import (
	"context"

	"gozero-ddd/internal/domain/repository"
)

// MemoryUnitOfWork 内存工作单元实现
// 用于内存存储模式，不需要真正的事务支持
type MemoryUnitOfWork struct{}

// NewMemoryUnitOfWork 创建内存工作单元
func NewMemoryUnitOfWork() *MemoryUnitOfWork {
	return &MemoryUnitOfWork{}
}

// 确保实现了接口
var _ repository.UnitOfWork = (*MemoryUnitOfWork)(nil)

// Begin 开始事务（内存模式下直接返回原上下文）
func (u *MemoryUnitOfWork) Begin(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

// Commit 提交事务（内存模式下无操作）
func (u *MemoryUnitOfWork) Commit(ctx context.Context) error {
	return nil
}

// Rollback 回滚事务（内存模式下无操作）
func (u *MemoryUnitOfWork) Rollback(ctx context.Context) error {
	return nil
}

// Transaction 在事务中执行函数（内存模式下直接执行）
func (u *MemoryUnitOfWork) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

