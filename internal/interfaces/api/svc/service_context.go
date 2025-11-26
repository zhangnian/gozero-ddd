package svc

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gozero-ddd/internal/application/command"
	"gozero-ddd/internal/application/eventhandler"
	"gozero-ddd/internal/application/query"
	"gozero-ddd/internal/domain/event"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/service"
	"gozero-ddd/internal/infrastructure/config"
	"gozero-ddd/internal/infrastructure/eventbus"
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

	// 工作单元（事务管理）
	UnitOfWork repository.UnitOfWork

	// 事件总线（领域事件发布与订阅）
	EventBus event.EventBus

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
	MergeKnowledgeBasesHandler *command.MergeKnowledgeBasesHandler

	// 查询处理器
	GetKnowledgeBaseHandler   *query.GetKnowledgeBaseHandler
	ListKnowledgeBasesHandler *query.ListKnowledgeBasesHandler
	ListDocumentsHandler      *query.ListDocumentsHandler
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	var db *gorm.DB
	var uow repository.UnitOfWork
	var kbRepo repository.KnowledgeBaseRepository
	var docRepo repository.DocumentRepository

	// 根据配置选择仓储实现
	if c.UseMemory {
		// 使用内存仓储（开发测试用）
		log.Println("📦 使用内存存储")
		kbRepo = persistence.NewMemoryKnowledgeBaseRepository()
		docRepo = persistence.NewMemoryDocumentRepository()
		// 内存模式下使用空的工作单元
		uow = persistence.NewMemoryUnitOfWork()
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

		// 创建工作单元（事务管理）
		uow = persistence.NewGormUnitOfWork(db)

		// 先创建文档仓储
		docRepo = persistence.NewGormDocumentRepository(db)
		// 知识库仓储需要文档仓储来加载关联数据
		kbRepo = persistence.NewGormKnowledgeBaseRepository(db, docRepo)
	}

	// ==================== 初始化领域事件系统 ====================
	// 创建事件总线
	evtBus := eventbus.NewMemoryEventBus()

	// 注册事件处理器
	// 这些处理器会在领域事件发布时被调用
	registerEventHandlers(evtBus)

	log.Println("📫 领域事件系统初始化完成")

	// 初始化领域服务
	knowledgeService := service.NewKnowledgeService(kbRepo, docRepo)

	return &ServiceContext{
		Config:     c,
		DB:         db,
		UnitOfWork: uow,
		EventBus:   evtBus,

		// 仓储
		KnowledgeBaseRepo: kbRepo,
		DocumentRepo:      docRepo,

		// 领域服务
		KnowledgeService: knowledgeService,

		// 命令处理器（注入事件发布器）
		CreateKnowledgeBaseHandler: command.NewCreateKnowledgeBaseHandler(knowledgeService, evtBus),
		UpdateKnowledgeBaseHandler: command.NewUpdateKnowledgeBaseHandler(kbRepo, evtBus),
		DeleteKnowledgeBaseHandler: command.NewDeleteKnowledgeBaseHandler(kbRepo, knowledgeService),
		AddDocumentHandler:         command.NewAddDocumentHandler(uow, kbRepo, docRepo, evtBus),
		RemoveDocumentHandler:      command.NewRemoveDocumentHandler(uow, kbRepo, docRepo),
		MergeKnowledgeBasesHandler: command.NewMergeKnowledgeBasesHandler(uow, kbRepo, docRepo),

		// 查询处理器
		GetKnowledgeBaseHandler:   query.NewGetKnowledgeBaseHandler(kbRepo, docRepo),
		ListKnowledgeBasesHandler: query.NewListKnowledgeBasesHandler(kbRepo),
		ListDocumentsHandler:      query.NewListDocumentsHandler(docRepo),
	}
}

// registerEventHandlers 注册所有事件处理器
// 在应用启动时调用，将处理器注册到事件总线
func registerEventHandlers(evtBus event.EventBus) {
	// ==================== 知识库相关事件处理器 ====================

	// 1. 注册知识库创建事件处理器
	kbCreatedHandler := eventhandler.NewKnowledgeBaseCreatedHandler()
	evtBus.Subscribe(kbCreatedHandler.EventName(), kbCreatedHandler)

	// 2. 注册知识库更新事件处理器
	kbUpdatedHandler := eventhandler.NewKnowledgeBaseUpdatedHandler()
	evtBus.Subscribe(kbUpdatedHandler.EventName(), kbUpdatedHandler)

	// ==================== 文档相关事件处理器 ====================

	// 3. 注册文档添加事件处理器
	docAddedHandler := eventhandler.NewDocumentAddedHandler()
	evtBus.Subscribe(docAddedHandler.EventName(), docAddedHandler)

	// 4. 注册文档删除事件处理器
	docRemovedHandler := eventhandler.NewDocumentRemovedHandler()
	evtBus.Subscribe(docRemovedHandler.EventName(), docRemovedHandler)

	// ==================== 全局事件处理器 ====================

	// 5. 注册审计日志处理器（全局处理器，处理所有事件）
	auditLogHandler := eventhandler.NewAuditLogHandler()
	evtBus.SubscribeAll(auditLogHandler)
}
