package svc

import (
	"gozero-ddd/internal/application/command"
	"gozero-ddd/internal/application/query"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/service"
	"gozero-ddd/internal/infrastructure/config"
	"gozero-ddd/internal/infrastructure/persistence"
)

// ServiceContext 服务上下文
// go-zero 使用 ServiceContext 来管理依赖注入
// 这是 go-zero 框架的核心设计模式之一
type ServiceContext struct {
	Config config.Config

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
	// 初始化仓储（这里使用内存实现，生产环境应使用数据库实现）
	kbRepo := persistence.NewMemoryKnowledgeBaseRepository()
	docRepo := persistence.NewMemoryDocumentRepository()

	// 初始化领域服务
	knowledgeService := service.NewKnowledgeService(kbRepo, docRepo)

	return &ServiceContext{
		Config: c,

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

