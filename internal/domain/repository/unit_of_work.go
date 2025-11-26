package repository

import "context"

// UnitOfWork 工作单元接口
// 定义在领域层，用于抽象事务操作
// 领域层不依赖具体的数据库技术，只定义接口
type UnitOfWork interface {
	// Begin 开始事务，返回带事务的上下文
	Begin(ctx context.Context) (context.Context, error)

	// Commit 提交事务
	Commit(ctx context.Context) error

	// Rollback 回滚事务
	Rollback(ctx context.Context) error

	// Transaction 在事务中执行函数（推荐使用）
	// 自动处理提交和回滚
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// TransactionalRepository 支持事务的仓储接口
// 仓储实现需要能够从上下文中获取事务连接
type TransactionalRepository interface {
	// WithTx 返回使用事务连接的仓储实例
	WithTx(ctx context.Context) interface{}
}

