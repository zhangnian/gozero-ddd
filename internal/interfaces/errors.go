package interfaces

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gozero-ddd/internal/domain"
	"gozero-ddd/internal/domain/valueobject"
)

// ToGrpcError 将领域错误转换为 gRPC 错误
// 根据错误类型返回适当的 gRPC 状态码
func ToGrpcError(err error) error {
	if err == nil {
		return nil
	}

	// 检查是否为"未找到"类型的错误
	if domain.IsNotFoundError(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	// 检查是否为验证错误
	if domain.IsValidationError(err) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// 检查是否为冲突错误
	if domain.IsConflictError(err) {
		return status.Error(codes.AlreadyExists, err.Error())
	}

	// 检查值对象验证错误
	if errors.Is(err, valueobject.ErrInvalidKnowledgeBaseID) ||
		errors.Is(err, valueobject.ErrInvalidDocumentID) ||
		errors.Is(err, valueobject.ErrEmptyID) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// 默认返回内部错误
	return status.Error(codes.Internal, err.Error())
}

// HTTPErrorCode 根据错误类型返回适当的 HTTP 状态码
func HTTPErrorCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// 检查是否为"未找到"类型的错误
	if domain.IsNotFoundError(err) {
		return http.StatusNotFound
	}

	// 检查是否为验证错误
	if domain.IsValidationError(err) {
		return http.StatusBadRequest
	}

	// 检查是否为冲突错误
	if domain.IsConflictError(err) {
		return http.StatusConflict
	}

	// 检查值对象验证错误
	if errors.Is(err, valueobject.ErrInvalidKnowledgeBaseID) ||
		errors.Is(err, valueobject.ErrInvalidDocumentID) ||
		errors.Is(err, valueobject.ErrEmptyID) {
		return http.StatusBadRequest
	}

	// 默认返回内部服务器错误
	return http.StatusInternalServerError
}

