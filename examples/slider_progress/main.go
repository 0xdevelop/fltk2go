//go:build ignore

package main

import (
	"runtime"

	"examples/slider_progress"
	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(760, 560, "Slider & Progress Example")
	root := win.RootView()

	slider_progress.BuildView(root)

	win.Show()
	fltk2go.Run()
}
