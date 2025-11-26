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

// ServiceContext 服务上下文
// go-zero 使用 ServiceContext 来管理依赖注入
// 这是 go-zero 框架的核心设计模式之一
type ServiceContext struct {
	Config config.Config

	// 数据库连接
	DB *gorm.DB

	// 仓储
	KnowledgeBaseRepo repository.KnowledgeBaseRepository
	DocumentRepo      repository.DocumentRepository

	// 领域服务
	KnowledgeService *service.KnowledgeService

	// 命令处理器
	CreateKnowledgeBaseHandler *command.CreateKnowledgeBaseHandler
	UpdateKnowledgeBaseHandler *command.UpdateKnowledgeBaseHandler
	DeleteKnowledgeBaseHandler *command.DeleteKnowledgeBaseHandler
	AddDocumentHandler         *command.AddDocumentHandler
	RemoveDocumentHandler      *command.RemoveDocumentHandler

	// 查询处理器
	GetKnowledgeBaseHandler   *query.GetKnowledgeBaseHandler
	ListKnowledgeBasesHandler *query.ListKnowledgeBasesHandler
	ListDocumentsHandler      *query.ListDocumentsHandler
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	var db *gorm.DB
	var kbRepo repository.KnowledgeBaseRepository
	var docRepo repository.DocumentRepository

	// 根据配置选择仓储实现
	if c.UseMemory {
		// 使用内存仓储（开发测试用）
		log.Println("📦 使用内存存储")
		kbRepo = persistence.NewMemoryKnowledgeBaseRepository()
		docRepo = persistence.NewMemoryDocumentRepository()
	} else {
		// 使用 GORM + MySQL（生产环境）
		log.Println("📦 使用 MySQL 存储 (GORM)")
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
			log.Println("🔄 自动迁移数据库表结构...")
			if err := db.AutoMigrate(&model.KnowledgeBaseModel{}, &model.DocumentModel{}); err != nil {
				log.Fatalf("❌ 数据库迁移失败: %v", err)
			}
		}

		// 先创建文档仓储
		docRepo = persistence.NewGormDocumentRepository(db)
		// 知识库仓储需要文档仓储来加载关联数据
		kbRepo = persistence.NewGormKnowledgeBaseRepository(db, docRepo)
	}

	// 初始化领域服务
	knowledgeService := service.NewKnowledgeService(kbRepo, docRepo)

	return &ServiceContext{
		Config: c,
		DB:     db,

		// 仓储
		KnowledgeBaseRepo: kbRepo,
		DocumentRepo:      docRepo,

		// 领域服务
		KnowledgeService: knowledgeService,

		// 命令处理器
		CreateKnowledgeBaseHandler: command.NewCreateKnowledgeBaseHandler(knowledgeService),
		UpdateKnowledgeBaseHandler: command.NewUpdateKnowledgeBaseHandler(kbRepo),
		DeleteKnowledgeBaseHandler: command.NewDeleteKnowledgeBaseHandler(kbRepo, knowledgeService),
		AddDocumentHandler:         command.NewAddDocumentHandler(kbRepo, docRepo),
		RemoveDocumentHandler:      command.NewRemoveDocumentHandler(kbRepo, docRepo),

		// 查询处理器
		GetKnowledgeBaseHandler:   query.NewGetKnowledgeBaseHandler(kbRepo, docRepo),
		ListKnowledgeBasesHandler: query.NewListKnowledgeBasesHandler(kbRepo),
		ListDocumentsHandler:      query.NewListDocumentsHandler(docRepo),
	}
}
