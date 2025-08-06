package main

import (
	"clipboard-monitor/clipboard"
	"clipboard-monitor/hotkey"
	"clipboard-monitor/keyboard"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	webview "github.com/webview/webview_go"
)

type ClipboardApp struct {
	w            webview.WebView
	monitor      *clipboard.Monitor
	ctx          context.Context
	cancel       context.CancelFunc
	hidden       bool // 窗口是否隐藏
	hotkeyMgr    *hotkey.HotkeyManager
	globalHotkey bool // 全局热键是否启用
}

func NewClipboardApp() *ClipboardApp {
	ctx, cancel := context.WithCancel(context.Background())

	return &ClipboardApp{
		monitor:   clipboard.NewMonitor(50), // Save last 50 records
		ctx:       ctx,
		cancel:    cancel,
		hotkeyMgr: hotkey.NewHotkeyManager(),
	}
}

func (ca *ClipboardApp) setupUI() error {
	debug := false

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

	ca.w.SetHtml(string(htmlContent))

	return nil
}

func (ca *ClipboardApp) bindFunctions() {
	// 绑定获取历史记录函数
	ca.w.Bind("getHistory", func() interface{} {
		history := ca.monitor.GetHistory()

		// 确保返回的是数组，即使为空
		if history == nil {
			return []interface{}{}
		}

		return history
	})

	// 绑定复制到剪贴板函数
	ca.w.Bind("copyToClipboardGo", func(content string) interface{} {
		err := ca.monitor.CopyToClipboard(content)
		if err != nil {
			return map[string]string{"error": err.Error()}
		}
		return map[string]bool{"success": true}
	})

	// 绑定清空历史函数
	ca.w.Bind("clearHistory", func() interface{} {
		ca.monitor.ClearHistory()
		return map[string]bool{"success": true}
	})

	// 绑定获取版本信息函数
	ca.w.Bind("getVersionInfo", func() interface{} {
		return GetVersionInfo()
	})

	// 绑定删除单个历史记录项函数
	ca.w.Bind("deleteHistoryItemGo", func(index int) interface{} {
		err := ca.monitor.DeleteHistoryItem(index)
		if err != nil {
			return map[string]string{"error": err.Error()}
		}
		return map[string]bool{"success": true}
	})

	// 绑定快捷键设置函数
	ca.w.Bind("saveHotkeySettings", func(config map[string]interface{}) interface{} {
		log.Printf("收到快捷键配置保存请求: %+v", config)
		// 这里可以保存快捷键配置到文件或注册表
		// 暂时返回成功，实际实现可以根据需要添加
		result := map[string]bool{"success": true}
		log.Printf("返回保存结果: %+v", result)
		return result
	})

	ca.w.Bind("getHotkeySettings", func() interface{} {
		log.Printf("收到获取快捷键配置请求")
		// 这里可以从文件或注册表读取快捷键配置
		// 暂时返回默认配置
		result := map[string]interface{}{
			"hotkey":  "Ctrl+Shift+V",
			"enabled": false,
		}
		log.Printf("返回快捷键配置: %+v", result)
		return result
	})

	ca.w.Bind("setGlobalHotkeyEnabled", func(enabled bool) interface{} {
		log.Printf("收到设置全局快捷键状态请求: %v", enabled)

		if enabled && !ca.globalHotkey {
			// 启用全局热键
			err := ca.hotkeyMgr.RegisterHotkey(func() {
				log.Printf("全局热键被触发，显示快速选择界面")
				// 显示快速选择界面
				ca.w.Eval(`
					if (typeof showQuickSelector === 'function') {
						showQuickSelector();
					}
				`)
			})
			if err != nil {
				log.Printf("注册全局热键失败: %v", err)
				return map[string]string{"error": "注册全局热键失败: " + err.Error()}
			}
			ca.globalHotkey = true
		} else if !enabled && ca.globalHotkey {
			// 禁用全局热键
			err := ca.hotkeyMgr.UnregisterHotkey()
			if err != nil {
				log.Printf("注销全局热键失败: %v", err)
				return map[string]string{"error": "注销全局热键失败: " + err.Error()}
			}
			ca.globalHotkey = false
		}

		result := map[string]bool{"success": true}
		log.Printf("返回设置结果: %+v", result)
		return result
	})

	// 绑定窗口控制函数
	ca.w.Bind("hideWindowGo", func() interface{} {
		log.Printf("隐藏窗口")
		ca.hidden = true
		// 注意：WebView可能不支持直接隐藏，这里只是标记状态
		return map[string]bool{"success": true}
	})

	ca.w.Bind("showWindowGo", func() interface{} {
		log.Printf("显示窗口")
		ca.hidden = false
		return map[string]bool{"success": true}
	})

	ca.w.Bind("minimizeWindow", func() interface{} {
		log.Printf("最小化窗口")
		// WebView的最小化功能有限，这里只是记录
		return map[string]bool{"success": true}
	})

	ca.w.Bind("toggleWindow", func() interface{} {
		log.Printf("切换窗口显示状态，当前状态: %v", ca.hidden)
		ca.hidden = !ca.hidden
		if ca.hidden {
			log.Printf("窗口已隐藏")
		} else {
			log.Printf("窗口已显示")
		}
		return map[string]interface{}{
			"success": true,
			"hidden":  ca.hidden,
		}
	})

	ca.w.Bind("setMinimizeToTrayEnabled", func(enabled bool) interface{} {
		log.Printf("设置最小化到托盘: %v", enabled)
		// 这里可以设置托盘相关的配置
		return map[string]bool{"success": true}
	})

	// 绑定直接粘贴功能
	ca.w.Bind("pasteContentGo", func(content string) interface{} {
		contentPreview := content
		if len(content) > 50 {
			contentPreview = content[:50] + "..."
		}
		log.Printf("执行直接粘贴: %s", contentPreview)

		// 先复制到剪贴板
		err := ca.monitor.CopyToClipboard(content)
		if err != nil {
			return map[string]string{"error": "复制到剪贴板失败: " + err.Error()}
		}

		// 等待一小段时间确保剪贴板更新
		time.Sleep(50 * time.Millisecond)

		// 模拟 Ctrl+V 按键
		go func() {
			// 稍微延迟以确保窗口切换完成
			time.Sleep(100 * time.Millisecond)

			// 发送 Ctrl+V 组合键
			err := keyboard.SendCtrlV()
			if err != nil {
				log.Printf("发送Ctrl+V失败: %v", err)
			} else {
				log.Printf("已发送Ctrl+V组合键")
			}
		}()

		return map[string]bool{"success": true}
	})

	// 绑定快速粘贴功能（全局热键触发）
	ca.w.Bind("quickPaste", func(index int) interface{} {
		log.Printf("快速粘贴第 %d 项", index)

		history := ca.monitor.GetHistory()
		if index < 0 || index >= len(history) {
			return map[string]string{"error": "索引超出范围"}
		}

		content := history[index].Content

		// 复制到剪贴板
		err := ca.monitor.CopyToClipboard(content)
		if err != nil {
			return map[string]string{"error": "复制失败: " + err.Error()}
		}

		// 发送Ctrl+V
		go func() {
			time.Sleep(100 * time.Millisecond)
			err := keyboard.SendCtrlV()
			if err != nil {
				log.Printf("发送Ctrl+V失败: %v", err)
			} else {
				log.Printf("已快速粘贴内容")
			}
		}()

		return map[string]bool{"success": true}
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ca *ClipboardApp) startMonitoring() {
	// Set new content callback
	ca.monitor.SetOnNewContent(func(entry clipboard.ClipboardEntry) {
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
		err := ca.monitor.Start(ca.ctx)
		if err != nil && err != context.Canceled {
			log.Printf("Monitor error: %v", err)
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
	if ca.globalHotkey {
		ca.hotkeyMgr.UnregisterHotkey()
	}
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
