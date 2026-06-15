package switchview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type UISwitch struct {
	v   view.UIView
	raw *fltk_bridge.ToggleButton
}

func NewUISwitch(r *foundation.Rect) *UISwitch {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 64, Height: 28}
	}

	raw := fltk_bridge.NewToggleButton(r.X, r.Y, r.Width, r.Height, "")
	s := &UISwitch{raw: raw}
	s.v.BindRaw(raw)
	return s
}

func (s *UISwitch) View() *view.UIView {
	if s == nil {
		return nil
	}
	return &s.v
}

func (s *UISwitch) Raw() *fltk_bridge.ToggleButton {
	if s == nil {
		return nil
	}
	return s.raw
}

func (s *UISwitch) SetOn(on bool) {
	if s != nil && s.raw != nil {
		s.raw.SetValue(on)
	}
}

func (s *UISwitch) SetValue(v bool) {
	s.SetOn(v)
}

func (s *UISwitch) IsOn() bool {
	if s == nil || s.raw == nil {
		return false
	}
	return s.raw.Value()
}

func (s *UISwitch) Value() bool {
	return s.IsOn()
}

// OnValueChanged invokes cb with the current on/off state.
func (s *UISwitch) OnValueChanged(cb func(bool)) {
	if s != nil && s.raw != nil {
		s.raw.SetCallback(func() {
			if cb != nil {
				cb(s.raw.Value())
			}
		})
	}
}

// OnChange is a convenience variant for callers that do not need the state.
func (s *UISwitch) OnChange(cb func()) {
	if s != nil && s.raw != nil {
		s.raw.SetCallback(func() {
			if cb != nil {
				cb()
			}
		})
	}
}

func (s *UISwitch) On(event fltk_bridge.Event, handler func(fltk_bridge.Event) bool) {
	if s != nil {
		s.v.On(event, handler)
	}
}
