//go:build ignore

package main

import (
	"runtime"

	"examples/auth"
	"github.com/0xdevelop/fltk2go"
	"github.com/0xdevelop/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(400, 360, "Auth Example")
	root := win.RootView()

	auth.BuildView(root)

	win.Show()
	fltk2go.Run()
}
