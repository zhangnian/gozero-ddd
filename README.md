# Go-Zero DDD 知识库管理系统

> 基于 go-zero 框架和 DDD（领域驱动设计）的知识库管理 Demo 项目

## 📁 项目结构

```
gozero-ddd/
├── cmd/                          # 应用入口
│   └── api/
│       └── main.go
├── internal/                     # 内部代码（DDD分层架构）
│   ├── domain/                   # 🔷 领域层 - DDD核心
│   │   ├── entity/              # 实体（具有唯一标识的对象）
│   │   ├── valueobject/         # 值对象（无唯一标识，通过属性值比较）
│   │   ├── repository/          # 仓储接口（抽象持久化操作）
│   │   ├── service/             # 领域服务（跨实体的业务逻辑）
│   │   └── event/               # 领域事件（领域内发生的事件）
│   ├── application/             # 🔶 应用层 - 编排领域对象
│   │   ├── command/             # 命令处理（写操作）
│   │   ├── query/               # 查询处理（读操作）
│   │   └── dto/                 # 数据传输对象
│   ├── infrastructure/          # 🔵 基础设施层 - 技术实现
│   │   ├── persistence/         # 持久化实现（MySQL + 内存）
│   │   │   └── model/          # 数据库模型
│   │   └── config/              # 配置管理
│   └── interfaces/              # 🟢 接口层 - 对外暴露
│       └── api/                 # HTTP API
│           ├── handler/         # 请求处理器
│           ├── middleware/      # 中间件
│           ├── svc/             # 服务上下文（依赖注入）
│           └── types/           # 请求/响应类型
├── api/                         # API 定义文件
│   └── knowledge.api
├── etc/                         # 配置文件
│   └── knowledge.yaml
├── scripts/                     # 脚本文件
│   └── init.sql                # 数据库初始化脚本
├── go.mod
├── Makefile
└── README.md
```

## 🏗️ DDD 分层架构说明

### 1. 领域层 (Domain Layer)
**职责**：包含核心业务逻辑，是整个应用的心脏

- **Entity（实体）**：具有唯一标识和生命周期的对象
- **Value Object（值对象）**：没有唯一标识，通过属性值来比较
- **Repository（仓储接口）**：定义持久化抽象，不关心具体实现
- **Domain Service（领域服务）**：处理跨实体的复杂业务逻辑
- **Domain Event（领域事件）**：记录领域内发生的重要事件

### 2. 应用层 (Application Layer)
**职责**：编排领域对象，协调业务流程

- **Command（命令）**：处理写操作，改变系统状态
- **Query（查询）**：处理读操作，获取数据
- **DTO（数据传输对象）**：在层之间传递数据

### 3. 基础设施层 (Infrastructure Layer)
**职责**：提供技术实现和外部服务集成

- **Persistence（持久化）**：实现仓储接口
  - `GormKnowledgeBaseRepository` - GORM + MySQL 实现
  - `GormDocumentRepository` - GORM + MySQL 实现
  - `MemoryKnowledgeBaseRepository` - 内存实现（测试用）
  - `MemoryDocumentRepository` - 内存实现（测试用）
- **Model（数据模型）**：GORM 数据库表映射模型
- **Config（配置）**：管理应用配置

### 4. 接口层 (Interfaces Layer)
**职责**：处理外部请求，适配不同的接入方式

- **API Handler（处理器）**：处理 HTTP 请求
- **Middleware（中间件）**：处理横切关注点
- **ServiceContext（服务上下文）**：依赖注入容器

## 🚀 快速开始

### 1. 安装依赖
```bash
go mod tidy
```

### 2. 初始化数据库

**方式一：使用 MySQL**

```bash
# 执行数据库初始化脚本
mysql -u root -p < scripts/init.sql
```

然后修改 `etc/knowledge.yaml` 中的数据库配置：

```yaml
UseMemory: false
MySQL:
  DataSource: "root:your_password@tcp(127.0.0.1:3306)/knowledge_db?charset=utf8mb4&parseTime=True&loc=Local"
```

**方式二：使用内存存储（开发测试）**

修改 `etc/knowledge.yaml`：

```yaml
UseMemory: true
```

### 3. 运行项目
```bash
# 方式一：使用 make
make run

# 方式二：直接运行
go run cmd/api/main.go -f etc/knowledge.yaml
```

### 4. 访问 API
```bash
# 创建知识库
curl -X POST http://localhost:8888/api/v1/knowledge \
  -H "Content-Type: application/json" \
  -d '{"name": "技术文档", "description": "技术相关的知识库"}'

# 获取知识库列表
curl http://localhost:8888/api/v1/knowledge

# 获取单个知识库
curl http://localhost:8888/api/v1/knowledge/{id}

# 更新知识库
curl -X PUT http://localhost:8888/api/v1/knowledge/{id} \
  -H "Content-Type: application/json" \
  -d '{"name": "新名称", "description": "新描述"}'

# 删除知识库
curl -X DELETE http://localhost:8888/api/v1/knowledge/{id}

# 添加文档
curl -X POST http://localhost:8888/api/v1/knowledge/{id}/documents \
  -H "Content-Type: application/json" \
  -d '{"title": "Go语言入门", "content": "Go是一门简洁的语言...", "tags": ["go", "programming"]}'

# 获取文档列表
curl http://localhost:8888/api/v1/knowledge/{id}/documents

# 删除文档
curl -X DELETE http://localhost:8888/api/v1/knowledge/{id}/documents/{doc_id}
```

## 🗄️ 数据库设计

### knowledge_bases 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 知识库ID (UUID) |
| name | VARCHAR(255) | 知识库名称（唯一） |
| description | TEXT | 描述 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### documents 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 文档ID (UUID) |
| knowledge_base_id | VARCHAR(36) | 所属知识库ID（外键） |
| title | VARCHAR(500) | 文档标题 |
| content | LONGTEXT | 文档内容 |
| tags | JSON | 标签列表 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

## 🔑 核心设计原则

### 1. 依赖倒置原则
- 领域层不依赖任何外部层
- 外部层通过接口依赖领域层
- 仓储接口定义在领域层，实现在基础设施层

```
domain/repository/         <- 接口定义
infrastructure/persistence/ <- 具体实现（MySQL/Memory）
```

### 2. 聚合根设计
- `KnowledgeBase` 是聚合根，管理 `Document` 实体
- 所有对 `Document` 的操作都通过 `KnowledgeBase` 进行

### 3. 命令查询职责分离 (CQRS)
- **Command**：处理创建、更新、删除操作
- **Query**：处理查询操作

### 4. 数据模型与领域模型分离
- `model/` 目录下是数据库模型，负责 ORM 映射
- `entity/` 目录下是领域实体，包含业务逻辑
- 两者通过 `ToEntity()` 和 `FromEntity()` 方法转换

## 📝 go-zero 框架最佳实践

1. **API 定义**：使用 `.api` 文件定义接口规范
2. **配置管理**：使用 YAML 配置文件，支持多环境
3. **中间件**：实现认证、日志等横切关注点
4. **错误处理**：统一的错误响应格式
5. **依赖注入**：通过 ServiceContext 管理依赖
6. **数据库访问**：使用 GORM ORM 框架
7. **自动迁移**：支持 GORM AutoMigrate 自动建表

## 📚 参考资料

- [go-zero 官方文档](https://go-zero.dev/)
- [领域驱动设计](https://domainlanguage.com/ddd/)
- [CQRS 模式](https://martinfowler.com/bliki/CQRS.html)
