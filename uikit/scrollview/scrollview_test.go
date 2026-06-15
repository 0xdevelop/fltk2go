package scrollview_test

import (
	"testing"

	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/scrollview"
	"github.com/0xYeah/fltk2go/uikit/view"
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

func TestUIScrollView_Creation(t *testing.T) {
	sv := scrollview.NewUIScrollView(&foundation.Rect{X: 10, Y: 10, Width: 200, Height: 300})
	if sv == nil {
		t.Fatal("Expected UIScrollView to be created")
	}
	if sv.Raw() == nil {
		t.Fatal("Expected Raw() to return underlying Scroll widget")
	}
	if sv.View() == nil {
		t.Fatal("Expected View() to return UIView")
	}
}

func TestUIScrollView_ScrollType(t *testing.T) {
	sv := scrollview.NewUIScrollView(nil)
	sv.SetScrollType(scrollview.ScrollVerticalAlways)
	// It's hard to verify inner state without getters in fltk_bridge, but we can verify it doesn't crash
}

func TestUIScrollView_AddRemoveSubview(t *testing.T) {
	sv := scrollview.NewUIScrollView(nil)
	child := newMockViewable()

	sv.AddSubview(child)

	// fltk_bridge.Scroll 继承了 Group，ChildCount 应该为 1 (如果有自带滚动条可能是大于1，但至少大于0)
	count := sv.Raw().ChildCount()
	if count <= 0 {
		t.Fatalf("Expected ChildCount > 0 after AddSubview, got %d", count)
	}

	sv.RemoveSubview(child)
}
