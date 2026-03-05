package tableview

import "github.com/0xYeah/fltk2go/fltk_bridge"

// BridgeTable is the minimal interface TableView needs from the underlying FLTK table.
type BridgeTable interface {
	SetRows(rows int)
	Redraw()
	SetDrawCellHandler(fn func(row int, x, y, w, h int))
	SetEventHandler(fn func(row int) bool)
	GetSelectedRow() int
	Widget() fltk_bridge.Widget
}

func newBridgeTable(x, y, w, h int) (BridgeTable, error) {
	return newBridgeTableImpl(x, y, w, h), nil
}
