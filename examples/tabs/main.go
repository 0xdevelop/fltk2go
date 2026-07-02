//go:build ignore

package main

import (
	"runtime"

	"examples/tabs"
	"github.com/0xdevelop/fltk2go"
	"github.com/0xdevelop/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(800, 520, "Tabs Example")
	root := win.RootView()

	tabs.BuildView(root)

	win.Show()
	fltk2go.Run()
}
