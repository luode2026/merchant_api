.PHONY: help admin app gen tidy deps clean

help:
	@echo "可用命令:"
	@echo "  make admin      - 启动 Admin 端服务"
	@echo "  make app        - 启动 App 端服务"
	@echo "  make gen        - 根据数据库表生成模型文件"
	@echo "  make tidy       - 整理依赖"
	@echo "  make deps       - 下载依赖"
	@echo "  make clean      - 清理构建文件"

# 启动 Admin 端
admin:
	@echo "启动 Admin 端服务..."
	go run cmd/admin/main.go

# 启动 App 端
app:
	@echo "启动 App 端服务..."
	go run cmd/app/main.go

# 生成模型文件 (使用示例: make gen table=users)
gen:
	@echo "根据数据库表生成模型文件..."
	go run cmd/generator/main.go -table="$(table)"

# 整理依赖
tidy:
	go mod tidy

# 下载依赖
deps:
	go mod download

# 清理
clean:
	rm -rf bin/
	go clean
