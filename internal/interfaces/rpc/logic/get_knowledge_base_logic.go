package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"gozero-ddd/internal/application/dto"
	"gozero-ddd/internal/application/query"
	"gozero-ddd/internal/interfaces"
	"gozero-ddd/internal/interfaces/rpc/pb"
	"gozero-ddd/internal/interfaces/rpc/svc"
)

// GetKnowledgeBaseLogic 获取知识库逻辑
// 在 go-zero 中，每个 RPC 方法对应一个 Logic 结构
// Logic 负责协调应用层（Command/Query Handler）完成业务逻辑
type GetKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewGetKnowledgeBaseLogic 创建获取知识库逻辑
func NewGetKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetKnowledgeBaseLogic {
	return &GetKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetKnowledgeBase 获取知识库详情
// 演示：如何在 gRPC 服务中使用 DDD 的 Query Handler
// 流程：gRPC Request -> Logic -> Query Handler -> Repository -> Domain Entity -> DTO -> gRPC Response
func (l *GetKnowledgeBaseLogic) GetKnowledgeBase(req *pb.GetKnowledgeBaseRequest) (*pb.GetKnowledgeBaseResponse, error) {
	l.Logger.Infof("📥 [gRPC] GetKnowledgeBase 请求: id=%s, includeDocuments=%v", req.Id, req.IncludeDocuments)

	// 构建查询对象（CQRS 模式中的 Query）
	qry := &query.GetKnowledgeBaseQuery{
		ID:               req.Id,
		IncludeDocuments: req.IncludeDocuments,
	}

	// 通过应用层容器访问查询处理器
	// Query Handler 负责：
	// - 验证参数格式
	// - 通过仓储获取领域实体
	// - 将领域实体转换为 DTO
	result, err := l.svcCtx.App.Queries.GetKnowledgeBase.Handle(l.ctx, qry)
	if err != nil {
		l.Logger.Errorf("❌ 获取知识库失败: %v", err)
		// 使用统一的错误转换函数
		return nil, interfaces.ToGrpcError(err)
	}

	// 将 DTO 转换为 gRPC 响应
	// 注意：这里进行了 DTO -> Protobuf 的转换
	// 这种转换保持了各层之间的解耦
	resp := &pb.GetKnowledgeBaseResponse{
		KnowledgeBase: convertToProtoKnowledgeBase(result),
	}

	l.Logger.Infof("✅ [gRPC] GetKnowledgeBase 成功: name=%s", result.Name)
	return resp, nil
}

// convertToProtoKnowledgeBase 将 DTO 转换为 Protobuf 消息
// 这个转换函数放在接口层，因为它是接口层特有的转换逻辑
func convertToProtoKnowledgeBase(d *dto.KnowledgeBaseDTO) *pb.KnowledgeBase {
	kb := &pb.KnowledgeBase{
		Id:            d.ID,
		Name:          d.Name,
		Description:   d.Description,
		DocumentCount: int32(d.DocumentCount),
		CreatedAt:     d.CreatedAt.Unix(),
		UpdatedAt:     d.UpdatedAt.Unix(),
	}

	// 转换文档列表
	if len(d.Documents) > 0 {
		kb.Documents = make([]*pb.Document, len(d.Documents))
		for i, doc := range d.Documents {
			kb.Documents[i] = &pb.Document{
				Id:              doc.ID,
				KnowledgeBaseId: doc.KnowledgeBaseID,
				Title:           doc.Title,
				Content:         doc.Content,
				Tags:            doc.Tags,
				CreatedAt:       doc.CreatedAt.Unix(),
				UpdatedAt:       doc.UpdatedAt.Unix(),
			}
		}
	}

	return kb
}
