package navigationbar

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/colors"
	"github.com/0xYeah/fltk2go/uikit/view"
)

// UIBarButtonItem represents an item on the navigation bar (e.g., a button).
type UIBarButtonItem struct {
	View view.Viewable
}

// NewUIBarButtonItem creates a standard button item.
func NewUIBarButtonItem(title string, handler func()) *UIBarButtonItem {
	btn := button.NewUIButton(nil, title)
	btn.Raw().SetBox(fltk_bridge.FLAT_BOX)
	// Use a transparent or matching background color if needed
	btn.SetBackgroundColor(colors.Background.Rgb) // Match the default nav bar background color
	btn.OnTouchUpInside(handler)
	return &UIBarButtonItem{View: btn}
}

// NewUIBarButtonItemWithCustomView creates an item with a custom view.
func NewUIBarButtonItemWithCustomView(v view.Viewable) *UIBarButtonItem {
	return &UIBarButtonItem{View: v}
}

// UINavigationItem represents the content displayed on the navigation bar.
type UINavigationItem struct {
	Title               string
	LeftBarButtonItems  []*UIBarButtonItem
	RightBarButtonItems []*UIBarButtonItem
}

// NewUINavigationItem creates a new navigation item with a title.
func NewUINavigationItem(title string) *UINavigationItem {
	return &UINavigationItem{
		Title: title,
	}
}
