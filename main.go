package main

import (
	"clipboard-monitor/clipboard"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	webview "github.com/webview/webview_go"
)

type ClipboardApp struct {
	w       webview.WebView
	monitor *clipboard.Monitor
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewClipboardApp() *ClipboardApp {
	ctx, cancel := context.WithCancel(context.Background())

	return &ClipboardApp{
		monitor: clipboard.NewMonitor(50), // Save last 50 records
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (ca *ClipboardApp) setupUI() error {
	debug := false // 启用调试模式
	log.Printf("WebView 调试模式: %v", debug)

	// 创建 WebView
	w := webview.New(debug)
	ca.w = w

	// 设置窗口属性
	w.SetTitle("剪贴板监控器")
	w.SetSize(800, 600, webview.HintNone)

	// 绑定 Go 函数到 JavaScript
	ca.bindFunctions()

	// 加载外部 HTML 文件
	err := ca.loadHTMLFile()
	if err != nil {
		return fmt.Errorf("failed to load HTML file: %v", err)
	}

	return nil
}

func (ca *ClipboardApp) loadHTMLFile() error {
	// 获取当前执行文件的目录
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	execDir := filepath.Dir(execPath)

	// 构建 web 目录路径
	webDir := filepath.Join(execDir, "web")
	indexPath := filepath.Join(webDir, "index.html")

	// 检查文件是否存在
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// 如果在执行文件目录找不到，尝试当前工作目录
		wd, _ := os.Getwd()
		webDir = filepath.Join(wd, "web")
		indexPath = filepath.Join(webDir, "index.html")

		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			return fmt.Errorf("web files not found in %s or %s", filepath.Join(execDir, "web"), webDir)
		}
	}

	// 读取 HTML 文件内容
	htmlContent, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read HTML file: %v", err)
	}

	log.Printf("Successfully loaded HTML file from: %s", indexPath)
	ca.w.SetHtml(string(htmlContent))

	return nil
}

func (ca *ClipboardApp) bindFunctions() {
	// 绑定获取历史记录函数
	ca.w.Bind("getHistory", func() interface{} {
		history := ca.monitor.GetHistory()
		log.Printf("获取历史记录: %d 条", len(history))
		for i, entry := range history {
			log.Printf("  [%d] %s: %s", i, entry.Timestamp.Format("15:04:05"), entry.Content[:min(30, len(entry.Content))])
		}

		// 确保返回的是数组，即使为空
		if history == nil {
			log.Println("历史记录为 nil，返回空数组")
			return []interface{}{}
		}

		log.Printf("返回历史记录数组，长度: %d", len(history))
		return history
	})

	// 绑定复制到剪贴板函数
	ca.w.Bind("copyToClipboardGo", func(content string) interface{} {
		err := ca.monitor.CopyToClipboard(content)
		if err != nil {
			log.Printf("复制失败: %v", err)
			return map[string]string{"error": err.Error()}
		}
		log.Printf("复制成功: %s", content[:min(50, len(content))])
		return map[string]bool{"success": true}
	})

	// 绑定清空历史函数
	ca.w.Bind("clearHistory", func() interface{} {
		ca.monitor.ClearHistory()
		log.Println("历史记录已清空")
		return map[string]bool{"success": true}
	})

	// 绑定获取版本信息函数
	ca.w.Bind("getVersionInfo", func() interface{} {
		return GetVersionInfo()
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ca *ClipboardApp) startMonitoring() {
	log.Println("开始启动剪贴板监控...")

	// Set new content callback
	ca.monitor.SetOnNewContent(func(entry clipboard.ClipboardEntry) {
		log.Printf("检测到新剪贴板内容: %s", entry.Content[:min(50, len(entry.Content))])
		// 通知前端更新
		if ca.w != nil {
			// 可以通过 JavaScript 更新状态
			ca.w.Eval(fmt.Sprintf(`
				if (typeof updateStatus === 'function') {
					updateStatus('新内容检测到: %s');
				}
			`, entry.Timestamp.Format("15:04:05")))
		}
	})

	// Start monitoring in background
	go func() {
		log.Println("剪贴板监控协程已启动")
		err := ca.monitor.Start(ca.ctx)
		if err != nil && err != context.Canceled {
			log.Printf("Monitor error: %v", err)
		} else {
			log.Println("剪贴板监控正常结束")
		}
	}()
}

func (ca *ClipboardApp) Run() error {
	// 设置 UI
	err := ca.setupUI()
	if err != nil {
		return fmt.Errorf("failed to setup UI: %v", err)
	}

	// 开始监控
	ca.startMonitoring()

	// 运行 WebView
	ca.w.Run()

	// 清理资源
	ca.cancel()
	ca.w.Destroy()

	return nil
}

// initConsole initializes console encoding
func initConsole() {
	if runtime.GOOS == "windows" {
		// Set environment variables for GBK
		os.Setenv("LANG", "zh_CN.GBK")
		os.Setenv("LC_ALL", "zh_CN.GBK")
		os.Setenv("LC_CTYPE", "zh_CN.GBK")
	}

	// Call platform-specific initialization
	initPlatformSpecific()
}

func main() {
	// Initialize console encoding
	initConsole()

	app := NewClipboardApp()
	err := app.Run()
	if err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
