package group

import (
	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

// UIGroup is a lightweight UIKit-style container backed by Fl_Group.
// It gives applications an explicit panel/container object that can host
// child views via AddSubview while still allowing absolute FLTK layouts.
type UIGroup struct {
	v   view.UIView
	raw *fltk_bridge.Group
}

func NewUIGroup(r *foundation.Rect) *UIGroup {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 100, Height: 100}
	}

	raw := fltk_bridge.NewGroup(r.X, r.Y, r.Width, r.Height)
	raw.End()

	g := &UIGroup{raw: raw}
	g.v.BindRaw(raw)
	g.v.SetAutomationRole("group")
	return g
}

func (g *UIGroup) View() *view.UIView {
	if g == nil {
		return nil
	}
	return &g.v
}

func (g *UIGroup) Raw() *fltk_bridge.Group {
	if g == nil {
		return nil
	}
	return g.raw
}

func (g *UIGroup) AddSubview(child view.Viewable) {
	if g == nil {
		return
	}
	g.v.AddSubview(child)
}

func (g *UIGroup) SetBackgroundColor(rgb uint) {
	if g != nil && g.raw != nil {
		g.raw.SetColor(fltk_bridge.Color(rgb))
		g.raw.SetBox(fltk_bridge.FLAT_BOX)
		g.raw.Redraw()
	}
}

func (g *UIGroup) SetAutomationID(id string) *UIGroup {
	if g != nil {
		g.v.SetAutomationID(id)
	}
	return g
}

func (g *UIGroup) SetAutomationName(name string) *UIGroup {
	if g != nil {
		g.v.SetAutomationName(name)
	}
	return g
}
