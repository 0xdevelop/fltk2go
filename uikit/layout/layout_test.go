package layout

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"testing"
)

func TestRectIntersects(t *testing.T) {
	r1 := Rect{0, 0, 100, 100}
	r2 := Rect{50, 50, 100, 100}
	if !r1.Intersects(r2) {
		t.Errorf("Expected r1 to intersect r2")
	}

	r3 := Rect{200, 200, 10, 10}
	if r1.Intersects(r3) {
		t.Errorf("Expected r1 not to intersect r3")
	}
}

func TestContainerNodeMeasureLayout(t *testing.T) {
	c := NewContainerNode()
	c.Padding = 10
	c.Spacing = 5

	t1 := NewTextNode("Hello", fltk_bridge.HELVETICA, 12, fltk_bridge.Color(0))
	t1.LineHeight = 16
	t2 := NewTextNode("World", fltk_bridge.HELVETICA, 12, fltk_bridge.Color(0))
	t2.LineHeight = 16

	c.AddChild(t1)
	c.AddChild(t2)

	w, h := c.Measure(100)
	if h <= 0 {
		t.Errorf("Expected positive height, got %d", h)
	}
	if w <= 0 {
		t.Errorf("Expected positive width, got %d", w)
	}

	c.Layout(0, 0)
	bounds := c.Bounds()
	if bounds.X != 0 || bounds.Y != 0 {
		t.Errorf("Expected container at (0,0), got %v", bounds)
	}
	if bounds.W != w || bounds.H != h {
		t.Errorf("Expected bounds size (%d,%d), got (%d,%d)", w, h, bounds.W, bounds.H)
	}
}
