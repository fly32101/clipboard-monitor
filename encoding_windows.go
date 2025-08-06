//go:build windows
// +build windows

package main

import (
	"syscall"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleCP       = kernel32.NewProc("SetConsoleCP")
	procSetConsoleOutputCP = kernel32.NewProc("SetConsoleOutputCP")
)

// SetWindowsConsoleUTF8 设置 Windows 控制台为 UTF-8 编码
func SetWindowsConsoleUTF8() {
	// 设置控制台输入编码为 UTF-8 (65001)
	procSetConsoleCP.Call(uintptr(65001))
	// 设置控制台输出编码为 UTF-8 (65001)
	procSetConsoleOutputCP.Call(uintptr(65001))
}

// SetWindowsAppUTF8 设置 Windows 应用程序为 UTF-8 编码
func SetWindowsAppUTF8() {
	// 设置当前线程的代码页为 UTF-8
	kernel32.NewProc("SetThreadLocale").Call(uintptr(0x0804)) // 中文简体
}

// initPlatformSpecific 平台特定初始化
func initPlatformSpecific() {
	SetWindowsConsoleUTF8()
	SetWindowsAppUTF8()
}
