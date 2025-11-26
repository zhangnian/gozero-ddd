package dto

// MergeResultDTO 合并知识库结果DTO
type MergeResultDTO struct {
	SourceID       string `json:"source_id"`        // 源知识库ID
	SourceName     string `json:"source_name"`      // 源知识库名称
	TargetID       string `json:"target_id"`        // 目标知识库ID
	TargetName     string `json:"target_name"`      // 目标知识库名称
	DocumentsMoved int    `json:"documents_moved"`  // 移动的文档数量
	SourceDeleted  bool   `json:"source_deleted"`   // 源知识库是否已删除
}

