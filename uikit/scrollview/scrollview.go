package scrollview

import (
	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

type UIScrollViewScrollType = fltk_bridge.ScrollType

var (
	ScrollHorizontal       = fltk_bridge.SCROLL_HORIZONTAL
	ScrollVertical         = fltk_bridge.SCROLL_VERTICAL
	ScrollBoth             = fltk_bridge.SCROLL_BOTH
	ScrollHorizontalAlways = fltk_bridge.SCROLL_HORIZONTAL_ALWAYS
	ScrollVerticalAlways   = fltk_bridge.SCROLL_VERTICAL_ALWAYS
	ScrollBothAlways       = fltk_bridge.SCROLL_BOTH_ALWAYS
)

type UIScrollView struct {
	v   view.UIView
	raw *fltk_bridge.Scroll
}

func NewUIScrollView(r *foundation.Rect) *UIScrollView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 240, Height: 160}
	}

	raw := fltk_bridge.NewScroll(r.X, r.Y, r.Width, r.Height)
	raw.SetType(fltk_bridge.SCROLL_BOTH)

	s := &UIScrollView{raw: raw}
	s.v.BindRaw(raw)
	s.v.BindHost(raw)
	return s
}

func (s *UIScrollView) View() *view.UIView {
	if s == nil {
		return nil
	}
	return &s.v
}

func (s *UIScrollView) Raw() *fltk_bridge.Scroll {
	if s == nil {
		return nil
	}
	return s.raw
}

func (s *UIScrollView) AddSubview(child view.Viewable) {
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

func (s *UIScrollView) ScrollTo(x, y int) {
	if s != nil && s.raw != nil {
		s.raw.ScrollTo(x, y)
	}
}

func (s *UIScrollView) ContentOffset() (int, int) {
	if s == nil || s.raw == nil {
		return 0, 0
	}
	return s.raw.XPosition(), s.raw.YPosition()
}

func (s *UIScrollView) XPosition() int {
	if s == nil || s.raw == nil {
		return 0
	}
	return s.raw.XPosition()
}

func (s *UIScrollView) YPosition() int {
	if s == nil || s.raw == nil {
		return 0
	}
	return s.raw.YPosition()
}

func (s *UIScrollView) SetScrollType(t fltk_bridge.ScrollType) {
	if s != nil && s.raw != nil {
		s.raw.SetType(t)
	}
}

func (s *UIScrollView) RemoveSubview(child view.Viewable) {
	if s == nil || s.raw == nil || child == nil {
		return
	}
	cv := child.View()
	if cv == nil || cv.Raw() == nil {
		return
	}
	s.raw.Remove(cv.Raw())
}
