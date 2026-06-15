package main

import (
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	"github.com/0xYeah/fltk2go/uikit/input"
	"github.com/0xYeah/fltk2go/uikit/label"
	"github.com/0xYeah/fltk2go/uikit/view"
	"github.com/0xYeah/fltk2go/uikit/window"
)

const (
	Primary     uint = 0x2563EB00
	PrimaryDark uint = 0x1E293B00
	MutedText   uint = 0x64748B00
	Panel       uint = 0xF8FAFC00
	Card        uint = 0xFFFFFF00
	Field       uint = 0xF1F5F900
	Border      uint = 0xCBD5E100
	Orange      uint = 0xF9731600
	White       uint = 0xFFFFFFFF
)

func fieldLabel(root *view.UIView, x, y int, text string) {
	l := label.NewUILabel(&foundation.Rect{X: x, Y: y, Width: 300, Height: 20}, text)
	l.SetFontSize(12)
	l.SetTextColor(MutedText)
	l.SetFont(fltk_bridge.HELVETICA_BOLD)
	root.AddSubview(l)
}

func main() {
	runtime.LockOSThread()

	win := window.NewUIWindow(900, 640, "Input Example")
	root := win.RootView()

	background := label.NewUILabel(&foundation.Rect{X: 0, Y: 0, Width: 900, Height: 640}, "")
	background.SetBackgroundColor(Panel)
	background.SetFrame(fltk_bridge.FLAT_BOX)
	root.AddSubview(background)

	title := label.NewUILabel(&foundation.Rect{X: 56, Y: 30, Width: 788, Height: 34}, "Input Playground")
	title.SetFontSize(24)
	title.SetFont(fltk_bridge.HELVETICA_BOLD)
	title.SetTextColor(PrimaryDark)
	root.AddSubview(title)

	description := label.NewUILabel(&foundation.Rect{X: 56, Y: 68, Width: 788, Height: 26}, "Readable form spacing, explicit labels, and immediate preview feedback.")
	description.SetFontSize(14)
	description.SetTextColor(MutedText)
	root.AddSubview(description)

	formCard := label.NewUILabel(&foundation.Rect{X: 56, Y: 120, Width: 390, Height: 420}, "")
	formCard.SetBackgroundColor(Card)
	formCard.SetFrame(fltk_bridge.ROUNDED_BOX)
	root.AddSubview(formCard)

	previewCard := label.NewUILabel(&foundation.Rect{X: 482, Y: 120, Width: 362, Height: 420}, "")
	previewCard.SetBackgroundColor(Card)
	previewCard.SetFrame(fltk_bridge.ROUNDED_BOX)
	root.AddSubview(previewCard)

	formTitle := label.NewUILabel(&foundation.Rect{X: 84, Y: 146, Width: 320, Height: 28}, "Form fields")
	formTitle.SetFontSize(17)
	formTitle.SetFont(fltk_bridge.HELVETICA_BOLD)
	formTitle.SetTextColor(PrimaryDark)
	root.AddSubview(formTitle)

	fieldLabel(root, 84, 190, "TEXT")
	textInput := input.New(84, 212, 320, 38, "")
	textInput.SetFontSize(14)
	textInput.SetBackgroundColor(Field)
	root.AddSubview(textInput)

	fieldLabel(root, 84, 264, "INTEGER")
	intInput := input.NewWithType(84, 286, 320, 38, "", input.IntInput)
	intInput.SetFontSize(14)
	intInput.SetBackgroundColor(Field)
	root.AddSubview(intInput)

	fieldLabel(root, 84, 338, "FLOAT")
	floatInput := input.NewWithType(84, 360, 320, 38, "", input.FloatInput)
	floatInput.SetFontSize(14)
	floatInput.SetBackgroundColor(Field)
	root.AddSubview(floatInput)

	fieldLabel(root, 84, 412, "PASSWORD")
	secretInput := input.New(84, 434, 320, 38, "")
	secretInput.SetFontSize(14)
	secretInput.SetBackgroundColor(Field)
	root.AddSubview(secretInput)

	fieldLabel(root, 510, 190, "MULTILINE NOTE")
	multilineInput := input.New(510, 212, 296, 94, "")
	multilineInput.SetFontSize(14)
	multilineInput.SetBackgroundColor(Field)
	root.AddSubview(multilineInput)

	previewTitle := label.NewUILabel(&foundation.Rect{X: 510, Y: 146, Width: 296, Height: 28}, "Live preview")
	previewTitle.SetFontSize(17)
	previewTitle.SetFont(fltk_bridge.HELVETICA_BOLD)
	previewTitle.SetTextColor(PrimaryDark)
	root.AddSubview(previewTitle)

	displayLabel := label.NewUILabel(&foundation.Rect{X: 510, Y: 332, Width: 296, Height: 144}, "Input values will appear here")
	displayLabel.SetFontSize(14)
	displayLabel.SetAlignment(fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_INSIDE)
	displayLabel.SetBackgroundColor(0xF8FAFC00)
	displayLabel.SetFrame(fltk_bridge.ROUNDED_BOX)
	displayLabel.SetTextColor(PrimaryDark)
	root.AddSubview(displayLabel)

	updateBtn := button.NewUIButton(&foundation.Rect{X: 84, Y: 560, Width: 170, Height: 44}, "Update preview")
	updateBtn.SetBackgroundColor(Primary)
	updateBtn.SetTitleColor(White)
	root.AddSubview(updateBtn)

	clearBtn := button.NewUIButton(&foundation.Rect{X: 274, Y: 560, Width: 132, Height: 44}, "Clear")
	clearBtn.SetBackgroundColor(Orange)
	clearBtn.SetTitleColor(White)
	root.AddSubview(clearBtn)

	updateBtn.OnTouchUpInside(func() {
		displayText := "Text: " + textInput.Text() + "\n\n"
		displayText += "Integer: " + intInput.Text() + "\n"
		displayText += "Float: " + floatInput.Text() + "\n"
		displayText += "Password: " + secretInput.Text() + "\n\n"
		displayText += "Note: " + multilineInput.Text()
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
