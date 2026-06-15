//go:build ignore

package main

import (
	"runtime"

	"examples/navigationbar"
	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(800, 600, "Example")
	root := win.RootView()

	navigationbar.BuildView(root)

	win.Show()
	fltk2go.Run()
}
