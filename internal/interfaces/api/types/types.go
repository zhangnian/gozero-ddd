package types

// ========== 知识库相关请求/响应 ==========

// CreateKnowledgeBaseRequest 创建知识库请求
type CreateKnowledgeBaseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,optional"`
}

// UpdateKnowledgeBaseRequest 更新知识库请求
type UpdateKnowledgeBaseRequest struct {
	ID          string `path:"id"`
	Name        string `json:"name"`
	Description string `json:"description,optional"`
}

// GetKnowledgeBaseRequest 获取知识库请求
type GetKnowledgeBaseRequest struct {
	ID               string `path:"id"`
	IncludeDocuments bool   `form:"include_documents,optional"`
}

// DeleteKnowledgeBaseRequest 删除知识库请求
type DeleteKnowledgeBaseRequest struct {
	ID string `path:"id"`
}

// MergeKnowledgeBasesRequest 合并知识库请求
type MergeKnowledgeBasesRequest struct {
	SourceID string `json:"source_id"` // 源知识库ID（将被删除）
	TargetID string `json:"target_id"` // 目标知识库ID（保留）
}

// KnowledgeBaseResponse 知识库响应
type KnowledgeBaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ========== 文档相关请求/响应 ==========

// AddDocumentRequest 添加文档请求
type AddDocumentRequest struct {
	KnowledgeBaseID string   `path:"id"`
	Title           string   `json:"title"`
	Content         string   `json:"content"`
	Tags            []string `json:"tags,optional"`
}

// RemoveDocumentRequest 删除文档请求
type RemoveDocumentRequest struct {
	KnowledgeBaseID string `path:"id"`
	DocumentID      string `path:"doc_id"`
}

// ListDocumentsRequest 列出文档请求
type ListDocumentsRequest struct {
	KnowledgeBaseID string `path:"id"`
}

// DocumentResponse 文档响应
type DocumentResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ========== 通用响应 ==========

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
	}
}
