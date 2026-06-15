package splitview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/uikit/view"
)

// SplitView 分割视图组件
type SplitView struct {
	// 底层FLTK Flex容器
	raw *fltk_bridge.Flex

	// 左/上视图
	leftView view.Viewable

	// 右/下视图
	rightView view.Viewable

	// 基础视图
	v view.UIView
}

// Orientation 分割方向
type Orientation int

const (
	// Horizontal 水平分割
	Horizontal Orientation = iota
	// Vertical 垂直分割
	Vertical
)

// New 创建一个新的分割视图
func New(x, y, width, height int, orientation Orientation) *SplitView {
	flex := fltk_bridge.NewFlex(x, y, width, height)

	// 设置分割方向
	if orientation == Horizontal {
		flex.SetType(fltk_bridge.ROW)
	} else {
		flex.SetType(fltk_bridge.COLUMN)
	}

	// 设置间距
	flex.SetGap(5)

	sv := &SplitView{
		raw: flex,
	}

	// 绑定底层widget到view
	sv.v.BindRaw(flex)
	sv.v.BindHost(flex)

	return sv
}

// SetLeftView 设置左/上视图
func (sv *SplitView) SetLeftView(v view.Viewable) {
	if sv == nil {
		return
	}
	// 清除所有子视图
	sv.clearViews()

	// 保存左视图
	sv.leftView = v

	// 重新添加所有视图
	sv.addViews()
}

// SetRightView 设置右/下视图
func (sv *SplitView) SetRightView(v view.Viewable) {
	if sv == nil {
		return
	}
	// 清除所有子视图
	sv.clearViews()

	// 保存右视图
	sv.rightView = v

	// 重新添加所有视图
	sv.addViews()
}

// SetPosition 设置分割位置（0.0-1.0之间的值）
// 注意：由于使用Flex实现，此方法仅在设置固定大小时有效
func (sv *SplitView) SetPosition(pos float64) {
	// 这里可以根据需要实现固定大小的逻辑
}

// Position 获取分割位置
func (sv *SplitView) Position() float64 {
	// 这里可以根据需要实现获取位置的逻辑
	return 0.5
}

// SetResizable 设置是否可调整大小
func (sv *SplitView) SetResizable(resizable bool) {
	// Flex布局默认是可调整大小的
}

// Resizable 获取是否可调整大小
func (sv *SplitView) Resizable() bool {
	return true
}

// SetLeftViewFixed 设置左视图为固定大小
func (sv *SplitView) SetLeftViewFixed(size int) {
	if sv == nil || sv.raw == nil || sv.leftView == nil || sv.leftView.View() == nil || sv.leftView.View().Raw() == nil {
		return
	}
	if sv.leftView != nil {
		sv.raw.Fixed(sv.leftView.View().Raw(), size)
		sv.raw.End()
	}
}

// SetRightViewFixed 设置右视图为固定大小
func (sv *SplitView) SetRightViewFixed(size int) {
	if sv == nil || sv.raw == nil || sv.rightView == nil || sv.rightView.View() == nil || sv.rightView.View().Raw() == nil {
		return
	}
	if sv.rightView != nil {
		sv.raw.Fixed(sv.rightView.View().Raw(), size)
		sv.raw.End()
	}
}

// View 返回基础视图，实现view.Viewable接口
func (sv *SplitView) View() *view.UIView {
	if sv == nil {
		return nil
	}
	return &sv.v
}

// Raw 返回底层FLTK Flex容器
func (sv *SplitView) Raw() *fltk_bridge.Flex {
	if sv == nil {
		return nil
	}
	return sv.raw
}

// 清除所有子视图
func (sv *SplitView) clearViews() {
	if sv == nil {
		return
	}
	// 由于FLTK的Flex容器不支持直接移除子视图
	// 我们需要重新创建Flex容器
	if sv.raw != nil {
		// 保存当前位置和大小
		x := sv.raw.X()
		y := sv.raw.Y()
		w := sv.raw.W()
		h := sv.raw.H()

		// 保存当前方向
		var orientation Orientation
		if sv.raw.Type() == uint8(fltk_bridge.ROW) {
			orientation = Horizontal
		} else {
			orientation = Vertical
		}

		// 创建新的Flex容器
		sv.raw = fltk_bridge.NewFlex(x, y, w, h)

		// 设置方向
		if orientation == Horizontal {
			sv.raw.SetType(fltk_bridge.ROW)
		} else {
			sv.raw.SetType(fltk_bridge.COLUMN)
		}

		// 设置间距
		sv.raw.SetGap(5)

		// 重新绑定底层widget到view
		sv.v.BindRaw(sv.raw)
		sv.v.BindHost(sv.raw)
	}
}

// 添加所有子视图
func (sv *SplitView) addViews() {
	if sv == nil {
		return
	}
	if sv.raw != nil {
		// 添加左视图
		if sv.leftView != nil && sv.leftView.View() != nil && sv.leftView.View().Raw() != nil {
			sv.raw.Add(sv.leftView.View().Raw())
			sv.leftView.View().BindHost(sv.raw)
		}

		// 添加右视图
		if sv.rightView != nil && sv.rightView.View() != nil && sv.rightView.View().Raw() != nil {
			sv.raw.Add(sv.rightView.View().Raw())
			sv.rightView.View().BindHost(sv.raw)
		}

		// 结束Flex布局
		sv.raw.End()
	}
}
