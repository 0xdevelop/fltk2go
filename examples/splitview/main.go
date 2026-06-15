//go:build ignore

package main

import (
	"runtime"

	"examples/splitview"
	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(1000, 700, "Split View Example")
	root := win.RootView()

	splitview.BuildView(root)

	win.Show()
	fltk2go.Run()
}
