package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gozero-ddd/internal/application/command"
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
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(400, err.Error()))
		return
	}

	cmd := &command.MergeKnowledgeBasesCommand{
		SourceID: req.SourceID,
		TargetID: req.TargetID,
	}

	result, err := h.svcCtx.MergeKnowledgeBasesHandler.Handle(r.Context(), cmd)
	if err != nil {
		// 根据错误类型返回不同的状态码
		switch err {
		case command.ErrSourceKnowledgeBaseNotFound, command.ErrTargetKnowledgeBaseNotFound:
			httpx.WriteJson(w, http.StatusNotFound, types.NewErrorResponse(404, err.Error()))
		case command.ErrCannotMergeSameKnowledgeBase:
			httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(400, err.Error()))
		default:
			httpx.WriteJson(w, http.StatusInternalServerError, types.NewErrorResponse(500, err.Error()))
		}
		return
	}

	httpx.WriteJson(w, http.StatusOK, types.NewSuccessResponse(result))
}

