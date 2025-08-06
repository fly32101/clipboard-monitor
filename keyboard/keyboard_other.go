//go:build !windows

package keyboard

import "fmt"

// SendCtrlV 发送 Ctrl+V 组合键 (非Windows平台暂不支持)
func SendCtrlV() error {
	return fmt.Errorf("keyboard simulation not supported on this platform")
}

// SendKey 发送单个按键 (非Windows平台暂不支持)
func SendKey(vkCode int) error {
	return fmt.Errorf("keyboard simulation not supported on this platform")
}

// IsKeyPressed 检查按键是否被按下 (非Windows平台暂不支持)
func IsKeyPressed(vkCode int) bool {
	return false
}
