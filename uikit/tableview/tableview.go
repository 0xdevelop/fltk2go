package tableview

import "github.com/0xYeah/fltk2go/uikit/view"

type TableView struct {
	table      BridgeTable
	v          view.UIView
	customDraw func(row, x, y, w, h int)

	dataSource DataSource
	delegate   Delegate

	defaultRowHeight int

	reusePool map[string][]*TableViewCell
	visible   map[int]*TableViewCell
}

func New(x, y, w, h int) (*TableView, error) {
	bt, err := newBridgeTable(x, y, w, h)
	if err != nil {
		return nil, err
	}

	tv := &TableView{
		table:            bt,
		defaultRowHeight: 24,
		reusePool:        map[string][]*TableViewCell{},
		visible:          map[int]*TableViewCell{},
	}
	tv.v.BindRaw(bt.Widget())

	tv.table.SetDrawCellHandler(tv.onDrawCell)
	tv.table.SetEventHandler(tv.onEvent)

	return tv, nil
}

// View implements view.Viewable — enables root.AddSubview(tv).
func (tv *TableView) View() *view.UIView { return &tv.v }

// Raw returns the underlying BridgeTable (e.g. for win.Raw().Add(tv.Raw().Widget())).
func (tv *TableView) Raw() BridgeTable { return tv.table }

func (tv *TableView) SetDataSource(ds DataSource) { tv.dataSource = ds }
func (tv *TableView) SetDelegate(d Delegate)      { tv.delegate = d }

func (tv *TableView) SetDefaultRowHeight(h int) {
	if h > 0 {
		tv.defaultRowHeight = h
	}
}

// SetCustomDraw sets a custom cell-drawing function called for every visible row.
// When set, it replaces the default DataSource-driven cell drawing.
func (tv *TableView) SetCustomDraw(fn func(row, x, y, w, h int)) {
	tv.customDraw = fn
}

// GetSelectedRow returns the 0-based index of the selected row, or -1 if none.
func (tv *TableView) GetSelectedRow() int {
	return tv.table.GetSelectedRow()
}

func (tv *TableView) Dequeue(reuseID string) *TableViewCell {
	list := tv.reusePool[reuseID]
	if n := len(list); n > 0 {
		c := list[n-1]
		tv.reusePool[reuseID] = list[:n-1]
		c.PrepareForReuse()
		return c
	}
	return NewCell(reuseID)
}

func (tv *TableView) Enqueue(c *TableViewCell) {
	if c == nil || c.ReuseID == "" {
		return
	}
	tv.reusePool[c.ReuseID] = append(tv.reusePool[c.ReuseID], c)
}

func (tv *TableView) ReloadData() {
	if tv.dataSource == nil {
		tv.table.SetRows(0)
		tv.table.Redraw()
		return
	}

	for row, cell := range tv.visible {
		_ = row
		tv.Enqueue(cell)
	}
	tv.visible = map[int]*TableViewCell{}

	rows := tv.dataSource.NumberOfRows(tv)
	if rows < 0 {
		rows = 0
	}
	tv.table.SetRows(rows)
	tv.table.Redraw()
}

// ── callbacks ──────────────────────────────────────────────────────────────

func (tv *TableView) onDrawCell(row int, x, y, w, h int) {
	if tv.customDraw != nil {
		tv.customDraw(row, x, y, w, h)
		return
	}
	if tv.dataSource == nil {
		return
	}
	cell, ok := tv.visible[row]
	if !ok {
		cell = tv.dataSource.CellForRow(tv, row)
		if cell == nil {
			return
		}
		cell.row = row
		tv.visible[row] = cell
	}
	_, _, _, _ = x, y, w, h
}

func (tv *TableView) onEvent(row int) bool {
	if tv.delegate != nil && row >= 0 {
		tv.delegate.DidSelectRow(tv, row)
		return true
	}
	return false
}
