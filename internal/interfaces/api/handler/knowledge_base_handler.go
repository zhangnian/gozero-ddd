package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"gozero-ddd/internal/application/command"
	"gozero-ddd/internal/application/query"
	"gozero-ddd/internal/interfaces/api/svc"
	"gozero-ddd/internal/interfaces/api/types"
)

// KnowledgeBaseHandler 知识库处理器
type KnowledgeBaseHandler struct {
	svcCtx *svc.ServiceContext
}

// NewKnowledgeBaseHandler 创建知识库处理器
func NewKnowledgeBaseHandler(svcCtx *svc.ServiceContext) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{svcCtx: svcCtx}
}

// Create 创建知识库
func (h *KnowledgeBaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req types.CreateKnowledgeBaseRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(400, err.Error()))
		return
	}

	cmd := &command.CreateKnowledgeBaseCommand{
		Name:        req.Name,
		Description: req.Description,
	}

	result, err := h.svcCtx.CreateKnowledgeBaseHandler.Handle(r.Context(), cmd)
	if err != nil {
		httpx.WriteJson(w, http.StatusInternalServerError, types.NewErrorResponse(500, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusCreated, types.NewSuccessResponse(result))
}

// Update 更新知识库
func (h *KnowledgeBaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateKnowledgeBaseRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(400, err.Error()))
		return
	}

	cmd := &command.UpdateKnowledgeBaseCommand{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	result, err := h.svcCtx.UpdateKnowledgeBaseHandler.Handle(r.Context(), cmd)
	if err != nil {
		httpx.WriteJson(w, http.StatusInternalServerError, types.NewErrorResponse(500, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusOK, types.NewSuccessResponse(result))
}

// Get 获取知识库
func (h *KnowledgeBaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	var req types.GetKnowledgeBaseRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(400, err.Error()))
		return
	}

	qry := &query.GetKnowledgeBaseQuery{
		ID:               req.ID,
		IncludeDocuments: req.IncludeDocuments,
	}

	result, err := h.svcCtx.GetKnowledgeBaseHandler.Handle(r.Context(), qry)
	if err != nil {
		httpx.WriteJson(w, http.StatusNotFound, types.NewErrorResponse(404, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusOK, types.NewSuccessResponse(result))
}

// List 列出所有知识库
func (h *KnowledgeBaseHandler) List(w http.ResponseWriter, r *http.Request) {
	qry := &query.ListKnowledgeBasesQuery{}

	result, err := h.svcCtx.ListKnowledgeBasesHandler.Handle(r.Context(), qry)
	if err != nil {
		httpx.WriteJson(w, http.StatusInternalServerError, types.NewErrorResponse(500, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusOK, types.NewSuccessResponse(result))
}

// Delete 删除知识库
func (h *KnowledgeBaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var req types.DeleteKnowledgeBaseRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, types.NewErrorResponse(400, err.Error()))
		return
	}

	cmd := &command.DeleteKnowledgeBaseCommand{
		ID: req.ID,
	}

	if err := h.svcCtx.DeleteKnowledgeBaseHandler.Handle(r.Context(), cmd); err != nil {
		httpx.WriteJson(w, http.StatusInternalServerError, types.NewErrorResponse(500, err.Error()))
		return
	}

	httpx.WriteJson(w, http.StatusOK, types.NewSuccessResponse(nil))
}

