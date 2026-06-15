package label

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type UILabel struct {
	v   view.UIView
	raw *fltk_bridge.Box
}

func NewUILabel(r *foundation.Rect, text string) *UILabel {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 100, Height: 30}
	}

	b := fltk_bridge.NewBox(fltk_bridge.NO_BOX, r.X, r.Y, r.Width, r.Height, text)

	l := &UILabel{
		raw: b,
	}
	l.v.BindRaw(b)
	return l
}

func (l *UILabel) View() *view.UIView {
	if l == nil {
		return nil
	}
	return &l.v
}

func (l *UILabel) SetText(s string) {
	if l != nil && l.raw != nil {
		l.raw.SetLabel(s)
	}
}

func (l *UILabel) SetFontSize(px int) {
	if l != nil && l.raw != nil {
		l.raw.SetLabelSize(px)
	}
}

func (l *UILabel) SetTextColor(c uint) {
	if l != nil && l.raw != nil {
		l.raw.SetLabelColor(fltk_bridge.Color(c))
	}
}

func (l *UILabel) SetFont(font fltk_bridge.Font) {
	if l != nil && l.raw != nil {
		l.raw.SetLabelFont(font)
	}
}

func (l *UILabel) SetAlignment(align fltk_bridge.Align) {
	if l != nil && l.raw != nil {
		l.raw.SetAlign(align)
	}
}

func (l *UILabel) SetFrame(boxType fltk_bridge.BoxType) {
	if l != nil && l.raw != nil {
		l.raw.SetBox(boxType)
	}
}

func (l *UILabel) SetBackgroundColor(rgb uint) {
	if l != nil && l.raw != nil {
		l.raw.SetColor(fltk_bridge.Color(rgb))
	}
}

func (l *UILabel) Raw() *fltk_bridge.Box {
	if l == nil {
		return nil
	}
	return l.raw
}
