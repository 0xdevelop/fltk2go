//go:build ignore

package main

import (
	"runtime"

	"examples/input"
	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(900, 640, "Input Example")
	root := win.RootView()

	input.BuildView(root)

	win.Show()
	fltk2go.Run()
}
