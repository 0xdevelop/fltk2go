package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/fltk_bridge"
)

const (
	winW      = 1040
	winH      = 680
	initSplit = 300

	ink       = fltk_bridge.Color(0x0F172A00)
	slate     = fltk_bridge.Color(0x1E293B00)
	muted     = fltk_bridge.Color(0x33415500)
	panel     = fltk_bridge.Color(0xF8FAFC00)
	panelSoft = fltk_bridge.Color(0xEEF2F700)
	line      = fltk_bridge.Color(0xCBD5E100)
	accent    = fltk_bridge.Color(0x2563EB00)
	white     = fltk_bridge.Color(0xFFFFFF00)
)

type example struct {
	title    string
	subtitle string
	html     string
	dir      string
}

var examples = []example{
	{
		title:    "Counter",
		subtitle: "UIButton + UILabel basics",
		html: `<h2><font color="#0F172A">Counter</font></h2>
<p><font color="#475569">A compact state-management demo with a single primary action and immediate label feedback.</font></p>
<h3>What it shows</h3>
<ul>
<li>UIButton click callback via <code>OnTouchUpInside</code></li>
<li>UILabel text updates with <code>SetText</code></li>
<li><code>runtime.LockOSThread()</code> for GUI thread affinity</li>
</ul>`,
		dir: "./counter",
	},
	{
		title:    "Comprehensive",
		subtitle: "Buttons, state, table rows",
		html: `<h2><font color="#0F172A">Comprehensive</font></h2>
<p><font color="#475569">A polished component gallery with button states, semantic actions, and a custom-drawn table.</font></p>
<h3>What it shows</h3>
<ul>
<li>System / checkbox / radio / toggle buttons</li>
<li>TableView DataSource / Delegate pattern</li>
<li>Dynamic add / delete / reload workflows</li>
<li>Custom row drawing with visual hierarchy</li>
</ul>`,
		dir: "./comprehensive",
	},
	{
		title:    "Input",
		subtitle: "Forms with preview feedback",
		html: `<h2><font color="#0F172A">Input</font></h2>
<p><font color="#475569">Form controls arranged as a readable two-column workflow: inputs on the left, live preview on the right.</font></p>
<h3>What it shows</h3>
<ul>
<li>Text, integer, float, password, and multiline inputs</li>
<li>Readable labels instead of placeholder-only context</li>
<li>Display summary updated from field values</li>
</ul>`,
		dir: "./input",
	},
	{
		title:    "SplitView",
		subtitle: "Resizable master-detail layout",
		html: `<h2><font color="#0F172A">SplitView</font></h2>
<p><font color="#475569">Horizontal and vertical split views with clearer pane boundaries and cross-panel updates.</font></p>
<h3>What it shows</h3>
<ul>
<li>Horizontal / vertical split modes</li>
<li>SetLeftView / SetRightView panel composition</li>
<li>Fixed left pane sizing</li>
<li>Cross-panel selection and detail updates</li>
</ul>`,
		dir: "./splitview",
	},
	{
		title:    "TableView",
		subtitle: "Server list management",
		html: `<h2><font color="#0F172A">TableView</font></h2>
<p><font color="#475569">A server management list with table rows, status semantics, and add/remove actions.</font></p>
<h3>What it shows</h3>
<ul>
<li>DataSource.NumberOfRows / CellForRow</li>
<li>Delegate.DidSelectRow / RowHeight</li>
<li>Custom row drawing</li>
<li>Dynamic server records</li>
</ul>`,
		dir: "./tableview",
	},
	{
		title:    "Slider & Progress",
		subtitle: "Live controls with metrics",
		html: `<h2><font color="#0F172A">Slider &amp; Progress</font></h2>
<p><font color="#475569">Two spacious control cards with real-time progress feedback and one clear action cluster.</font></p>
<h3>What it shows</h3>
<ul>
<li>UIKit-style UISlider and UIProgressView wrappers</li>
<li>Value-changed callbacks</li>
<li>Reset / 50% / Max semantic actions</li>
<li>Readable spacing for desktop and remote sessions</li>
</ul>`,
		dir: "./slider_progress",
	},
	{
		title:    "Tabs",
		subtitle: "Grouped settings panels",
		html: `<h2><font color="#0F172A">Tabs</font></h2>
<p><font color="#475569">Tabbed controls demonstrating FLTK begin/end ownership with settings-like interaction patterns.</font></p>
<h3>What it shows</h3>
<ul>
<li>Tabs + Group begin/end ownership</li>
<li>Choice dropdowns</li>
<li>Spinner and value slider controls</li>
<li>Multi-control output updates</li>
</ul>`,
		dir: "./tabs",
	},
	{
		title:    "TableView Demo",
		subtitle: "JSON-backed server table",
		html: `<h2><font color="#0F172A">TableView Demo</font></h2>
<p><font color="#475569">A complete TableView scenario that loads server data from JSON and refreshes it in-place.</font></p>
<h3>What it shows</h3>
<ul>
<li>JSON parsing from <code>servers.json</code></li>
<li>ServerTableDataSource / ServerTableDelegate</li>
<li>tableview.New list creation</li>
<li>Refresh action with local data reload</li>
</ul>
<p><font color="#64748B"><i>The launcher runs this from its own folder so the JSON file resolves correctly.</i></font></p>`,
		dir: "./tableview_demo",
	},
}

func main() {
	runtime.LockOSThread()

	win := fltk_bridge.NewWindow(winW, winH, "FLTK2Go — Examples Launcher")
	win.SetColor(panel)

	// Tile covers the whole window and provides the draggable split handle.
	tile := fltk_bridge.NewTile(0, 0, winW, winH)
	tile.SetColor(panel)

	// ── Left panel: example list ──────────────────────────────────────────
	leftGrp := fltk_bridge.NewGroup(0, 0, initSplit, winH)
	leftGrp.SetColor(slate)
	leftGrp.SetBox(fltk_bridge.FLAT_BOX)

	brand := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, 0, 0, initSplit, 86, "  FLTK2Go\n  Examples")
	brand.SetColor(ink)
	brand.SetLabelColor(white)
	brand.SetLabelFont(fltk_bridge.HELVETICA_BOLD)
	brand.SetLabelSize(18)
	brand.SetAlign(fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_INSIDE)

	listHint := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, 0, 86, initSplit, 34, "  Select a demo")
	listHint.SetColor(slate)
	listHint.SetLabelColor(fltk_bridge.Color(0xCBD5E100))
	listHint.SetLabelSize(12)
	listHint.SetAlign(fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_INSIDE)

	browser := fltk_bridge.NewHoldBrowser(16, 124, initSplit-32, winH-148)
	browser.SetBox(fltk_bridge.ROUNDED_BOX)
	browser.SetColor(muted)
	browser.SetSelectionColor(accent)
	browser.SetLabelColor(white)
	for i, e := range examples {
		browser.Add(e.title)
		if i == 0 {
			browser.SetValue(1)
		}
	}
	browser.End() // HoldBrowser is a Group subclass; restore leftGrp as current.

	leftGrp.Resizable(browser)
	leftGrp.End()

	// ── Right panel: preview ──────────────────────────────────────────────
	rW := winW - initSplit
	rightGrp := fltk_bridge.NewGroup(initSplit, 0, rW, winH)
	rightGrp.SetColor(panel)
	rightGrp.SetBox(fltk_bridge.FLAT_BOX)

	topPad := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, initSplit, 0, rW, 24, "")
	topPad.SetColor(panel)

	titleBar := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, initSplit+32, 24, rW-64, 72,
		"Select an example")
	titleBar.SetColor(panel)
	titleBar.SetLabelColor(ink)
	titleBar.SetLabelFont(fltk_bridge.HELVETICA_BOLD)
	titleBar.SetLabelSize(22)
	titleBar.SetAlign(fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_INSIDE)

	card := fltk_bridge.NewBox(fltk_bridge.ROUNDED_BOX, initSplit+32, 112, rW-64, winH-210, "")
	card.SetColor(white)

	helpView := fltk_bridge.NewHelpView(initSplit+54, 134, rW-108, winH-254)
	helpView.SetValue(examples[0].html)
	helpView.TextSize(14)

	footer := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, initSplit+32, winH-82, rW-64, 1, "")
	footer.SetColor(line)

	runBtn := fltk_bridge.NewButton(initSplit+32, winH-64, 220, 44,
		"Run selected demo")
	runBtn.SetBox(fltk_bridge.ROUNDED_BOX)
	runBtn.SetColor(accent)
	runBtn.SetLabelColor(white)
	runBtn.SetLabelFont(fltk_bridge.HELVETICA_BOLD)
	runBtn.SetLabelSize(14)

	status := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, initSplit+272, winH-64, rW-304, 44, "Ready — examples run with go run from their own folders")
	status.SetColor(panel)
	status.SetLabelColor(fltk_bridge.Color(0x64748B00))
	status.SetLabelSize(12)
	status.SetAlign(fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_INSIDE)

	rightGrp.Resizable(card)
	rightGrp.End()

	tile.End()

	selectExample := func(idx int) {
		if idx < 0 || idx >= len(examples) {
			return
		}
		e := examples[idx]
		titleBar.SetLabel(e.title)
		titleBar.Redraw()
		helpView.SetValue(e.html)
		status.SetLabel("Ready — " + e.subtitle)
		status.Redraw()
	}
	selectExample(0)

	// ── Callbacks ─────────────────────────────────────────────────────────
	browser.SetCallback(func() {
		selectExample(browser.Value() - 1) // FLTK browser is 1-based
	})

	runBtn.SetCallback(func() {
		idx := browser.Value() - 1
		if idx < 0 || idx >= len(examples) {
			return
		}
		wd, _ := os.Getwd()
		exampleDir := filepath.Join(wd, examples[idx].dir)
		cmd := exec.Command("go", "run", ".")
		cmd.Dir = exampleDir
		_ = cmd.Start()
		status.SetLabel("Launching — " + examples[idx].title)
		status.Redraw()
	})

	win.Show()
	fltk2go.Run()
}
