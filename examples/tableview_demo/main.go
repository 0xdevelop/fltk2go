package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/label"
	"github.com/0xYeah/fltk2go/uikit/tableview"
	"github.com/0xYeah/fltk2go/uikit/window"
)

const (
	BLUE   uint = 0x42A5F500
	GRAY   uint = 0x75757500
	WHITE  uint = 0xFFFFFF00
	GREEN  uint = 0x66BB6A00
	RED    uint = 0xE5393500
	ORANGE uint = 0xFFA72600
)

// Server 服务器数据结构
type Server struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Status string `json:"status"`
}

// ServerList 服务器列表结构
type ServerList struct {
	Servers []Server `json:"servers"`
}

// ServerTableDataSource 实现了 tableview.DataSource 接口
type ServerTableDataSource struct {
	servers []Server
}

// ServerTableDelegate 实现了 tableview.Delegate 接口
type ServerTableDelegate struct{}

// NumberOfRows 返回表格行数
func (ds *ServerTableDataSource) NumberOfRows(tv *tableview.TableView) int {
	return len(ds.servers)
}

// CellForRow 返回指定行的单元格
func (ds *ServerTableDataSource) CellForRow(tv *tableview.TableView, row int) *tableview.TableViewCell {
	if row < 0 || row >= len(ds.servers) {
		return nil
	}

	cell := tv.Dequeue("server_cell")
	if cell == nil {
		cell = tableview.NewCell("server_cell")
	}

	return cell
}

// DidSelectRow 处理行选择事件
func (delegate *ServerTableDelegate) DidSelectRow(tv *tableview.TableView, row int) {
	fmt.Printf("选中了第 %d 行\n", row)
}

// RowHeight 返回行高
func (delegate *ServerTableDelegate) RowHeight(tv *tableview.TableView, row int) int {
	return 0 // 使用默认行高
}

func main() {
	// 将当前 goroutine 绑定到当前操作系统线程
	runtime.LockOSThread()

	// 加载服务器数据
	servers, err := loadServers("servers.json")
	if err != nil {
		fmt.Printf("加载服务器数据失败: %v\n", err)
		return
	}

	// 创建主窗口
	win := window.NewUIWindow(800, 600, "表格视图演示")
	root := win.RootView()

	// 标题
	title := label.NewUILabel(&foundation.Rect{X: 20, Y: 20, Width: 760, Height: 40}, "服务器列表")
	title.SetFontSize(24)
	title.SetTextColor(GRAY)
	root.AddSubview(title)

	// 刷新按钮
	refreshBtn := button.NewUIButton(&foundation.Rect{X: 680, Y: 20, Width: 100, Height: 40}, "刷新数据")
	refreshBtn.SetBackgroundColor(BLUE)
	root.AddSubview(refreshBtn)

	// 创建表格视图
	tv, err := tableview.New(20, 70, 760, 500)
	if err != nil {
		fmt.Printf("创建表格失败: %v\n", err)
		return
	}

	// 添加表格到窗口
	win.Raw().Add(tv.Raw().Widget())

	// 设置数据源和代理
	dataSource := &ServerTableDataSource{servers: servers.Servers}
	delegate := &ServerTableDelegate{}
	tv.SetDataSource(dataSource)
	tv.SetDelegate(delegate)

	// 初始日志
	fmt.Println("表格视图演示程序启动成功")
	fmt.Printf("加载服务器数据完成，共 %d 台服务器\n", len(servers.Servers))

	// 刷新按钮事件处理
	refreshBtn.OnTouchUpInside(func() {
		fmt.Println("刷新服务器数据")
		// 模拟刷新数据
		newServers, err := loadServers("servers.json")
		if err == nil {
			dataSource.servers = newServers.Servers
			tv.ReloadData()
			fmt.Printf("刷新完成，共 %d 台服务器\n", len(newServers.Servers))
		} else {
			fmt.Printf("刷新失败: %v\n", err)
		}
	})

	// 显示窗口并运行
	win.Show()
	fltk2go.Run()
}

// loadServers 从JSON文件加载服务器数据
func loadServers(filename string) (*ServerList, error) {
	// 读取文件
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析JSON
	var serverList ServerList
	err = json.Unmarshal(data, &serverList)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return &serverList, nil
}
