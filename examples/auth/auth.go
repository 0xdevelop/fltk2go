package auth

import (
	"fmt"

	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/loginview"
	"github.com/0xdevelop/fltk2go/uikit/navigationbar"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

func BuildView(parent *view.UIView) view.Viewable {
	// 添加顶部导航栏
	navBar := navigationbar.NewUINavigationBar(&foundation.Rect{X: 0, Y: 0, Width: 400, Height: 44})
	navItem := navigationbar.NewUINavigationItem("Login")

	// 添加导航栏右侧帮助按钮
	helpBtn := navigationbar.NewUIBarButtonItem("Help", func() {
		fmt.Println("Help clicked")
	})
	navItem.RightBarButtonItems = []*navigationbar.UIBarButtonItem{helpBtn}

	navBar.SetItem(navItem)
	parent.AddSubview(navBar)

	// 添加 LoginView 到导航栏下方
	lv := loginview.NewLoginView(&foundation.Rect{X: 0, Y: 44, Width: 400, Height: 316})
	lv.OnLoginClick(func(username, password string) {
		fmt.Printf("Login Clicked! Username: %s, Password: %s\n", username, password)
	})

	parent.AddSubview(lv)

	return nil
}
