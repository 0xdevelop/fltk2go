package loginview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/textfield"
	"github.com/0xYeah/fltk2go/uikit/view"
)

// LoginView 复合组件
type LoginView struct {
	v   view.UIView
	raw *fltk_bridge.Group

	usernameInput *textfield.UITextField
	passwordInput *textfield.UISecretTextField
	loginButton   *button.UIButton

	onLoginClick func(username, password string)
}

// NewLoginView 创建一个新的登录表单
func NewLoginView(r *foundation.Rect) *LoginView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 400, Height: 300}
	}
	lv := &LoginView{}
	lv.raw = fltk_bridge.NewGroup(r.X, r.Y, r.Width, r.Height, "")
	lv.raw.End()
	lv.v.BindRaw(lv.raw)

	// 计算居中布局
	inputWidth := r.Width - 40
	if inputWidth > 300 {
		inputWidth = 300
	}
	if inputWidth < 100 {
		inputWidth = r.Width
	}
	inputHeight := 36
	spacing := 20

	startX := r.X + (r.Width-inputWidth)/2
	startY := r.Y + (r.Height-(inputHeight*3+spacing*2))/2
	if startY < r.Y {
		startY = r.Y
	}

	// 账户输入框
	lv.usernameInput = textfield.NewUITextField(startX, startY, inputWidth, inputHeight, "Username")
	lv.v.AddSubview(lv.usernameInput.View())

	// 密码输入框
	lv.passwordInput = textfield.NewUISecretTextField(startX, startY+inputHeight+spacing, inputWidth, inputHeight, "Password")
	lv.v.AddSubview(lv.passwordInput.View())

	// 登录按钮
	lv.loginButton = button.NewUIButton(&foundation.Rect{X: startX, Y: startY + (inputHeight+spacing)*2, Width: inputWidth, Height: inputHeight}, "Login")
	lv.v.AddSubview(lv.loginButton.View())

	// 绑定点击事件
	lv.loginButton.OnTouchUpInside(func() {
		if lv.onLoginClick != nil {
			lv.onLoginClick(lv.usernameInput.Text(), lv.passwordInput.Text())
		}
	})

	return lv
}

// View 返回基础视图
func (lv *LoginView) View() *view.UIView {
	return &lv.v
}

// OnLoginClick 设置登录点击事件回调
func (lv *LoginView) OnLoginClick(callback func(username, password string)) {
	lv.onLoginClick = callback
}

// Username 获取用户名
func (lv *LoginView) Username() string {
	return lv.usernameInput.Text()
}

// Password 获取密码
func (lv *LoginView) Password() string {
	return lv.passwordInput.Text()
}

// SetUsername 设置用户名
func (lv *LoginView) SetUsername(username string) {
	lv.usernameInput.SetText(username)
}

// SetPassword 设置密码
func (lv *LoginView) SetPassword(password string) {
	lv.passwordInput.SetText(password)
}
