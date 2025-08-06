//go:build windows

package keyboard

import (
	"syscall"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procKeybd_event = user32.NewProc("keybd_event")
	procGetKeyState = user32.NewProc("GetKeyState")
)

// Virtual key codes
const (
	VK_CONTROL = 0x11
	VK_V       = 0x56
)

// Key event flags
const (
	KEYEVENTF_KEYUP = 0x0002
)

// SendCtrlV 发送 Ctrl+V 组合键
func SendCtrlV() error {
	// 按下 Ctrl 键
	procKeybd_event.Call(
		uintptr(VK_CONTROL),
		0,
		0,
		0,
	)

	// 按下 V 键
	procKeybd_event.Call(
		uintptr(VK_V),
		0,
		0,
		0,
	)

	// 释放 V 键
	procKeybd_event.Call(
		uintptr(VK_V),
		0,
		uintptr(KEYEVENTF_KEYUP),
		0,
	)

	// 释放 Ctrl 键
	procKeybd_event.Call(
		uintptr(VK_CONTROL),
		0,
		uintptr(KEYEVENTF_KEYUP),
		0,
	)

	return nil
}

// SendKey 发送单个按键
func SendKey(vkCode int) error {
	// 按下键
	procKeybd_event.Call(
		uintptr(vkCode),
		0,
		0,
		0,
	)

	// 释放键
	procKeybd_event.Call(
		uintptr(vkCode),
		0,
		uintptr(KEYEVENTF_KEYUP),
		0,
	)

	return nil
}

// IsKeyPressed 检查按键是否被按下
func IsKeyPressed(vkCode int) bool {
	ret, _, _ := procGetKeyState.Call(uintptr(vkCode))
	return (ret & 0x8000) != 0
}
