package valueobject

import "github.com/google/uuid"

// KnowledgeBaseID 知识库ID值对象
// 值对象是不可变的，没有唯一标识，通过值来比较
type KnowledgeBaseID string

// NewKnowledgeBaseID 创建新的知识库ID
func NewKnowledgeBaseID() KnowledgeBaseID {
	return KnowledgeBaseID(uuid.New().String())
}

// KnowledgeBaseIDFromString 从字符串创建知识库ID
func KnowledgeBaseIDFromString(s string) KnowledgeBaseID {
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

// DocumentIDFromString 从字符串创建文档ID
func DocumentIDFromString(s string) DocumentID {
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

