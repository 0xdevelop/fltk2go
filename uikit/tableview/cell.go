package tableview

import (
	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/uikit/textlayout"
)

// TableViewCell：UITableViewCell 的最小抽象
type TableViewCell struct {
	ReuseID string

	row int
	col int

	Text      string
	Font      fltk_bridge.Font
	FontSize  int
	TextColor fltk_bridge.Color
	Align     fltk_bridge.Align
	Selected  bool

	preparedText *textlayout.PreparedText
}

func NewCell(reuseID string) *TableViewCell {
	return &TableViewCell{
		ReuseID:   reuseID,
		Font:      fltk_bridge.HELVETICA,
		FontSize:  14,
		TextColor: fltk_bridge.Color(0),
		Align:     fltk_bridge.ALIGN_CENTER | fltk_bridge.ALIGN_CLIP,
	}
}

// PrepareForReuse：复用前清理状态
func (c *TableViewCell) PrepareForReuse() {
	c.Text = ""
	c.preparedText = nil
	c.Selected = false
}

func (c *TableViewCell) SetText(text string) {
	c.Text = text
	if text != "" {
		c.preparedText = textlayout.Prepare(text, c.Font, c.FontSize)
	} else {
		c.preparedText = nil
	}
}

func (c *TableViewCell) Row() int { return c.row }
func (c *TableViewCell) Col() int { return c.col }
