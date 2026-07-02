package loginview

import (
	"fmt"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/loginview"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

func BuildView(parent *view.UIView) view.Viewable {

	lv := loginview.NewLoginView(&foundation.Rect{X: 0, Y: 0, Width: 400, Height: 300})
	lv.OnLoginClick(func(username, password string) {
		fmt.Printf("Login Clicked! Username: %s, Password: %s\n", username, password)
	})

	parent.AddSubview(lv)

	return nil
}
