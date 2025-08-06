package clipboard

import (
	"context"
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {
	monitor := NewMonitor(10)
	if monitor == nil {
		t.Fatal("NewMonitor returned nil")
	}

	if monitor.maxHistory != 10 {
		t.Errorf("Expected maxHistory to be 10, got %d", monitor.maxHistory)
	}

	if len(monitor.history) != 0 {
		t.Errorf("Expected empty history, got %d items", len(monitor.history))
	}
}

func TestAddToHistory(t *testing.T) {
	monitor := NewMonitor(3)

	// 添加第一个条目
	entry1 := ClipboardEntry{
		Content:   "test content 1",
		Timestamp: time.Now(),
	}
	monitor.addToHistory(entry1)

	history := monitor.GetHistory()
	if len(history) != 1 {
		t.Errorf("Expected 1 item in history, got %d", len(history))
	}

	if history[0].Content != "test content 1" {
		t.Errorf("Expected 'test content 1', got '%s'", history[0].Content)
	}

	// 添加更多条目
	entry2 := ClipboardEntry{
		Content:   "test content 2",
		Timestamp: time.Now(),
	}
	monitor.addToHistory(entry2)

	entry3 := ClipboardEntry{
		Content:   "test content 3",
		Timestamp: time.Now(),
	}
	monitor.addToHistory(entry3)

	history = monitor.GetHistory()
	if len(history) != 3 {
		t.Errorf("Expected 3 items in history, got %d", len(history))
	}

	// 检查顺序（最新的应该在前面）
	if history[0].Content != "test content 3" {
		t.Errorf("Expected newest item first, got '%s'", history[0].Content)
	}

	// 添加第四个条目，应该移除最老的
	entry4 := ClipboardEntry{
		Content:   "test content 4",
		Timestamp: time.Now(),
	}
	monitor.addToHistory(entry4)

	history = monitor.GetHistory()
	if len(history) != 3 {
		t.Errorf("Expected 3 items in history after overflow, got %d", len(history))
	}

	// 检查最老的条目是否被移除
	found := false
	for _, entry := range history {
		if entry.Content == "test content 1" {
			found = true
			break
		}
	}
	if found {
		t.Error("Oldest entry should have been removed")
	}
}

func TestClearHistory(t *testing.T) {
	monitor := NewMonitor(10)

	// 添加一些条目
	entry := ClipboardEntry{
		Content:   "test content",
		Timestamp: time.Now(),
	}
	monitor.addToHistory(entry)

	if len(monitor.GetHistory()) == 0 {
		t.Error("History should not be empty before clearing")
	}

	// 清空历史
	monitor.ClearHistory()

	if len(monitor.GetHistory()) != 0 {
		t.Error("History should be empty after clearing")
	}
}

func TestSetOnNewContent(t *testing.T) {
	monitor := NewMonitor(10)

	called := false
	var receivedEntry ClipboardEntry

	monitor.SetOnNewContent(func(entry ClipboardEntry) {
		called = true
		receivedEntry = entry
	})

	// 模拟添加新内容
	entry := ClipboardEntry{
		Content:   "test callback",
		Timestamp: time.Now(),
	}

	// 直接调用回调函数进行测试
	monitor.mu.RLock()
	callback := monitor.onNewContent
	monitor.mu.RUnlock()

	if callback != nil {
		callback(entry)
	}

	if !called {
		t.Error("Callback should have been called")
	}

	if receivedEntry.Content != "test callback" {
		t.Errorf("Expected 'test callback', got '%s'", receivedEntry.Content)
	}
}

func TestMonitorContext(t *testing.T) {
	monitor := NewMonitor(10)

	// 创建一个可以取消的上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 立即取消上下文
	cancel()

	// 启动监听器应该立即返回
	err := monitor.Start(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}
