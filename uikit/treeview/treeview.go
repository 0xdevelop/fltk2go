package treeview

import (
	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

type TreeDataSource interface {
	// RootNode 返回根节点的路径名称
	RootNode() string
	// Children 返回指定 path 下的所有子节点相对名称列表
	Children(path string) []string
	// HasChildren 判断指定 path 是否有子节点，用于懒加载展示折叠图标
	HasChildren(path string) bool
}

type UITreeView struct {
	v          view.UIView
	raw        *fltk_bridge.Tree
	dataSource TreeDataSource

	onItemSelected func(path string)
}

func NewUITreeView(r *foundation.Rect) *UITreeView {
	if r == nil {
		r = &foundation.Rect{X: 0, Y: 0, Width: 200, Height: 300}
	}
	t := fltk_bridge.NewTree(r.X, r.Y, r.Width, r.Height)
	t.SetShowRoot(true)
	t.SetConnectorStyle(fltk_bridge.TreeConnectorDotted)
	t.SetSelectMode(fltk_bridge.TreeSelectSingle)

	tv := &UITreeView{raw: t}
	tv.v.BindRaw(t)

	// 绑定回调，处理懒加载与选中事件
	t.SetCallback(func() { tv.handleTreeCallback() })

	return tv
}

func (tv *UITreeView) View() *view.UIView     { return &tv.v }
func (tv *UITreeView) Raw() *fltk_bridge.Tree { return tv.raw }

func (tv *UITreeView) SetDataSource(ds TreeDataSource) {
	tv.dataSource = ds
	tv.ReloadData()
}

func (tv *UITreeView) ReloadData() {
	if tv.raw == nil || tv.dataSource == nil {
		return
	}
	tv.raw.Clear()
	root := tv.dataSource.RootNode()
	if root == "" {
		root = "Root"
	}

	tv.loadChildren(root)
}

func (tv *UITreeView) loadChildren(path string) {
	children := tv.dataSource.Children(path)
	for _, child := range children {
		childFullPath := path + "/" + child

		tv.raw.Add(childFullPath)
		// 懒加载：如果有子节点，添加一个 dummy 节点以显示展开图标
		if tv.dataSource.HasChildren(childFullPath) {
			tv.raw.Add(childFullPath + "/(loading...)")
		}
	}
}

func (tv *UITreeView) handleTreeCallback() {
	if tv.raw == nil {
		return
	}

	reason := tv.raw.CallbackReason()
	item := tv.raw.CallbackItem()
	if !item.IsValid() {
		return
	}

	path := tv.raw.ItemPathname(item)

	switch reason {
	case fltk_bridge.TreeReasonSelected:
		if tv.onItemSelected != nil {
			tv.onItemSelected(path)
		}
	case fltk_bridge.TreeReasonOpened:
		// 懒加载处理
		if tv.dataSource != nil {
			// 清理占位的 loading 节点
			tv.raw.ClearChildren(item)
			tv.loadChildren(path)
		}
	}
}

func (tv *UITreeView) OnItemSelected(cb func(path string)) {
	tv.onItemSelected = cb
}
