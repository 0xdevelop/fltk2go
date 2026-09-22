package tableview

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/uikit/textlayout"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

type TableColumn struct {
	Identifier string
	Title      string
	Width      int
	Align      fltk_bridge.Align
}

// TableContextMenuState identifies the data row that requested a native context
// menu. Selected reports whether it was already selected before a pointer
// request. TableView selects and publishes a right-clicked row before invoking
// the owner callback; keyboard requests retain the current selection.
type TableContextMenuState struct {
	Row      int
	Selected bool
}

// TableKeyEvent is the native keyboard state offered to an owner before the
// table applies its built-in Enter and selection-navigation behavior.
type TableKeyEvent struct {
	Key   int
	Text  string
	State int
}

type TableView struct {
	table      BridgeTable
	v          view.UIView
	customDraw func(ctx fltk_bridge.TableContext, row, col, x, y, w, h int)

	dataSource        DataSource
	delegate          Delegate
	onActivate        func(row int)
	onContextMenu     func(TableContextMenuState)
	onHeaderClick     func(column int)
	onKey             func(TableKeyEvent) bool
	headerClickActive bool
	selectedRow       int

	columns []TableColumn

	defaultRowHeight int
	headerHeight     int
	emptyMessage     string

	reusePool   map[string][]*TableViewCell
	visible     map[string]*TableViewCell
	drawnInPage map[string]bool
}

func New(x, y, w, h int) (*TableView, error) {
	bt, err := newBridgeTable(x, y, w, h)
	if err != nil {
		return nil, err
	}
	return newWithBridgeTable(bt), nil
}

func newWithBridgeTable(bt BridgeTable) *TableView {
	tv := &TableView{
		table:            bt,
		selectedRow:      -1,
		defaultRowHeight: 24,
		headerHeight:     24,
		reusePool:        map[string][]*TableViewCell{},
		visible:          map[string]*TableViewCell{},
		drawnInPage:      map[string]bool{},
	}

	if bt != nil {
		if raw := bt.Widget(); raw != nil {
			tv.v.BindRaw(raw)
		}
		bt.SetDrawCellHandler(tv.onDrawCell)
		bt.SetEventHandler(tv.onEvent)
	}
	tv.v.SetAutomationRole("table").SetAutomationValueHandler(func() (string, bool) {
		return strconv.Itoa(tv.GetSelectedRow()), true
	})
	tv.v.On(fltk_bridge.KEYDOWN, func(fltk_bridge.Event) bool {
		return tv.handleKeyEvent(TableKeyEvent{
			Key: fltk_bridge.EventKey(), Text: fltk_bridge.EventText(), State: fltk_bridge.EventState(),
		})
	})

	return tv
}

func (tv *TableView) AddColumn(c TableColumn) {
	if tv == nil || tv.table == nil {
		return
	}
	if c.Width == 0 {
		c.Width = 100
	}
	if c.Align == 0 {
		c.Align = fltk_bridge.ALIGN_CENTER | fltk_bridge.ALIGN_CLIP
	}
	tv.columns = append(tv.columns, c)
	tv.table.SetColumnCount(len(tv.columns))
	tv.table.SetColumnWidth(len(tv.columns)-1, c.Width)
	tv.table.AllowColumnResizing()
	tv.table.EnableColumnHeaders()
	tv.table.SetColumnHeaderHeight(tv.headerHeight)
}

// View implements view.Viewable — enables root.AddSubview(tv).
func (tv *TableView) View() *view.UIView {
	if tv == nil {
		return nil
	}
	return &tv.v
}

// Raw returns the underlying BridgeTable (e.g. for win.Raw().Add(tv.Raw().Widget())).
func (tv *TableView) Raw() BridgeTable {
	if tv == nil {
		return nil
	}
	return tv.table
}

func (tv *TableView) SetDataSource(ds DataSource) {
	if tv != nil {
		tv.dataSource = ds
	}
}

func (tv *TableView) SetDelegate(d Delegate) {
	if tv != nil {
		tv.delegate = d
	}
}

func (tv *TableView) SetDefaultRowHeight(h int) {
	if tv != nil && h > 0 {
		tv.defaultRowHeight = h
	}
}

func (tv *TableView) SetHeaderHeight(h int) {
	if tv != nil && tv.table != nil && h > 0 {
		tv.headerHeight = h
		tv.table.SetColumnHeaderHeight(h)
	}
}

// SetBackgroundColor sets the table viewport background, including the empty
// area below the last row.
func (tv *TableView) SetBackgroundColor(color fltk_bridge.Color) {
	if tv != nil && tv.table != nil {
		tv.table.SetBackgroundColor(color)
		tv.table.Redraw()
	}
}

// SetEmptyMessage configures centered supporting text for an empty table and
// mirrors that state into semantic automation metadata.
func (tv *TableView) SetEmptyMessage(message string) {
	if tv == nil {
		return
	}
	tv.emptyMessage = strings.TrimSpace(message)
	if tv.emptyMessage == "" {
		tv.v.ClearAutomationProperty("emptyMessage")
	} else {
		tv.v.SetAutomationProperty("emptyMessage", tv.emptyMessage)
	}
	if tv.table != nil {
		tv.table.Redraw()
	}
}

func (tv *TableView) emptyMessageForDrawing() string {
	if tv == nil || tv.emptyMessage == "" || tv.dataSource == nil || tv.dataSource.NumberOfRows(tv) != 0 {
		return ""
	}
	return tv.emptyMessage
}

func (tv *TableView) drawEmptyMessage(x, y, w, h int) {
	message := tv.emptyMessageForDrawing()
	if message == "" || w <= 0 || h <= 0 {
		return
	}
	fltk_bridge.PushClip(x, y, w, h)
	fltk_bridge.SetDrawColor(fltk_bridge.Color(0x7A849300))
	fltk_bridge.SetDrawFont(fltk_bridge.HELVETICA, 14)
	fltk_bridge.Draw(message, x+16, y, w-32, h, fltk_bridge.ALIGN_CENTER|fltk_bridge.ALIGN_CLIP)
	fltk_bridge.PopClip()
}

// SetCustomDraw sets a custom cell-drawing function called for every visible row.
// When set, it replaces the default DataSource-driven cell drawing.
func (tv *TableView) SetCustomDraw(fn func(ctx fltk_bridge.TableContext, row, col, x, y, w, h int)) {
	if tv != nil {
		tv.customDraw = fn
	}
}

// GetSelectedRow returns the 0-based index of the selected row, or -1 if none.
func (tv *TableView) GetSelectedRow() int {
	if tv == nil || tv.table == nil {
		return -1
	}
	if tv.selectedRow >= 0 {
		return tv.selectedRow
	}
	return tv.table.GetSelectedRow()
}

// SelectRow selects, reveals, and publishes a data row through the delegate.
// It returns false when the row is outside the current data source.
func (tv *TableView) SelectRow(row int) bool {
	if tv == nil || tv.table == nil || tv.dataSource == nil || row < 0 || row >= tv.dataSource.NumberOfRows(tv) {
		return false
	}
	tv.table.SelectRow(row)
	tv.table.ScrollToRow(row)
	tv.selectedRow = row
	if tv.delegate != nil {
		tv.delegate.DidSelectRow(tv, row)
	}
	tv.table.Redraw()
	return true
}

// OnActivate registers the primary row action used by native double-click,
// Enter, and semantic debug automation.
func (tv *TableView) OnActivate(handler func(row int)) {
	if tv == nil {
		return
	}
	tv.onActivate = handler
	if handler == nil {
		tv.v.OnAutomationClick(nil)
		return
	}
	tv.v.OnAutomationClick(func() error {
		tv.ActivateSelected()
		return nil
	})
}

// OnContextMenu registers an application-owned menu request for native
// right-clicks and the standard Shift+F10/Menu keyboard gestures. The table owns
// hit-testing and row selection while the application owns labels, enablement,
// and actions.
func (tv *TableView) OnContextMenu(handler func(TableContextMenuState)) {
	if tv != nil {
		tv.onContextMenu = handler
	}
}

// OnColumnHeaderClick registers a primary-click callback for column headers.
// Header interaction never changes the selected data row; consumers can use
// the stable zero-based column index to implement sorting or filtering.
func (tv *TableView) OnColumnHeaderClick(handler func(column int)) {
	if tv != nil {
		tv.onHeaderClick = handler
	}
}

// OnKey registers an owner keyboard policy. The handler runs before built-in
// table navigation and should return true only when it consumes the event.
// This lets products add row commands while preserving native navigation.
func (tv *TableView) OnKey(handler func(TableKeyEvent) bool) {
	if tv != nil {
		tv.onKey = handler
	}
}

// ActivateSelected invokes the primary action for the selected row.
func (tv *TableView) ActivateSelected() bool {
	if tv == nil || tv.onActivate == nil {
		return false
	}
	row := tv.GetSelectedRow()
	if row < 0 {
		return false
	}
	tv.onActivate(row)
	return true
}

// RequestContextMenuSelected requests the owner menu for the current valid row
// without republishing selection or invoking the row's primary action.
func (tv *TableView) RequestContextMenuSelected() bool {
	if tv == nil || tv.onContextMenu == nil || tv.dataSource == nil {
		return false
	}
	row := tv.GetSelectedRow()
	if row < 0 || row >= tv.dataSource.NumberOfRows(tv) {
		return false
	}
	tv.onContextMenu(TableContextMenuState{Row: row, Selected: true})
	return true
}

func isContextMenuKey(event TableKeyEvent) bool {
	modifiers := event.State & (fltk_bridge.SHIFT | fltk_bridge.CTRL | fltk_bridge.ALT | fltk_bridge.META)
	return (event.Key == fltk_bridge.F10 && modifiers == fltk_bridge.SHIFT) ||
		(event.Key == fltk_bridge.MENU && modifiers == 0)
}

func (tv *TableView) handleKey(key int) bool {
	return tv.handleKeyEvent(TableKeyEvent{Key: key})
}

func (tv *TableView) handleKeyEvent(event TableKeyEvent) bool {
	if tv != nil && tv.onKey != nil && tv.onKey(event) {
		return true
	}
	if isContextMenuKey(event) {
		return tv.RequestContextMenuSelected()
	}
	if tv == nil || tv.dataSource == nil {
		return false
	}
	rows := tv.dataSource.NumberOfRows(tv)
	if rows <= 0 {
		return false
	}
	selected := tv.GetSelectedRow()
	switch event.Key {
	case fltk_bridge.ENTER_KEY:
		return tv.ActivateSelected()
	case fltk_bridge.UP:
		if selected < 0 {
			selected = rows
		}
		return tv.SelectRow(max(0, selected-1))
	case fltk_bridge.DOWN:
		return tv.SelectRow(min(rows-1, selected+1))
	case fltk_bridge.HOME:
		return tv.SelectRow(0)
	case fltk_bridge.END:
		return tv.SelectRow(rows - 1)
	default:
		return false
	}
}

func (tv *TableView) Dequeue(reuseID string) *TableViewCell {
	if tv == nil {
		return NewCell(reuseID)
	}
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
	if tv == nil || c == nil || c.ReuseID == "" {
		return
	}
	tv.reusePool[c.ReuseID] = append(tv.reusePool[c.ReuseID], c)
}

func (tv *TableView) ReloadData() {
	if tv == nil || tv.table == nil {
		return
	}
	if tv.dataSource == nil {
		tv.table.SetRows(0)
		tv.table.Redraw()
		return
	}

	for _, cell := range tv.visible {
		tv.Enqueue(cell)
	}
	tv.visible = map[string]*TableViewCell{}

	rows := tv.dataSource.NumberOfRows(tv)
	if rows < 0 {
		rows = 0
	}
	if tv.selectedRow >= rows {
		tv.selectedRow = -1
	}
	tv.table.SetRows(rows)
	tv.table.Redraw()
}

// ── callbacks ──────────────────────────────────────────────────────────────

func (tv *TableView) onDrawCell(ctx fltk_bridge.TableContext, row, col int, x, y, w, h int) {
	if tv == nil {
		return
	}
	if ctx == fltk_bridge.ContextEmpty {
		tv.drawEmptyMessage(x, y, w, h)
		return
	}
	if tv.customDraw != nil {
		tv.customDraw(ctx, row, col, x, y, w, h)
		return
	}

	switch ctx {
	case fltk_bridge.ContextStartPage:
		fltk_bridge.SetDrawFont(fltk_bridge.HELVETICA, 14)
		tv.drawnInPage = make(map[string]bool)
	case fltk_bridge.ContextEndPage:
		// Cleanup invisible cells to prevent memory leak and allow reuse
		for key, cell := range tv.visible {
			if !tv.drawnInPage[key] {
				tv.Enqueue(cell)
				delete(tv.visible, key)
			}
		}
	case fltk_bridge.ContextColHeader:
		fltk_bridge.PushClip(x, y, w, h)
		fltk_bridge.DrawBox(fltk_bridge.THIN_UP_BOX, x, y, w, h, fltk_bridge.Color(0xDDDDDD00))
		fltk_bridge.SetDrawColor(fltk_bridge.Color(0))
		if col >= 0 && col < len(tv.columns) {
			fltk_bridge.Draw(tv.columns[col].Title, x, y, w, h, tv.columns[col].Align)
		}
		fltk_bridge.PopClip()
	case fltk_bridge.ContextCell:
		if tv.dataSource == nil {
			return
		}
		fltk_bridge.PushClip(x, y, w, h)

		cell := tv.cellFor(row, col)

		bgColor := fltk_bridge.Color(0xFFFFFF00)
		selected := tv.table != nil && tv.table.GetSelectedRow() == row
		if selected {
			bgColor = fltk_bridge.Color(0xBBDEFB00)
		} else if row%2 == 1 {
			bgColor = fltk_bridge.Color(0xF5F5F500)
		}
		fltk_bridge.DrawBox(fltk_bridge.FLAT_BOX, x, y, w, h, bgColor)

		if cell != nil {
			fltk_bridge.SetDrawColor(cell.TextColor)
			fltk_bridge.SetDrawFont(cell.Font, cell.FontSize)

			if cell.preparedText != nil {
				lineHeight := cell.FontSize + 4
				layoutResult := textlayout.Layout(cell.preparedText, w-10, lineHeight)

				startY := y + (h-layoutResult.Height)/2
				if startY < y {
					startY = y
				}

				for i, line := range layoutResult.Lines {
					fltk_bridge.Draw(line.Text, x+5, startY+(i*lineHeight), w-10, lineHeight, cell.Align)
				}
			} else if cell.Text != "" {
				fltk_bridge.Draw(cell.Text, x+5, y, w-10, h, cell.Align)
			}
		}

		fltk_bridge.SetDrawColor(fltk_bridge.Color(0xDDDDDD00))
		fltk_bridge.DrawRect(x, y, w, h)

		fltk_bridge.PopClip()
	}
}

func (tv *TableView) cellFor(row, col int) *TableViewCell {
	if tv == nil || tv.dataSource == nil {
		return nil
	}
	key := fmt.Sprintf("%d_%d", row, col)
	tv.drawnInPage[key] = true
	cell, ok := tv.visible[key]
	if ok {
		return cell
	}
	cell = tv.dataSource.CellForColumn(tv, row, col)
	if cell != nil {
		cell.row = row
		cell.col = col
		tv.visible[key] = cell
	}
	return cell
}

func (tv *TableView) onEvent(interaction TableInteraction) bool {
	if tv == nil {
		return false
	}
	if interaction.Context == fltk_bridge.ContextColHeader {
		if interaction.Button == fltk_bridge.RightMouse || interaction.Column < 0 || interaction.Column >= len(tv.columns) || tv.onHeaderClick == nil {
			return false
		}
		if interaction.Event != fltk_bridge.NO_EVENT && interaction.Event != fltk_bridge.PUSH {
			return true
		}
		if tv.headerClickActive {
			return true
		}
		tv.headerClickActive = true
		defer func() { tv.headerClickActive = false }()
		tv.onHeaderClick(interaction.Column)
		return true
	}
	if interaction.Row < 0 {
		return false
	}
	if interaction.Button == fltk_bridge.RightMouse {
		if tv.onContextMenu == nil || tv.dataSource == nil || interaction.Row >= tv.dataSource.NumberOfRows(tv) {
			return false
		}
	} else if tv.table != nil {
		// A native row click selects data but FLTK does not reliably transfer
		// keyboard focus to Fl_Table_Row. Claim it explicitly so owner key
		// handlers and built-in row navigation work immediately after a click.
		tv.table.TakeFocus()
	}
	selected := tv.GetSelectedRow() == interaction.Row
	tv.selectedRow = interaction.Row
	if tv.delegate != nil {
		tv.delegate.DidSelectRow(tv, interaction.Row)
	}
	if interaction.Button == fltk_bridge.RightMouse {
		tv.onContextMenu(TableContextMenuState{Row: interaction.Row, Selected: selected})
		return true
	}
	if interaction.Clicks > 0 && tv.onActivate != nil {
		tv.onActivate(interaction.Row)
	}
	return tv.delegate != nil || tv.onActivate != nil
}
