package fltk_bridge

/*
#include <stdlib.h>
#include "tree.h"
*/
import "C"
import "unsafe"

type Tree struct {
	Group
}

func NewTree(x, y, w, h int, text ...string) *Tree {
	t := &Tree{}
	initWidget(t, unsafe.Pointer(C.go_fltk_new_Tree(C.int(x), C.int(y), C.int(w), C.int(h), cStringOpt(text))))
	return t
}

func (t *Tree) SetShowRoot(show bool) {
	if show {
		C.go_fltk_Tree_set_show_root((*C.Fl_Tree)(t.ptr()), 1)
	} else {
		C.go_fltk_Tree_set_show_root((*C.Fl_Tree)(t.ptr()), 0)
	}
}

type TreeItem struct {
	ptr *C.Fl_Tree_Item
}

func (t *Tree) Add(path string) TreeItem {
	pathStr := C.CString(path)
	defer C.free(unsafe.Pointer(pathStr))
	itemPtr := C.go_fltk_Tree_add((*C.Fl_Tree)(t.ptr()), pathStr)
	return TreeItem{ptr: itemPtr}
}
func (t *Tree) Remove(item TreeItem) bool {
	return C.go_fltk_Tree_remove((*C.Fl_Tree)(t.ptr()), item.ptr) == 0
}
func (t *Tree) Clear() {
	C.go_fltk_Tree_clear((*C.Fl_Tree)(t.ptr()))
}
func (t *Tree) ClearChildren(item TreeItem) {
	C.go_fltk_Tree_clear_children((*C.Fl_Tree)(t.ptr()), item.ptr)
}

func (t TreeItem) SetWidget(w Widget) {
	C.go_fltk_Tree_Item_set_widget(t.ptr, w.getWidget().ptr())
}

type TreeItemDrawMode uint

var (
	TreeItemDrawDefault        = TreeItemDrawMode(C.go_FL_TREE_ITEM_DRAW_DEFAULT)
	TreeItemDrawLabelAndWidget = TreeItemDrawMode(C.go_FL_TREE_ITEM_DRAW_LABEL_AND_WIDGET)
	TreeItemHeightFromWidget   = TreeItemDrawMode(C.go_FL_TREE_ITEM_HEIGHT_FROM_WIDGET)
)

func (t *Tree) SetItemDrawMode(drawMode TreeItemDrawMode) {
	C.go_fltk_Tree_set_item_draw_mode((*C.Fl_Tree)(t.ptr()), C.uint(drawMode))
}

type TreeConnector int

var (
	TreeConnectorNone   = TreeConnector(C.go_FL_TREE_CONNECTOR_NONE)
	TreeConnectorDotted = TreeConnector(C.go_FL_TREE_CONNECTOR_DOTTED)
	TreeConnectorSolid  = TreeConnector(C.go_FL_TREE_CONNECTOR_SOLID)
)

func (t *Tree) SetConnectorStyle(style TreeConnector) {
	C.go_fltk_Tree_set_connector_style((*C.Fl_Tree)(t.ptr()), C.int(style))
}

type TreeSelect int

var (
	TreeSelectNone            = TreeSelect(C.go_FL_TREE_SELECT_NONE)
	TreeSelectSingle          = TreeSelect(C.go_FL_TREE_SELECT_SINGLE)
	TreeSelectMulti           = TreeSelect(C.go_FL_TREE_SELECT_MULTI)
	TreeSelectSingleDraggable = TreeSelect(C.go_FL_TREE_SELECT_SINGLE_DRAGGABLE)
)

func (t *Tree) SetSelectMode(selectMode TreeSelect) {
	C.go_fltk_Tree_set_select_mode((*C.Fl_Tree)(t.ptr()), C.int(selectMode))
}

type TreeReason int

const (
	TreeReasonNone       TreeReason = 0 // FL_TREE_REASON_NONE
	TreeReasonSelected   TreeReason = 1 // FL_TREE_REASON_SELECTED
	TreeReasonDeselected TreeReason = 2 // FL_TREE_REASON_DESELECTED
	TreeReasonOpened     TreeReason = 3 // FL_TREE_REASON_OPENED
	TreeReasonClosed     TreeReason = 4 // FL_TREE_REASON_CLOSED
	TreeReasonDragged    TreeReason = 5 // FL_TREE_REASON_DRAGGED
)

func (t *Tree) CallbackItem() TreeItem {
	return TreeItem{ptr: C.go_fltk_Tree_callback_item((*C.Fl_Tree)(t.ptr()))}
}

func (t *Tree) CallbackReason() TreeReason {
	return TreeReason(C.go_fltk_Tree_callback_reason((*C.Fl_Tree)(t.ptr())))
}

func (t *Tree) ItemPathname(item TreeItem) string {
	if item.ptr == nil {
		return ""
	}
	return C.GoString(C.go_fltk_Tree_item_pathname((*C.Fl_Tree)(t.ptr()), item.ptr))
}

func (t TreeItem) IsValid() bool {
	return t.ptr != nil
}
