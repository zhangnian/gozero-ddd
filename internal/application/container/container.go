package container

import (
	"log"

	"gozero-ddd/internal/application/command"
	"gozero-ddd/internal/application/query"
	"gozero-ddd/internal/domain/event"
	"gozero-ddd/internal/domain/repository"
	"gozero-ddd/internal/domain/service"
)

// InfraDependencies 基础设施层依赖接口
// 定义应用层所需的基础设施层依赖
// 通过接口隔离，应用层不直接依赖基础设施层的具体实现
type InfraDependencies interface {
	GetUnitOfWork() repository.UnitOfWork
	GetEventBus() event.EventPublisher
	GetKnowledgeBaseRepo() repository.KnowledgeBaseRepository
	GetDocumentRepo() repository.DocumentRepository
	GetKnowledgeService() *service.KnowledgeService
}

// ApplicationContainer 应用层容器
// 负责管理所有应用层的组件：命令处理器、查询处理器
// 这是接口层唯一应该依赖的容器
type ApplicationContainer struct {
	// 命令处理器（写操作）
	Commands *CommandHandlers

	// 查询处理器（读操作）
	Queries *QueryHandlers
}

// CommandHandlers 命令处理器集合
// CQRS 模式中的 Command 端
type CommandHandlers struct {
	CreateKnowledgeBase *command.CreateKnowledgeBaseHandler
	UpdateKnowledgeBase *command.UpdateKnowledgeBaseHandler
	DeleteKnowledgeBase *command.DeleteKnowledgeBaseHandler
	AddDocument         *command.AddDocumentHandler
	RemoveDocument      *command.RemoveDocumentHandler
	MergeKnowledgeBases *command.MergeKnowledgeBasesHandler
}

// QueryHandlers 查询处理器集合
// CQRS 模式中的 Query 端
type QueryHandlers struct {
	GetKnowledgeBase   *query.GetKnowledgeBaseHandler
	ListKnowledgeBases *query.ListKnowledgeBasesHandler
	ListDocuments      *query.ListDocumentsHandler
}

// NewApplicationContainer 创建应用层容器
// 参数为基础设施层依赖，实现依赖注入
func NewApplicationContainer(deps InfraDependencies) *ApplicationContainer {
	container := &ApplicationContainer{
		Commands: &CommandHandlers{},
		Queries:  &QueryHandlers{},
	}

	// 初始化命令处理器
	container.initCommandHandlers(deps)

	// 初始化查询处理器
	container.initQueryHandlers(deps)

	log.Println("✅ [Application] 应用层容器初始化完成")

	return container
}

// initCommandHandlers 初始化所有命令处理器
func (c *ApplicationContainer) initCommandHandlers(deps InfraDependencies) {
	uow := deps.GetUnitOfWork()
	eventBus := deps.GetEventBus()
	kbRepo := deps.GetKnowledgeBaseRepo()
	docRepo := deps.GetDocumentRepo()
	kbService := deps.GetKnowledgeService()

	// 创建知识库
	c.Commands.CreateKnowledgeBase = command.NewCreateKnowledgeBaseHandler(kbService, eventBus)

	// 更新知识库
	c.Commands.UpdateKnowledgeBase = command.NewUpdateKnowledgeBaseHandler(kbRepo, eventBus)

	// 删除知识库
	c.Commands.DeleteKnowledgeBase = command.NewDeleteKnowledgeBaseHandler(kbRepo, kbService)

	// 添加文档
	c.Commands.AddDocument = command.NewAddDocumentHandler(uow, kbRepo, docRepo, eventBus)

	// 删除文档
	c.Commands.RemoveDocument = command.NewRemoveDocumentHandler(uow, kbRepo, docRepo)

	// 合并知识库
	c.Commands.MergeKnowledgeBases = command.NewMergeKnowledgeBasesHandler(uow, kbRepo, docRepo)

	log.Println("📝 [Application] 命令处理器初始化完成")
}

// initQueryHandlers 初始化所有查询处理器
func (c *ApplicationContainer) initQueryHandlers(deps InfraDependencies) {
	kbRepo := deps.GetKnowledgeBaseRepo()
	docRepo := deps.GetDocumentRepo()

	// 获取知识库详情
	c.Queries.GetKnowledgeBase = query.NewGetKnowledgeBaseHandler(kbRepo, docRepo)

	// 列出所有知识库
	c.Queries.ListKnowledgeBases = query.NewListKnowledgeBasesHandler(kbRepo)

	// 列出文档
	c.Queries.ListDocuments = query.NewListDocumentsHandler(docRepo)

	log.Println("🔍 [Application] 查询处理器初始化完成")
}
