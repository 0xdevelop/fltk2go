package splitview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/uikit/view"
)

// SplitView 分割视图组件
type SplitView struct {
	raw         *fltk_bridge.Group
	leftView    view.Viewable
	rightView   view.Viewable
	divider     *view.UIView
	v           view.UIView
	orientation Orientation

	splitPos    int // 当前分割位置（相对于容器的像素偏移）
	minSplit    int
	dividerSize int

	isDragging      bool
	lastMouse       int
	initialSplitPos int
}

// Orientation 分割方向
type Orientation int

const (
	// Horizontal 水平分割 (左右)
	Horizontal Orientation = iota
	// Vertical 垂直分割 (上下)
	Vertical
)

// New 创建一个新的分割视图
func New(x, y, width, height int, orientation Orientation) *SplitView {
	grp := fltk_bridge.NewGroup(x, y, width, height)
	grp.End()

	sv := &SplitView{
		raw:         grp,
		orientation: orientation,
		dividerSize: 6,
		minSplit:    50,
	}

	if orientation == Horizontal {
		sv.splitPos = width / 2
	} else {
		sv.splitPos = height / 2
	}

	sv.v.BindRaw(grp)

	// 设置自己为可变大小组件，避免 FLTK 自动按比例缩放子控件
	grp.Resizable(grp)

	// 创建 Divider (利用 Box 捕获事件)
	divBox := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, 0, 0, 0, 0, "")
	divBox.SetColor(fltk_bridge.Color(0xDDDDDD00))

	sv.divider = &view.UIView{}
	sv.divider.BindRaw(divBox)
	sv.divider.BindHost(grp)

	// 闭包事件流绑定：处理拖拽事件
	sv.divider.On(fltk_bridge.PUSH, sv.onPush)
	sv.divider.On(fltk_bridge.DRAG, sv.onDrag)
	sv.divider.On(fltk_bridge.RELEASE, sv.onRelease)
	sv.divider.On(fltk_bridge.ENTER, sv.onEnter)
	sv.divider.On(fltk_bridge.LEAVE, sv.onLeave)

	grp.Add(divBox)

	grp.SetResizeHandler(func() {
		sv.layout()
	})

	sv.layout()

	return sv
}

// SetLeftView 设置左/上视图
func (sv *SplitView) SetLeftView(v view.Viewable) {
	if sv == nil || sv.raw == nil {
		return
	}
	if sv.leftView != nil && sv.leftView.View() != nil && sv.leftView.View().Raw() != nil {
		sv.raw.Remove(sv.leftView.View().Raw())
	}
	sv.leftView = v
	if v != nil && v.View() != nil && v.View().Raw() != nil {
		sv.raw.Add(v.View().Raw())
		v.View().BindHost(sv.raw)
	}
	sv.layout()
}

// SetRightView 设置右/下视图
func (sv *SplitView) SetRightView(v view.Viewable) {
	if sv == nil || sv.raw == nil {
		return
	}
	if sv.rightView != nil && sv.rightView.View() != nil && sv.rightView.View().Raw() != nil {
		sv.raw.Remove(sv.rightView.View().Raw())
	}
	sv.rightView = v
	if v != nil && v.View() != nil && v.View().Raw() != nil {
		sv.raw.Add(v.View().Raw())
		v.View().BindHost(sv.raw)
	}
	sv.layout()
}

// SetLeftViewFixed 设置左侧固定大小
func (sv *SplitView) SetLeftViewFixed(size int) {
	if sv == nil {
		return
	}
	sv.splitPos = size
	sv.layout()
}

// SetRightViewFixed 设置右侧/下侧固定大小
func (sv *SplitView) SetRightViewFixed(size int) {
	if sv == nil || sv.raw == nil {
		return
	}
	if sv.orientation == Horizontal {
		sv.splitPos = sv.raw.W() - size - sv.dividerSize
	} else {
		sv.splitPos = sv.raw.H() - size - sv.dividerSize
	}
	sv.layout()
}

// SetPosition sets the split position as a ratio between 0 and 1.
func (sv *SplitView) SetPosition(pos float64) {
	if sv == nil || sv.raw == nil {
		return
	}
	if pos < 0 {
		pos = 0
	}
	if pos > 1 {
		pos = 1
	}
	if sv.orientation == Horizontal {
		sv.splitPos = int(float64(sv.raw.W()) * pos)
	} else {
		sv.splitPos = int(float64(sv.raw.H()) * pos)
	}
	sv.layout()
}

// Position returns the current split position as a ratio between 0 and 1.
func (sv *SplitView) Position() float64 {
	if sv == nil || sv.raw == nil {
		return 0
	}
	if sv.orientation == Horizontal {
		if sv.raw.W() == 0 {
			return 0
		}
		return float64(sv.splitPos) / float64(sv.raw.W())
	}
	if sv.raw.H() == 0 {
		return 0
	}
	return float64(sv.splitPos) / float64(sv.raw.H())
}

// SetResizable keeps API compatibility; this implementation is always draggable.
func (sv *SplitView) SetResizable(bool) {}

// Resizable reports whether the split view divider can be dragged.
func (sv *SplitView) Resizable() bool { return sv != nil }

// View 返回基础视图，实现 view.Viewable 接口
func (sv *SplitView) View() *view.UIView {
	if sv == nil {
		return nil
	}
	return &sv.v
}

// Raw 返回底层容器
func (sv *SplitView) Raw() fltk_bridge.Widget {
	if sv == nil {
		return nil
	}
	return sv.raw
}

// On 委托给底层的 UIView 绑定事件
func (sv *SplitView) On(event fltk_bridge.Event, handler func(fltk_bridge.Event) bool) {
	if sv == nil {
		return
	}
	sv.v.On(event, handler)
}

// AddSubview 委托给底层的 UIView
func (sv *SplitView) AddSubview(child view.Viewable) {
	if sv == nil {
		return
	}
	sv.v.AddSubview(child)
}

func (sv *SplitView) layout() {
	if sv.raw == nil {
		return
	}

	w := sv.raw.W()
	h := sv.raw.H()
	x := sv.raw.X()
	y := sv.raw.Y()

	var maxAllowed int
	if sv.orientation == Horizontal {
		maxAllowed = w - sv.minSplit - sv.dividerSize
	} else {
		maxAllowed = h - sv.minSplit - sv.dividerSize
	}

	if maxAllowed < sv.minSplit {
		maxAllowed = sv.minSplit
	}

	if sv.splitPos < sv.minSplit {
		sv.splitPos = sv.minSplit
	}
	if sv.splitPos > maxAllowed {
		sv.splitPos = maxAllowed
	}

	if sv.orientation == Horizontal {
		if sv.leftView != nil && sv.leftView.View() != nil && sv.leftView.View().Raw() != nil {
			if r, ok := sv.leftView.View().Raw().(interface{ Resize(x, y, w, h int) }); ok {
				r.Resize(x, y, sv.splitPos, h)
			}
		}
		if sv.divider != nil {
			if r, ok := sv.divider.Raw().(interface{ Resize(x, y, w, h int) }); ok {
				r.Resize(x+sv.splitPos, y, sv.dividerSize, h)
			}
		}
		if sv.rightView != nil && sv.rightView.View() != nil && sv.rightView.View().Raw() != nil {
			if r, ok := sv.rightView.View().Raw().(interface{ Resize(x, y, w, h int) }); ok {
				r.Resize(x+sv.splitPos+sv.dividerSize, y, w-sv.splitPos-sv.dividerSize, h)
			}
		}
	} else {
		if sv.leftView != nil && sv.leftView.View() != nil && sv.leftView.View().Raw() != nil {
			if r, ok := sv.leftView.View().Raw().(interface{ Resize(x, y, w, h int) }); ok {
				r.Resize(x, y, w, sv.splitPos)
			}
		}
		if sv.divider != nil {
			if r, ok := sv.divider.Raw().(interface{ Resize(x, y, w, h int) }); ok {
				r.Resize(x, y+sv.splitPos, w, sv.dividerSize)
			}
		}
		if sv.rightView != nil && sv.rightView.View() != nil && sv.rightView.View().Raw() != nil {
			if r, ok := sv.rightView.View().Raw().(interface{ Resize(x, y, w, h int) }); ok {
				r.Resize(x, y+sv.splitPos+sv.dividerSize, w, h-sv.splitPos-sv.dividerSize)
			}
		}
	}

	sv.raw.Redraw()
}

// 事件流处理

func (sv *SplitView) onPush(e fltk_bridge.Event) bool {
	if sv == nil {
		return false
	}
	sv.isDragging = true
	if sv.orientation == Horizontal {
		sv.lastMouse = fltk_bridge.EventXRoot()
	} else {
		sv.lastMouse = fltk_bridge.EventYRoot()
	}
	sv.initialSplitPos = sv.splitPos
	return true // 必须返回 true 以便接收后续 DRAG
}

func (sv *SplitView) onDrag(e fltk_bridge.Event) bool {
	if sv == nil {
		return false
	}
	if !sv.isDragging {
		return false
	}

	var currentMouse int
	if sv.orientation == Horizontal {
		currentMouse = fltk_bridge.EventXRoot()
	} else {
		currentMouse = fltk_bridge.EventYRoot()
	}

	delta := currentMouse - sv.lastMouse
	sv.splitPos = sv.initialSplitPos + delta

	sv.layout()

	// 实时触发父容器重绘，确保丝滑无闪烁
	if sv.raw.Parent() != nil {
		sv.raw.Parent().Redraw()
	}
	return true
}

func (sv *SplitView) onRelease(e fltk_bridge.Event) bool {
	if sv == nil {
		return false
	}
	sv.isDragging = false
	return true
}

func (sv *SplitView) onEnter(e fltk_bridge.Event) bool {
	if sv == nil {
		return false
	}
	// 可选：鼠标进入时改变 Divider 颜色或 Cursor
	if sv.divider != nil && sv.divider.Raw() != nil {
		if r, ok := sv.divider.Raw().(interface {
			SetColor(fltk_bridge.Color)
			Redraw()
		}); ok {
			r.SetColor(fltk_bridge.Color(0xCCCCCC00)) // slightly darker
			r.Redraw()
		}
	}
	return true
}

func (sv *SplitView) onLeave(e fltk_bridge.Event) bool {
	if sv == nil {
		return false
	}
	if !sv.isDragging && sv.divider != nil && sv.divider.Raw() != nil {
		if r, ok := sv.divider.Raw().(interface {
			SetColor(fltk_bridge.Color)
			Redraw()
		}); ok {
			r.SetColor(fltk_bridge.Color(0xDDDDDD00)) // restore
			r.Redraw()
		}
	}
	return true
}
