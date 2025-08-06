# 剪贴板监听器

[![Build and Test](https://github.com/YOUR_USERNAME/YOUR_REPO_NAME/workflows/Build%20and%20Test/badge.svg)](https://github.com/YOUR_USERNAME/YOUR_REPO_NAME/actions)
[![Release](https://github.com/YOUR_USERNAME/YOUR_REPO_NAME/workflows/Release/badge.svg)](https://github.com/YOUR_USERNAME/YOUR_REPO_NAME/actions)
[![Latest Release](https://img.shields.io/github/v/release/YOUR_USERNAME/YOUR_REPO_NAME)](https://github.com/YOUR_USERNAME/YOUR_REPO_NAME/releases/latest)

一个用 Go 语言开发的跨平台剪贴板监听器，具有图形用户界面。

## 功能特性

- 🔍 **实时监听**: 自动检测剪贴板内容变化
- 📋 **历史记录**: 保存最近 50 条剪贴板历史
- 🖱️ **快速复制**: 双击历史记录条目即可复制
- 🗑️ **清空功能**: 一键清空所有历史记录
- 🔄 **手动刷新**: 手动刷新历史记录列表
- 🌐 **跨平台**: 支持 Windows、macOS 和 Linux

## 技术栈

- **Go 1.21+**: 主要编程语言
- **Fyne v2**: GUI 框架
- **atotto/clipboard**: 剪贴板操作库

## 下载和安装

### 方式一：下载预编译版本（推荐）

前往 [Releases 页面](https://github.com/YOUR_USERNAME/YOUR_REPO_NAME/releases/latest) 下载适合你操作系统的版本：

- **Windows**: `clipboard-monitor-windows-amd64.exe`
- **macOS (Intel)**: `clipboard-monitor-macos-amd64`
- **macOS (Apple Silicon)**: `clipboard-monitor-macos-arm64`
- **Linux**: `clipboard-monitor-linux-amd64`

下载后直接运行即可，无需安装其他依赖。

### 方式二：从源码编译

#### 前置要求

- Go 1.21 或更高版本
- Git
- 系统依赖：
  - **Linux**: `libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libxxf86vm-dev`
  - **Windows**: 无额外要求
  - **macOS**: 无额外要求

#### 编译步骤

1. 克隆或下载项目到本地
2. 进入项目目录
3. 安装依赖：
   ```bash
   go mod tidy
   ```
4. 运行程序：
   ```bash
   go run main.go
   ```

#### 编译可执行文件

```bash
# 编译当前平台
go build -ldflags="-s -w" -o clipboard-monitor

# 交叉编译 Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o clipboard-monitor.exe

# 交叉编译 macOS (Intel)
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o clipboard-monitor-mac-amd64

# 交叉编译 macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-s -w" -o clipboard-monitor-mac-arm64

# 交叉编译 Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o clipboard-monitor-linux
```

> **注意**: 交叉编译需要相应平台的 CGO 工具链。推荐使用 GitHub Actions 自动构建。

## 使用说明

1. **启动程序**: 运行程序后会显示主界面
2. **监听剪贴板**: 程序会自动开始监听剪贴板变化
3. **查看历史**: 所有剪贴板内容变化都会显示在列表中
4. **复制内容**: 双击任意历史记录条目可将其复制到剪贴板
5. **清空历史**: 点击"清空历史"按钮可清除所有记录
6. **刷新列表**: 点击"刷新"按钮手动更新显示

## 界面说明

- **状态栏**: 显示当前程序状态和最后操作时间
- **历史列表**: 显示剪贴板历史记录，包含时间戳和内容预览
- **操作按钮**: 
  - 清空历史: 清除所有历史记录
  - 刷新: 手动刷新列表显示
  - 关于: 显示程序信息

## 项目结构

```
clipboard-monitor/
├── main.go              # 主程序入口和 GUI 界面
├── clipboard/
│   └── monitor.go       # 剪贴板监听核心逻辑
├── go.mod              # Go 模块文件
└── README.md           # 项目说明文档
```

## 自动构建和发布

本项目使用 GitHub Actions 进行自动构建和发布：

### 构建工作流

- **持续集成**: 每次推送到主分支或创建 PR 时自动运行测试和构建
- **多平台构建**: 自动在 Windows、macOS、Linux 上测试构建
- **代码质量检查**: 自动运行 `go fmt`、`go vet` 和测试

### 发布工作流

- **自动发布**: 推送版本标签（如 `v1.0.0`）时自动创建 GitHub Release
- **跨平台二进制**: 自动构建 Windows、macOS (Intel/ARM)、Linux 版本
- **优化构建**: 使用 `-ldflags="-s -w"` 减小二进制文件大小

### 创建新版本

1. 更新版本号并提交代码
2. 创建并推送标签：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. GitHub Actions 将自动创建 Release 并上传二进制文件

## 注意事项

- 程序需要访问系统剪贴板的权限
- 在某些系统上可能需要授予相应的权限
- 长文本内容在列表中会被截断显示，但完整内容会被保存
- 程序会过滤空内容，只记录有效的剪贴板变化

## 中文显示问题

如果遇到中文乱码问题，请参考：[中文显示问题解决方案.md](./中文显示问题解决方案.md)

**快速解决方案**：
1. Windows 10/11：启用"Beta 版：使用 Unicode UTF-8 提供全球语言支持"
2. 或使用提供的 `build-windows.bat` 重新构建
3. 或右键程序 → 属性 → 兼容性 → 设置兼容模式

## 许可证

本项目采用 MIT 许可证。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。
