package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"gozero-ddd/internal/application/command"
	"gozero-ddd/internal/interfaces"
	"gozero-ddd/internal/interfaces/rpc/pb"
	"gozero-ddd/internal/interfaces/rpc/svc"
)

// CreateKnowledgeBaseLogic 创建知识库逻辑
// 演示：在 gRPC 服务中使用 DDD 的 Command Handler
type CreateKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewCreateKnowledgeBaseLogic 创建逻辑实例
func NewCreateKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateKnowledgeBaseLogic {
	return &CreateKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateKnowledgeBase 创建知识库
// 演示：如何在 gRPC 服务中使用 DDD 的 Command Handler
// 流程：gRPC Request -> Logic -> Command Handler -> Domain Service -> Repository -> Domain Entity
//
// DDD 分层职责说明：
// 1. 接口层（本文件）：接收请求，调用应用层，转换响应和错误
// 2. 应用层（Command Handler）：编排业务用例，参数验证
// 3. 领域层（Domain Service）：核心业务逻辑，如"名称不能重复"
// 4. 基础设施层（Repository）：数据持久化
func (l *CreateKnowledgeBaseLogic) CreateKnowledgeBase(req *pb.CreateKnowledgeBaseRequest) (*pb.CreateKnowledgeBaseResponse, error) {
	l.Logger.Infof("📥 [gRPC] CreateKnowledgeBase 请求: name=%s", req.Name)

	// 构建命令对象（CQRS 模式中的 Command）
	// Command 是一个值对象，代表一个写操作的意图
	cmd := &command.CreateKnowledgeBaseCommand{
		Name:        req.Name,
		Description: req.Description,
	}

	// 调用应用层的 Command Handler
	// Command Handler 职责：
	// - 验证参数格式
	// - 调用领域服务执行业务逻辑
	// - 领域服务会验证业务规则（如名称唯一性）
	// - 通过仓储持久化领域实体
	// - 返回 DTO（而非领域实体，保护领域层封装）
	result, err := l.svcCtx.CreateKnowledgeBaseHandler.Handle(l.ctx, cmd)
	if err != nil {
		l.Logger.Errorf("❌ 创建知识库失败: %v", err)
		// 使用统一的错误转换函数
		return nil, interfaces.ToGrpcError(err)
	}

	// 将 DTO 转换为 gRPC 响应
	resp := &pb.CreateKnowledgeBaseResponse{
		KnowledgeBase: &pb.KnowledgeBase{
			Id:            result.ID,
			Name:          result.Name,
			Description:   result.Description,
			DocumentCount: int32(result.DocumentCount),
			CreatedAt:     result.CreatedAt.Unix(),
			UpdatedAt:     result.UpdatedAt.Unix(),
		},
	}

	l.Logger.Infof("✅ [gRPC] CreateKnowledgeBase 成功: id=%s, name=%s", result.ID, result.Name)
	return resp, nil
}
