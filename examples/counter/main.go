//go:build ignore

package main

import (
	"runtime"

	"examples/counter"
	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(600, 400, "Counter")
	root := win.RootView()

	counter.BuildView(root)

	win.Show()
	fltk2go.Run()
}
