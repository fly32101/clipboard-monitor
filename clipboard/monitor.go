package clipboard

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/atotto/clipboard"
)

// ClipboardEntry 表示剪贴板条目
type ClipboardEntry struct {
	Content   string
	Timestamp time.Time
}

// Monitor 剪贴板监听器
type Monitor struct {
	mu           sync.RWMutex
	history      []ClipboardEntry
	lastContent  string
	maxHistory   int
	onNewContent func(entry ClipboardEntry)
}

// NewMonitor 创建新的剪贴板监听器
func NewMonitor(maxHistory int) *Monitor {
	return &Monitor{
		history:    make([]ClipboardEntry, 0),
		maxHistory: maxHistory,
	}
}

// SetOnNewContent 设置新内容回调函数
func (m *Monitor) SetOnNewContent(callback func(entry ClipboardEntry)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onNewContent = callback
}

// Start 开始监听剪贴板
func (m *Monitor) Start(ctx context.Context) error {
	// 获取初始剪贴板内容
	initialContent, err := clipboard.ReadAll()
	if err == nil && initialContent != "" {
		m.lastContent = initialContent
		entry := ClipboardEntry{
			Content:   initialContent,
			Timestamp: time.Now(),
		}
		m.addToHistory(entry)
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			content, err := clipboard.ReadAll()
			if err != nil {
				continue
			}

			m.mu.RLock()
			lastContent := m.lastContent
			m.mu.RUnlock()

			if content != lastContent && content != "" {
				entry := ClipboardEntry{
					Content:   content,
					Timestamp: time.Now(),
				}

				m.mu.Lock()
				m.lastContent = content
				m.addToHistory(entry)
				callback := m.onNewContent
				m.mu.Unlock()

				if callback != nil {
					callback(entry)
				}
			}
		}
	}
}

// addToHistory 添加到历史记录（支持去重）
func (m *Monitor) addToHistory(entry ClipboardEntry) {
	// 查找是否已存在相同内容
	for i, existingEntry := range m.history {
		if existingEntry.Content == entry.Content {
			// 找到重复内容，更新时间戳并移动到顶部
			m.history[i].Timestamp = entry.Timestamp
			// 将该项移动到列表顶部
			updatedEntry := m.history[i]
			m.history = append(m.history[:i], m.history[i+1:]...)
			m.history = append([]ClipboardEntry{updatedEntry}, m.history...)
			return
		}
	}

	// 没有找到重复内容，添加新项
	m.history = append([]ClipboardEntry{entry}, m.history...)
	if len(m.history) > m.maxHistory {
		m.history = m.history[:m.maxHistory]
	}
}

// GetHistory 获取历史记录
func (m *Monitor) GetHistory() []ClipboardEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回副本以避免并发问题
	history := make([]ClipboardEntry, len(m.history))
	copy(history, m.history)
	return history
}

// ClearHistory 清空历史记录
func (m *Monitor) ClearHistory() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = make([]ClipboardEntry, 0)
}

// CopyToClipboard 复制内容到剪贴板
func (m *Monitor) CopyToClipboard(content string) error {
	return clipboard.WriteAll(content)
}

// DeleteHistoryItem 删除指定索引的历史记录项
func (m *Monitor) DeleteHistoryItem(index int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if index < 0 || index >= len(m.history) {
		return fmt.Errorf("index out of range: %d", index)
	}

	// 删除指定索引的项目
	m.history = append(m.history[:index], m.history[index+1:]...)
	return nil
}
