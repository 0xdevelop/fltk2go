//go:build ignore

package main

import (
	"runtime"

	"examples/tableview_demo"
	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(800, 600, "表格视图演示")
	root := win.RootView()

	tableview_demo.BuildView(root)

	win.Show()
	fltk2go.Run()
}
