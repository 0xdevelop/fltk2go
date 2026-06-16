package label

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/textlayout"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type LabelStyle struct {
	Font      fltk_bridge.Font
	FontSize  int
	TextColor uint
	Align     fltk_bridge.Align
	BgColor   uint
	BoxType   fltk_bridge.BoxType
}

func DefaultLabelStyle() LabelStyle {
	return LabelStyle{
		Font:      fltk_bridge.HELVETICA,
		FontSize:  14,
		TextColor: 0,
		Align:     fltk_bridge.ALIGN_CENTER | fltk_bridge.ALIGN_CLIP,
		BgColor:   0,
		BoxType:   fltk_bridge.NO_BOX,
	}
}

type UILabel struct {
	v   view.UIView
	raw *fltk_bridge.Box

	text         string
	preparedText *textlayout.PreparedText

	style LabelStyle
}

func NewUILabel(r *foundation.Rect, text string) *UILabel {
	return NewUILabelWithOptions(r, text, DefaultLabelStyle())
}

func NewUILabelWithOptions(r *foundation.Rect, text string, style LabelStyle) *UILabel {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 100, Height: 30}
	}

	b := fltk_bridge.NewBox(fltk_bridge.NO_BOX, r.X, r.Y, r.Width, r.Height, "")

	l := &UILabel{
		raw:   b,
		style: style,
	}
	l.v.BindRaw(b)
	l.v.SetAutomationRole("text").SetAutomationName(text).SetAutomationValueHandler(func() (string, bool) {
		return l.text, true
	})

	l.SetText(text)

	b.SetDrawHandler(func(drawSuper func()) {
		// Draw background box if any
		if l.style.BoxType != fltk_bridge.NO_BOX {
			fltk_bridge.DrawBox(l.style.BoxType, l.raw.X(), l.raw.Y(), l.raw.W(), l.raw.H(), fltk_bridge.Color(l.style.BgColor))
		} else if l.style.BgColor != 0 {
			fltk_bridge.DrawBox(fltk_bridge.FLAT_BOX, l.raw.X(), l.raw.Y(), l.raw.W(), l.raw.H(), fltk_bridge.Color(l.style.BgColor))
		}

		if l.preparedText == nil {
			return
		}

		// Draw text using textlayout
		x, y, w, h := l.raw.X(), l.raw.Y(), l.raw.W(), l.raw.H()
		fltk_bridge.PushClip(x, y, w, h)

		fltk_bridge.SetDrawColor(fltk_bridge.Color(l.style.TextColor))
		fltk_bridge.SetDrawFont(l.style.Font, l.style.FontSize)

		lineHeight := l.style.FontSize + 4
		layoutResult := textlayout.Layout(l.preparedText, w, lineHeight)

		startY := y
		if (l.style.Align & fltk_bridge.ALIGN_BOTTOM) != 0 {
			startY = y + h - layoutResult.Height
		} else if (l.style.Align & fltk_bridge.ALIGN_TOP) != 0 {
			startY = y
		} else { // center vertical
			startY = y + (h-layoutResult.Height)/2
		}

		if startY < y {
			startY = y
		}

		for i, line := range layoutResult.Lines {
			fltk_bridge.Draw(line.Text, x, startY+(i*lineHeight), w, lineHeight, l.style.Align)
		}

		fltk_bridge.PopClip()
	})

	return l
}

func (l *UILabel) View() *view.UIView {
	if l == nil {
		return nil
	}
	return &l.v
}

func (l *UILabel) SetText(s string) {
	if l == nil {
		return
	}
	l.text = s
	l.v.SetAutomationName(s)
	if s != "" {
		l.preparedText = textlayout.Prepare(s, l.style.Font, l.style.FontSize)
	} else {
		l.preparedText = nil
	}
	if l.raw != nil {
		l.raw.Redraw()
	}
}

func (l *UILabel) ApplyStyle(style LabelStyle) {
	if l == nil {
		return
	}
	l.style = style
	if l.text != "" {
		l.preparedText = textlayout.Prepare(l.text, l.style.Font, l.style.FontSize)
	}
	if l.raw != nil {
		l.raw.Redraw()
	}
}

func (l *UILabel) Style() LabelStyle {
	if l == nil {
		return LabelStyle{}
	}
	return l.style
}

func (l *UILabel) SetFontSize(px int) {
	if l == nil {
		return
	}
	l.style.FontSize = px
	if l.text != "" {
		l.preparedText = textlayout.Prepare(l.text, l.style.Font, l.style.FontSize)
	}
	if l.raw != nil {
		l.raw.Redraw()
	}
}

func (l *UILabel) SetTextColor(c uint) {
	if l == nil {
		return
	}
	l.style.TextColor = c
	if l.raw != nil {
		l.raw.Redraw()
	}
}

func (l *UILabel) SetFont(font fltk_bridge.Font) {
	if l == nil {
		return
	}
	l.style.Font = font
	if l.text != "" {
		l.preparedText = textlayout.Prepare(l.text, l.style.Font, l.style.FontSize)
	}
	if l.raw != nil {
		l.raw.Redraw()
	}
}

func (l *UILabel) SetAlignment(align fltk_bridge.Align) {
	if l == nil {
		return
	}
	l.style.Align = align
	if l.raw != nil {
		l.raw.Redraw()
	}
}

func (l *UILabel) SetFrame(boxType fltk_bridge.BoxType) {
	if l == nil {
		return
	}
	l.style.BoxType = boxType
	if l.raw != nil {
		l.raw.Redraw()
	}
}

func (l *UILabel) SetBackgroundColor(rgb uint) {
	if l == nil {
		return
	}
	l.style.BgColor = rgb
	if l.raw != nil {
		l.raw.Redraw()
	}
}

func (l *UILabel) Raw() *fltk_bridge.Box {
	if l == nil {
		return nil
	}
	return l.raw
}
