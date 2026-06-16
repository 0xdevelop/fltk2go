package progress

import (
	"fmt"

	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type ProgressStyle struct {
	BoxType        fltk_bridge.BoxType
	Color          uint
	SelectionColor uint
}

func DefaultProgressStyle() ProgressStyle {
	return ProgressStyle{
		BoxType:        fltk_bridge.ROUND_DOWN_BOX,
		Color:          0xE0E0E000,
		SelectionColor: 0x4CAF5000,
	}
}

type UIProgressView struct {
	v     view.UIView
	raw   *fltk_bridge.Progress
	style ProgressStyle
}

type UIProgress = UIProgressView

func NewUIProgressView(r *foundation.Rect) *UIProgressView {
	return NewUIProgressWithOptions(r, DefaultProgressStyle())
}

func NewUIProgress(r *foundation.Rect) *UIProgressView {
	return NewUIProgressWithOptions(r, DefaultProgressStyle())
}

func NewUIProgressWithOptions(r *foundation.Rect, style ProgressStyle) *UIProgressView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 160, Height: 20}
	}

	raw := fltk_bridge.NewProgress(r.X, r.Y, r.Width, r.Height)
	raw.SetMinimum(0)
	raw.SetMaximum(1)
	raw.SetValue(0)

	p := &UIProgressView{raw: raw, style: style}
	p.v.BindRaw(raw)
	p.v.SetAutomationRole("progressbar").SetAutomationValueHandler(func() (string, bool) {
		return fmt.Sprintf("%g", p.Value()), true
	})
	p.ApplyStyle(style)
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

func (p *UIProgressView) SetMinimum(v float64) {
	p.SetMinimumValue(v)
}

func (p *UIProgressView) SetMaximumValue(v float64) {
	if p != nil && p.raw != nil {
		p.raw.SetMaximum(v)
	}
}

func (p *UIProgressView) SetMaximum(v float64) {
	p.SetMaximumValue(v)
}

func (p *UIProgressView) SetValue(v float64) {
	p.SetProgress(v)
}

func (p *UIProgressView) Value() float64 {
	return p.Progress()
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

func (p *UIProgressView) ApplyStyle(style ProgressStyle) {
	if p == nil || p.raw == nil {
		return
	}
	p.style = style
	p.raw.SetBox(style.BoxType)
	p.raw.SetColor(fltk_bridge.Color(style.Color))
	p.raw.SetSelectionColor(fltk_bridge.Color(style.SelectionColor))
	p.raw.Redraw()
}

func (p *UIProgressView) Style() ProgressStyle {
	if p == nil {
		return ProgressStyle{}
	}
	return p.style
}
