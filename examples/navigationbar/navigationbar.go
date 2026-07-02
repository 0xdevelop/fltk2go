package navigationbar

import (
	"fmt"
	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/button"
	"github.com/0xdevelop/fltk2go/uikit/navigationbar"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

func BuildView(parent *view.UIView) view.Viewable {

	win := fltk_bridge.NewWindow(800, 600, "NavigationBar Example (Desktop Style)")

	// Create a navigation bar spanning the full width
	navBar := navigationbar.NewUINavigationBar(&foundation.Rect{X: 0, Y: 0, Width: 800, Height: 44})

	// Create navigation items
	item1 := navigationbar.NewUINavigationItem("Dashboard")

	backBtn := navigationbar.NewUIBarButtonItem("Back", func() {
		fmt.Println("Back clicked")
	})
	item1.LeftBarButtonItems = []*navigationbar.UIBarButtonItem{backBtn}

	addBtn := navigationbar.NewUIBarButtonItem("Settings", func() {
		fmt.Println("Settings clicked")

		// Push a new item when Settings is clicked to simulate navigation
		newItem := navigationbar.NewUINavigationItem("System Preferences")

		backBtn2 := navigationbar.NewUIBarButtonItem("Done", func() {
			navBar.SetItem(item1)
		})
		newItem.LeftBarButtonItems = []*navigationbar.UIBarButtonItem{backBtn2}

		saveBtn := navigationbar.NewUIBarButtonItem("Apply", func() {
			fmt.Println("Applied")
			navBar.SetItem(item1)
		})
		newItem.RightBarButtonItems = []*navigationbar.UIBarButtonItem{saveBtn}

		navBar.SetItem(newItem)
	})
	item1.RightBarButtonItems = []*navigationbar.UIBarButtonItem{addBtn}

	navBar.SetItem(item1)

	// Make the window resizable, leaving the navigation bar attached to the top
	contentGroup := fltk_bridge.NewGroup(0, 44, 800, 600-44, "")
	contentBtn := button.NewUIButton(&foundation.Rect{X: 300, Y: 250, Width: 200, Height: 44}, "Change Background")
	contentBtn.OnTouchUpInside(func() {
		navBar.SetBackgroundColor(0xFFE0B200) // Light orange
	})
	contentGroup.Add(contentBtn.Raw())
	contentGroup.End()

	win.Add(contentGroup)
	win.Add(navBar.View().Raw())

	// Make the contentGroup resizable, so the nav bar stays at the top
	win.Resizable(contentGroup)

	win.End()

	return nil
}
