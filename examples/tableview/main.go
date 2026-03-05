package main

import (
	"runtime"
	"strconv"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/tableview"
	"github.com/0xYeah/fltk2go/uikit/window"
)

const (
	BLUE   uint = 0x42A5F500
	GRAY   uint = 0x75757500
	ORANGE uint = 0xFFA72600
)

type Server struct {
	ID       string
	Name     string
	IP       string
	Location string
	Status   string
}

type ServerDataSource struct {
	servers []Server
}

func NewServerDataSource() *ServerDataSource {
	return &ServerDataSource{servers: []Server{
		{ID: "1", Name: "Server 1", IP: "192.168.1.1", Location: "Beijing", Status: "Running"},
		{ID: "2", Name: "Server 2", IP: "192.168.1.2", Location: "Shanghai", Status: "Stopped"},
		{ID: "3", Name: "Server 3", IP: "192.168.1.3", Location: "Guangzhou", Status: "Running"},
		{ID: "4", Name: "Server 4", IP: "192.168.1.4", Location: "Shenzhen", Status: "Running"},
		{ID: "5", Name: "Server 5", IP: "192.168.1.5", Location: "Chengdu", Status: "Stopped"},
	}}
}

func (ds *ServerDataSource) NumberOfRows(tv *tableview.TableView) int { return len(ds.servers) }

func (ds *ServerDataSource) CellForRow(tv *tableview.TableView, row int) *tableview.TableViewCell {
	cell := tv.Dequeue("serverCell")
	if cell == nil {
		cell = tableview.NewCell("serverCell")
	}
	_ = ds.servers[row]
	return cell
}

type ServerDelegate struct{}

func (d *ServerDelegate) DidSelectRow(tv *tableview.TableView, row int) {
	println("Selected row:", row)
}
func (d *ServerDelegate) RowHeight(tv *tableview.TableView, row int) int { return 40 }

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(800, 600, "Server Management")
	root := win.RootView()

	tv, _ := tableview.New(50, 50, 700, 400)
	dataSource := NewServerDataSource()
	delegate := &ServerDelegate{}
	tv.SetDataSource(dataSource)
	tv.SetDelegate(delegate)

	tv.SetCustomDraw(func(row int, x, y, w, h int) {
		server := dataSource.servers[row]
		println("Drawing row", row, ":", server.Name)
	})

	addBtn := button.NewUIButton(&foundation.Rect{X: 50, Y: 500, Width: 120, Height: 36}, "Add Server")
	addBtn.SetBackgroundColor(BLUE)

	removeBtn := button.NewUIButton(&foundation.Rect{X: 200, Y: 500, Width: 120, Height: 36}, "Remove Server")
	removeBtn.SetBackgroundColor(ORANGE)

	refreshBtn := button.NewUIButton(&foundation.Rect{X: 350, Y: 500, Width: 120, Height: 36}, "Refresh")
	refreshBtn.SetBackgroundColor(GRAY)

	addBtn.OnTouchUpInside(func() {
		newServer := Server{
			ID:       strconv.Itoa(len(dataSource.servers) + 1),
			Name:     "Server " + strconv.Itoa(len(dataSource.servers)+1),
			IP:       "192.168.1." + strconv.Itoa(len(dataSource.servers)+1),
			Location: "New Location",
			Status:   "Running",
		}
		dataSource.servers = append(dataSource.servers, newServer)
		tv.ReloadData()
	})
	removeBtn.OnTouchUpInside(func() {
		selectedRow := tv.GetSelectedRow()
		if selectedRow >= 0 && selectedRow < len(dataSource.servers) {
			dataSource.servers = append(dataSource.servers[:selectedRow], dataSource.servers[selectedRow+1:]...)
			tv.ReloadData()
		}
	})
	refreshBtn.OnTouchUpInside(func() { tv.ReloadData() })

	root.AddSubview(tv)
	root.AddSubview(addBtn)
	root.AddSubview(removeBtn)
	root.AddSubview(refreshBtn)

	win.Show()
	fltk2go.Run()
}
