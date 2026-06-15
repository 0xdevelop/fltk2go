package stepper

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type UIStepper struct {
	v   view.UIView
	raw *fltk_bridge.Spinner
}

func NewUIStepper(r *foundation.Rect) *UIStepper {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 100, Height: 30}
	}

	spinner := fltk_bridge.NewSpinner(r.X, r.Y, r.Width, r.Height, "")
	spinner.SetType(fltk_bridge.SPINNER_FLOAT_INPUT)
	spinner.SetMinimum(0)
	spinner.SetMaximum(100)
	spinner.SetStep(1.0)
	spinner.SetValue(0)

	s := &UIStepper{raw: spinner}
	s.v.BindRaw(spinner)

	return s
}

func (s *UIStepper) View() *view.UIView { return &s.v }

func (s *UIStepper) Raw() *fltk_bridge.Spinner { return s.raw }

func (s *UIStepper) SetRange(min, max float64) {
	s.raw.SetMinimum(min)
	s.raw.SetMaximum(max)
}

func (s *UIStepper) SetStep(step float64) {
	s.raw.SetStep(step)
}

func (s *UIStepper) Value() float64 {
	return s.raw.Value()
}

func (s *UIStepper) SetValue(val float64) {
	s.raw.SetValue(val)
}

func (s *UIStepper) OnValueChanged(cb func(val float64)) {
	s.raw.SetCallback(func() {
		cb(s.raw.Value())
	})
}

// On 绑定事件
func (s *UIStepper) On(event fltk_bridge.Event, handler func(fltk_bridge.Event) bool) {
	s.v.On(event, handler)
}
