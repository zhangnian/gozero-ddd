package routes

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"gozero-ddd/internal/interfaces/api/handler"
	"gozero-ddd/internal/interfaces/api/middleware"
	"gozero-ddd/internal/interfaces/api/svc"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(server *rest.Server, svcCtx *svc.ServiceContext) {
	// 创建处理器
	kbHandler := handler.NewKnowledgeBaseHandler(svcCtx)
	docHandler := handler.NewDocumentHandler(svcCtx)
	mergeHandler := handler.NewMergeHandler(svcCtx)

	// 创建中间件
	loggingMiddleware := middleware.NewLoggingMiddleware()

	// 注册知识库相关路由
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{loggingMiddleware.Handle},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/knowledge",
					Handler: kbHandler.Create,
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/knowledge",
					Handler: kbHandler.List,
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/knowledge/:id",
					Handler: kbHandler.Get,
				},
				{
					Method:  http.MethodPut,
					Path:    "/api/v1/knowledge/:id",
					Handler: kbHandler.Update,
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/knowledge/:id",
					Handler: kbHandler.Delete,
				},
				// 合并知识库（事务演示）
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/knowledge/merge",
					Handler: mergeHandler.MergeKnowledgeBases,
				},
			}...,
		),
	)

	// 注册文档相关路由
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{loggingMiddleware.Handle},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/knowledge/:id/documents",
					Handler: docHandler.Add,
				},
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/knowledge/:id/documents",
					Handler: docHandler.List,
				},
				{
					Method:  http.MethodDelete,
					Path:    "/api/v1/knowledge/:id/documents/:doc_id",
					Handler: docHandler.Remove,
				},
			}...,
		),
	)
}
