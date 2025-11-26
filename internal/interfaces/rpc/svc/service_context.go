package svc

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gozero-ddd/internal/application/command"
	"gozero-ddd/internal/application/query"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/service"
	"gozero-ddd/internal/infrastructure/config"
	"gozero-ddd/internal/infrastructure/persistence"
	"gozero-ddd/internal/infrastructure/persistence/model"
)

// ServiceContext gRPC 服务上下文
// go-zero 的依赖注入容器，管理所有服务依赖
// 与 REST API 的 ServiceContext 类似，但专门用于 gRPC 服务
type ServiceContext struct {
	Config config.RpcConfig

	// 数据库连接
	DB *gorm.DB

	// 工作单元（事务管理）
	UnitOfWork repository.UnitOfWork

	// 仓储层 - 负责数据持久化
	KnowledgeBaseRepo repository.KnowledgeBaseRepository
	DocumentRepo      repository.DocumentRepository

	// 领域服务 - 处理跨实体的业务逻辑
	KnowledgeService *service.KnowledgeService

	// 命令处理器 - 处理写操作（CQRS 模式中的 Command）
	CreateKnowledgeBaseHandler *command.CreateKnowledgeBaseHandler

	// 查询处理器 - 处理读操作（CQRS 模式中的 Query）
	GetKnowledgeBaseHandler *query.GetKnowledgeBaseHandler
}

// NewServiceContext 创建 gRPC 服务上下文
// 初始化所有依赖，实现依赖注入
func NewServiceContext(c config.RpcConfig) *ServiceContext {
	var db *gorm.DB
	var uow repository.UnitOfWork
	var kbRepo repository.KnowledgeBaseRepository
	var docRepo repository.DocumentRepository

	// 根据配置选择仓储实现（策略模式）
	if c.UseMemory {
		// 使用内存仓储（开发测试用）
		log.Println("📦 [gRPC] 使用内存存储")
		kbRepo = persistence.NewMemoryKnowledgeBaseRepository()
		docRepo = persistence.NewMemoryDocumentRepository()
		uow = persistence.NewMemoryUnitOfWork()
	} else {
		// 使用 GORM + MySQL（生产环境）
		log.Println("📦 [gRPC] 使用 MySQL 存储 (GORM)")
		if c.MySQL.DataSource == "" {
			log.Fatal("❌ MySQL DataSource 未配置")
		}

		// 创建 GORM 数据库连接
		var err error
		db, err = gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("❌ 连接数据库失败: %v", err)
		}

		// 自动迁移表结构（开发环境使用）
		if c.MySQL.AutoMigrate {
			log.Println("🔄 [gRPC] 自动迁移数据库表结构...")
			if err := db.AutoMigrate(&model.KnowledgeBaseModel{}, &model.DocumentModel{}); err != nil {
				log.Fatalf("❌ 数据库迁移失败: %v", err)
			}
		}

		// 创建工作单元（事务管理）
		uow = persistence.NewGormUnitOfWork(db)

		// 创建仓储实例
		docRepo = persistence.NewGormDocumentRepository(db)
		kbRepo = persistence.NewGormKnowledgeBaseRepository(db, docRepo)
	}

	// 初始化领域服务
	knowledgeService := service.NewKnowledgeService(kbRepo, docRepo)

	return &ServiceContext{
		Config:     c,
		DB:         db,
		UnitOfWork: uow,

		// 仓储
		KnowledgeBaseRepo: kbRepo,
		DocumentRepo:      docRepo,

		// 领域服务
		KnowledgeService: knowledgeService,

		// 命令处理器 - 用于 CreateKnowledgeBase RPC
		CreateKnowledgeBaseHandler: command.NewCreateKnowledgeBaseHandler(knowledgeService),

		// 查询处理器 - 用于 GetKnowledgeBase RPC
		GetKnowledgeBaseHandler: query.NewGetKnowledgeBaseHandler(kbRepo, docRepo),
	}
}
