package main

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/fltk_bridge"
)

const (
	winW      = 940
	winH      = 620
	initSplit = 240
)

type example struct {
	title string
	html  string
	dir   string
}

var examples = []example{
	{
		title: "Counter",
		html: `<h2>Counter</h2>
<p>基础按钮点击计数器。</p>
<b>演示要点：</b>
<ul>
<li>UIButton 点击回调 OnTouchUpInside</li>
<li>UILabel 文本动态更新 SetText</li>
<li>runtime.LockOSThread() 线程绑定</li>
</ul>`,
		dir: "./counter",
	},
	{
		title: "Comprehensive",
		html: `<h2>Comprehensive</h2>
<p>综合控件示例，展示多种按钮类型与 TableView 联动。</p>
<b>演示要点：</b>
<ul>
<li>SystemButton / CheckboxButton / RadioButton / ToggleButton</li>
<li>TableView DataSource / Delegate 模式</li>
<li>动态添加 / 删除行，ReloadData</li>
<li>SetCustomDraw 自定义行绘制</li>
</ul>`,
		dir: "./comprehensive",
	},
	{
		title: "Input",
		html: `<h2>Input</h2>
<p>各类输入框示例。</p>
<b>演示要点：</b>
<ul>
<li>普通文本输入 input.New</li>
<li>整数输入 input.IntInput</li>
<li>浮点输入 input.FloatInput</li>
<li>密码框 / 多行文本框</li>
<li>Text() 读取输入内容</li>
</ul>`,
		dir: "./input",
	},
	{
		title: "SplitView",
		html: `<h2>SplitView</h2>
<p>水平与垂直分割视图布局示例。</p>
<b>演示要点：</b>
<ul>
<li>splitview.New 水平/垂直模式</li>
<li>SetLeftView / SetRightView 设置子面板</li>
<li>SetLeftViewFixed 固定左侧宽度</li>
<li>跨面板联动更新</li>
</ul>`,
		dir: "./splitview",
	},
	{
		title: "TableView",
		html: `<h2>TableView</h2>
<p>服务器管理列表，演示 DataSource / Delegate 模式。</p>
<b>演示要点：</b>
<ul>
<li>DataSource.NumberOfRows / CellForRow</li>
<li>Delegate.DidSelectRow / RowHeight</li>
<li>SetCustomDraw 自定义行绘制</li>
<li>动态添加 / 删除服务器记录</li>
</ul>`,
		dir: "./tableview",
	},
	{
		title: "Slider &amp; Progress",
		html: `<h2>Slider &amp; Progress</h2>
<p>滑块与进度条实时联动控制。</p>
<b>演示要点：</b>
<ul>
<li>fltk_bridge.NewSlider 水平滑块</li>
<li>fltk_bridge.NewValueSlider 带数值显示</li>
<li>fltk_bridge.NewProgress 进度条</li>
<li>WhenChanged 实时回调</li>
<li>Reset / 50% / Max 三挡按钮</li>
</ul>`,
		dir: "./slider_progress",
	},
	{
		title: "Tabs",
		html: `<h2>Tabs</h2>
<p>选项卡容器示例，展示 FLTK begin/end 自动归属机制。</p>
<b>演示要点：</b>
<ul>
<li>fltk_bridge.NewTabs + NewGroup begin/end</li>
<li>Choice 下拉颜色选择器</li>
<li>Spinner 整数调节</li>
<li>ValueSlider 浮点调节</li>
<li>多控件联动更新显示</li>
</ul>`,
		dir: "./tabs",
	},
	{
		title: "TableView Demo",
		html: `<h2>TableView Demo</h2>
<p>从 JSON 文件加载服务器数据的完整 TableView 演示。</p>
<b>演示要点：</b>
<ul>
<li>encoding/json 解析 servers.json</li>
<li>ServerTableDataSource / ServerTableDelegate</li>
<li>tableview.New 创建列表</li>
<li>刷新按钮重新加载数据</li>
</ul>
<p><i>注：需在 tableview_demo/ 目录下运行，JSON 文件与可执行文件同级。</i></p>`,
		dir: "./tableview_demo",
	},
}

func main() {
	runtime.LockOSThread()

	win := fltk_bridge.NewWindow(winW, winH, "FLTK2Go — Examples Launcher")

	// Tile covers the whole window and provides the draggable split handle.
	tile := fltk_bridge.NewTile(0, 0, winW, winH)

	// ── Left panel: example list ──────────────────────────────────────────
	leftGrp := fltk_bridge.NewGroup(0, 0, initSplit, winH)

	listHdr := fltk_bridge.NewBox(fltk_bridge.UP_BOX, 0, 0, initSplit, 32, "Examples")
	listHdr.SetLabelFont(fltk_bridge.HELVETICA_BOLD)
	listHdr.SetLabelSize(13)

	browser := fltk_bridge.NewHoldBrowser(0, 32, initSplit, winH-32)
	for _, e := range examples {
		browser.Add(e.title)
	}
	browser.End() // HoldBrowser is a Group subclass; restore leftGrp as current.

	leftGrp.Resizable(browser)
	leftGrp.End()

	// ── Right panel: preview ──────────────────────────────────────────────
	rW := winW - initSplit
	rightGrp := fltk_bridge.NewGroup(initSplit, 0, rW, winH)

	titleBar := fltk_bridge.NewBox(fltk_bridge.UP_BOX, initSplit, 0, rW, 44,
		"  Select an example")
	titleBar.SetLabelFont(fltk_bridge.HELVETICA_BOLD)
	titleBar.SetLabelSize(15)
	titleBar.SetAlign(fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_INSIDE)

	helpView := fltk_bridge.NewHelpView(initSplit, 44, rW, winH-44-58)
	helpView.SetValue(
		`<p><font color="#999999">← 从左侧列表选择一个示例，查看说明并运行。</font></p>`)
	helpView.TextSize(13)

	runBtn := fltk_bridge.NewButton(initSplit+10, winH-48, 240, 38,
		"Run Selected Example")
	runBtn.SetColor(fltk_bridge.Color(0x42A5F500))
	runBtn.SetLabelColor(fltk_bridge.Color(0xFFFFFF00))
	runBtn.SetLabelFont(fltk_bridge.HELVETICA_BOLD)
	runBtn.SetLabelSize(13)

	rightGrp.Resizable(helpView)
	rightGrp.End()

	tile.End()

	// ── Callbacks ─────────────────────────────────────────────────────────
	browser.SetCallback(func() {
		idx := browser.Value() - 1 // FLTK browser is 1-based
		if idx < 0 || idx >= len(examples) {
			return
		}
		e := examples[idx]
		titleBar.SetLabel("  " + e.title)
		titleBar.Redraw()
		helpView.SetValue(e.html)
	})

	runBtn.SetCallback(func() {
		idx := browser.Value() - 1
		if idx < 0 || idx >= len(examples) {
			return
		}
		wd, _ := os.Getwd()
		cmd := exec.Command("go", "run", examples[idx].dir)
		cmd.Dir = wd
		_ = cmd.Start()
	})

	win.Show()
	fltk2go.Run()
}
