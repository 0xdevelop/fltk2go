package textview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type UITextView struct {
	v      view.UIView
	raw    *fltk_bridge.TextEditor
	buffer *fltk_bridge.TextBuffer
}

func NewUITextView(r *foundation.Rect) *UITextView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 240, Height: 120}
	}

	raw := fltk_bridge.NewTextEditor(r.X, r.Y, r.Width, r.Height)
	buffer := fltk_bridge.NewTextBuffer()
	raw.SetBuffer(buffer)
	raw.SetWrapMode(fltk_bridge.WRAP_AT_BOUNDS)

	t := &UITextView{raw: raw, buffer: buffer}
	t.v.BindRaw(raw)
	return t
}

func (t *UITextView) View() *view.UIView {
	if t == nil {
		return nil
	}
	return &t.v
}

func (t *UITextView) Raw() *fltk_bridge.TextEditor {
	if t == nil {
		return nil
	}
	return t.raw
}

func (t *UITextView) TextBuffer() *fltk_bridge.TextBuffer {
	if t == nil {
		return nil
	}
	return t.buffer
}

func (t *UITextView) SetText(text string) {
	if t != nil && t.buffer != nil {
		t.buffer.SetText(text)
	}
}

func (t *UITextView) Text() string {
	if t == nil || t.buffer == nil {
		return ""
	}
	return t.buffer.Text()
}

func (t *UITextView) Append(text string) {
	if t != nil && t.buffer != nil {
		t.buffer.Append(text)
	}
}

func (t *UITextView) AppendText(text string) {
	t.Append(text)
}

func (t *UITextView) SetWrapAtBounds() {
	if t != nil && t.raw != nil {
		t.raw.SetWrapMode(fltk_bridge.WRAP_AT_BOUNDS)
	}
}

func (t *UITextView) SetFontSize(size int) {
	if t != nil && t.raw != nil {
		t.raw.SetTextSize(size)
	}
}

func (t *UITextView) SetTextColor(rgb uint) {
	if t != nil && t.raw != nil {
		t.raw.SetTextColor(fltk_bridge.Color(rgb))
	}
}

func (t *UITextView) OnTextChanged(cb func()) {
	if t == nil || t.buffer == nil {
		return
	}
	t.buffer.AddModifyCallback(func(int, int, int, int, string) {
		if cb != nil {
			cb()
		}
	})
}

func (t *UITextView) On(event fltk_bridge.Event, handler func(fltk_bridge.Event) bool) {
	if t != nil {
		t.v.On(event, handler)
	}
}
