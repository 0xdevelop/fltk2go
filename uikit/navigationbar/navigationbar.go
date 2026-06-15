package navigationbar

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/colors"
	"github.com/0xYeah/fltk2go/uikit/label"
	"github.com/0xYeah/fltk2go/uikit/view"
)

// UINavigationBar represents an iOS-style top navigation bar.
type UINavigationBar struct {
	v   view.UIView
	raw *fltk_bridge.Group

	titleLabel *label.UILabel
	bottomLine *fltk_bridge.Box

	item *UINavigationItem

	bgColor   uint
	lineColor uint

	activeViews []view.Viewable
}

// NewUINavigationBar creates a new navigation bar with a given frame.
func NewUINavigationBar(r *foundation.Rect) *UINavigationBar {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 320, Height: 44}
	}

	g := fltk_bridge.NewGroup(r.X, r.Y, r.Width, r.Height, "")
	g.SetBox(fltk_bridge.FLAT_BOX)

	bar := &UINavigationBar{
		raw:       g,
		bgColor:   colors.Background.Rgb,
		lineColor: colors.Gray0.Rgb,
	}
	bar.v.BindRaw(g)
	g.SetColor(fltk_bridge.Color(bar.bgColor))

	// Title label
	titleRect := &foundation.Rect{X: r.X + 80, Y: r.Y, Width: r.Width - 160, Height: r.Height - 1}
	lbl := label.NewUILabel(titleRect, "")
	lbl.SetAlignment(fltk_bridge.ALIGN_CENTER | fltk_bridge.ALIGN_INSIDE)
	lbl.SetFont(fltk_bridge.HELVETICA_BOLD)
	lbl.SetBackgroundColor(bar.bgColor)
	bar.titleLabel = lbl

	// Add bottom line
	lineRect := &foundation.Rect{X: r.X, Y: r.Y + r.Height - 1, Width: r.Width, Height: 1}
	line := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, lineRect.X, lineRect.Y, lineRect.Width, lineRect.Height, "")
	line.SetColor(fltk_bridge.Color(bar.lineColor))
	bar.bottomLine = line

	g.Add(lbl.Raw())
	g.Add(line)

	// Make the title label the resizable component so buttons stick to the sides
	g.Resizable(lbl.Raw())

	g.End()

	return bar
}

// View returns the underlying UIView.
func (b *UINavigationBar) View() *view.UIView { return &b.v }

// SetBackgroundColor sets the background color of the navigation bar.
func (b *UINavigationBar) SetBackgroundColor(color uint) {
	b.bgColor = color
	if b.raw != nil {
		b.raw.SetColor(fltk_bridge.Color(color))
		b.titleLabel.SetBackgroundColor(color)
		b.raw.Redraw()
	}
}

// SetBottomLineColor sets the color of the bottom divider line.
func (b *UINavigationBar) SetBottomLineColor(color uint) {
	b.lineColor = color
	if b.bottomLine != nil {
		b.bottomLine.SetColor(fltk_bridge.Color(color))
		b.bottomLine.Redraw()
	}
}

// SetItem sets the navigation item to be displayed (title, buttons).
func (b *UINavigationBar) SetItem(item *UINavigationItem) {
	b.item = item

	// Remove previously active views
	for _, v := range b.activeViews {
		if v != nil && v.View() != nil && v.View().Raw() != nil {
			b.raw.Remove(v.View().Raw())
		}
	}
	b.activeViews = nil

	if item == nil {
		b.titleLabel.SetText("")
		b.raw.Redraw()
		return
	}

	b.titleLabel.SetText(item.Title)

	// Layout Left Bar Button Items
	currentX := b.raw.X() + 8
	for _, lItem := range item.LeftBarButtonItems {
		if lItem.View == nil {
			continue
		}
		cv := lItem.View.View()
		if cv == nil || cv.Raw() == nil {
			continue
		}

		if sw, ok := cv.Raw().(interface {
			W() int
			H() int
			Resize(x, y, w, h int)
		}); ok {
			w := sw.W()
			h := sw.H()
			// Fit within bar height if too tall, else center vertically
			if h > b.raw.H()-2 {
				h = b.raw.H() - 8
			}
			y := b.raw.Y() + (b.raw.H()-h)/2
			sw.Resize(currentX, y, w, h)
			b.raw.Add(cv.Raw())
			b.activeViews = append(b.activeViews, lItem.View)
			currentX += w + 8
		}
	}

	// Layout Right Bar Button Items
	currentX = b.raw.X() + b.raw.W() - 8
	for i := len(item.RightBarButtonItems) - 1; i >= 0; i-- {
		rItem := item.RightBarButtonItems[i]
		if rItem.View == nil {
			continue
		}
		cv := rItem.View.View()
		if cv == nil || cv.Raw() == nil {
			continue
		}

		if sw, ok := cv.Raw().(interface {
			W() int
			H() int
			Resize(x, y, w, h int)
		}); ok {
			w := sw.W()
			h := sw.H()
			if h > b.raw.H()-2 {
				h = b.raw.H() - 8
			}
			currentX -= w
			y := b.raw.Y() + (b.raw.H()-h)/2
			sw.Resize(currentX, y, w, h)
			b.raw.Add(cv.Raw())
			b.activeViews = append(b.activeViews, rItem.View)
			currentX -= 8
		}
	}

	b.raw.Redraw()
}

// PushItem is an alias for SetItem for now.
func (b *UINavigationBar) PushItem(item *UINavigationItem) {
	b.SetItem(item)
}
