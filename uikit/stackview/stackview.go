package stackview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type Axis int

const (
	AxisVertical Axis = iota
	AxisHorizontal
)

type UIStackView struct {
	v   view.UIView
	raw *fltk_bridge.Flex
}

func NewUIStackView(r *foundation.Rect, axis Axis) *UIStackView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 240, Height: 160}
	}

	raw := fltk_bridge.NewFlex(r.X, r.Y, r.Width, r.Height)
	s := &UIStackView{raw: raw}
	s.v.BindRaw(raw)
	s.v.BindHost(raw)
	s.SetAxis(axis)
	return s
}

func (s *UIStackView) View() *view.UIView {
	if s == nil {
		return nil
	}
	return &s.v
}

func (s *UIStackView) Raw() *fltk_bridge.Flex {
	if s == nil {
		return nil
	}
	return s.raw
}

func (s *UIStackView) AddArrangedSubview(child view.Viewable) {
	if s == nil || s.raw == nil || child == nil {
		return
	}
	cv := child.View()
	if cv == nil || cv.Raw() == nil {
		return
	}
	s.raw.Add(cv.Raw())
	cv.BindHost(s.raw)
}

func (s *UIStackView) SetAxis(axis Axis) {
	if s == nil || s.raw == nil {
		return
	}
	if axis == AxisHorizontal {
		s.raw.SetType(fltk_bridge.ROW)
		return
	}
	s.raw.SetType(fltk_bridge.COLUMN)
}

func (s *UIStackView) SetSpacing(spacing int) {
	if s != nil && s.raw != nil {
		s.raw.SetGap(spacing)
		s.raw.SetSpacing(spacing)
	}
}

func (s *UIStackView) SetMargin(margin int) {
	if s != nil && s.raw != nil {
		s.raw.SetMargin(margin)
	}
}

func (s *UIStackView) SetFixedSize(child view.Viewable, size int) {
	if s == nil || s.raw == nil || child == nil {
		return
	}
	cv := child.View()
	if cv == nil || cv.Raw() == nil {
		return
	}
	s.raw.Fixed(cv.Raw(), size)
}

func (s *UIStackView) Layout() {
	if s != nil && s.raw != nil {
		s.raw.Layout()
	}
}

func (s *UIStackView) End() {
	if s != nil && s.raw != nil {
		s.raw.End()
	}
}
