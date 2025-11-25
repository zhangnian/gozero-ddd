.PHONY: build run clean test tidy

# 应用名称
APP_NAME=knowledge-api

# 构建目录
BUILD_DIR=build

# Go 编译参数
GO=go
GOFLAGS=-ldflags="-s -w"

# 默认目标
all: build

# 下载依赖
tidy:
	$(GO) mod tidy

# 构建应用
build: tidy
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

# 运行应用
run: tidy
	$(GO) run ./cmd/api -f etc/knowledge.yaml

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
	@echo "  make build  - 构建应用"
	@echo "  make run    - 运行应用"
	@echo "  make test   - 运行测试"
	@echo "  make tidy   - 下载依赖"
	@echo "  make clean  - 清理构建产物"
	@echo "  make fmt    - 格式化代码"
	@echo "  make lint   - 代码检查"

