package progress

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type UIProgressView struct {
	v   view.UIView
	raw *fltk_bridge.Progress
}

func NewUIProgressView(r *foundation.Rect) *UIProgressView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 160, Height: 20}
	}

	raw := fltk_bridge.NewProgress(r.X, r.Y, r.Width, r.Height)
	raw.SetMinimum(0)
	raw.SetMaximum(1)
	raw.SetValue(0)

	p := &UIProgressView{raw: raw}
	p.v.BindRaw(raw)
	return p
}

func (p *UIProgressView) View() *view.UIView {
	if p == nil {
		return nil
	}
	return &p.v
}

func (p *UIProgressView) Raw() *fltk_bridge.Progress {
	if p == nil {
		return nil
	}
	return p.raw
}

func (p *UIProgressView) SetProgress(v float64) {
	if p != nil && p.raw != nil {
		p.raw.SetValue(v)
		p.raw.Redraw()
	}
}

func (p *UIProgressView) Progress() float64 {
	if p == nil || p.raw == nil {
		return 0
	}
	return p.raw.Value()
}

func (p *UIProgressView) SetMinimumValue(v float64) {
	if p != nil && p.raw != nil {
		p.raw.SetMinimum(v)
	}
}

func (p *UIProgressView) SetMaximumValue(v float64) {
	if p != nil && p.raw != nil {
		p.raw.SetMaximum(v)
	}
}

func (p *UIProgressView) SetTrackColor(rgb uint) {
	if p != nil && p.raw != nil {
		p.raw.SetColor(fltk_bridge.Color(rgb))
	}
}

func (p *UIProgressView) SetProgressTintColor(rgb uint) {
	if p != nil && p.raw != nil {
		p.raw.SetSelectionColor(fltk_bridge.Color(rgb))
	}
}
