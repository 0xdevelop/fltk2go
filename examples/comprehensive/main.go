package main

import (
	"runtime"
	"strconv"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/label"
	"github.com/0xYeah/fltk2go/uikit/tableview"
	"github.com/0xYeah/fltk2go/uikit/window"
)

const (
	BLUE   uint = 0x42A5F500
	GRAY   uint = 0x75757500
	ORANGE uint = 0xFFA72600
	GREEN  uint = 0x4CAF5000
	RED    uint = 0xF4433600
	WHITE  uint = 0xFFFFFFFF
)

type Item struct {
	ID    string
	Name  string
	Value string
}

type ItemDataSource struct {
	items []Item
}

func NewItemDataSource() *ItemDataSource {
	return &ItemDataSource{items: []Item{
		{ID: "1", Name: "Item 1", Value: "Value 1"},
		{ID: "2", Name: "Item 2", Value: "Value 2"},
		{ID: "3", Name: "Item 3", Value: "Value 3"},
		{ID: "4", Name: "Item 4", Value: "Value 4"},
		{ID: "5", Name: "Item 5", Value: "Value 5"},
	}}
}

func (ds *ItemDataSource) NumberOfRows(tv *tableview.TableView) int { return len(ds.items) }

func (ds *ItemDataSource) CellForRow(tv *tableview.TableView, row int) *tableview.TableViewCell {
	cell := tv.Dequeue("itemCell")
	if cell == nil {
		cell = tableview.NewCell("itemCell")
	}
	_ = ds.items[row]
	return cell
}

type ItemDelegate struct{}

func (d *ItemDelegate) DidSelectRow(tv *tableview.TableView, row int) {
	print("Selected row:", row, "\n")
}
func (d *ItemDelegate) RowHeight(tv *tableview.TableView, row int) int { return 40 }

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(1000, 700, "FLTK2Go Comprehensive Example")
	root := win.RootView()

	title := label.NewUILabel(&foundation.Rect{X: 50, Y: 20, Width: 900, Height: 40}, "FLTK2Go UI Components Example")
	title.SetFontSize(24)
	title.SetFont(fltk_bridge.HELVETICA_BOLD)
	title.SetAlignment(fltk_bridge.ALIGN_CENTER)
	root.AddSubview(title)

	description := label.NewUILabel(&foundation.Rect{X: 50, Y: 70, Width: 900, Height: 30}, "This example demonstrates various UI components including buttons, labels, and table view.")
	description.SetFontSize(14)
	description.SetTextColor(GRAY)
	root.AddSubview(description)

	buttonTitle := label.NewUILabel(&foundation.Rect{X: 50, Y: 120, Width: 400, Height: 30}, "Button Examples:")
	buttonTitle.SetFontSize(16)
	buttonTitle.SetFont(fltk_bridge.HELVETICA_BOLD)
	root.AddSubview(buttonTitle)

	systemBtn := button.NewUIButton(&foundation.Rect{X: 50, Y: 160, Width: 120, Height: 36}, "System Button")
	systemBtn.SetBackgroundColor(BLUE)
	systemBtn.SetTitleColor(WHITE)
	root.AddSubview(systemBtn)

	checkBtn := button.NewUIButtonWithType(&foundation.Rect{X: 200, Y: 160, Width: 120, Height: 36}, "Checkbox", button.CheckboxButton)
	root.AddSubview(checkBtn)

	radioBtn := button.NewUIButtonWithType(&foundation.Rect{X: 350, Y: 160, Width: 120, Height: 36}, "Radio", button.RadioButton)
	root.AddSubview(radioBtn)

	toggleBtn := button.NewUIButtonWithType(&foundation.Rect{X: 500, Y: 160, Width: 120, Height: 36}, "Toggle", button.ToggleButton)
	root.AddSubview(toggleBtn)

	tableTitle := label.NewUILabel(&foundation.Rect{X: 50, Y: 220, Width: 900, Height: 30}, "Table View Example:")
	tableTitle.SetFontSize(16)
	tableTitle.SetFont(fltk_bridge.HELVETICA_BOLD)
	root.AddSubview(tableTitle)

	tv, _ := tableview.New(50, 260, 900, 300)
	dataSource := NewItemDataSource()
	delegate := &ItemDelegate{}
	tv.SetDataSource(dataSource)
	tv.SetDelegate(delegate)

	tv.SetCustomDraw(func(row int, x, y, w, h int) {
		item := dataSource.items[row]
		fltk_bridge.SetDrawColor(fltk_bridge.WHITE)
		fltk_bridge.DrawRectfWithColor(x, y, w, h, fltk_bridge.WHITE)
		fltk_bridge.SetDrawColor(fltk_bridge.LIGHT2)
		fltk_bridge.DrawRect(x, y, w, h)
		fltk_bridge.SetDrawColor(fltk_bridge.BLACK)
		fltk_bridge.SetDrawFont(fltk_bridge.HELVETICA, 14)
		fltk_bridge.Draw(item.Name, x+10, y, w/2-10, h, fltk_bridge.ALIGN_LEFT|fltk_bridge.ALIGN_INSIDE)
		fltk_bridge.Draw(item.Value, x+w/2, y, w/2-10, h, fltk_bridge.ALIGN_LEFT|fltk_bridge.ALIGN_INSIDE)
	})
	root.AddSubview(tv)

	addBtn := button.NewUIButton(&foundation.Rect{X: 50, Y: 590, Width: 120, Height: 36}, "Add Item")
	addBtn.SetBackgroundColor(GREEN)
	addBtn.SetTitleColor(WHITE)

	removeBtn := button.NewUIButton(&foundation.Rect{X: 200, Y: 590, Width: 120, Height: 36}, "Remove Item")
	removeBtn.SetBackgroundColor(RED)
	removeBtn.SetTitleColor(WHITE)

	refreshBtn := button.NewUIButton(&foundation.Rect{X: 350, Y: 590, Width: 120, Height: 36}, "Refresh")
	refreshBtn.SetBackgroundColor(GRAY)
	refreshBtn.SetTitleColor(WHITE)

	addBtn.OnTouchUpInside(func() {
		newItem := Item{
			ID:    strconv.Itoa(len(dataSource.items) + 1),
			Name:  "Item " + strconv.Itoa(len(dataSource.items)+1),
			Value: "Value " + strconv.Itoa(len(dataSource.items)+1),
		}
		dataSource.items = append(dataSource.items, newItem)
		tv.ReloadData()
	})
	removeBtn.OnTouchUpInside(func() {
		selectedRow := tv.GetSelectedRow()
		if selectedRow >= 0 && selectedRow < len(dataSource.items) {
			dataSource.items = append(dataSource.items[:selectedRow], dataSource.items[selectedRow+1:]...)
			tv.ReloadData()
		}
	})
	refreshBtn.OnTouchUpInside(func() { tv.ReloadData() })

	root.AddSubview(addBtn)
	root.AddSubview(removeBtn)
	root.AddSubview(refreshBtn)

	win.Show()
	fltk2go.Run()
}
