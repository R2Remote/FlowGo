.PHONY: run build test clean migrate

# 运行项目
run:
	go run cmd/server/main.go

# 构建项目
build:
	go build -o bin/server cmd/server/main.go

# 运行测试
test:
	go test ./...

# 清理构建文件
clean:
	rm -rf bin/

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 安装依赖
deps:
	go mod download
	go mod tidy

# 数据库迁移（需要根据实际情况实现）
migrate:
	@echo "Database migration not implemented yet"
