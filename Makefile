.PHONY: build run test lint clean docker docker-up docker-down

# 变量定义
APP_NAME=shortlink-service
BUILD_DIR=bin
MAIN_PATH=cmd/server/main.go

# 构建
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(APP_NAME)"

# 运行
run:
	@echo "Starting $(APP_NAME)..."
	go run $(MAIN_PATH)

# 测试
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Test completed. Coverage report: coverage.html"

# 测试（简洁模式）
test-quick:
	go test -race ./...

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# 清理
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean completed"

# Docker 构建
docker:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

# Docker 启动
docker-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d

# Docker 停止
docker-down:
	@echo "Stopping services..."
	docker-compose down

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Dependencies installed"

# 帮助信息
help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests with coverage"
	@echo "  make test-quick  - Run tests (quick mode)"
	@echo "  make lint        - Run linter"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make docker      - Build Docker image"
	@echo "  make docker-up   - Start services with docker-compose"
	@echo "  make docker-down - Stop docker-compose services"
	@echo "  make deps        - Install dependencies"
	@echo "  make help        - Show this help message"
