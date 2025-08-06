package hotkey

import (
	"log"
	"sync"

	"github.com/robotn/gohook"
)

// Manager 热键管理器
type Manager struct {
	mu        sync.RWMutex
	callbacks map[string]func()
	running   bool
}

// NewManager 创建新的热键管理器
func NewManager() *Manager {
	return &Manager{
		callbacks: make(map[string]func()),
		running:   false,
	}
}

// RegisterHotkey 注册热键
func (m *Manager) RegisterHotkey(key string, callback func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callbacks[key] = callback
}

// Start 开始监听热键
func (m *Manager) Start() {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.mu.Unlock()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Hotkey manager panic: %v", r)
			}
		}()

		// 监听键盘事件
		evChan := hook.Start()
		defer hook.End()

		for ev := range evChan {
			if ev.Kind == hook.KeyDown {
				m.handleKeyEvent(ev)
			}
		}
	}()
}

// Stop 停止监听热键
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		hook.End()
		m.running = false
	}
}

// handleKeyEvent 处理键盘事件
func (m *Manager) handleKeyEvent(ev hook.Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 检查是否是注册的热键组合
	// Ctrl+Shift+V 组合键
	if ev.Mask&hook.MaskCtrl != 0 && ev.Mask&hook.MaskShift != 0 && ev.Keychar == 'v' {
		if callback, exists := m.callbacks["ctrl+shift+v"]; exists {
			go callback()
		}
	}

	// Ctrl+` 组合键 (反引号)
	if ev.Mask&hook.MaskCtrl != 0 && ev.Rawcode == 192 { // 192 是反引号的键码
		if callback, exists := m.callbacks["ctrl+`"]; exists {
			go callback()
		}
	}
}

// IsRunning 检查是否正在运行
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}
