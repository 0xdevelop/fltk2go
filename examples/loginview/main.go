//go:build ignore

package main

import (
	"runtime"

	"examples/loginview"
	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/uikit/window"
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(400, 300, "Login Example")
	root := win.RootView()

	loginview.BuildView(root)

	win.Show()
	fltk2go.Run()
}
