package main

import (
	"fmt"
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/label"
	"github.com/0xYeah/fltk2go/uikit/window"
)

const (
	BLUE   uint = 0x42A5F500
	GRAY   uint = 0x75757500
	GREEN  uint = 0x4CAF5000
	RED    uint = 0xF4433600
	WHITE  uint = 0xFFFFFFFF
	YELLOW uint = 0xFFC10700
	PURPLE uint = 0x9C27B000
)

type colorEntry struct {
	name  string
	color uint
}

var colorPalette = []colorEntry{
	{"Red", RED},
	{"Green", GREEN},
	{"Blue", BLUE},
	{"Yellow", YELLOW},
	{"Purple", PURPLE},
}

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(800, 520, "Tabs Example")
	root := win.RootView()

	// Title must be created before Tabs so it auto-adds to Window.
	title := label.NewUILabel(&foundation.Rect{X: 50, Y: 12, Width: 700, Height: 36}, "FLTK2Go Tabs Example")
	title.SetFontSize(22)
	title.SetFont(fltk_bridge.HELVETICA_BOLD)
	title.SetAlignment(fltk_bridge.ALIGN_CENTER)
	root.AddSubview(title)

	// Tabs constructor calls begin(); subsequent widgets auto-add to Tabs.
	tabs := fltk_bridge.NewTabs(20, 58, 760, 420)

	// ── Tab 1: Color Picker ──────────────────────────────────────
	group1 := fltk_bridge.NewGroup(20, 88, 760, 390, "  Colors  ")

	previewBox := fltk_bridge.NewBox(fltk_bridge.FLAT_BOX, 310, 120, 160, 100)
	previewBox.SetColor(fltk_bridge.Color(RED))

	colorNameBox := fltk_bridge.NewBox(fltk_bridge.NO_BOX, 490, 148, 230, 28, "Red")
	colorNameBox.SetLabelSize(15)
	colorNameBox.SetLabelFont(fltk_bridge.HELVETICA_BOLD)

	fltk_bridge.NewBox(fltk_bridge.NO_BOX, 50, 118, 220, 28, "Choose Color:")

	choice := fltk_bridge.NewChoice(50, 152, 220, 32)
	for _, c := range colorPalette {
		c := c
		choice.Add(c.name, func() {
			previewBox.SetColor(fltk_bridge.Color(c.color))
			previewBox.Redraw()
			colorNameBox.SetLabel(c.name)
			colorNameBox.Redraw()
		})
	}
	choice.SetValue(0)

	descBox := fltk_bridge.NewBox(fltk_bridge.NO_BOX, 50, 260, 680, 28,
		"Select a color from the dropdown to see the preview.")
	descBox.SetLabelSize(13)
	descBox.SetLabelColor(fltk_bridge.Color(GRAY))

	group1.End()

	// ── Tab 2: Number Controls ───────────────────────────────────
	group2 := fltk_bridge.NewGroup(20, 88, 760, 390, "  Numbers  ")

	fltk_bridge.NewBox(fltk_bridge.NO_BOX, 50, 118, 150, 28, "Count:")
	spinner := fltk_bridge.NewSpinner(210, 115, 130, 32)
	spinner.SetMinimum(1)
	spinner.SetMaximum(999)
	spinner.SetValue(10)
	spinner.SetStep(1)
	spinner.SetType(fltk_bridge.SPINNER_INT_INPUT)

	fltk_bridge.NewBox(fltk_bridge.NO_BOX, 50, 168, 150, 28, "Scale:")
	scaleSlider := fltk_bridge.NewValueSlider(210, 165, 530, 32, "Scale")
	scaleSlider.SetType(fltk_bridge.HOR_NICE_SLIDER)
	scaleSlider.SetMinimum(0.1)
	scaleSlider.SetMaximum(5.0)
	scaleSlider.SetValue(1.0)
	scaleSlider.SetStep(0.1)
	scaleSlider.SetTextSize(13)
	scaleSlider.SetCallbackCondition(fltk_bridge.WhenChanged)

	resultBox := fltk_bridge.NewBox(fltk_bridge.ENGRAVED_BOX, 50, 222, 680, 60,
		"Count: 10   Scale: 1.0x   Result: 10.0")
	resultBox.SetLabelSize(14)
	resultBox.SetColor(fltk_bridge.Color(0xF5F5F500))

	updateResult := func() {
		count := int(spinner.Value())
		scale := scaleSlider.Value()
		result := float64(count) * scale
		resultBox.SetLabel(fmt.Sprintf("Count: %d   Scale: %.1fx   Result: %.1f", count, scale, result))
		resultBox.Redraw()
	}
	spinner.SetCallback(updateResult)
	scaleSlider.SetCallback(updateResult)

	group2.End()

	// ── Tab 3: About ─────────────────────────────────────────────
	group3 := fltk_bridge.NewGroup(20, 88, 760, 390, "  About  ")

	appNameBox := fltk_bridge.NewBox(fltk_bridge.NO_BOX, 50, 120, 680, 40, "FLTK2Go")
	appNameBox.SetLabelSize(28)
	appNameBox.SetLabelFont(fltk_bridge.HELVETICA_BOLD)
	appNameBox.SetAlign(fltk_bridge.ALIGN_CENTER)

	versionBox := fltk_bridge.NewBox(fltk_bridge.NO_BOX, 50, 168, 680, 28, "Version 0.0.1")
	versionBox.SetLabelSize(15)
	versionBox.SetLabelColor(fltk_bridge.Color(GRAY))
	versionBox.SetAlign(fltk_bridge.ALIGN_CENTER)

	descAboutBox := fltk_bridge.NewBox(fltk_bridge.NO_BOX, 50, 210, 680, 56,
		"A Go binding for the FLTK GUI toolkit.\nSimple, fast, and cross-platform.")
	descAboutBox.SetLabelSize(14)
	descAboutBox.SetAlign(fltk_bridge.ALIGN_CENTER | fltk_bridge.ALIGN_INSIDE)

	repoBox := fltk_bridge.NewBox(fltk_bridge.NO_BOX, 50, 278, 680, 28,
		"github.com/0xYeah/fltk2go")
	repoBox.SetLabelSize(13)
	repoBox.SetLabelColor(fltk_bridge.Color(BLUE))
	repoBox.SetAlign(fltk_bridge.ALIGN_CENTER)

	group3.End()

	tabs.End()

	_, _, _, _, _ = descBox, appNameBox, versionBox, descAboutBox, repoBox

	win.Show()
	fltk2go.Run()
}
