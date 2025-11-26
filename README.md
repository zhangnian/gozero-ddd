# Go-Zero DDD 知识库管理系统

> 基于 go-zero 框架和 DDD（领域驱动设计）的知识库管理 Demo 项目

## 📁 项目结构

```
gozero-ddd/
├── cmd/                          # 应用入口
│   ├── api/
│   │   └── main.go              # REST API 入口
│   └── rpc/
│       └── main.go              # gRPC 服务入口
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
│       ├── api/                 # HTTP REST API
│       │   ├── handler/         # 请求处理器
│       │   ├── middleware/      # 中间件
│       │   ├── svc/             # 服务上下文（依赖注入）
│       │   └── types/           # 请求/响应类型
│       └── rpc/                 # gRPC 服务
│           ├── pb/              # Protocol Buffer 生成代码
│           ├── server/          # gRPC 服务实现
│           ├── logic/           # 业务逻辑
│           └── svc/             # 服务上下文
├── rpc/                         # Proto 文件定义
│   └── knowledge.proto
├── api/                         # API 定义文件
│   └── knowledge.api
├── etc/                         # 配置文件
│   ├── knowledge.yaml           # REST API 配置
│   └── knowledge-rpc.yaml       # gRPC 服务配置
├── examples/                    # 示例代码
│   └── grpc_client/            # gRPC 客户端示例
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

**启动 REST API 服务（端口 8888）**
```bash
# 方式一：使用 make
make run

# 方式二：直接运行
go run cmd/api/main.go -f etc/knowledge.yaml
```

**启动 gRPC 服务（端口 9999）**
```bash
# 方式一：使用 make
make run-rpc

# 方式二：直接运行
go run cmd/rpc/main.go -f etc/knowledge-rpc.yaml
```

### 4. 访问 REST API
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

### 5. 访问 gRPC 接口

本项目提供了两个 gRPC 接口来演示 go-zero + DDD 中 gRPC 的正确使用方式：

**使用 grpcurl 测试（需要先安装）**
```bash
# 创建知识库（演示 Command 操作）
grpcurl -plaintext \
  -d '{"name":"gRPC测试知识库","description":"通过gRPC创建"}' \
  localhost:9999 knowledge.KnowledgeService/CreateKnowledgeBase

# 获取知识库详情（演示 Query 操作）
grpcurl -plaintext \
  -d '{"id":"<知识库ID>","include_documents":true}' \
  localhost:9999 knowledge.KnowledgeService/GetKnowledgeBase
```

**使用 Go 客户端示例**
```bash
# 先启动 gRPC 服务
make run-rpc

# 在另一个终端运行客户端示例
go run examples/grpc_client/main.go
```

## 🔄 gRPC + DDD 架构说明

### gRPC 请求处理流程

```
gRPC Request 
  → Server (实现 gRPC 接口) 
  → Logic (业务逻辑协调) 
  → Command/Query Handler (应用层) 
  → Domain Service (领域服务) 
  → Repository (仓储) 
  → Database
```

### gRPC 分层职责

| 层级 | 目录 | 职责 |
|------|------|------|
| Proto 定义 | `rpc/` | 定义 gRPC 接口和消息 |
| PB 代码 | `interfaces/rpc/pb/` | Protocol Buffer 生成代码 |
| Server 层 | `interfaces/rpc/server/` | 实现 gRPC 接口，创建 Logic |
| Logic 层 | `interfaces/rpc/logic/` | 协调业务逻辑，调用应用层 |
| 应用层 | `application/command/query/` | CQRS 模式的命令/查询处理 |
| 领域层 | `domain/` | 核心业务逻辑 |
| 基础设施层 | `infrastructure/` | 数据持久化实现 |

### 关键设计原则

1. **Logic 是请求级别的**：每个 gRPC 请求创建一个新的 Logic 实例
2. **复用应用层**：gRPC 和 REST API 共享相同的 Command/Query Handler
3. **DTO 转换在接口层**：Protobuf ↔ DTO 的转换发生在 Logic 层
4. **领域层不知道传输协议**：领域实体和服务与 gRPC/REST 无关

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

## 🔄 事务管理最佳实践

### 工作单元模式 (Unit of Work)

事务管理是 DDD 中的重要话题。本项目采用**工作单元模式**来抽象事务操作：

```
domain/repository/unit_of_work.go     <- 接口定义（领域层）
infrastructure/persistence/gorm_unit_of_work.go  <- GORM实现（基础设施层）
```

### 事务使用原则

1. **事务在应用层控制**：领域层不应该知道事务的存在
2. **通过上下文传递事务**：使用 `context.Context` 传递事务连接
3. **仓储自动感知事务**：仓储实现从上下文获取事务连接

### 代码示例

```go
// 应用层：在事务中执行多个操作
func (h *MergeKnowledgeBasesHandler) Handle(ctx context.Context, cmd *MergeKnowledgeBasesCommand) error {
    // 使用工作单元执行事务
    return h.unitOfWork.Transaction(ctx, func(txCtx context.Context) error {
        // 所有操作都在同一个事务中
        // 使用 txCtx 调用仓储方法
        
        docs, _ := h.docRepo.FindByKnowledgeBaseID(txCtx, sourceID)
        
        for _, doc := range docs {
            h.docRepo.Save(txCtx, newDoc)  // 使用事务连接
            h.docRepo.Delete(txCtx, doc.ID())
        }
        
        h.kbRepo.Delete(txCtx, sourceID)
        
        return nil  // 返回 nil 自动提交，返回 error 自动回滚
    })
}
```

### 合并知识库 API（事务演示）

```bash
# 将知识库A的所有文档移动到知识库B，然后删除知识库A
# 此操作在事务中执行，要么全部成功，要么全部失败
curl -X POST http://localhost:8888/api/v1/knowledge/merge \
  -H "Content-Type: application/json" \
  -d '{"source_id": "知识库A的ID", "target_id": "知识库B的ID"}'
```

## 🔑 核心设计原则

### 1. 依赖倒置原则
- 领域层不依赖任何外部层
- 外部层通过接口依赖领域层
- 仓储接口定义在领域层，实现在基础设施层

```
domain/repository/         <- 接口定义（包括 UnitOfWork）
infrastructure/persistence/ <- 具体实现（GORM/Memory）
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
