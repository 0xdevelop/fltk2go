package tabview_test

import (
	"testing"

	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/tabview"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

// mockViewable 实现一个简单的 Viewable
type mockViewable struct {
	v view.UIView
}

func newMockViewable() *mockViewable {
	m := &mockViewable{}
	raw := fltk_bridge.NewGroup(0, 0, 100, 100, "")
	m.v.BindRaw(raw)
	return m
}

func (m *mockViewable) View() *view.UIView {
	return &m.v
}

func TestUITabView_Creation(t *testing.T) {
	tv := tabview.NewUITabView(&foundation.Rect{X: 10, Y: 10, Width: 400, Height: 300})
	if tv == nil {
		t.Fatal("Expected UITabView to be created")
	}
	if tv.Raw() == nil {
		t.Fatal("Expected Raw() to return underlying Group widget")
	}
	if tv.View() == nil {
		t.Fatal("Expected View() to return UIView")
	}
}

func TestUITabView_AddTab(t *testing.T) {
	tv := tabview.NewUITabView(nil)

	child1 := newMockViewable()
	child2 := newMockViewable()

	tv.AddTab("Tab 1", child1)
	tv.AddTab("Tab 2", child2)

	// Since we added two tabs, activeIndex should be 0 by default
	// and there should be multiple children in the tab bar and content area
	if tv.Raw().ChildCount() < 2 {
		t.Fatalf("Expected UITabView to have children")
	}
}

func TestUITabView_SelectTab(t *testing.T) {
	tv := tabview.NewUITabView(nil)

	child1 := newMockViewable()
	child2 := newMockViewable()

	tv.AddTab("Tab 1", child1)
	tv.AddTab("Tab 2", child2)

	called := false
	tv.OnTabChanged(func(idx int) {
		if idx == 1 {
			called = true
		}
	})

	tv.SelectTab(1)
	if !called {
		t.Fatal("Expected OnTabChanged callback to be triggered for index 1")
	}
}
