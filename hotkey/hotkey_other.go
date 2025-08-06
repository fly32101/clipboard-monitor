//go:build !windows

package hotkey

import "fmt"

// HotkeyManager 全局热键管理器 (非Windows平台)
type HotkeyManager struct{}

// NewHotkeyManager 创建新的热键管理器
func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{}
}

// RegisterHotkey 注册全局热键 (非Windows平台暂不支持)
func (hm *HotkeyManager) RegisterHotkey(callback func()) error {
	return fmt.Errorf("global hotkey not supported on this platform")
}

// UnregisterHotkey 注销热键 (非Windows平台暂不支持)
func (hm *HotkeyManager) UnregisterHotkey() error {
	return nil
}

// IsRegistered 检查热键是否已注册
func (hm *HotkeyManager) IsRegistered() bool {
	return false
}
