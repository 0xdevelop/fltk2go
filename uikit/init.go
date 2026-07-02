package uikit

import (
	"github.com/0xdevelop/fltk2go/lang"
)

// Init 初始化UIKit
// 可以选择设置语言偏好
func Init(opts ...InitOption) {
	// 应用初始化选项
	for _, opt := range opts {
		opt()
	}
}

// InitOption 初始化选项类型
type InitOption func()

// WithLanguage 设置语言偏好
// l: 语言类型，支持 lang.LangEnglish 和 lang.LangChinese
func WithLanguage(l lang.Language) InitOption {
	return func() {
		lang.SetLanguage(l)
	}
}

// WithEnglish 设置使用英语
func WithEnglish() InitOption {
	return func() {
		lang.SetLanguage(lang.LangEnglish)
	}
}

// WithChinese 设置使用中文
func WithChinese() InitOption {
	return func() {
		lang.SetLanguage(lang.LangChinese)
	}
}
