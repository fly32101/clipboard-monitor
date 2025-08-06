package main

// Version 信息
const (
	Version   = "1.0.0"
	BuildDate = "2024-01-01"
	GitCommit = "unknown"
)

// GetVersionInfo 返回版本信息字符串
func GetVersionInfo() string {
	return "Clipboard Monitor v" + Version
}
