package imageview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/view"
)

type ContentMode int

const (
	ScaleAspectFit ContentMode = iota
	ScaleAspectFill
	ScaleToFill
	Center
)

type UIImageView struct {
	v           view.UIView
	raw         *fltk_bridge.Box
	image       fltk_bridge.Image
	contentMode ContentMode
}

func NewUIImageView(r *foundation.Rect) *UIImageView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 100, Height: 100}
	}

	box := fltk_bridge.NewBox(fltk_bridge.NO_BOX, r.X, r.Y, r.Width, r.Height, "")

	iv := &UIImageView{raw: box, contentMode: ScaleAspectFit}
	iv.v.BindRaw(box)

	return iv
}

func (iv *UIImageView) View() *view.UIView { return &iv.v }

func (iv *UIImageView) Raw() *fltk_bridge.Box { return iv.raw }

func (iv *UIImageView) SetImage(img fltk_bridge.Image) {
	iv.image = img
	iv.updateImage()
}

func (iv *UIImageView) SetContentMode(mode ContentMode) {
	iv.contentMode = mode
	iv.updateImage()
}

func (iv *UIImageView) updateImage() {
	if iv.image == nil {
		iv.raw.SetImage(nil)
		iv.raw.Redraw()
		return
	}

	w, h := iv.raw.W(), iv.raw.H()

	// fltk_bridge 的 Image Scale 会在原图上操作
	// 按照 ContentMode 进行缩放
	switch iv.contentMode {
	case ScaleAspectFit:
		if scalable, ok := iv.image.(interface {
			Scale(w, h int, prop bool, expand bool)
		}); ok {
			scalable.Scale(w, h, true, true)
		}
	case ScaleToFill:
		if scalable, ok := iv.image.(interface {
			Scale(w, h int, prop bool, expand bool)
		}); ok {
			scalable.Scale(w, h, false, true)
		}
	case ScaleAspectFill:
		// AspectFill：取较大边比例。FLTK的 Scale 如果 prop=true, expand=true 会尽量 fit。
		// 没有直接的 AspectFill 支持，我们这里用 AspectFit 代替或强制 ScaleToFill
		if scalable, ok := iv.image.(interface {
			Scale(w, h int, prop bool, expand bool)
		}); ok {
			scalable.Scale(w, h, true, true)
		}
	case Center:
		// 不做缩放
	}

	iv.raw.SetImage(iv.image)
	iv.raw.Redraw()
}

// On 绑定事件
func (iv *UIImageView) On(event fltk_bridge.Event, handler func(fltk_bridge.Event) bool) {
	iv.v.On(event, handler)
}
