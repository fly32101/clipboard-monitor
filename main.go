package main

import (
	"clipboard-monitor/clipboard"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
	"time"
)

type ClipboardApp struct {
	app         fyne.App
	window      fyne.Window
	monitor     *clipboard.Monitor
	historyList *widget.List
	statusLabel *widget.Label
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewClipboardApp() *ClipboardApp {
	myApp := app.New()
	myApp.SetIcon(nil)

	window := myApp.NewWindow("剪贴板监听器")
	window.Resize(fyne.NewSize(600, 400))

	ctx, cancel := context.WithCancel(context.Background())

	return &ClipboardApp{
		app:     myApp,
		window:  window,
		monitor: clipboard.NewMonitor(50), // 保存最近50条记录
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (ca *ClipboardApp) setupUI() {
	// 状态标签
	ca.statusLabel = widget.NewLabel("状态: 准备就绪")

	// 历史记录列表
	ca.historyList = widget.NewList(
		func() int {
			return len(ca.monitor.GetHistory())
		},
		func() fyne.CanvasObject {
			// 创建一个简单的标签作为列表项
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapWord
			return label
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			history := ca.monitor.GetHistory()
			if id >= len(history) {
				return
			}

			entry := history[id]
			label := obj.(*widget.Label)

			// 格式化显示内容
			content := entry.Content
			if len(content) > 80 {
				content = content[:80] + "..."
			}
			// 替换换行符为空格以便显示
			content = strings.ReplaceAll(content, "\n", " ")
			content = strings.ReplaceAll(content, "\r", " ")

			// 组合时间和内容
			displayText := fmt.Sprintf("[%s] %s",
				entry.Timestamp.Format("15:04:05"), content)
			label.SetText(displayText)
		},
	)

	// 双击复制功能
	ca.historyList.OnSelected = func(id widget.ListItemID) {
		history := ca.monitor.GetHistory()
		if id < len(history) {
			entry := history[id]
			err := ca.monitor.CopyToClipboard(entry.Content)
			if err != nil {
				dialog.ShowError(err, ca.window)
			} else {
				ca.statusLabel.SetText(fmt.Sprintf("已复制: %s",
					time.Now().Format("15:04:05")))
			}
		}
	}

	// 按钮
	clearBtn := widget.NewButton("清空历史", func() {
		dialog.ShowConfirm("确认", "确定要清空所有历史记录吗？",
			func(confirmed bool) {
				if confirmed {
					ca.monitor.ClearHistory()
					ca.historyList.Refresh()
					ca.statusLabel.SetText("历史记录已清空")
				}
			}, ca.window)
	})

	refreshBtn := widget.NewButton("刷新", func() {
		ca.historyList.Refresh()
		ca.statusLabel.SetText("列表已刷新")
	})

	aboutBtn := widget.NewButton("关于", func() {
		dialog.ShowInformation("关于",
			GetVersionInfo()+"\n\n功能:\n- 实时监听剪贴板变化\n- 显示历史记录\n- 双击条目可复制\n- 跨平台支持\n\n作者: AI Assistant\n构建版本: "+Version,
			ca.window)
	})

	// 布局
	buttonContainer := container.NewHBox(
		clearBtn,
		refreshBtn,
		aboutBtn,
	)

	content := container.NewVBox(
		ca.statusLabel,
		widget.NewSeparator(),
		widget.NewLabel("剪贴板历史记录 (双击复制):"),
		ca.historyList,
		widget.NewSeparator(),
		buttonContainer,
	)

	ca.window.SetContent(content)
}

func (ca *ClipboardApp) startMonitoring() {
	// 设置新内容回调
	ca.monitor.SetOnNewContent(func(entry clipboard.ClipboardEntry) {
		// 在UI线程中更新界面
		ca.statusLabel.SetText(fmt.Sprintf("检测到新内容: %s",
			entry.Timestamp.Format("15:04:05")))
		ca.historyList.Refresh()
	})

	// 在后台启动监听
	go func() {
		err := ca.monitor.Start(ca.ctx)
		if err != nil && err != context.Canceled {
			log.Printf("监听器错误: %v", err)
		}
	}()
}

func (ca *ClipboardApp) Run() {
	ca.setupUI()
	ca.startMonitoring()

	// 设置窗口关闭回调
	ca.window.SetCloseIntercept(func() {
		ca.cancel() // 停止监听
		ca.app.Quit()
	})

	ca.statusLabel.SetText("状态: 正在监听剪贴板...")
	ca.window.ShowAndRun()
}

func main() {
	app := NewClipboardApp()
	app.Run()
}
