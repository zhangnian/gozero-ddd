.PHONY: build build-rpc run run-rpc clean test tidy proto

# 应用名称
APP_NAME=knowledge-api
RPC_NAME=knowledge-rpc

# 构建目录
BUILD_DIR=build

# Go 编译参数
GO=go
GOFLAGS=-ldflags="-s -w"

# 默认目标
all: build build-rpc

# 下载依赖
tidy:
	$(GO) mod tidy

# ==================== REST API ====================

# 构建 REST API 应用
build: tidy
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

# 运行 REST API 应用
run: tidy
	$(GO) run ./cmd/api -f etc/knowledge.yaml

# ==================== gRPC 服务 ====================

# 构建 gRPC 服务
build-rpc: tidy
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(RPC_NAME) ./cmd/rpc

# 运行 gRPC 服务
run-rpc: tidy
	$(GO) run ./cmd/rpc -f etc/knowledge-rpc.yaml

# 生成 Proto 文件（需要安装 protoc 和 protoc-gen-go）
# 安装命令：
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
proto:
	protoc --go_out=. --go-grpc_out=. rpc/knowledge.proto

# ==================== 通用命令 ====================

# 运行测试
test:
	$(GO) test -v ./...

# 清理构建产物
clean:
	@rm -rf $(BUILD_DIR)

# 代码格式化
fmt:
	$(GO) fmt ./...

# 代码检查
lint:
	golangci-lint run ./...

# 生成 API 代码（可选，本项目手动编写）
# api:
#	goctl api go -api api/knowledge.api -dir .

# 帮助信息
help:
	@echo "可用命令:"
	@echo ""
	@echo "REST API:"
	@echo "  make build     - 构建 REST API 应用"
	@echo "  make run       - 运行 REST API 应用 (端口 8888)"
	@echo ""
	@echo "gRPC 服务:"
	@echo "  make build-rpc - 构建 gRPC 服务"
	@echo "  make run-rpc   - 运行 gRPC 服务 (端口 9999)"
	@echo "  make proto     - 生成 Proto 代码"
	@echo ""
	@echo "通用命令:"
	@echo "  make all       - 构建所有服务"
	@echo "  make test      - 运行测试"
	@echo "  make tidy      - 下载依赖"
	@echo "  make clean     - 清理构建产物"
	@echo "  make fmt       - 格式化代码"
	@echo "  make lint      - 代码检查"

