package main

import (
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/input"
	"github.com/0xYeah/fltk2go/uikit/label"
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

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(800, 600, "Input Example")
	root := win.RootView()

	title := label.NewUILabel(&foundation.Rect{X: 50, Y: 20, Width: 700, Height: 40}, "Input Example")
	title.SetFontSize(24)
	title.SetFont(fltk_bridge.HELVETICA_BOLD)
	title.SetAlignment(fltk_bridge.ALIGN_CENTER)
	root.AddSubview(title)

	textInput := input.New(100, 100, 300, 36, "Enter text here...")
	textInput.SetFontSize(14)
	root.AddSubview(textInput)

	intInput := input.NewWithType(100, 150, 300, 36, "Enter integer here...", input.IntInput)
	intInput.SetFontSize(14)
	root.AddSubview(intInput)

	floatInput := input.NewWithType(100, 200, 300, 36, "Enter float here...", input.FloatInput)
	floatInput.SetFontSize(14)
	root.AddSubview(floatInput)

	secretInput := input.New(100, 250, 300, 36, "Enter password here...")
	secretInput.SetFontSize(14)
	root.AddSubview(secretInput)

	multilineInput := input.New(100, 300, 300, 100, "Enter multiline text here...")
	multilineInput.SetFontSize(14)
	root.AddSubview(multilineInput)

	displayLabel := label.NewUILabel(&foundation.Rect{X: 450, Y: 100, Width: 250, Height: 200}, "Input values will appear here")
	displayLabel.SetFontSize(14)
	displayLabel.SetAlignment(fltk_bridge.ALIGN_LEFT)
	displayLabel.SetBackgroundColor(0xF5F5F500)
	displayLabel.SetFrame(fltk_bridge.ENGRAVED_BOX)
	root.AddSubview(displayLabel)

	updateBtn := button.NewUIButton(&foundation.Rect{X: 100, Y: 420, Width: 120, Height: 36}, "Update Display")
	updateBtn.SetBackgroundColor(BLUE)
	updateBtn.SetTitleColor(WHITE)
	root.AddSubview(updateBtn)

	clearBtn := button.NewUIButton(&foundation.Rect{X: 250, Y: 420, Width: 120, Height: 36}, "Clear All")
	clearBtn.SetBackgroundColor(ORANGE)
	clearBtn.SetTitleColor(WHITE)
	root.AddSubview(clearBtn)

	updateBtn.OnTouchUpInside(func() {
		displayText := "Text Input: " + textInput.Text() + "\n"
		displayText += "Int Input: " + intInput.Text() + "\n"
		displayText += "Float Input: " + floatInput.Text() + "\n"
		displayText += "Password Input: " + secretInput.Text() + "\n"
		displayText += "Multiline Input: " + multilineInput.Text()
		displayLabel.SetText(displayText)
	})

	clearBtn.OnTouchUpInside(func() {
		textInput.SetText("")
		intInput.SetText("")
		floatInput.SetText("")
		secretInput.SetText("")
		multilineInput.SetText("")
		displayLabel.SetText("Input values will appear here")
	})

	win.Show()
	fltk2go.Run()
}
