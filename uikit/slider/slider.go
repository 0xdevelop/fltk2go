package slider

import (
	"fmt"

	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

type SliderStyle struct {
	Type           fltk_bridge.SliderType
	BoxType        fltk_bridge.BoxType
	Color          uint
	SelectionColor uint
	TextSize       int
}

func DefaultSliderStyle() SliderStyle {
	return SliderStyle{
		Type:           fltk_bridge.HOR_NICE_SLIDER,
		BoxType:        fltk_bridge.ROUND_DOWN_BOX,
		Color:          0xE0E0E000,
		SelectionColor: 0x42A5F500,
		TextSize:       13,
	}
}

type UISlider struct {
	v     view.UIView
	raw   *fltk_bridge.Slider
	style SliderStyle
}

func NewUISlider(r *foundation.Rect) *UISlider {
	return NewUISliderWithOptions(r, DefaultSliderStyle())
}

func NewUISliderWithOptions(r *foundation.Rect, style SliderStyle) *UISlider {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 160, Height: 28}
	}

	raw := fltk_bridge.NewSlider(r.X, r.Y, r.Width, r.Height)
	raw.SetMinimum(0)
	raw.SetMaximum(1)
	raw.SetStep(0.01)

	s := &UISlider{raw: raw, style: style}
	s.v.BindRaw(raw)
	s.v.SetAutomationRole("slider").SetAutomationValueHandler(func() (string, bool) {
		return fmt.Sprintf("%g", s.Value()), true
	})
	s.ApplyStyle(style)
	return s
}

func (s *UISlider) View() *view.UIView {
	if s == nil {
		return nil
	}
	return &s.v
}

func (s *UISlider) Raw() *fltk_bridge.Slider {
	if s == nil {
		return nil
	}
	return s.raw
}

func (s *UISlider) SetMinimumValue(v float64) {
	if s != nil && s.raw != nil {
		s.raw.SetMinimum(v)
	}
}

func (s *UISlider) SetMinimum(v float64) {
	s.SetMinimumValue(v)
}

func (s *UISlider) SetMaximumValue(v float64) {
	if s != nil && s.raw != nil {
		s.raw.SetMaximum(v)
	}
}

func (s *UISlider) SetMaximum(v float64) {
	s.SetMaximumValue(v)
}

func (s *UISlider) SetValue(v float64) {
	if s != nil && s.raw != nil {
		s.raw.SetValue(v)
	}
}

func (s *UISlider) Value() float64 {
	if s == nil || s.raw == nil {
		return 0
	}
	return s.raw.Value()
}

func (s *UISlider) SetStep(v float64) {
	if s != nil && s.raw != nil {
		s.raw.SetStep(v)
	}
}

// OnValueChanged invokes cb with the current slider value whenever FLTK reports
// a value-change event.
func (s *UISlider) OnValueChanged(cb func(float64)) {
	if s != nil && s.raw != nil {
		s.raw.SetCallbackCondition(fltk_bridge.WhenChanged)
		s.raw.SetCallback(func() {
			if cb != nil {
				cb(s.raw.Value())
			}
		})
	}
}

// OnChange is a convenience variant for callers that do not need the value.
func (s *UISlider) OnChange(cb func()) {
	if s != nil && s.raw != nil {
		s.raw.SetCallbackCondition(fltk_bridge.WhenChanged)
		s.raw.SetCallback(func() {
			if cb != nil {
				cb()
			}
		})
	}
}

func (s *UISlider) SetType(t fltk_bridge.SliderType) {
	if s != nil && s.raw != nil {
		s.raw.SetType(t)
	}
}

func (s *UISlider) SetVertical(vertical bool) {
	if s == nil || s.raw == nil {
		return
	}
	if vertical {
		s.raw.SetType(fltk_bridge.VERT_NICE_SLIDER)
		return
	}
	s.raw.SetType(fltk_bridge.HOR_NICE_SLIDER)
}

func (s *UISlider) ApplyStyle(style SliderStyle) {
	if s == nil || s.raw == nil {
		return
	}
	s.style = style
	s.raw.SetType(style.Type)
	s.raw.SetBox(style.BoxType)
	s.raw.SetColor(fltk_bridge.Color(style.Color))
	s.raw.SetSelectionColor(fltk_bridge.Color(style.SelectionColor))
	if textSizer, ok := any(s.raw).(interface{ SetTextSize(int) }); ok {
		textSizer.SetTextSize(style.TextSize)
	}
	s.raw.Redraw()
}

func (s *UISlider) Style() SliderStyle {
	if s == nil {
		return SliderStyle{}
	}
	return s.style
}
