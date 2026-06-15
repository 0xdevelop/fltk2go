package tabview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/view"
)

// UITabView 封装基于现代 UI 理念的分段控制器/多标签视图
type UITabView struct {
	v           view.UIView
	raw         *fltk_bridge.Group
	tabBar      *fltk_bridge.Group
	contentArea *fltk_bridge.Group
	highlight   *fltk_bridge.Box

	tabs        []*tabItem
	activeIndex int

	onTabChanged func(index int)
}

type tabItem struct {
	btn     *button.UIButton
	content view.Viewable
}

// NewUITabView 创建一个新的分段标签容器
func NewUITabView(r *foundation.Rect) *UITabView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 400, Height: 300}
	}

	raw := fltk_bridge.NewGroup(r.X, r.Y, r.Width, r.Height, "")

	// 顶部 TabBar 区域 (高度固定 40)
	tabBarHeight := 40
	tabBar := fltk_bridge.NewGroup(r.X, r.Y, r.Width, tabBarHeight, "")
	tabBar.SetBox(fltk_bridge.FLAT_BOX)
	tabBar.SetColor(fltk_bridge.ColorFromRgb(245, 245, 245)) // 浅灰色背景

	// 底部高亮指示器
	highlight := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, r.X, r.Y+tabBarHeight-3, 0, 3, "")
	highlight.SetColor(fltk_bridge.ColorFromRgb(0, 122, 255)) // iOS 蓝
	tabBar.Add(highlight)

	tabBar.End()

	// 内容区域
	contentArea := fltk_bridge.NewGroup(r.X, r.Y+tabBarHeight, r.Width, r.Height-tabBarHeight, "")
	contentArea.SetBox(fltk_bridge.FLAT_BOX)
	contentArea.SetColor(fltk_bridge.ColorFromRgb(255, 255, 255))
	contentArea.End()

	raw.End()

	tv := &UITabView{
		raw:         raw,
		tabBar:      tabBar,
		contentArea: contentArea,
		highlight:   highlight,
		tabs:        make([]*tabItem, 0),
		activeIndex: -1,
	}
	tv.v.BindRaw(raw)
	return tv
}

// View 实现 Viewable 接口
func (tv *UITabView) View() *view.UIView {
	return &tv.v
}

// Raw 返回底层容器
func (tv *UITabView) Raw() *fltk_bridge.Group {
	return tv.raw
}

// AddTab 添加一个新标签页
func (tv *UITabView) AddTab(title string, content view.Viewable) {
	idx := len(tv.tabs)

	// 创建标签按钮
	btn := button.NewUIButton(&foundation.Rect{X: 0, Y: 0, Width: 100, Height: 40}, title)
	btn.Raw().SetBox(fltk_bridge.FLAT_BOX)
	btn.SetBackgroundColor(uint(fltk_bridge.ColorFromRgb(245, 245, 245)))
	btn.SetTitleColor(uint(fltk_bridge.ColorFromRgb(51, 51, 51)))

	btn.OnTouchUpInside(func() {
		tv.SelectTab(idx)
	})

	tv.tabBar.Add(btn.Raw())

	// 添加内容视图
	if content != nil && content.View() != nil {
		contentWidget := content.View().Raw()
		if contentWidget != nil {
			if rw, ok := contentWidget.(interface {
				Resize(x, y, w, h int)
				Show()
				Hide()
			}); ok {
				rw.Resize(tv.contentArea.X(), tv.contentArea.Y(), tv.contentArea.W(), tv.contentArea.H())
				rw.Hide()
			}
			tv.contentArea.Add(contentWidget)
		}
	}

	tv.tabs = append(tv.tabs, &tabItem{btn: btn, content: content})

	// 重新布局
	tv.relayoutTabs()

	// 默认选中第一个
	if tv.activeIndex == -1 {
		tv.SelectTab(0)
	}
}

// relayoutTabs 重新排列顶部 Tab 按钮
func (tv *UITabView) relayoutTabs() {
	count := len(tv.tabs)
	if count == 0 {
		return
	}
	width := tv.tabBar.W() / count

	for i, item := range tv.tabs {
		item.btn.Raw().Resize(tv.tabBar.X()+i*width, tv.tabBar.Y(), width, tv.tabBar.H()-3) // 留出 3px 给底部高亮线
	}

	// 保证高亮线始终在最上层
	tv.tabBar.Remove(tv.highlight)
	tv.tabBar.Add(tv.highlight)

	// 立即更新高亮线位置
	if tv.activeIndex >= 0 && tv.activeIndex < count {
		tv.highlight.Resize(tv.tabBar.X()+tv.activeIndex*width, tv.tabBar.Y()+tv.tabBar.H()-3, width, 3)
	}
}

// SelectTab 选中指定的标签页
func (tv *UITabView) SelectTab(index int) {
	if index < 0 || index >= len(tv.tabs) || index == tv.activeIndex {
		return
	}

	// 恢复旧的样式，隐藏内容
	if tv.activeIndex >= 0 {
		oldItem := tv.tabs[tv.activeIndex]
		oldItem.btn.SetTitleColor(uint(fltk_bridge.ColorFromRgb(51, 51, 51)))
		if oldItem.content != nil && oldItem.content.View() != nil {
			if rw, ok := oldItem.content.View().Raw().(interface{ Hide() }); ok {
				rw.Hide()
			}
		}
	}

	// 应用新样式，显示内容
	tv.activeIndex = index
	newItem := tv.tabs[index]
	newItem.btn.SetTitleColor(uint(fltk_bridge.ColorFromRgb(0, 122, 255)))
	if newItem.content != nil && newItem.content.View() != nil {
		if rw, ok := newItem.content.View().Raw().(interface{ Show() }); ok {
			rw.Show()
		}
	}

	// 触发高亮移动动画
	tv.animateHighlight(index)

	tv.raw.Redraw()

	// 触发回调
	if tv.onTabChanged != nil {
		tv.onTabChanged(index)
	}
}

// animateHighlight 简易平滑移动高亮线条
func (tv *UITabView) animateHighlight(targetIndex int) {
	count := len(tv.tabs)
	if count == 0 {
		return
	}
	targetWidth := tv.tabBar.W() / count
	targetX := tv.tabBar.X() + targetIndex*targetWidth
	startX := tv.highlight.X()

	if startX == targetX {
		tv.highlight.Resize(startX, tv.highlight.Y(), targetWidth, 3)
		return
	}

	tv.highlight.Resize(startX, tv.highlight.Y(), targetWidth, 3)

	steps := 10
	stepX := (targetX - startX) / steps
	if stepX == 0 {
		tv.highlight.Resize(targetX, tv.highlight.Y(), targetWidth, 3)
		return
	}

	currentStep := 0
	var animFunc func()
	animFunc = func() {
		currentStep++
		if currentStep >= steps {
			fltk_bridge.Lock()
			tv.highlight.Resize(targetX, tv.highlight.Y(), targetWidth, 3)
			tv.tabBar.Redraw()
			fltk_bridge.Unlock()
			fltk_bridge.AwakeNullMessage()
		} else {
			fltk_bridge.Lock()
			tv.highlight.Resize(startX+currentStep*stepX, tv.highlight.Y(), targetWidth, 3)
			tv.tabBar.Redraw()
			fltk_bridge.Unlock()
			fltk_bridge.AwakeNullMessage()
			fltk_bridge.AddTimeout(0.016, animFunc) // 约 60FPS
		}
	}
	fltk_bridge.AddTimeout(0.016, animFunc)
}

// OnTabChanged 注册切换标签的回调
func (tv *UITabView) OnTabChanged(cb func(index int)) {
	tv.onTabChanged = cb
}
