package valueobject

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidKnowledgeBaseID = errors.New("invalid knowledge base ID format")
	ErrInvalidDocumentID      = errors.New("invalid document ID format")
	ErrEmptyID                = errors.New("ID cannot be empty")
)

// KnowledgeBaseID 知识库ID值对象
// 值对象是不可变的，没有唯一标识，通过值来比较
type KnowledgeBaseID string

// NewKnowledgeBaseID 创建新的知识库ID
func NewKnowledgeBaseID() KnowledgeBaseID {
	return KnowledgeBaseID(uuid.New().String())
}

// KnowledgeBaseIDFromString 从字符串创建知识库ID（带验证）
func KnowledgeBaseIDFromString(s string) (KnowledgeBaseID, error) {
	if s == "" {
		return "", ErrEmptyID
	}
	if _, err := uuid.Parse(s); err != nil {
		return "", ErrInvalidKnowledgeBaseID
	}
	return KnowledgeBaseID(s), nil
}

// MustKnowledgeBaseIDFromString 从字符串创建知识库ID（不验证，用于从数据库重建）
// 仅在确定数据来源可靠时使用（如从数据库读取）
func MustKnowledgeBaseIDFromString(s string) KnowledgeBaseID {
	return KnowledgeBaseID(s)
}

// String 转换为字符串
func (id KnowledgeBaseID) String() string {
	return string(id)
}

// IsEmpty 判断ID是否为空
func (id KnowledgeBaseID) IsEmpty() bool {
	return string(id) == ""
}

// DocumentID 文档ID值对象
type DocumentID string

// NewDocumentID 创建新的文档ID
func NewDocumentID() DocumentID {
	return DocumentID(uuid.New().String())
}

// DocumentIDFromString 从字符串创建文档ID（带验证）
func DocumentIDFromString(s string) (DocumentID, error) {
	if s == "" {
		return "", ErrEmptyID
	}
	if _, err := uuid.Parse(s); err != nil {
		return "", ErrInvalidDocumentID
	}
	return DocumentID(s), nil
}

// MustDocumentIDFromString 从字符串创建文档ID（不验证，用于从数据库重建）
// 仅在确定数据来源可靠时使用（如从数据库读取）
func MustDocumentIDFromString(s string) DocumentID {
	return DocumentID(s)
}

// String 转换为字符串
func (id DocumentID) String() string {
	return string(id)
}

// IsEmpty 判断ID是否为空
func (id DocumentID) IsEmpty() bool {
	return string(id) == ""
}
