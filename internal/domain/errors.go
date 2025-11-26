package domain

import "errors"

// 领域层统一错误定义
// 这些错误可以被应用层和接口层识别并转换为适当的响应

var (
	// 知识库相关错误
	ErrKnowledgeBaseNotFound   = errors.New("knowledge base not found")
	ErrKnowledgeBaseNameExists = errors.New("knowledge base name already exists")
	ErrKnowledgeBaseNameEmpty  = errors.New("knowledge base name cannot be empty")

	// 文档相关错误
	ErrDocumentNotFound     = errors.New("document not found")
	ErrDocumentTitleEmpty   = errors.New("document title cannot be empty")
	ErrDocumentContentEmpty = errors.New("document content cannot be empty")

	// 操作相关错误
	ErrCannotMergeSameKnowledgeBase = errors.New("cannot merge knowledge base with itself")
)

// DomainError 领域错误接口
// 用于标识领域层产生的错误，便于上层进行错误类型判断
type DomainError interface {
	error
	IsDomainError() bool
}

// domainError 领域错误实现
type domainError struct {
	message string
	code    string
}

func (e *domainError) Error() string {
	return e.message
}

func (e *domainError) IsDomainError() bool {
	return true
}

func (e *domainError) Code() string {
	return e.code
}

// NewDomainError 创建领域错误
func NewDomainError(code, message string) *domainError {
	return &domainError{
		code:    code,
		message: message,
	}
}

// IsDomainError 判断是否为领域错误
func IsDomainError(err error) bool {
	var de DomainError
	return errors.As(err, &de)
}

// IsNotFoundError 判断是否为"未找到"类型的错误
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrKnowledgeBaseNotFound) ||
		errors.Is(err, ErrDocumentNotFound)
}

// IsValidationError 判断是否为验证错误
func IsValidationError(err error) bool {
	return errors.Is(err, ErrKnowledgeBaseNameEmpty) ||
		errors.Is(err, ErrDocumentTitleEmpty) ||
		errors.Is(err, ErrDocumentContentEmpty)
}

// IsConflictError 判断是否为冲突错误
func IsConflictError(err error) bool {
	return errors.Is(err, ErrKnowledgeBaseNameExists) ||
		errors.Is(err, ErrCannotMergeSameKnowledgeBase)
}

