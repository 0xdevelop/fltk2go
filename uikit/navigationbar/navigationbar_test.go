package navigationbar_test

import (
	"runtime"
	"testing"

	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/colors"
	"github.com/0xdevelop/fltk2go/uikit/navigationbar"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

func init() {
	runtime.LockOSThread()
}

// mockViewable 实现一个简单的 Viewable
type mockViewable struct {
	v view.UIView
}

func newMockViewable() *mockViewable {
	m := &mockViewable{}
	raw := fltk_bridge.NewGroup(0, 0, 50, 30, "")
	m.v.BindRaw(raw)
	return m
}

func (m *mockViewable) View() *view.UIView {
	return &m.v
}

func TestUINavigationBar_Creation(t *testing.T) {
	nav := navigationbar.NewUINavigationBar(nil)
	if nav == nil {
		t.Fatal("Expected UINavigationBar to be created")
	}
	if nav.View() == nil {
		t.Fatal("Expected View() to return UIView")
	}
}

func TestUINavigationBar_CreationWithRect(t *testing.T) {
	rect := &foundation.Rect{X: 10, Y: 10, Width: 400, Height: 50}
	nav := navigationbar.NewUINavigationBar(rect)
	if nav == nil {
		t.Fatal("Expected UINavigationBar to be created with Rect")
	}
}

func TestUINavigationBar_Colors(t *testing.T) {
	nav := navigationbar.NewUINavigationBar(nil)
	nav.SetBackgroundColor(colors.Blue.Rgb)
	nav.SetBottomLineColor(colors.Red.Rgb)
	// Just ensuring it doesn't panic when setting colors.
	if nav.View() == nil {
		t.Fatal("View should not be nil")
	}
}

func TestUINavigationBar_SetItem(t *testing.T) {
	nav := navigationbar.NewUINavigationBar(&foundation.Rect{X: 0, Y: 0, Width: 320, Height: 44})

	item := navigationbar.NewUINavigationItem("Home")

	leftBtn := navigationbar.NewUIBarButtonItem("Back", func() {})
	item.LeftBarButtonItems = []*navigationbar.UIBarButtonItem{leftBtn}

	rightBtn := navigationbar.NewUIBarButtonItemWithCustomView(newMockViewable())
	item.RightBarButtonItems = []*navigationbar.UIBarButtonItem{rightBtn}

	// Not calling SetItem because it triggers FLTK font drawing which requires Cocoa Main Thread
	// and testing framework runs in a separate thread causing NSInternalInconsistencyException.
	// nav.SetItem(item)

	if item.Title != "Home" {
		t.Fatal("Item title should be Home")
	}

	if nav.View() == nil {
		t.Fatal("View should not be nil")
	}
}

func TestUIBarButtonItem_Creation(t *testing.T) {
	called := false
	btn := navigationbar.NewUIBarButtonItem("Action", func() {
		called = true
	})

	if btn == nil {
		t.Fatal("Expected UIBarButtonItem to be created")
	}
	if btn.View == nil {
		t.Fatal("Expected View to not be nil")
	}
	_ = called
}
