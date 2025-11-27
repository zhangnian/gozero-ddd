package svc

import (
	"log"

	appcontainer "gozero-ddd/internal/application/container"
	"gozero-ddd/internal/infrastructure/config"
	infracontainer "gozero-ddd/internal/infrastructure/container"
)

// ServiceContext 服务上下文
// go-zero 使用 ServiceContext 来管理依赖注入
//
// 重构后的 ServiceContext 遵循 DDD 分层原则：
// - 接口层只依赖应用层（通过 App 容器访问）
// - 不直接暴露基础设施层组件（DB、仓储等）
// - 不直接暴露领域层组件（领域服务等）
type ServiceContext struct {
	Config config.Config

	// 应用层容器 - 接口层唯一应该访问的入口
	// 包含所有的 Command Handler 和 Query Handler
	App *appcontainer.ApplicationContainer

	// 基础设施容器 - 内部持有，用于资源管理（如关闭数据库连接）
	// 注意：这里使用小写字母开头，表示不对外暴露
	infra *infracontainer.InfrastructureContainer
}

// NewServiceContext 创建服务上下文
// 按照分层顺序初始化各层容器
func NewServiceContext(c config.Config) *ServiceContext {
	log.Println("🚀 [ServiceContext] 开始初始化服务上下文...")

	// 1. 创建基础设施层容器
	// 负责：数据库连接、仓储实现、事件总线、领域服务
	infra := infracontainer.NewInfrastructureContainer(&configAdapter{c})

	// 2. 创建应用层容器
	// 负责：Command Handler、Query Handler
	// 依赖：基础设施层容器
	app := appcontainer.NewApplicationContainer(infra)

	log.Println("✅ [ServiceContext] 服务上下文初始化完成")

	return &ServiceContext{
		Config: c,
		App:    app,
		infra:  infra,
	}
}

// Close 关闭服务上下文，释放资源
func (ctx *ServiceContext) Close() error {
	if ctx.infra != nil {
		return ctx.infra.Close()
	}
	return nil
}

// configAdapter 配置适配器
// 将 config.Config 适配为 InfraConfig 接口
type configAdapter struct {
	config.Config
}

func (a *configAdapter) GetMySQLDataSource() string {
	return a.MySQL.DataSource
}

func (a *configAdapter) IsAutoMigrate() bool {
	return a.MySQL.AutoMigrate
}
