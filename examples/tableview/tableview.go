package tableview

import (
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/button"
	"github.com/0xdevelop/fltk2go/uikit/tableview"
	"github.com/0xdevelop/fltk2go/uikit/view"
	"strconv"
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

func (ds *ServerDataSource) CellForColumn(tv *tableview.TableView, row, col int) *tableview.TableViewCell {
	cell := tv.Dequeue("serverCell")
	if cell == nil {
		cell = tableview.NewCell("serverCell")
	}
	server := ds.servers[row]
	switch col {
	case 0:
		cell.SetText(server.ID)
	case 1:
		cell.SetText(server.Name)
	case 2:
		cell.SetText(server.IP)
	case 3:
		cell.SetText(server.Location)
	case 4:
		cell.SetText(server.Status)
	}
	return cell
}

type ServerDelegate struct{}

func (d *ServerDelegate) DidSelectRow(tv *tableview.TableView, row int) {
	println("Selected row:", row)
}
func (d *ServerDelegate) RowHeight(tv *tableview.TableView, row int) int { return 40 }

func BuildView(parent *view.UIView) view.Viewable {

	tv, _ := tableview.New(50, 50, 700, 400)
	dataSource := NewServerDataSource()
	delegate := &ServerDelegate{}
	tv.AddColumn(tableview.TableColumn{Identifier: "ID", Title: "ID", Width: 50})
	tv.AddColumn(tableview.TableColumn{Identifier: "Name", Title: "Server Name", Width: 150})
	tv.AddColumn(tableview.TableColumn{Identifier: "IP", Title: "IP Address", Width: 150})
	tv.AddColumn(tableview.TableColumn{Identifier: "Location", Title: "Location", Width: 150})
	tv.AddColumn(tableview.TableColumn{Identifier: "Status", Title: "Status", Width: 100})

	tv.SetDataSource(dataSource)
	tv.SetDelegate(delegate)

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

	parent.AddSubview(tv)
	parent.AddSubview(addBtn)
	parent.AddSubview(removeBtn)
	parent.AddSubview(refreshBtn)

	return nil
}
