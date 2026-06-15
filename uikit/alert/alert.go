package alert

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/label"
	"github.com/0xYeah/fltk2go/uikit/view"
	"github.com/0xYeah/fltk2go/uikit/window"
)

type UIAlert struct {
	win       *window.UIWindow
	content   view.Viewable
	onConfirm func()
	onCancel  func()
}

// NewUIAlert 创建一个标准的模态对话框
func NewUIAlert(title, message string) *UIAlert {
	win := window.NewUIWindow(300, 150, title)
	win.Raw().SetModal()

	alert := &UIAlert{
		win: win,
	}

	// 默认的消息内容
	msgLabel := label.NewUILabel(&foundation.Rect{X: 20, Y: 20, Width: 260, Height: 60}, message)
	msgLabel.Raw().SetAlign(fltk_bridge.ALIGN_CENTER | fltk_bridge.ALIGN_INSIDE | fltk_bridge.ALIGN_WRAP)
	win.RootView().AddSubview(msgLabel)

	// 确认按钮
	confirmBtn := button.NewUIButton(&foundation.Rect{X: 160, Y: 100, Width: 80, Height: 30}, "OK")
	confirmBtn.OnTouchUpInside(func() {
		if alert.onConfirm != nil {
			alert.onConfirm()
		}
		alert.Hide()
	})
	win.RootView().AddSubview(confirmBtn)

	// 取消按钮
	cancelBtn := button.NewUIButton(&foundation.Rect{X: 60, Y: 100, Width: 80, Height: 30}, "Cancel")
	cancelBtn.OnTouchUpInside(func() {
		if alert.onCancel != nil {
			alert.onCancel()
		}
		alert.Hide()
	})
	win.RootView().AddSubview(cancelBtn)

	return alert
}

// NewUIModal 创建一个嵌入自定义 UIView 的模态对话框
func NewUIModal(title string, width, height int, customView view.Viewable) *UIAlert {
	win := window.NewUIWindow(width, height, title)
	win.Raw().SetModal()

	alert := &UIAlert{
		win:     win,
		content: customView,
	}

	if customView != nil {
		win.RootView().AddSubview(customView)
	}

	return alert
}

func (a *UIAlert) OnConfirm(cb func()) {
	a.onConfirm = cb
}

func (a *UIAlert) OnCancel(cb func()) {
	a.onCancel = cb
}

func (a *UIAlert) Show() {
	if a.win != nil {
		a.win.Show()
	}
}

func (a *UIAlert) Hide() {
	if a.win != nil && a.win.Raw() != nil {
		a.win.Raw().Hide()
	}
}
