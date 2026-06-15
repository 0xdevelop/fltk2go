package layout

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/uikit/textlayout"
	"github.com/0xYeah/fltk2go/uikit/view"
)

// Rect represents a bounding box in 2D space.
type Rect struct {
	X, Y, W, H int
}

// Intersects checks if two rectangles overlap.
func (r Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.W && r.X+r.W > other.X &&
		r.Y < other.Y+other.H && r.Y+r.H > other.Y
}

// Node is an element in the Virtual Layout Tree.
// It decouples Measure, Layout, and Draw phases for high performance.
type Node interface {
	Measure(maxWidth int) (int, int)
	Layout(x, y int)
	Draw(viewport Rect)
	Bounds() Rect
}

// TextNode is a leaf node that renders text using cached text segments.
type TextNode struct {
	Text       string
	Font       fltk_bridge.Font
	FontSize   int
	Color      fltk_bridge.Color
	Align      fltk_bridge.Align
	LineHeight int

	preparedText *textlayout.PreparedText
	layoutResult *textlayout.LayoutResult
	rect         Rect
}

func NewTextNode(text string, font fltk_bridge.Font, fontSize int, color fltk_bridge.Color) *TextNode {
	return &TextNode{
		Text:     text,
		Font:     font,
		FontSize: fontSize,
		Color:    color,
		Align:    fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_TOP,
	}
}

func (n *TextNode) Measure(maxWidth int) (int, int) {
	if n.preparedText == nil {
		n.preparedText = textlayout.Prepare(n.Text, n.Font, n.FontSize)
	}
	lineHeight := n.LineHeight
	if lineHeight == 0 {
		lineHeight = n.FontSize + 4
	}
	n.layoutResult = textlayout.Layout(n.preparedText, maxWidth, lineHeight)

	maxW := 0
	for _, l := range n.layoutResult.Lines {
		if l.Width > maxW {
			maxW = l.Width
		}
	}

	return maxW, n.layoutResult.Height
}

func (n *TextNode) Layout(x, y int) {
	w, h := 0, 0
	if n.layoutResult != nil {
		h = n.layoutResult.Height
		for _, l := range n.layoutResult.Lines {
			if l.Width > w {
				w = l.Width
			}
		}
	}
	n.rect = Rect{X: x, Y: y, W: w, H: h}
}

func (n *TextNode) Draw(viewport Rect) {
	if !n.rect.Intersects(viewport) {
		return // Skip rendering if outside viewport (Dirty Rect concept)
	}
	if n.layoutResult == nil {
		return
	}

	fltk_bridge.SetDrawColor(n.Color)
	fltk_bridge.SetDrawFont(n.Font, n.FontSize)

	lineHeight := n.LineHeight
	if lineHeight == 0 {
		lineHeight = n.FontSize + 4
	}

	for i, line := range n.layoutResult.Lines {
		lineY := n.rect.Y + i*lineHeight
		// Render only visible lines within the dirty rect
		if lineY+lineHeight >= viewport.Y && lineY <= viewport.Y+viewport.H {
			fltk_bridge.Draw(line.Text, n.rect.X, lineY, n.rect.W, lineHeight, n.Align)
		}
	}
}

func (n *TextNode) Bounds() Rect {
	return n.rect
}

// ViewNode is a leaf node that wraps an existing Viewable component.
// It manages the position and size of the underlying FLTK widget during layout.
type ViewNode struct {
	View view.Viewable

	// Optional fixed size constraints. If 0, it uses the widget's current size.
	FixedWidth  int
	FixedHeight int

	rect Rect
}

func NewViewNode(v view.Viewable, w, h int) *ViewNode {
	return &ViewNode{
		View:        v,
		FixedWidth:  w,
		FixedHeight: h,
	}
}

func (n *ViewNode) Measure(maxWidth int) (int, int) {
	w, h := n.FixedWidth, n.FixedHeight
	if n.View != nil && n.View.View() != nil && n.View.View().Raw() != nil {
		if getter, ok := n.View.View().Raw().(interface {
			W() int
			H() int
		}); ok {
			if w <= 0 {
				w = getter.W()
			}
			if h <= 0 {
				h = getter.H()
			}
		}
	}
	return w, h
}

func (n *ViewNode) Layout(x, y int) {
	w, h := n.Measure(0) // Unconstrained measure to get the preferred size
	n.rect = Rect{X: x, Y: y, W: w, H: h}

	if n.View != nil && n.View.View() != nil && n.View.View().Raw() != nil {
		if resizable, ok := n.View.View().Raw().(interface{ Resize(x, y, w, h int) }); ok {
			resizable.Resize(x, y, w, h)
		}
	}
}

func (n *ViewNode) Draw(viewport Rect) {
	// For ViewNode, FLTK handles drawing the underlying widget natively.
	// We don't need to manually draw it here, but we could toggle its visibility
	// based on the dirty rect for extreme optimization (optional).
	if n.View != nil && n.View.View() != nil && n.View.View().Raw() != nil {
		if !n.rect.Intersects(viewport) {
			// n.View.View().Raw().Hide() // Optional
		} else {
			// n.View.View().Raw().Show() // Optional
		}
	}
}

func (n *ViewNode) Bounds() Rect {
	return n.rect
}

// ContainerNode holds multiple child nodes and manages their layout.
type ContainerNode struct {
	Children []Node
	BgColor  fltk_bridge.Color
	Padding  int
	Spacing  int

	rect Rect

	IsHorizontal bool
}

func NewContainerNode() *ContainerNode {
	return &ContainerNode{
		BgColor: fltk_bridge.Color(0),
	}
}

func (c *ContainerNode) AddChild(node Node) {
	c.Children = append(c.Children, node)
}

func (c *ContainerNode) Measure(maxWidth int) (int, int) {
	w := 0
	h := c.Padding * 2

	availWidth := maxWidth - c.Padding*2
	if availWidth < 0 {
		availWidth = 0
	}

	if c.IsHorizontal {
		w = c.Padding * 2
		maxH := 0
		for i, child := range c.Children {
			cw, ch := child.Measure(availWidth)
			w += cw
			if i > 0 {
				w += c.Spacing
			}
			if ch > maxH {
				maxH = ch
			}
		}
		h += maxH
	} else {
		maxW := 0
		for i, child := range c.Children {
			cw, ch := child.Measure(availWidth)
			h += ch
			if i > 0 {
				h += c.Spacing
			}
			if cw > maxW {
				maxW = cw
			}
		}
		w = c.Padding*2 + maxW
	}

	return w, h
}

func (c *ContainerNode) Layout(x, y int) {
	c.rect.X = x
	c.rect.Y = y

	cx := x + c.Padding
	cy := y + c.Padding

	maxW, maxH := 0, 0

	for _, child := range c.Children {
		child.Layout(cx, cy)
		cb := child.Bounds()
		if c.IsHorizontal {
			cx += cb.W + c.Spacing
			if cb.H > maxH {
				maxH = cb.H
			}
		} else {
			cy += cb.H + c.Spacing
			if cb.W > maxW {
				maxW = cb.W
			}
		}
	}

	if c.IsHorizontal {
		c.rect.W = cx - x + c.Padding
		if len(c.Children) > 0 {
			c.rect.W -= c.Spacing
		}
		c.rect.H = maxH + c.Padding*2
	} else {
		c.rect.H = cy - y + c.Padding
		if len(c.Children) > 0 {
			c.rect.H -= c.Spacing
		}
		c.rect.W = maxW + c.Padding*2
	}
}

func (c *ContainerNode) Draw(viewport Rect) {
	if !c.rect.Intersects(viewport) {
		return
	}

	if c.BgColor != 0 {
		fltk_bridge.DrawBox(fltk_bridge.FLAT_BOX, c.rect.X, c.rect.Y, c.rect.W, c.rect.H, c.BgColor)
	}

	for _, child := range c.Children {
		child.Draw(viewport)
	}
}

func (c *ContainerNode) Bounds() Rect {
	return c.rect
}

// Engine acts as the top-level renderer for a Virtual Layout Tree.
type Engine struct {
	Root Node
}

func NewEngine(root Node) *Engine {
	return &Engine{
		Root: root,
	}
}

// Render recalculates bounds if needed and draws only the nodes within the viewport.
// Uses FLTK's PushClip and PopClip for hardware-level dirty rect culling.
func (e *Engine) Render(viewport Rect) {
	if e.Root == nil {
		return
	}
	fltk_bridge.PushClip(viewport.X, viewport.Y, viewport.W, viewport.H)
	defer fltk_bridge.PopClip()

	e.Root.Draw(viewport)
}

// LayoutView wraps an Engine into a standard Viewable component,
// so it can receive events via On() and be added via AddSubview().
type LayoutView struct {
	raw    *fltk_bridge.Box
	v      view.UIView
	engine *Engine
}

// NewLayoutView creates a LayoutView.
func NewLayoutView(x, y, w, h int, engine *Engine) *LayoutView {
	// Create a transparent box to act as the host widget for the layout.
	// We use NO_BOX so it doesn't draw its own background, we'll draw via the engine.
	box := fltk_bridge.NewBox(fltk_bridge.NO_BOX, x, y, w, h, "")

	lv := &LayoutView{
		raw:    box,
		engine: engine,
	}

	lv.v.BindRaw(box)

	box.SetDrawHandler(func(defaultDraw func()) {
		// Custom draw: render the layout engine
		if lv.engine != nil {
			// Update layout bounds to match the widget
			bx, by, bw, bh := box.X(), box.Y(), box.W(), box.H()

			// If the root is a node, measure and layout
			if lv.engine.Root != nil {
				// We don't re-measure and re-layout every frame in a real app unless dirty,
				// but for simplicity here we layout if size changes.
				// In a full implementation we'd check if bounds changed.
				lv.engine.Root.Measure(bw)
				lv.engine.Root.Layout(bx, by)
			}

			lv.engine.Render(Rect{X: bx, Y: by, W: bw, H: bh})
		}
	})

	return lv
}

func (lv *LayoutView) View() *view.UIView {
	return &lv.v
}

func (lv *LayoutView) Raw() fltk_bridge.Widget {
	return lv.raw
}

// On delegates event binding to the underlying UIView
func (lv *LayoutView) On(event fltk_bridge.Event, handler func(fltk_bridge.Event) bool) {
	lv.v.On(event, handler)
}
