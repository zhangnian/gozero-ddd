package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gozero-ddd/internal/application/command"
	"gozero-ddd/internal/application/query"
	"gozero-ddd/internal/interfaces"
	"gozero-ddd/internal/interfaces/api/svc"
	"gozero-ddd/internal/interfaces/api/types"
)

// DocumentHandler 文档处理器
type DocumentHandler struct {
	svcCtx *svc.ServiceContext
}

// NewDocumentHandler 创建文档处理器
func NewDocumentHandler(svcCtx *svc.ServiceContext) *DocumentHandler {
	return &DocumentHandler{svcCtx: svcCtx}
}

// Add 添加文档
func (h *DocumentHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req types.AddDocumentRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	cmd := &command.AddDocumentCommand{
		KnowledgeBaseID: req.KnowledgeBaseID,
		Title:           req.Title,
		Content:         req.Content,
		Tags:            req.Tags,
	}

	// 通过应用层容器访问命令处理器
	result, err := h.svcCtx.App.Commands.AddDocument.Handle(r.Context(), cmd)
	if err != nil {
		code := interfaces.HTTPErrorCode(err)
		httpx.WriteJson(w, code, types.NewErrorResponse(code, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusCreated, types.NewSuccessResponse(result))
}

// List 列出文档
func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	var req types.ListDocumentsRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	qry := &query.ListDocumentsQuery{
		KnowledgeBaseID: req.KnowledgeBaseID,
	}

	// 通过应用层容器访问查询处理器
	result, err := h.svcCtx.App.Queries.ListDocuments.Handle(r.Context(), qry)
	if err != nil {
		code := interfaces.HTTPErrorCode(err)
		httpx.WriteJson(w, code, types.NewErrorResponse(code, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusOK, types.NewSuccessResponse(result))
}

// Remove 删除文档
func (h *DocumentHandler) Remove(w http.ResponseWriter, r *http.Request) {
	var req types.RemoveDocumentRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	cmd := &command.RemoveDocumentCommand{
		KnowledgeBaseID: req.KnowledgeBaseID,
		DocumentID:      req.DocumentID,
	}

	// 通过应用层容器访问命令处理器
	if err := h.svcCtx.App.Commands.RemoveDocument.Handle(r.Context(), cmd); err != nil {
		code := interfaces.HTTPErrorCode(err)
		httpx.WriteJson(w, code, types.NewErrorResponse(code, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusOK, types.NewSuccessResponse(nil))
}
