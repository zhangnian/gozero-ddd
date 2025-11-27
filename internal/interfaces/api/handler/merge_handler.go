package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gozero-ddd/internal/application/command"
	"gozero-ddd/internal/interfaces"
	"gozero-ddd/internal/interfaces/api/svc"
	"gozero-ddd/internal/interfaces/api/types"
)

// MergeHandler 合并处理器
type MergeHandler struct {
	svcCtx *svc.ServiceContext
}

// NewMergeHandler 创建合并处理器
func NewMergeHandler(svcCtx *svc.ServiceContext) *MergeHandler {
	return &MergeHandler{svcCtx: svcCtx}
}

// MergeKnowledgeBases 合并知识库
// POST /api/v1/knowledge/merge
// 将源知识库的所有文档移动到目标知识库，然后删除源知识库
// 此操作在事务中执行，保证原子性
func (h *MergeHandler) MergeKnowledgeBases(w http.ResponseWriter, r *http.Request) {
	var req types.MergeKnowledgeBasesRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	cmd := &command.MergeKnowledgeBasesCommand{
		SourceID: req.SourceID,
		TargetID: req.TargetID,
	}

	// 通过应用层容器访问命令处理器
	result, err := h.svcCtx.App.Commands.MergeKnowledgeBases.Handle(r.Context(), cmd)
	if err != nil {
		// 使用统一的错误转换函数
		code := interfaces.HTTPErrorCode(err)
		httpx.WriteJson(w, code, types.NewErrorResponse(code, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusOK, types.NewSuccessResponse(result))
}
