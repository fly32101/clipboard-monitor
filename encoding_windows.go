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

// initPlatformSpecific 平台特定初始化
func initPlatformSpecific() {
	SetWindowsConsoleUTF8()
}
