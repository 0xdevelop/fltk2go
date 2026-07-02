package textfield

import (
	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

// UISecretTextField 密码输入框组件 (UIKit 风格封装)
type UISecretTextField struct {
	raw *fltk_bridge.SecretInput
	v   view.UIView
}

// NewUISecretTextField 创建一个新的密码输入框
func NewUISecretTextField(x, y, width, height int, placeholder string) *UISecretTextField {
	input := fltk_bridge.NewSecretInput(x, y, width, height)
	if placeholder != "" {
		input.SetLabel(placeholder)
		input.SetAlign(fltk_bridge.ALIGN_INSIDE | fltk_bridge.ALIGN_LEFT)
	}

	// 默认样式调整，使其更贴近 iOS 风格
	input.SetBox(fltk_bridge.ROUNDED_BOX)
	input.SetColor(fltk_bridge.Color(0xFFFFFF00)) // 假设默认白底

	tf := &UISecretTextField{
		raw: input,
	}

	// 绑定到底层视图
	tf.v.BindRaw(input)

	return tf
}

// View 返回基础视图，实现 view.Viewable 接口
func (tf *UISecretTextField) View() *view.UIView {
	return &tf.v
}

// Raw 返回底层的 FLTK 组件
func (tf *UISecretTextField) Raw() fltk_bridge.Widget {
	return tf.raw
}

// Text 获取输入框内容
func (tf *UISecretTextField) Text() string {
	if tf.raw != nil {
		return tf.raw.Value()
	}
	return ""
}

// SetText 设置输入框内容
func (tf *UISecretTextField) SetText(text string) {
	if tf.raw != nil {
		tf.raw.SetValue(text)
	}
}

// Placeholder 获取占位符文本
func (tf *UISecretTextField) Placeholder() string {
	if tf.raw != nil {
		return tf.raw.Label()
	}
	return ""
}

// SetPlaceholder 设置占位符文本
func (tf *UISecretTextField) SetPlaceholder(placeholder string) {
	if tf.raw != nil {
		tf.raw.SetLabel(placeholder)
		tf.raw.SetAlign(fltk_bridge.ALIGN_INSIDE | fltk_bridge.ALIGN_LEFT)
	}
}

// SetFontSize 设置字体大小
func (tf *UISecretTextField) SetFontSize(size int) {
	if tf.raw != nil {
		tf.raw.SetLabelSize(size)
	}
}

// SetFont 设置字体
func (tf *UISecretTextField) SetFont(font fltk_bridge.Font) {
	if tf.raw != nil {
		tf.raw.SetLabelFont(font)
	}
}

// SetTextColor 设置文本颜色
func (tf *UISecretTextField) SetTextColor(color uint) {
	if tf.raw != nil {
		tf.raw.SetLabelColor(fltk_bridge.Color(color))
	}
}

// SetBackgroundColor 设置背景颜色
func (tf *UISecretTextField) SetBackgroundColor(color uint) {
	if tf.raw != nil {
		tf.raw.SetColor(fltk_bridge.Color(color))
	}
}

// SetEnabled 设置是否可用
func (tf *UISecretTextField) SetEnabled(enabled bool) {
	if tf.raw != nil {
		if enabled {
			tf.raw.Activate()
		} else {
			tf.raw.Deactivate()
		}
	}
}

// IsEnabled 获取是否可用
func (tf *UISecretTextField) IsEnabled() bool {
	if tf.raw != nil {
		return tf.raw.IsActive()
	}
	return false
}

// OnChange 设置文本内容变化时的回调函数
func (tf *UISecretTextField) OnChange(callback func()) {
	if tf.raw != nil {
		tf.raw.SetCallback(callback)
		tf.raw.SetCallbackCondition(fltk_bridge.WhenChanged)
	}
}

// SetCornerRadius 模拟 iOS 设置圆角 (通过 FLTK 盒模型模拟)
func (tf *UISecretTextField) SetCornerRadius() {
	if tf.raw != nil {
		tf.raw.SetBox(fltk_bridge.ROUNDED_BOX)
	}
}

// SetBorderStyle 模拟 iOS 设置边框样式
func (tf *UISecretTextField) SetBorderStyle(box fltk_bridge.BoxType) {
	if tf.raw != nil {
		tf.raw.SetBox(box)
	}
}
