//go:build ignore

package main

import (
	"runtime"

	"examples/tableview"
	"github.com/0xdevelop/fltk2go"
	"github.com/0xdevelop/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(800, 600, "Server Management")
	root := win.RootView()

	tableview.BuildView(root)

	win.Show()
	fltk2go.Run()
}
