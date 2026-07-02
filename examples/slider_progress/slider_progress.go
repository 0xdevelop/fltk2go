package slider_progress

import (
	"fmt"

	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit"
	"github.com/0xdevelop/fltk2go/uikit/button"
	"github.com/0xdevelop/fltk2go/uikit/label"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

const (
	SliderBlue       uint = 0x2563EB00
	SliderGreen      uint = 0x22C55E00
	SliderRed        uint = 0xEF444400
	SliderInk        uint = 0x0F172A00
	SliderMuted      uint = 0x64748B00
	SliderPanel      uint = 0xF8FAFC00
	SliderCard       uint = 0xFFFFFF00
	SliderTrack      uint = 0xE2E8F000
	SliderWhite      uint = 0xFFFFFFFF
	SliderButtonGray uint = 0x33415500
)

func sectionTitle(x, y int, title, subtitle string, color uint) (*label.UILabel, *label.UILabel) {
	t := label.NewUILabel(&foundation.Rect{X: x, Y: y, Width: 330, Height: 28}, title)
	t.SetFontSize(17)
	t.SetFont(fltk_bridge.HELVETICA_BOLD)
	t.SetTextColor(SliderInk)
	s := label.NewUILabel(&foundation.Rect{X: x, Y: y + 30, Width: 330, Height: 24}, subtitle)
	s.SetFontSize(13)
	s.SetTextColor(color)
	return t, s
}

func BuildView(parent *view.UIView) view.Viewable {
	background := label.NewUILabel(&foundation.Rect{X: 0, Y: 0, Width: 760, Height: 560}, "")
	background.SetBackgroundColor(SliderPanel)
	background.SetFrame(fltk_bridge.FLAT_BOX)
	parent.AddSubview(background)

	title := label.NewUILabel(&foundation.Rect{X: 48, Y: 26, Width: 664, Height: 34}, "Slider & Progress Controls")
	title.SetFontSize(24)
	title.SetFont(fltk_bridge.HELVETICA_BOLD)
	title.SetTextColor(SliderInk)
	parent.AddSubview(title)

	description := label.NewUILabel(&foundation.Rect{X: 48, Y: 64, Width: 664, Height: 24}, "Spacious control cards with live state feedback and clear action hierarchy.")
	description.SetFontSize(14)
	description.SetTextColor(SliderMuted)
	parent.AddSubview(description)

	volumeCard := label.NewUILabel(&foundation.Rect{X: 48, Y: 116, Width: 664, Height: 142}, "")
	volumeCard.SetBackgroundColor(SliderCard)
	volumeCard.SetFrame(fltk_bridge.ROUNDED_BOX)
	parent.AddSubview(volumeCard)

	volumeTitle, volumeLabel := sectionTitle(78, 142, "Volume", "50", SliderBlue)
	volumeLabel.View().SetAutomationID("slider.volume.label").SetAutomationName("Volume value label")
	parent.AddSubview(volumeTitle)
	parent.AddSubview(volumeLabel)

	volumeSlider := uikit.NewUISlider(&foundation.Rect{X: 78, Y: 196, Width: 604, Height: 34})
	volumeSlider.View().SetAutomationID("slider.volume").SetAutomationName("Volume slider")
	volumeSlider.SetMinimumValue(0)
	volumeSlider.SetMaximumValue(100)
	volumeSlider.SetValue(50)
	volumeSlider.SetStep(1)
	parent.AddSubview(volumeSlider)

	volumeBar := uikit.NewUIProgressView(&foundation.Rect{X: 78, Y: 230, Width: 604, Height: 16})
	volumeBar.View().SetAutomationID("slider.volume.progress").SetAutomationName("Volume progress")
	volumeBar.SetMinimumValue(0)
	volumeBar.SetMaximumValue(100)
	volumeBar.SetProgress(50)
	volumeBar.SetTrackColor(SliderTrack)
	volumeBar.SetProgressTintColor(SliderBlue)
	parent.AddSubview(volumeBar)

	volumeSlider.OnValueChanged(func(val float64) {
		volumeBar.SetProgress(val)
		volumeLabel.SetText(fmt.Sprintf("%.0f", val))
	})

	brightnessCard := label.NewUILabel(&foundation.Rect{X: 48, Y: 278, Width: 664, Height: 142}, "")
	brightnessCard.SetBackgroundColor(SliderCard)
	brightnessCard.SetFrame(fltk_bridge.ROUNDED_BOX)
	parent.AddSubview(brightnessCard)

	brightnessTitle, brightnessLabel := sectionTitle(78, 304, "Brightness", "75%", SliderGreen)
	brightnessLabel.View().SetAutomationID("slider.brightness.label").SetAutomationName("Brightness value label")
	parent.AddSubview(brightnessTitle)
	parent.AddSubview(brightnessLabel)

	brightnessSlider := uikit.NewUISlider(&foundation.Rect{X: 78, Y: 358, Width: 604, Height: 34})
	brightnessSlider.View().SetAutomationID("slider.brightness").SetAutomationName("Brightness slider")
	brightnessSlider.SetMinimumValue(0)
	brightnessSlider.SetMaximumValue(100)
	brightnessSlider.SetValue(75)
	brightnessSlider.SetStep(1)
	parent.AddSubview(brightnessSlider)

	brightnessBar := uikit.NewUIProgressView(&foundation.Rect{X: 78, Y: 392, Width: 604, Height: 16})
	brightnessBar.View().SetAutomationID("slider.brightness.progress").SetAutomationName("Brightness progress")
	brightnessBar.SetMinimumValue(0)
	brightnessBar.SetMaximumValue(100)
	brightnessBar.SetProgress(75)
	brightnessBar.SetTrackColor(SliderTrack)
	brightnessBar.SetProgressTintColor(SliderGreen)
	parent.AddSubview(brightnessBar)

	brightnessSlider.OnValueChanged(func(val float64) {
		brightnessBar.SetProgress(val)
		brightnessLabel.SetText(fmt.Sprintf("%.0f%%", val))
	})

	resetBtn := button.NewUIButton(&foundation.Rect{X: 48, Y: 460, Width: 142, Height: 44}, "Reset")
	resetBtn.View().SetAutomationID("slider.reset").SetAutomationName("Reset sliders")
	resetBtn.SetBackgroundColor(SliderRed)
	resetBtn.SetTitleColor(SliderWhite)
	parent.AddSubview(resetBtn)

	halfBtn := button.NewUIButton(&foundation.Rect{X: 210, Y: 460, Width: 142, Height: 44}, "Set 50%")
	halfBtn.View().SetAutomationID("slider.set_half").SetAutomationName("Set sliders to 50 percent")
	halfBtn.SetBackgroundColor(SliderButtonGray)
	halfBtn.SetTitleColor(SliderWhite)
	parent.AddSubview(halfBtn)

	maxBtn := button.NewUIButton(&foundation.Rect{X: 372, Y: 460, Width: 142, Height: 44}, "Max")
	maxBtn.View().SetAutomationID("slider.max").SetAutomationName("Max sliders")
	maxBtn.SetBackgroundColor(SliderGreen)
	maxBtn.SetTitleColor(SliderWhite)
	parent.AddSubview(maxBtn)

	note := label.NewUILabel(&foundation.Rect{X: 534, Y: 456, Width: 178, Height: 56}, "44px+ targets\n20px button gaps")
	note.SetFontSize(12)
	note.SetTextColor(SliderMuted)
	note.SetAlignment(fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_INSIDE)
	parent.AddSubview(note)

	resetBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(0)
		volumeBar.SetProgress(0)
		volumeLabel.SetText("0")
		brightnessSlider.SetValue(0)
		brightnessBar.SetProgress(0)
		brightnessLabel.SetText("0%")
	})
	halfBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(50)
		volumeBar.SetProgress(50)
		volumeLabel.SetText("50")
		brightnessSlider.SetValue(50)
		brightnessBar.SetProgress(50)
		brightnessLabel.SetText("50%")
	})
	maxBtn.OnTouchUpInside(func() {
		volumeSlider.SetValue(100)
		volumeBar.SetProgress(100)
		volumeLabel.SetText("100")
		brightnessSlider.SetValue(100)
		brightnessBar.SetProgress(100)
		brightnessLabel.SetText("100%")
	})

	return nil
}
