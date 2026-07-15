# XINFRA Makefile

.PHONY: help build run test docker clean

# 默认目标
help:
	@echo "XINFRA - 统一运维入口平台"
	@echo ""
	@echo "可用命令:"
	@echo "  make help          - 显示此帮助信息"
	@echo "  make build         - 构建前后端项目"
	@echo "  make run           - 启动开发服务器"
	@echo "  make test          - 运行测试"
	@echo "  make docker        - 构建 Docker 镜像"
	@echo "  make clean         - 清理构建产物"
	@echo "  make swagger       - 生成 Swagger API 文档"
	@echo ""

# 构建前后端项目
build: build-frontend build-server

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

build-server:
	@echo "Building server..."
	cd server && go build -o bin/xinfra ./cmd/server

# 启动开发服务器（前后端并行）
run:
	@echo "Starting development servers..."
	@trap 'kill 0' EXIT; \
	cd frontend && npm run dev & \
	cd server && go run ./cmd/server & \
	wait

# 运行测试
test: test-frontend test-server

test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm run test

test-server:
	@echo "Running server tests..."
	cd server && go test ./...

# 构建 Docker 镜像
docker: docker-frontend docker-server

docker-frontend:
	@echo "Building frontend Docker image..."
	docker build -t xinfra-frontend:latest -f deploy/docker/Dockerfile.frontend .

docker-server:
	@echo "Building server Docker image..."
	docker build -t xinfra-server:latest -f deploy/docker/Dockerfile.server .

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	rm -rf frontend/dist frontend/node_modules
	rm -rf server/bin server/vendor
	find . -type d -name vendor -exec rm -rf {} + 2>/dev/null || true

# 本地开发环境
dev-up:
	@echo "Starting local development environment..."
	docker-compose up -d

dev-down:
	@echo "Stopping local development environment..."
	docker-compose down

# Swagger 文档
swagger:
	@echo "Generating Swagger documentation..."
	cd server && swag init -g internal/router/router.go -o ./docs
	@echo "Swagger docs generated at server/docs/"

# 数据库迁移
migrate-up:
	@echo "Running database migrations..."
	cd server && go run ./cmd/server migrate up

migrate-down:
	@echo "Rolling back database migrations..."
	cd server && go run ./cmd/server migrate down
