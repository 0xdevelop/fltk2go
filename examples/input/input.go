package input

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/button"
	uiinput "github.com/0xYeah/fltk2go/uikit/input"
	"github.com/0xYeah/fltk2go/uikit/label"
	"github.com/0xYeah/fltk2go/uikit/view"
)

const (
	InputPrimary     uint = 0x2563EB00
	InputPrimaryDark uint = 0x1E293B00
	InputMutedText   uint = 0x64748B00
	InputPanel       uint = 0xF8FAFC00
	InputCard        uint = 0xFFFFFF00
	InputField       uint = 0xF1F5F900
	InputOrange      uint = 0xF9731600
	InputWhite       uint = 0xFFFFFFFF
)

func fieldLabel(parent *view.UIView, x, y int, text string) {
	l := label.NewUILabel(&foundation.Rect{X: x, Y: y, Width: 300, Height: 20}, text)
	l.SetFontSize(12)
	l.SetTextColor(InputMutedText)
	l.SetFont(fltk_bridge.HELVETICA_BOLD)
	parent.AddSubview(l)
}

func BuildView(parent *view.UIView) view.Viewable {
	background := label.NewUILabel(&foundation.Rect{X: 0, Y: 0, Width: 900, Height: 640}, "")
	background.SetBackgroundColor(InputPanel)
	background.SetFrame(fltk_bridge.FLAT_BOX)
	parent.AddSubview(background)

	title := label.NewUILabel(&foundation.Rect{X: 56, Y: 30, Width: 788, Height: 34}, "Input Playground")
	title.SetFontSize(24)
	title.SetFont(fltk_bridge.HELVETICA_BOLD)
	title.SetTextColor(InputPrimaryDark)
	parent.AddSubview(title)

	description := label.NewUILabel(&foundation.Rect{X: 56, Y: 68, Width: 788, Height: 26}, "Readable form spacing, explicit labels, and immediate preview feedback.")
	description.SetFontSize(14)
	description.SetTextColor(InputMutedText)
	parent.AddSubview(description)

	formCard := label.NewUILabel(&foundation.Rect{X: 56, Y: 120, Width: 390, Height: 420}, "")
	formCard.SetBackgroundColor(InputCard)
	formCard.SetFrame(fltk_bridge.ROUNDED_BOX)
	parent.AddSubview(formCard)

	previewCard := label.NewUILabel(&foundation.Rect{X: 482, Y: 120, Width: 362, Height: 420}, "")
	previewCard.SetBackgroundColor(InputCard)
	previewCard.SetFrame(fltk_bridge.ROUNDED_BOX)
	parent.AddSubview(previewCard)

	formTitle := label.NewUILabel(&foundation.Rect{X: 84, Y: 146, Width: 320, Height: 28}, "Form fields")
	formTitle.SetFontSize(17)
	formTitle.SetFont(fltk_bridge.HELVETICA_BOLD)
	formTitle.SetTextColor(InputPrimaryDark)
	parent.AddSubview(formTitle)

	fieldLabel(parent, 84, 190, "TEXT")
	textInput := uiinput.New(84, 212, 320, 38, "")
	textInput.SetFontSize(14)
	textInput.SetBackgroundColor(InputField)
	parent.AddSubview(textInput)

	fieldLabel(parent, 84, 264, "INTEGER")
	intInput := uiinput.NewWithType(84, 286, 320, 38, "", uiinput.IntInput)
	intInput.SetFontSize(14)
	intInput.SetBackgroundColor(InputField)
	parent.AddSubview(intInput)

	fieldLabel(parent, 84, 338, "FLOAT")
	floatInput := uiinput.NewWithType(84, 360, 320, 38, "", uiinput.FloatInput)
	floatInput.SetFontSize(14)
	floatInput.SetBackgroundColor(InputField)
	parent.AddSubview(floatInput)

	fieldLabel(parent, 84, 412, "PASSWORD")
	secretInput := uiinput.New(84, 434, 320, 38, "")
	secretInput.SetFontSize(14)
	secretInput.SetBackgroundColor(InputField)
	parent.AddSubview(secretInput)

	fieldLabel(parent, 510, 190, "MULTILINE NOTE")
	multilineInput := uiinput.New(510, 212, 296, 94, "")
	multilineInput.SetFontSize(14)
	multilineInput.SetBackgroundColor(InputField)
	parent.AddSubview(multilineInput)

	previewTitle := label.NewUILabel(&foundation.Rect{X: 510, Y: 146, Width: 296, Height: 28}, "Live preview")
	previewTitle.SetFontSize(17)
	previewTitle.SetFont(fltk_bridge.HELVETICA_BOLD)
	previewTitle.SetTextColor(InputPrimaryDark)
	parent.AddSubview(previewTitle)

	displayLabel := label.NewUILabel(&foundation.Rect{X: 510, Y: 332, Width: 296, Height: 144}, "Input values will appear here")
	displayLabel.SetFontSize(14)
	displayLabel.SetAlignment(fltk_bridge.ALIGN_LEFT | fltk_bridge.ALIGN_INSIDE)
	displayLabel.SetBackgroundColor(0xF8FAFC00)
	displayLabel.SetFrame(fltk_bridge.ROUNDED_BOX)
	displayLabel.SetTextColor(InputPrimaryDark)
	parent.AddSubview(displayLabel)

	updateBtn := button.NewUIButton(&foundation.Rect{X: 84, Y: 560, Width: 170, Height: 44}, "Update preview")
	updateBtn.SetBackgroundColor(InputPrimary)
	updateBtn.SetTitleColor(InputWhite)
	parent.AddSubview(updateBtn)

	clearBtn := button.NewUIButton(&foundation.Rect{X: 274, Y: 560, Width: 132, Height: 44}, "Clear")
	clearBtn.SetBackgroundColor(InputOrange)
	clearBtn.SetTitleColor(InputWhite)
	parent.AddSubview(clearBtn)

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

	return nil
}
