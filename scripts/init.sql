-- 知识库管理系统数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS knowledge_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE knowledge_db;

-- 知识库表
CREATE TABLE IF NOT EXISTS knowledge_bases (
    id VARCHAR(36) PRIMARY KEY COMMENT '知识库ID (UUID)',
    name VARCHAR(255) NOT NULL COMMENT '知识库名称',
    description TEXT COMMENT '知识库描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    -- 索引
    UNIQUE KEY uk_name (name),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识库表';

-- 文档表
CREATE TABLE IF NOT EXISTS documents (
    id VARCHAR(36) PRIMARY KEY COMMENT '文档ID (UUID)',
    knowledge_base_id VARCHAR(36) NOT NULL COMMENT '所属知识库ID',
    title VARCHAR(500) NOT NULL COMMENT '文档标题',
    content LONGTEXT NOT NULL COMMENT '文档内容',
    tags JSON COMMENT '标签列表 (JSON数组)',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    -- 索引
    KEY idx_knowledge_base_id (knowledge_base_id),
    KEY idx_created_at (created_at),
    
    -- 外键约束
    CONSTRAINT fk_documents_knowledge_base 
        FOREIGN KEY (knowledge_base_id) 
        REFERENCES knowledge_bases(id) 
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档表';

-- 插入示例数据（可选）
-- INSERT INTO knowledge_bases (id, name, description) VALUES
-- (UUID(), '技术文档', '技术相关的知识库'),
-- (UUID(), '产品手册', '产品使用说明');


