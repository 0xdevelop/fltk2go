package menubar

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type MenuItem struct {
	Title    string
	Shortcut int
	Flags    int
	Callback func()
	Children []MenuItem
}

type UIMenuBar struct {
	v   view.UIView
	raw *fltk_bridge.MenuBar
}

func NewUIMenuBar(r *foundation.Rect) *UIMenuBar {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 800, Height: 30}
	}
	mb := fltk_bridge.NewMenuBar(r.X, r.Y, r.Width, r.Height)

	m := &UIMenuBar{raw: mb}
	m.v.BindRaw(mb)

	return m
}

func (m *UIMenuBar) View() *view.UIView        { return &m.v }
func (m *UIMenuBar) Raw() *fltk_bridge.MenuBar { return m.raw }

func (m *UIMenuBar) SetMenu(items []MenuItem) {
	m.raw.Clear()
	m.buildMenu(items, "")
}

func (m *UIMenuBar) buildMenu(items []MenuItem, prefix string) {
	for _, item := range items {
		path := prefix + item.Title
		if len(item.Children) > 0 {
			// Submenu
			flags := item.Flags | fltk_bridge.SUBMENU
			m.raw.AddEx(path, item.Shortcut, item.Callback, flags)
			m.buildMenu(item.Children, path+"/")
		} else {
			m.raw.AddEx(path, item.Shortcut, item.Callback, item.Flags)
		}
	}
}

type UIContextMenu struct {
	v   view.UIView
	raw *fltk_bridge.MenuButton
}

func NewUIContextMenu(r *foundation.Rect) *UIContextMenu {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 0, Height: 0}
	}
	mb := fltk_bridge.NewMenuButton(r.X, r.Y, r.Width, r.Height)
	mb.SetType(fltk_bridge.POPUP3) // 右键弹出

	m := &UIContextMenu{raw: mb}
	m.v.BindRaw(mb)

	return m
}

func (m *UIContextMenu) View() *view.UIView           { return &m.v }
func (m *UIContextMenu) Raw() *fltk_bridge.MenuButton { return m.raw }

func (m *UIContextMenu) SetMenu(items []MenuItem) {
	m.raw.Clear()
	m.buildMenu(items, "")
}

func (m *UIContextMenu) buildMenu(items []MenuItem, prefix string) {
	for _, item := range items {
		path := prefix + item.Title
		if len(item.Children) > 0 {
			flags := item.Flags | fltk_bridge.SUBMENU
			m.raw.AddEx(path, item.Shortcut, item.Callback, flags)
			m.buildMenu(item.Children, path+"/")
		} else {
			m.raw.AddEx(path, item.Shortcut, item.Callback, item.Flags)
		}
	}
}

func (m *UIContextMenu) Popup() {
	m.raw.Popup()
}
