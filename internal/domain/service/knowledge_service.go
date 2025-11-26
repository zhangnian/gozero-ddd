package service

import (
	"context"

	"gozero-ddd/internal/domain"
	"gozero-ddd/internal/domain/entity"
	"gozero-ddd/internal/domain/repository"
)

// KnowledgeService 知识库领域服务
// 领域服务处理跨实体的业务逻辑，或不适合放在实体中的业务逻辑
type KnowledgeService struct {
	kbRepo  repository.KnowledgeBaseRepository
	docRepo repository.DocumentRepository
}

// NewKnowledgeService 创建知识库领域服务
func NewKnowledgeService(
	kbRepo repository.KnowledgeBaseRepository,
	docRepo repository.DocumentRepository,
) *KnowledgeService {
	return &KnowledgeService{
		kbRepo:  kbRepo,
		docRepo: docRepo,
	}
}

// CreateKnowledgeBase 创建知识库
// 包含业务规则验证：名称不能重复
func (s *KnowledgeService) CreateKnowledgeBase(ctx context.Context, name, description string) (*entity.KnowledgeBase, error) {
	// 检查名称是否已存在
	exists, err := s.kbRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrKnowledgeBaseNameExists
	}

	// 创建知识库实体
	kb, err := entity.NewKnowledgeBase(name, description)
	if err != nil {
		return nil, err
	}

	// 持久化
	if err := s.kbRepo.Save(ctx, kb); err != nil {
		return nil, err
	}

	return kb, nil
}

// DeleteKnowledgeBase 删除知识库及其所有文档
// 这是一个跨聚合的操作，适合放在领域服务中
func (s *KnowledgeService) DeleteKnowledgeBase(ctx context.Context, kb *entity.KnowledgeBase) error {
	// 先删除所有文档
	if err := s.docRepo.DeleteByKnowledgeBaseID(ctx, kb.ID()); err != nil {
		return err
	}

	// 再删除知识库
	return s.kbRepo.Delete(ctx, kb.ID())
}
