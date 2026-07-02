package counter

import (
	"strconv"

	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/button"
	"github.com/0xdevelop/fltk2go/uikit/label"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

const (
	BLUE uint = 0x42A5F500
	GRAY uint = 0x75757500
)

func BuildView(parent *view.UIView) view.Viewable {
	title := label.NewUILabel(&foundation.Rect{X: 20, Y: 20, Width: 560, Height: 40}, "Clicked 0 count")
	title.View().SetAutomationID("counter.title").SetAutomationName("Counter value")
	title.SetFontSize(20)
	title.SetTextColor(GRAY)

	btn := button.NewUIButton(&foundation.Rect{X: 20, Y: 80, Width: 160, Height: 44}, "点我 +1")
	btn.View().SetAutomationID("counter.increment").SetAutomationName("Increment counter")
	btn.SetBackgroundColor(BLUE)

	count := 0
	btn.OnTouchUpInside(func() {
		count++
		title.SetText("Clicked " + strconv.Itoa(count) + " count")
	})

	parent.AddSubview(title)
	parent.AddSubview(btn)

	return nil
}
