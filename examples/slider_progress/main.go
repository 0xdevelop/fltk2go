package main

import (
	"fmt"
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/label"
	"github.com/0xYeah/fltk2go/uikit/window"
)

const (
	BLUE  uint = 0x42A5F500
	GRAY  uint = 0x75757500
	GREEN uint = 0x4CAF5000
	RED   uint = 0xF4433600
	WHITE uint = 0xFFFFFFFF
)

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(700, 460, "Slider & Progress Example")
	root := win.RootView()
	rawWin := win.Raw()

	title := label.NewUILabel(&foundation.Rect{X: 50, Y: 20, Width: 600, Height: 40}, "Slider & Progress Controls")
	title.SetFontSize(22)
	title.SetFont(fltk_bridge.HELVETICA_BOLD)
	title.SetAlignment(fltk_bridge.ALIGN_CENTER)
	root.AddSubview(title)

	volumeLabel := label.NewUILabel(&foundation.Rect{X: 50, Y: 80, Width: 300, Height: 28}, "Volume: 50")
	volumeLabel.SetFontSize(15)
	root.AddSubview(volumeLabel)

	volumeSlider := fltk_bridge.NewSlider(50, 115, 580, 34, "Volume")
	volumeSlider.SetType(fltk_bridge.HOR_NICE_SLIDER)
	volumeSlider.SetMinimum(0)
	volumeSlider.SetMaximum(100)
	volumeSlider.SetValue(50)
	volumeSlider.SetStep(1)
	volumeSlider.SetCallbackCondition(fltk_bridge.WhenChanged)
	rawWin.Add(volumeSlider)

	volumeBar := fltk_bridge.NewProgress(50, 160, 580, 28, "")
	volumeBar.SetMinimum(0)
	volumeBar.SetMaximum(100)
	volumeBar.SetValue(50)
	volumeBar.SetColor(fltk_bridge.Color(0xE0E0E000))
	volumeBar.SetSelectionColor(fltk_bridge.Color(BLUE))
	rawWin.Add(volumeBar)

	volumeSlider.SetCallback(func() {
		val := volumeSlider.Value()
		volumeBar.SetValue(val)
		volumeLabel.SetText(fmt.Sprintf("Volume: %.0f", val))
		volumeBar.Redraw()
	})

	brightnessLabel := label.NewUILabel(&foundation.Rect{X: 50, Y: 210, Width: 300, Height: 28}, "Brightness:")
	brightnessLabel.SetFontSize(15)
	root.AddSubview(brightnessLabel)

	brightnessSlider := fltk_bridge.NewValueSlider(50, 245, 580, 34, "Brightness")
	brightnessSlider.SetType(fltk_bridge.HOR_NICE_SLIDER)
	brightnessSlider.SetMinimum(0)
	brightnessSlider.SetMaximum(100)
	brightnessSlider.SetValue(75)
	brightnessSlider.SetStep(1)
	brightnessSlider.SetTextSize(13)
	rawWin.Add(brightnessSlider)

	brightnessBar := fltk_bridge.NewProgress(50, 290, 580, 28, "")
	brightnessBar.SetMinimum(0)
	brightnessBar.SetMaximum(100)
	brightnessBar.SetValue(75)
	brightnessBar.SetColor(fltk_bridge.Color(0xE0E0E000))
	brightnessBar.SetSelectionColor(fltk_bridge.Color(GREEN))
	rawWin.Add(brightnessBar)

	brightnessSlider.SetCallbackCondition(fltk_bridge.WhenChanged)
	brightnessSlider.SetCallback(func() {
		val := brightnessSlider.Value()
		brightnessBar.SetValue(val)
		brightnessLabel.SetText(fmt.Sprintf("Brightness: %.0f%%", val))
		brightnessBar.Redraw()
	})

	resetBtn := button.NewUIButton(&foundation.Rect{X: 50, Y: 360, Width: 120, Height: 36}, "Reset All")
	resetBtn.SetBackgroundColor(RED)
	resetBtn.SetTitleColor(WHITE)
	resetBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(0)
		volumeBar.SetValue(0)
		volumeLabel.SetText("Volume: 0")
		brightnessSlider.SetValue(0)
		brightnessBar.SetValue(0)
		brightnessLabel.SetText("Brightness: 0%")
		volumeSlider.Redraw()
		volumeBar.Redraw()
		brightnessSlider.Redraw()
		brightnessBar.Redraw()
	})
	root.AddSubview(resetBtn)

	halfBtn := button.NewUIButton(&foundation.Rect{X: 200, Y: 360, Width: 120, Height: 36}, "50%")
	halfBtn.SetBackgroundColor(BLUE)
	halfBtn.SetTitleColor(WHITE)
	halfBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(50)
		volumeBar.SetValue(50)
		volumeLabel.SetText("Volume: 50")
		brightnessSlider.SetValue(50)
		brightnessBar.SetValue(50)
		brightnessLabel.SetText("Brightness: 50%")
		volumeSlider.Redraw()
		volumeBar.Redraw()
		brightnessSlider.Redraw()
		brightnessBar.Redraw()
	})
	root.AddSubview(halfBtn)

	maxBtn := button.NewUIButton(&foundation.Rect{X: 350, Y: 360, Width: 120, Height: 36}, "Max")
	maxBtn.SetBackgroundColor(GREEN)
	maxBtn.SetTitleColor(WHITE)
	maxBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(100)
		volumeBar.SetValue(100)
		volumeLabel.SetText("Volume: 100")
		brightnessSlider.SetValue(100)
		brightnessBar.SetValue(100)
		brightnessLabel.SetText("Brightness: 100%")
		volumeSlider.Redraw()
		volumeBar.Redraw()
		brightnessSlider.Redraw()
		brightnessBar.Redraw()
	})
	root.AddSubview(maxBtn)

	win.Show()
	fltk2go.Run()
}
