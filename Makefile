.PHONY: build run test clean deps

# 构建项目
build:
	go build -o bin/edu-course-client main.go

# 运行项目
run:
	go run main.go

# 运行测试
test:
	go test ./...

# 清理构建文件
clean:
	rm -rf bin/

# 下载依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run
