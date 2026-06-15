package dropdown

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type UIDropdown struct {
	v       view.UIView
	raw     *fltk_bridge.Choice
	options []string
}

func NewUIDropdown(r *foundation.Rect) *UIDropdown {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 120, Height: 30}
	}

	choice := fltk_bridge.NewChoice(r.X, r.Y, r.Width, r.Height, "")

	dd := &UIDropdown{
		raw: choice,
	}
	dd.v.BindRaw(choice)

	return dd
}

func (dd *UIDropdown) View() *view.UIView { return &dd.v }

func (dd *UIDropdown) Raw() *fltk_bridge.Choice { return dd.raw }

func (dd *UIDropdown) SetOptions(options []string) {
	dd.raw.Clear()
	dd.options = options
	for _, opt := range options {
		dd.raw.Add(opt, nil)
	}
	if len(options) > 0 {
		dd.raw.SetValue(0)
	}
}

func (dd *UIDropdown) Options() []string {
	return dd.options
}

func (dd *UIDropdown) SelectedIndex() int {
	return dd.raw.Value()
}

func (dd *UIDropdown) SelectedOption() string {
	idx := dd.SelectedIndex()
	if idx >= 0 && idx < len(dd.options) {
		return dd.options[idx]
	}
	return ""
}

func (dd *UIDropdown) SetSelectedIndex(index int) {
	if index >= 0 && index < len(dd.options) {
		dd.raw.SetValue(index)
	}
}

func (dd *UIDropdown) OnSelectionChanged(cb func(index int, option string)) {
	dd.raw.SetCallback(func() {
		idx := dd.SelectedIndex()
		opt := ""
		if idx >= 0 && idx < len(dd.options) {
			opt = dd.options[idx]
		}
		cb(idx, opt)
	})
}

// On 绑定事件
func (dd *UIDropdown) On(event fltk_bridge.Event, handler func(fltk_bridge.Event) bool) {
	dd.v.On(event, handler)
}
