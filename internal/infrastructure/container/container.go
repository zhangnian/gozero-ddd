package container

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gozero-ddd/internal/application/eventhandler"
	"gozero-ddd/internal/domain/event"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/service"
	"gozero-ddd/internal/infrastructure/eventbus"
	"gozero-ddd/internal/infrastructure/persistence"
	"gozero-ddd/internal/infrastructure/persistence/model"
)

// InfraConfig 基础设施配置接口
// 定义基础设施层初始化所需的配置
type InfraConfig interface {
	IsUseMemory() bool
	GetMySQLDataSource() string
	IsAutoMigrate() bool
}

// InfrastructureContainer 基础设施层容器
// 负责管理所有基础设施层的组件：数据库、仓储、事件总线等
// 这些组件对上层（应用层）是透明的，上层只依赖接口
type InfrastructureContainer struct {
	// 数据库连接（内部使用，不对外暴露）
	db *gorm.DB

	// 工作单元（事务管理）
	UnitOfWork repository.UnitOfWork

	// 事件总线
	EventBus event.EventBus

	// 仓储接口（注意：这里是接口类型，不是具体实现）
	KnowledgeBaseRepo repository.KnowledgeBaseRepository
	DocumentRepo      repository.DocumentRepository

	// 领域服务（领域层，但由基础设施层组装）
	KnowledgeService *service.KnowledgeService
}

// NewInfrastructureContainer 创建基础设施层容器
// 负责初始化所有基础设施组件
func NewInfrastructureContainer(cfg InfraConfig) *InfrastructureContainer {
	container := &InfrastructureContainer{}

	// 1. 初始化存储层（仓储和工作单元）
	container.initStorage(cfg)

	// 2. 初始化事件总线
	container.initEventBus()

	// 3. 初始化领域服务
	container.initDomainServices()

	return container
}

// initStorage 初始化存储层
func (c *InfrastructureContainer) initStorage(cfg InfraConfig) {
	if cfg.IsUseMemory() {
		// 使用内存仓储（开发测试用）
		log.Println("📦 [Infrastructure] 使用内存存储")
		c.KnowledgeBaseRepo = persistence.NewMemoryKnowledgeBaseRepository()
		c.DocumentRepo = persistence.NewMemoryDocumentRepository()
		c.UnitOfWork = persistence.NewMemoryUnitOfWork()
	} else {
		// 使用 GORM + MySQL（生产环境）
		log.Println("📦 [Infrastructure] 使用 MySQL 存储 (GORM)")

		dataSource := cfg.GetMySQLDataSource()
		if dataSource == "" {
			log.Fatal("❌ MySQL DataSource 未配置")
		}

		// 创建 GORM 数据库连接
		var err error
		c.db, err = gorm.Open(mysql.Open(dataSource), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("❌ 连接数据库失败: %v", err)
		}

		// 自动迁移表结构（开发环境使用）
		if cfg.IsAutoMigrate() {
			log.Println("🔄 [Infrastructure] 自动迁移数据库表结构...")
			if err := c.db.AutoMigrate(&model.KnowledgeBaseModel{}, &model.DocumentModel{}); err != nil {
				log.Fatalf("❌ 数据库迁移失败: %v", err)
			}
		}

		// 创建工作单元（事务管理）
		c.UnitOfWork = persistence.NewGormUnitOfWork(c.db)

		// 创建仓储实例
		c.DocumentRepo = persistence.NewGormDocumentRepository(c.db)
		c.KnowledgeBaseRepo = persistence.NewGormKnowledgeBaseRepository(c.db, c.DocumentRepo)
	}

	log.Println("✅ [Infrastructure] 存储层初始化完成")
}

// initEventBus 初始化事件总线
func (c *InfrastructureContainer) initEventBus() {
	// 创建事件总线
	c.EventBus = eventbus.NewMemoryEventBus()

	// 注册事件处理器
	c.registerEventHandlers()

	log.Println("✅ [Infrastructure] 事件总线初始化完成")
}

// registerEventHandlers 注册所有事件处理器
func (c *InfrastructureContainer) registerEventHandlers() {
	// 知识库创建事件处理器
	kbCreatedHandler := eventhandler.NewKnowledgeBaseCreatedHandler()
	c.EventBus.Subscribe(kbCreatedHandler.EventName(), kbCreatedHandler)

	// 知识库更新事件处理器
	kbUpdatedHandler := eventhandler.NewKnowledgeBaseUpdatedHandler()
	c.EventBus.Subscribe(kbUpdatedHandler.EventName(), kbUpdatedHandler)

	// 文档添加事件处理器
	docAddedHandler := eventhandler.NewDocumentAddedHandler()
	c.EventBus.Subscribe(docAddedHandler.EventName(), docAddedHandler)

	// 文档删除事件处理器
	docRemovedHandler := eventhandler.NewDocumentRemovedHandler()
	c.EventBus.Subscribe(docRemovedHandler.EventName(), docRemovedHandler)

	// 审计日志处理器（全局处理器，处理所有事件）
	auditLogHandler := eventhandler.NewAuditLogHandler()
	c.EventBus.SubscribeAll(auditLogHandler)

	log.Println("📫 [Infrastructure] 事件处理器注册完成")
}

// initDomainServices 初始化领域服务
func (c *InfrastructureContainer) initDomainServices() {
	c.KnowledgeService = service.NewKnowledgeService(c.KnowledgeBaseRepo, c.DocumentRepo)
	log.Println("✅ [Infrastructure] 领域服务初始化完成")
}

// Close 关闭基础设施资源
func (c *InfrastructureContainer) Close() error {
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// ==================== 实现 InfraDependencies 接口 ====================
// 这些方法用于向应用层提供依赖，而不是直接暴露内部字段

// GetUnitOfWork 获取工作单元
func (c *InfrastructureContainer) GetUnitOfWork() repository.UnitOfWork {
	return c.UnitOfWork
}

// GetEventBus 获取事件发布器
func (c *InfrastructureContainer) GetEventBus() event.EventPublisher {
	return c.EventBus
}

// GetKnowledgeBaseRepo 获取知识库仓储
func (c *InfrastructureContainer) GetKnowledgeBaseRepo() repository.KnowledgeBaseRepository {
	return c.KnowledgeBaseRepo
}

// GetDocumentRepo 获取文档仓储
func (c *InfrastructureContainer) GetDocumentRepo() repository.DocumentRepository {
	return c.DocumentRepo
}

// GetKnowledgeService 获取知识库领域服务
func (c *InfrastructureContainer) GetKnowledgeService() *service.KnowledgeService {
	return c.KnowledgeService
}
