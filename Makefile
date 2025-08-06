# 剪贴板监听器 Makefile

# 变量定义
APP_NAME = clipboard-monitor
VERSION = 1.0.0
BUILD_TIME = $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS = -s -w -X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# 默认目标
.PHONY: all
all: build

# 构建当前平台
.PHONY: build
build:
	CGO_ENABLED=1 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME) .

# 运行程序
.PHONY: run
run:
	go run .

# 清理构建文件
.PHONY: clean
clean:
	rm -f $(APP_NAME) $(APP_NAME).exe
	rm -f $(APP_NAME)-*

# 运行测试
.PHONY: test
test:
	go test -v ./...

# 代码格式化
.PHONY: fmt
fmt:
	go fmt ./...

# 代码检查
.PHONY: vet
vet:
	go vet ./...

# 下载依赖
.PHONY: deps
deps:
	go mod download
	go mod tidy

# 修复依赖问题
.PHONY: fix-deps
fix-deps:
	@echo "修复 Go 模块依赖问题..."
	go env -w GOPROXY=https://goproxy.cn,direct
	go env -w GOSUMDB=sum.golang.google.cn
	go clean -modcache
	rm -f go.sum
	go mod tidy
	@echo "依赖修复完成！"

# 跨平台构建
.PHONY: build-all
build-all: build-windows build-macos build-linux

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="$(LDFLAGS) -H windowsgui" -o $(APP_NAME)-windows-amd64.exe .

.PHONY: build-macos
build-macos:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-macos-amd64 .
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-macos-arm64 .

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-linux-amd64 .

# 创建发布包
.PHONY: package
package: build-all
	mkdir -p dist
	cp $(APP_NAME)-windows-amd64.exe dist/
	cp $(APP_NAME)-macos-amd64 dist/
	cp $(APP_NAME)-macos-arm64 dist/
	cp $(APP_NAME)-linux-amd64 dist/
	cp README.md dist/

# 帮助信息
.PHONY: help
help:
	@echo "可用的 make 目标:"
	@echo "  build        - 构建当前平台的可执行文件"
	@echo "  run          - 运行程序"
	@echo "  test         - 运行测试"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 代码检查"
	@echo "  deps         - 下载和整理依赖"
	@echo "  fix-deps     - 修复依赖问题"
	@echo "  clean        - 清理构建文件"
	@echo "  build-all    - 构建所有平台"
	@echo "  build-windows- 构建 Windows 版本"
	@echo "  build-macos  - 构建 macOS 版本"
	@echo "  build-linux  - 构建 Linux 版本"
	@echo "  package      - 创建发布包"
	@echo "  help         - 显示此帮助信息"
