package main

import (
	"fmt"
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit"
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

	title := label.NewUILabel(&foundation.Rect{X: 50, Y: 20, Width: 600, Height: 40}, "Slider & Progress Controls")
	title.SetFontSize(22)
	title.SetFont(fltk_bridge.HELVETICA_BOLD)
	title.SetAlignment(fltk_bridge.ALIGN_CENTER)
	root.AddSubview(title)

	volumeLabel := label.NewUILabel(&foundation.Rect{X: 50, Y: 80, Width: 300, Height: 28}, "Volume: 50")
	volumeLabel.SetFontSize(15)
	root.AddSubview(volumeLabel)

	volumeSlider := uikit.NewUISlider(&foundation.Rect{X: 50, Y: 115, Width: 580, Height: 34})
	volumeSlider.SetMinimumValue(0)
	volumeSlider.SetMaximumValue(100)
	volumeSlider.SetValue(50)
	volumeSlider.SetStep(1)
	root.AddSubview(volumeSlider)

	volumeBar := uikit.NewUIProgressView(&foundation.Rect{X: 50, Y: 160, Width: 580, Height: 28})
	volumeBar.SetMinimumValue(0)
	volumeBar.SetMaximumValue(100)
	volumeBar.SetProgress(50)
	volumeBar.SetTrackColor(0xE0E0E000)
	volumeBar.SetProgressTintColor(BLUE)
	root.AddSubview(volumeBar)

	volumeSlider.OnValueChanged(func(val float64) {
		volumeBar.SetProgress(val)
		volumeLabel.SetText(fmt.Sprintf("Volume: %.0f", val))
	})

	brightnessLabel := label.NewUILabel(&foundation.Rect{X: 50, Y: 210, Width: 300, Height: 28}, "Brightness:")
	brightnessLabel.SetFontSize(15)
	root.AddSubview(brightnessLabel)

	brightnessSlider := uikit.NewUISlider(&foundation.Rect{X: 50, Y: 245, Width: 580, Height: 34})
	brightnessSlider.SetMinimumValue(0)
	brightnessSlider.SetMaximumValue(100)
	brightnessSlider.SetValue(75)
	brightnessSlider.SetStep(1)
	root.AddSubview(brightnessSlider)

	brightnessBar := uikit.NewUIProgressView(&foundation.Rect{X: 50, Y: 290, Width: 580, Height: 28})
	brightnessBar.SetMinimumValue(0)
	brightnessBar.SetMaximumValue(100)
	brightnessBar.SetProgress(75)
	brightnessBar.SetTrackColor(0xE0E0E000)
	brightnessBar.SetProgressTintColor(GREEN)
	root.AddSubview(brightnessBar)

	brightnessSlider.OnValueChanged(func(val float64) {
		brightnessBar.SetProgress(val)
		brightnessLabel.SetText(fmt.Sprintf("Brightness: %.0f%%", val))
	})

	resetBtn := button.NewUIButton(&foundation.Rect{X: 50, Y: 360, Width: 120, Height: 36}, "Reset All")
	resetBtn.SetBackgroundColor(RED)
	resetBtn.SetTitleColor(WHITE)
	resetBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(0)
		volumeBar.SetProgress(0)
		volumeLabel.SetText("Volume: 0")
		brightnessSlider.SetValue(0)
		brightnessBar.SetProgress(0)
		brightnessLabel.SetText("Brightness: 0%")
	})
	root.AddSubview(resetBtn)

	halfBtn := button.NewUIButton(&foundation.Rect{X: 200, Y: 360, Width: 120, Height: 36}, "50%")
	halfBtn.SetBackgroundColor(BLUE)
	halfBtn.SetTitleColor(WHITE)
	halfBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(50)
		volumeBar.SetProgress(50)
		volumeLabel.SetText("Volume: 50")
		brightnessSlider.SetValue(50)
		brightnessBar.SetProgress(50)
		brightnessLabel.SetText("Brightness: 50%")
	})
	root.AddSubview(halfBtn)

	maxBtn := button.NewUIButton(&foundation.Rect{X: 350, Y: 360, Width: 120, Height: 36}, "Max")
	maxBtn.SetBackgroundColor(GREEN)
	maxBtn.SetTitleColor(WHITE)
	maxBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(100)
		volumeBar.SetProgress(100)
		volumeLabel.SetText("Volume: 100")
		brightnessSlider.SetValue(100)
		brightnessBar.SetProgress(100)
		brightnessLabel.SetText("Brightness: 100%")
	})
	root.AddSubview(maxBtn)

	win.Show()
	fltk2go.Run()
}
