package lang

import (
	"os"
	"runtime"
	"strings"
)

// Language 表示支持的语言类型
type Language string

const (
	// LangEnglish 英语
	LangEnglish Language = "en"
	// LangChinese 中文
	LangChinese Language = "zh"
)

var (
	// currentLang 当前使用的语言
	currentLang Language
	// langDetected 是否已检测过语言
	langDetected bool
)

// DetectLanguage 检测本地系统语言
func DetectLanguage() Language {
	if langDetected {
		return currentLang
	}

	var lang Language = LangEnglish

	// 根据不同操作系统检测语言
	switch runtime.GOOS {
	case "windows":
		// Windows 系统通过环境变量检测
		langID := os.Getenv("LANG")
		if langID == "" {
			langID = os.Getenv("LC_ALL")
		}
		if langID == "" {
			langID = os.Getenv("LC_MESSAGES")
		}
		if langID == "" {
			langID = os.Getenv("LANGUAGE")
		}
		if strings.HasPrefix(langID, "zh") {
			lang = LangChinese
		}
	case "darwin":
		// macOS 系统通过环境变量检测
		langID := os.Getenv("LANG")
		if strings.HasPrefix(langID, "zh") {
			lang = LangChinese
		}
	case "linux":
		// Linux 系统通过环境变量检测
		langID := os.Getenv("LANG")
		if strings.HasPrefix(langID, "zh") {
			lang = LangChinese
		}
	}

	currentLang = lang
	langDetected = true
	return lang
}

// GetCurrentLanguage 获取当前使用的语言
func GetCurrentLanguage() Language {
	return DetectLanguage()
}

// SetLanguage 手动设置语言偏好
func SetLanguage(lang Language) {
	currentLang = lang
	langDetected = true
}

// IsChinese 检查当前是否使用中文
func IsChinese() bool {
	return DetectLanguage() == LangChinese
}

// GetIncludePath 根据语言获取头文件路径
func GetIncludePath() string {
	if IsChinese() {
		return "libs/fltk/include_cn"
	}
	return "libs/fltk/include"
}
