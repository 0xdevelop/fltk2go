package tableview

import (
	"testing"

	"github.com/0xdevelop/fltk2go/fltk_bridge"
)

type fakeBridgeTable struct {
	rows       int
	redraws    int
	background fltk_bridge.Color
	draw       func(ctx fltk_bridge.TableContext, row, col int, x, y, w, h int)
	event      func(TableInteraction) bool
	selected   int
	scrolled   int
	focuses    int
}

func (f *fakeBridgeTable) SetRows(rows int) { f.rows = rows }
func (f *fakeBridgeTable) Redraw()          { f.redraws++ }
func (f *fakeBridgeTable) SetDrawCellHandler(fn func(ctx fltk_bridge.TableContext, row, col int, x, y, w, h int)) {
	f.draw = fn
}
func (f *fakeBridgeTable) SetEventHandler(fn func(TableInteraction) bool) { f.event = fn }
func (f *fakeBridgeTable) GetSelectedRow() int                            { return f.selected }
func (f *fakeBridgeTable) SelectRow(row int)                              { f.selected = row }
func (f *fakeBridgeTable) ScrollToRow(row int)                            { f.scrolled = row }
func (f *fakeBridgeTable) TakeFocus() bool                                { f.focuses++; return true }
func (f *fakeBridgeTable) Widget() fltk_bridge.Widget                     { return nil }
func (f *fakeBridgeTable) SetColumnCount(int)                             {}
func (f *fakeBridgeTable) SetColumnWidth(int, int)                        {}
func (f *fakeBridgeTable) AllowColumnResizing()                           {}
func (f *fakeBridgeTable) EnableColumnHeaders()                           {}
func (f *fakeBridgeTable) SetColumnHeaderHeight(int)                      {}
func (f *fakeBridgeTable) SetBackgroundColor(color fltk_bridge.Color)     { f.background = color }

func TestSetBackgroundColorForwardsToBridge(t *testing.T) {
	bridge := &fakeBridgeTable{}
	tv := newWithBridgeTable(bridge)
	want := fltk_bridge.Color(0x12345600)
	tv.SetBackgroundColor(want)
	if bridge.background != want {
		t.Fatalf("background = %#x, want %#x", bridge.background, want)
	}
}

func TestEmptyMessagePublishesSemanticStateAndTracksRows(t *testing.T) {
	bridge := &fakeBridgeTable{}
	tv := newWithBridgeTable(bridge)
	ds := &sliceDataSource{}
	tv.SetDataSource(ds)

	tv.SetEmptyMessage("No matching connections")
	if got := tv.View().AutomationSnapshot().Properties["emptyMessage"]; got != "No matching connections" {
		t.Fatalf("empty message automation property = %#v", got)
	}
	if got := tv.emptyMessageForDrawing(); got != "No matching connections" {
		t.Fatalf("empty message for zero rows = %q", got)
	}
	if bridge.redraws != 1 {
		t.Fatalf("setting empty message redrew %d times, want 1", bridge.redraws)
	}

	ds.rows = 1
	if got := tv.emptyMessageForDrawing(); got != "" {
		t.Fatalf("non-empty table drew empty message %q", got)
	}

	tv.SetEmptyMessage("")
	if _, ok := tv.View().AutomationSnapshot().Properties["emptyMessage"]; ok {
		t.Fatal("cleared empty message remains in automation metadata")
	}
}

func TestNativeSelectionGeometryUsesZeroBasedDataIndex(t *testing.T) {
	for raw, want := range map[int]int{-1: -1, 0: 0, 2: 2} {
		if got := normalizeSelectedRow(raw); got != want {
			t.Fatalf("normalizeSelectedRow(%d) = %d, want %d", raw, got, want)
		}
	}
}

func TestNativeRowInteractionTakesKeyboardFocus(t *testing.T) {
	bridge := &fakeBridgeTable{}
	tv := newWithBridgeTable(bridge)
	tv.SetDataSource(&sliceDataSource{rows: 1})
	tv.SetDelegate(&recordingDelegate{})

	if !tv.onEvent(TableInteraction{Context: fltk_bridge.ContextCell, Event: fltk_bridge.RELEASE, Row: 0, Button: fltk_bridge.LeftMouse}) {
		t.Fatal("row interaction was not handled")
	}
	if bridge.focuses != 1 {
		t.Fatalf("row interaction requested focus %d times, want 1", bridge.focuses)
	}

	tv.onEvent(TableInteraction{Context: fltk_bridge.ContextCell, Event: fltk_bridge.RELEASE, Row: 0, Button: fltk_bridge.RightMouse})
	if bridge.focuses != 1 {
		t.Fatalf("right-click changed keyboard focus count to %d", bridge.focuses)
	}
}

type sliceDataSource struct {
	rows  int
	cells map[int]*TableViewCell
}

func (s *sliceDataSource) NumberOfRows(*TableView) int { return s.rows }
func (s *sliceDataSource) CellForColumn(_ *TableView, row, col int) *TableViewCell {
	if s.cells == nil {
		s.cells = map[int]*TableViewCell{}
	}
	cell := NewCell("row")
	s.cells[row] = cell
	return cell
}

type recordingDelegate struct {
	selected []int
}

func (r *recordingDelegate) DidSelectRow(_ *TableView, row int) {
	r.selected = append(r.selected, row)
}
func (r *recordingDelegate) RowHeight(*TableView, int) int { return 0 }

func TestReloadDataUsesDataSourceAndClearsVisibleCells(t *testing.T) {
	bridge := &fakeBridgeTable{}
	tv := newWithBridgeTable(bridge)
	ds := &sliceDataSource{rows: 2}
	tv.SetDataSource(ds)

	tv.ReloadData()

	if bridge.rows != 2 {
		t.Fatalf("rows = %d, want 2", bridge.rows)
	}
	if bridge.redraws != 1 {
		t.Fatalf("redraws = %d, want 1", bridge.redraws)
	}

	tv.cellFor(1, 0)
	if got := tv.visible["1_0"]; got == nil || got.Row() != 1 {
		t.Fatalf("visible row was not cached with row index: %#v", got)
	}

	ds.rows = 1
	tv.ReloadData()
	if len(tv.visible) != 0 {
		t.Fatalf("visible cells were not cleared: %d", len(tv.visible))
	}
	if got := len(tv.reusePool["row"]); got != 1 {
		t.Fatalf("reuse pool size = %d, want 1", got)
	}
}

func TestReloadDataWithoutDataSourceClearsRows(t *testing.T) {
	bridge := &fakeBridgeTable{}
	tv := newWithBridgeTable(bridge)

	tv.ReloadData()

	if bridge.rows != 0 {
		t.Fatalf("rows = %d, want 0", bridge.rows)
	}
	if bridge.redraws != 1 {
		t.Fatalf("redraws = %d, want 1", bridge.redraws)
	}
}

func TestEventDelegateIgnoresNegativeRows(t *testing.T) {
	bridge := &fakeBridgeTable{selected: -1}
	tv := newWithBridgeTable(bridge)
	delegate := &recordingDelegate{}
	tv.SetDelegate(delegate)

	if bridge.event(TableInteraction{Row: -1}) {
		t.Fatal("event(-1) = true, want false")
	}
	if !bridge.event(TableInteraction{Row: 4}) {
		t.Fatal("event(4) = false, want true")
	}
	if len(delegate.selected) != 1 || delegate.selected[0] != 4 {
		t.Fatalf("selected rows = %#v, want [4]", delegate.selected)
	}
}

func TestDoubleClickActivatesSelectedRowAfterSelection(t *testing.T) {
	bridge := &fakeBridgeTable{selected: -1}
	tv := newWithBridgeTable(bridge)
	delegate := &recordingDelegate{}
	tv.SetDelegate(delegate)
	var activated []int
	tv.OnActivate(func(row int) { activated = append(activated, row) })

	if !bridge.event(TableInteraction{Row: 2, Clicks: 1}) {
		t.Fatal("double-click interaction was not handled")
	}
	if len(delegate.selected) != 1 || delegate.selected[0] != 2 {
		t.Fatalf("selection callbacks = %#v, want [2]", delegate.selected)
	}
	if len(activated) != 1 || activated[0] != 2 {
		t.Fatalf("activation callbacks = %#v, want [2]", activated)
	}
}

func TestRightClickRequestsContextMenuWithoutActivatingRow(t *testing.T) {
	bridge := &fakeBridgeTable{selected: -1}
	tv := newWithBridgeTable(bridge)
	tv.SetDataSource(&sliceDataSource{rows: 3})
	delegate := &recordingDelegate{}
	tv.SetDelegate(delegate)
	var requested []TableContextMenuState
	activated := -1
	tv.OnContextMenu(func(state TableContextMenuState) { requested = append(requested, state) })
	tv.OnActivate(func(row int) { activated = row })

	if !bridge.event(TableInteraction{Row: 2, Button: fltk_bridge.RightMouse, Clicks: 1}) {
		t.Fatal("right-click interaction was not handled")
	}
	if len(requested) != 1 || requested[0].Row != 2 || requested[0].Selected {
		t.Fatalf("context-menu requests = %#v, want row 2 previously unselected", requested)
	}
	if activated != -1 {
		t.Fatalf("right click activated row %d", activated)
	}
	if len(delegate.selected) != 1 || delegate.selected[0] != 2 {
		t.Fatalf("right click did not select its row: %#v", delegate.selected)
	}
}

func TestTableContextMenuRequiresHandlerAndValidRow(t *testing.T) {
	bridge := &fakeBridgeTable{selected: 0}
	tv := newWithBridgeTable(bridge)
	tv.SetDataSource(&sliceDataSource{rows: 1})
	if bridge.event(TableInteraction{Row: 0, Button: fltk_bridge.RightMouse}) {
		t.Fatal("right click without an owner handler was consumed")
	}
	tv.OnContextMenu(func(TableContextMenuState) {})
	if bridge.event(TableInteraction{Row: 1, Button: fltk_bridge.RightMouse}) {
		t.Fatal("out-of-range right click was consumed")
	}
}

func TestColumnHeaderClickPublishesColumnWithoutChangingSelection(t *testing.T) {
	bridge := &fakeBridgeTable{selected: 1}
	tv := newWithBridgeTable(bridge)
	tv.SetDataSource(&sliceDataSource{rows: 3})
	for index := 0; index < 3; index++ {
		tv.AddColumn(TableColumn{Identifier: string(rune('a' + index))})
	}
	delegate := &recordingDelegate{}
	tv.SetDelegate(delegate)
	var columns []int
	tv.OnColumnHeaderClick(func(column int) { columns = append(columns, column) })

	if !bridge.event(TableInteraction{Context: fltk_bridge.ContextColHeader, Row: -1, Column: 2, Button: fltk_bridge.LeftMouse}) {
		t.Fatal("column-header click was not handled")
	}
	if len(columns) != 1 || columns[0] != 2 {
		t.Fatalf("header callbacks = %#v, want [2]", columns)
	}
	if bridge.selected != 1 || len(delegate.selected) != 0 {
		t.Fatalf("header click changed row selection: native=%d callbacks=%#v", bridge.selected, delegate.selected)
	}
}

func TestColumnHeaderClickRejectsMissingHandlerInvalidColumnAndRightClick(t *testing.T) {
	bridge := &fakeBridgeTable{selected: 0}
	tv := newWithBridgeTable(bridge)
	tv.AddColumn(TableColumn{Identifier: "name"})
	if bridge.event(TableInteraction{Context: fltk_bridge.ContextColHeader, Row: -1, Column: 0}) {
		t.Fatal("header click without a handler was consumed")
	}
	tv.OnColumnHeaderClick(func(int) {})
	if bridge.event(TableInteraction{Context: fltk_bridge.ContextColHeader, Row: -1, Column: -1}) {
		t.Fatal("invalid header column was consumed")
	}
	if bridge.event(TableInteraction{Context: fltk_bridge.ContextColHeader, Row: -1, Column: 0, Button: fltk_bridge.RightMouse}) {
		t.Fatal("right-click header was treated as primary sorting action")
	}
}

func TestColumnHeaderClickSuppressesReentrantNativeCallbacks(t *testing.T) {
	bridge := &fakeBridgeTable{}
	tv := newWithBridgeTable(bridge)
	tv.AddColumn(TableColumn{Identifier: "name"})
	calls := 0
	tv.OnColumnHeaderClick(func(column int) {
		calls++
		if calls == 1 {
			bridge.event(TableInteraction{Context: fltk_bridge.ContextColHeader, Row: -1, Column: column})
		}
	})
	if !bridge.event(TableInteraction{Context: fltk_bridge.ContextColHeader, Row: -1, Column: 0}) {
		t.Fatal("outer header click was not handled")
	}
	if calls != 1 {
		t.Fatalf("reentrant header callbacks = %d, want 1", calls)
	}
}

func TestColumnHeaderClickIgnoresReleaseAfterPush(t *testing.T) {
	bridge := &fakeBridgeTable{}
	tv := newWithBridgeTable(bridge)
	tv.AddColumn(TableColumn{Identifier: "name"})
	calls := 0
	tv.OnColumnHeaderClick(func(int) { calls++ })
	for _, event := range []fltk_bridge.Event{fltk_bridge.PUSH, fltk_bridge.RELEASE} {
		if !bridge.event(TableInteraction{Context: fltk_bridge.ContextColHeader, Event: event, Row: -1, Column: 0}) {
			t.Fatalf("header event %v was not consumed", event)
		}
	}
	if calls != 1 {
		t.Fatalf("push/release header callbacks = %d, want 1", calls)
	}
}

func TestSelectRowClampsAndPublishesSemanticValue(t *testing.T) {
	bridge := &fakeBridgeTable{selected: -1}
	tv := newWithBridgeTable(bridge)
	tv.SetDataSource(&sliceDataSource{rows: 3})
	delegate := &recordingDelegate{}
	tv.SetDelegate(delegate)

	if !tv.SelectRow(1) {
		t.Fatal("SelectRow(1) = false, want true")
	}
	if bridge.selected != 1 || bridge.scrolled != 1 {
		t.Fatalf("bridge selected/scrolled = %d/%d, want 1/1", bridge.selected, bridge.scrolled)
	}
	if got := tv.View().AutomationSnapshot().Value; got != "1" {
		t.Fatalf("automation selected row = %q, want 1", got)
	}
	if tv.SelectRow(3) || tv.SelectRow(-1) {
		t.Fatal("out-of-range selection unexpectedly succeeded")
	}
}

func TestKeyboardNavigationAndActivation(t *testing.T) {
	bridge := &fakeBridgeTable{selected: -1}
	tv := newWithBridgeTable(bridge)
	tv.SetDataSource(&sliceDataSource{rows: 3})
	var activated int = -1
	tv.OnActivate(func(row int) { activated = row })

	if !tv.handleKey(fltk_bridge.DOWN) || bridge.selected != 0 {
		t.Fatalf("Down selected row %d, want 0", bridge.selected)
	}
	if !tv.handleKey(fltk_bridge.END) || bridge.selected != 2 {
		t.Fatalf("End selected row %d, want 2", bridge.selected)
	}
	if !tv.handleKey(fltk_bridge.ENTER_KEY) || activated != 2 {
		t.Fatalf("Enter activated row %d, want 2", activated)
	}
	if tv.handleKey('x') {
		t.Fatal("unhandled key unexpectedly consumed")
	}
}

func TestOwnerKeyboardHandlerRunsBeforeBuiltInNavigation(t *testing.T) {
	bridge := &fakeBridgeTable{selected: 1}
	tv := newWithBridgeTable(bridge)
	tv.SetDataSource(&sliceDataSource{rows: 3})
	var events []TableKeyEvent
	tv.OnKey(func(event TableKeyEvent) bool {
		events = append(events, event)
		return event.Key == ' '
	})

	space := TableKeyEvent{Key: ' ', Text: " ", State: fltk_bridge.SHIFT}
	if !tv.handleKeyEvent(space) {
		t.Fatal("owner-handled Space was not consumed")
	}
	if bridge.selected != 1 {
		t.Fatalf("owner-handled key changed selection to %d", bridge.selected)
	}
	if !tv.handleKeyEvent(TableKeyEvent{Key: fltk_bridge.DOWN}) || bridge.selected != 2 {
		t.Fatalf("unhandled Down did not use built-in navigation: selected=%d", bridge.selected)
	}
	if len(events) != 2 || events[0] != space || events[1].Key != fltk_bridge.DOWN {
		t.Fatalf("owner key events = %#v", events)
	}

	tv.OnKey(nil)
	if tv.handleKeyEvent(TableKeyEvent{Key: 'x'}) {
		t.Fatal("cleared owner handler still consumed an unknown key")
	}
}

func TestTableViewNilSafety(t *testing.T) {
	var tv *TableView
	tv.SetDataSource(nil)
	tv.SetDelegate(nil)
	tv.SetDefaultRowHeight(10)
	tv.SetCustomDraw(nil)
	tv.Enqueue(nil)
	tv.ReloadData()
	tv.onDrawCell(fltk_bridge.ContextCell, 0, 0, 0, 0, 0, 0)

	if tv.View() != nil {
		t.Fatal("nil TableView View() returned non-nil")
	}
	if tv.Raw() != nil {
		t.Fatal("nil TableView Raw() returned non-nil")
	}
	if tv.GetSelectedRow() != -1 {
		t.Fatal("nil TableView selected row should be -1")
	}
	if !tv.onEvent(TableInteraction{Row: 0}) {
		// expected false; keep branch explicit to ensure no panic above
		return
	}
	t.Fatal("nil TableView onEvent returned true")
}
