//go:build windows

package hotkey

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

var (
	user32                 = syscall.NewLazyDLL("user32.dll")
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procRegisterHotKey     = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey   = user32.NewProc("UnregisterHotKey")
	procGetMessage         = user32.NewProc("GetMessageW")
	procGetCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
)

// 修饰键常量
const (
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	MOD_WIN     = 0x0008
)

// 虚拟键码
const (
	VK_V = 0x56
)

// 消息类型
const (
	WM_HOTKEY = 0x0312
)

// MSG 结构体
type MSG struct {
	HWND    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// HotkeyManager 全局热键管理器
type HotkeyManager struct {
	registered bool
	callback   func()
	stopChan   chan bool
}

// NewHotkeyManager 创建新的热键管理器
func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{
		stopChan: make(chan bool),
	}
}

// RegisterHotkey 注册全局热键 (Ctrl+Shift+V)
func (hm *HotkeyManager) RegisterHotkey(callback func()) error {
	if hm.registered {
		return fmt.Errorf("hotkey already registered")
	}

	hm.callback = callback

	// 获取当前线程ID
	threadId, _, _ := procGetCurrentThreadId.Call()

	// 注册热键 Ctrl+Shift+V (ID=1)
	ret, _, err := procRegisterHotKey.Call(
		0,                              // NULL窗口句柄
		1,                              // 热键ID
		uintptr(MOD_CONTROL|MOD_SHIFT), // 修饰键
		uintptr(VK_V),                  // 虚拟键码
	)

	if ret == 0 {
		return fmt.Errorf("failed to register hotkey: %v", err)
	}

	hm.registered = true
	log.Printf("已注册全局热键: Ctrl+Shift+V (线程ID: %d)", threadId)

	// 启动消息循环
	go hm.messageLoop()

	return nil
}

// UnregisterHotkey 注销热键
func (hm *HotkeyManager) UnregisterHotkey() error {
	if !hm.registered {
		return nil
	}

	ret, _, err := procUnregisterHotKey.Call(0, 1)
	if ret == 0 {
		return fmt.Errorf("failed to unregister hotkey: %v", err)
	}

	hm.registered = false
	hm.stopChan <- true
	log.Printf("已注销全局热键")

	return nil
}

// messageLoop 消息循环
func (hm *HotkeyManager) messageLoop() {
	var msg MSG

	for {
		select {
		case <-hm.stopChan:
			return
		default:
			// 获取消息
			ret, _, _ := procGetMessage.Call(
				uintptr(unsafe.Pointer(&msg)),
				0,
				0,
				0,
			)

			if ret == 0 { // WM_QUIT
				return
			}

			if ret == ^uintptr(0) { // 错误
				continue
			}

			// 检查是否是热键消息
			if msg.Message == WM_HOTKEY && msg.WParam == 1 {
				log.Printf("检测到全局热键按下")
				if hm.callback != nil {
					go hm.callback()
				}
			}
		}
	}
}

// IsRegistered 检查热键是否已注册
func (hm *HotkeyManager) IsRegistered() bool {
	return hm.registered
}
