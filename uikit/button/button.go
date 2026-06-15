package button

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type UIButton struct {
	v   view.UIView
	raw *fltk_bridge.Button
}

func NewUIButton(r *foundation.Rect, title string) *UIButton {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 120, Height: 36}
	}

	btn := fltk_bridge.NewButton(r.X, r.Y, r.Width, r.Height, title)

	b := &UIButton{raw: btn}
	b.v.BindRaw(btn)
	return b
}

func (b *UIButton) View() *view.UIView {
	if b == nil {
		return nil
	}
	return &b.v
}

func (b *UIButton) SetTitle(s string) {
	if b != nil && b.raw != nil {
		b.raw.SetLabel(s)
	}
}

func (b *UIButton) SetBackgroundColor(rgb uint) {
	if b != nil && b.raw != nil {
		b.raw.SetColor(fltk_bridge.Color(rgb))
	}
}

func (b *UIButton) OnTouchUpInside(cb func()) {
	if b != nil && b.raw != nil {
		b.raw.SetCallback(cb)
	}
}

func (b *UIButton) SetTitleColor(rgb uint) {
	if b != nil && b.raw != nil {
		b.raw.SetLabelColor(fltk_bridge.Color(rgb))
	}
}

func (b *UIButton) Raw() *fltk_bridge.Button {
	if b == nil {
		return nil
	}
	return b.raw
}

// ButtonType 按钮类型
type ButtonType int

const (
	SystemButton   ButtonType = iota // 普通按钮
	CheckboxButton                   // 复选框按钮
	RadioButton                      // 单选按钮
	ToggleButton                     // 切换按钮
)

// NewUIButtonWithType 创建指定类型的按钮
func NewUIButtonWithType(r *foundation.Rect, title string, buttonType ButtonType) *UIButton {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 120, Height: 36}
	}

	var rawBtn *fltk_bridge.Button
	switch buttonType {
	case CheckboxButton:
		cb := fltk_bridge.NewCheckButton(r.X, r.Y, r.Width, r.Height, title)
		rawBtn = &cb.Button
	case RadioButton:
		rb := fltk_bridge.NewRadioButton(r.X, r.Y, r.Width, r.Height, title)
		rawBtn = &rb.Button
	case ToggleButton:
		tb := fltk_bridge.NewToggleButton(r.X, r.Y, r.Width, r.Height, title)
		rawBtn = &tb.Button
	default:
		rawBtn = fltk_bridge.NewButton(r.X, r.Y, r.Width, r.Height, title)
	}

	b := &UIButton{raw: rawBtn}
	b.v.BindRaw(rawBtn)
	return b
}
