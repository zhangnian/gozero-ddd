package server

import (
	"context"

	"gozero-ddd/internal/interfaces/rpc/logic"
	"gozero-ddd/internal/interfaces/rpc/pb"
	"gozero-ddd/internal/interfaces/rpc/svc"
)

// KnowledgeServer gRPC 服务实现
// 这是 gRPC 服务的入口点，实现了 pb.KnowledgeServiceServer 接口
// 在 go-zero 的 DDD 架构中，Server 层负责：
// 1. 实现 gRPC 接口
// 2. 创建并调用对应的 Logic
// 3. Logic 负责实际的业务逻辑协调
type KnowledgeServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedKnowledgeServiceServer
}

// NewKnowledgeServer 创建 gRPC 服务实例
func NewKnowledgeServer(svcCtx *svc.ServiceContext) *KnowledgeServer {
	return &KnowledgeServer{
		svcCtx: svcCtx,
	}
}

// GetKnowledgeBase 获取知识库详情
// 实现 pb.KnowledgeServiceServer 接口
// 每个请求创建一个新的 Logic 实例，确保并发安全
func (s *KnowledgeServer) GetKnowledgeBase(ctx context.Context, req *pb.GetKnowledgeBaseRequest) (*pb.GetKnowledgeBaseResponse, error) {
	// 创建 Logic 实例
	// Logic 是请求级别的，每个请求一个新实例
	l := logic.NewGetKnowledgeBaseLogic(ctx, s.svcCtx)
	return l.GetKnowledgeBase(req)
}

// CreateKnowledgeBase 创建知识库
// 实现 pb.KnowledgeServiceServer 接口
func (s *KnowledgeServer) CreateKnowledgeBase(ctx context.Context, req *pb.CreateKnowledgeBaseRequest) (*pb.CreateKnowledgeBaseResponse, error) {
	// 创建 Logic 实例
	l := logic.NewCreateKnowledgeBaseLogic(ctx, s.svcCtx)
	return l.CreateKnowledgeBase(req)
}

