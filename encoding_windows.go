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
	procSetThreadLocale    = kernel32.NewProc("SetThreadLocale")
)

// SetWindowsConsoleGBK 设置 Windows 控制台为 GBK 编码
func SetWindowsConsoleGBK() {
	// 设置控制台输入编码为 GBK (936)
	procSetConsoleCP.Call(uintptr(936))
	// 设置控制台输出编码为 GBK (936)
	procSetConsoleOutputCP.Call(uintptr(936))
}

// SetWindowsAppGBK 设置 Windows 应用程序为 GBK 编码
func SetWindowsAppGBK() {
	// 设置当前线程的区域设置为中文简体 GBK
	procSetThreadLocale.Call(uintptr(0x0804)) // 中文简体 (China)
}

// initPlatformSpecific 平台特定初始化
func initPlatformSpecific() {
	SetWindowsConsoleGBK()
	SetWindowsAppGBK()
}
