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
	"github.com/flopp/go-findfont"
	"log"
	"os"
	"runtime"
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

func init() {
	//设置中文字体:解决中文乱码问题
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func NewClipboardApp() *ClipboardApp {
	myApp := app.New()
	myApp.SetIcon(nil)

	window := myApp.NewWindow("Clipboard Monitor")
	window.Resize(fyne.NewSize(600, 400))

	ctx, cancel := context.WithCancel(context.Background())

	return &ClipboardApp{
		app:     myApp,
		window:  window,
		monitor: clipboard.NewMonitor(50), // Save last 50 records
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (ca *ClipboardApp) setupUI() {
	// Status label
	ca.statusLabel = widget.NewLabel("Status: Ready")

	// History list
	ca.historyList = widget.NewList(
		func() int {
			return len(ca.monitor.GetHistory())
		},
		func() fyne.CanvasObject {
			// Create a simple label as list item
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

			// Format display content
			content := entry.Content
			if len(content) > 80 {
				content = content[:80] + "..."
			}
			// Replace newlines with spaces for better display
			content = strings.ReplaceAll(content, "\n", " ")
			content = strings.ReplaceAll(content, "\r", " ")

			// Combine time and content
			displayText := fmt.Sprintf("[%s] %s",
				entry.Timestamp.Format("15:04:05"), content)
			label.SetText(displayText)
		},
	)

	// Double-click to copy functionality
	ca.historyList.OnSelected = func(id widget.ListItemID) {
		history := ca.monitor.GetHistory()
		if id < len(history) {
			entry := history[id]
			err := ca.monitor.CopyToClipboard(entry.Content)
			if err != nil {
				dialog.ShowError(err, ca.window)
			} else {
				ca.statusLabel.SetText(fmt.Sprintf("Copied: %s",
					time.Now().Format("15:04:05")))
			}
		}
	}

	// Buttons
	clearBtn := widget.NewButton("Clear History", func() {
		dialog.ShowConfirm("Confirmation", "Are you sure you want to clear all history records?",
			func(confirmed bool) {
				if confirmed {
					ca.monitor.ClearHistory()
					ca.historyList.Refresh()
					ca.statusLabel.SetText("History cleared")
				}
			}, ca.window)
	})

	refreshBtn := widget.NewButton("Refresh", func() {
		ca.historyList.Refresh()
		ca.statusLabel.SetText("List refreshed")
	})

	aboutBtn := widget.NewButton("About", func() {
		dialog.ShowInformation("About",
			GetVersionInfo()+"\n\nFeatures:\n- Real-time clipboard monitoring\n- Display history records\n- Double-click to copy\n- Cross-platform support\n\nAuthor: Fly \nBuild version: "+Version,
			ca.window)
	})

	// Layout
	buttonContainer := container.NewHBox(
		clearBtn,
		refreshBtn,
		aboutBtn,
	)

	content := container.NewVBox(
		ca.statusLabel,
		widget.NewSeparator(),
		widget.NewLabel("Clipboard History (Double-click to copy):"),
		ca.historyList,
		widget.NewSeparator(),
		buttonContainer,
	)

	ca.window.SetContent(content)
}

func (ca *ClipboardApp) startMonitoring() {
	// Set new content callback
	ca.monitor.SetOnNewContent(func(entry clipboard.ClipboardEntry) {
		// Update UI in the UI thread
		ca.statusLabel.SetText(fmt.Sprintf("New content detected: %s",
			entry.Timestamp.Format("15:04:05")))
		ca.historyList.Refresh()
	})

	// Start monitoring in background
	go func() {
		err := ca.monitor.Start(ca.ctx)
		if err != nil && err != context.Canceled {
			log.Printf("Monitor error: %v", err)
		}
	}()
}

func (ca *ClipboardApp) Run() {
	ca.setupUI()
	ca.startMonitoring()

	// Set window close callback
	ca.window.SetCloseIntercept(func() {
		ca.cancel() // Stop monitoring
		ca.app.Quit()
	})

	ca.statusLabel.SetText("Status: Monitoring clipboard...")
	ca.window.ShowAndRun()
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
	app.Run()
}
