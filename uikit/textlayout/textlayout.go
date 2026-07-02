package textlayout

import (
	"fmt"
	"strings"
	"sync"

	"github.com/0xdevelop/fltk2go/fltk_bridge"
)

// SegmentType identifies the type of a text segment.
type SegmentType int

const (
	WordSegment SegmentType = iota
	SpaceSegment
	BreakSegment
)

// Segment holds cached measurement data for a piece of text.
type Segment struct {
	Text  string
	Width int
	Type  SegmentType
}

// PreparedText represents a string that has been split and measured,
// allowing for extremely fast layout calculations without re-measuring.
// Inspired by pretext (https://github.com/chenglou/pretext).
type PreparedText struct {
	Segments []Segment
	Font     fltk_bridge.Font
	FontSize int
}

// LayoutLine represents a single computed line of text within a bounding box.
type LayoutLine struct {
	Text  string
	Width int
	Start int
	End   int // Exclusive
}

type LayoutResult struct {
	Height    int
	LineCount int
	Lines     []LayoutLine
}

var (
	measureCache = sync.Map{}
)

func getMeasure(text string, font fltk_bridge.Font, fontSize int) int {
	key := fmt.Sprintf("%d_%d_%s", font, fontSize, text)
	if val, ok := measureCache.Load(key); ok {
		return val.(int)
	}

	// fltk_bridge.SetDrawFont is called once in Prepare,
	// assuming it stays active during the loop
	w, _ := fltk_bridge.MeasureText(text, false)
	measureCache.Store(key, w)
	return w
}

// Prepare normalizes and breaks text into measurable segments (words, spaces, breaks),
// then measures each segment using FLTK's native font engine.
// This is the "expensive" step that should be cached.
func Prepare(text string, font fltk_bridge.Font, fontSize int) *PreparedText {
	fltk_bridge.SetDrawFont(font, fontSize)

	segments := []Segment{}
	var curWord strings.Builder

	// A simple tokenizer: breaks by spaces and newlines.
	// For production, this could be enhanced to support CJK, punctuation breaks, etc.
	for _, r := range text {
		if r == ' ' {
			if curWord.Len() > 0 {
				str := curWord.String()
				w := getMeasure(str, font, fontSize)
				segments = append(segments, Segment{Text: str, Width: w, Type: WordSegment})
				curWord.Reset()
			}
			w := getMeasure(" ", font, fontSize)
			segments = append(segments, Segment{Text: " ", Width: w, Type: SpaceSegment})
		} else if r == '\n' {
			if curWord.Len() > 0 {
				str := curWord.String()
				w := getMeasure(str, font, fontSize)
				segments = append(segments, Segment{Text: str, Width: w, Type: WordSegment})
				curWord.Reset()
			}
			segments = append(segments, Segment{Text: "\n", Width: 0, Type: BreakSegment})
		} else {
			curWord.WriteRune(r)
		}
	}
	if curWord.Len() > 0 {
		str := curWord.String()
		w := getMeasure(str, font, fontSize)
		segments = append(segments, Segment{Text: str, Width: w, Type: WordSegment})
	}

	return &PreparedText{
		Segments: segments,
		Font:     font,
		FontSize: fontSize,
	}
}

// Layout performs pure arithmetic line-breaking using the cached segment widths.
// This is incredibly fast (O(N) arithmetic) and skips the GUI engine entirely.
func Layout(prepared *PreparedText, maxWidth int, lineHeight int) *LayoutResult {
	lines := []LayoutLine{}

	var currentLine strings.Builder
	currentWidth := 0
	startIdx := 0

	for i, seg := range prepared.Segments {
		if seg.Type == BreakSegment {
			lines = append(lines, LayoutLine{
				Text:  currentLine.String(),
				Width: currentWidth,
				Start: startIdx,
				End:   i + 1,
			})
			currentLine.Reset()
			currentWidth = 0
			startIdx = i + 1
			continue
		}

		if currentWidth+seg.Width > maxWidth && currentWidth > 0 {
			// Wrap to next line
			lines = append(lines, LayoutLine{
				Text:  currentLine.String(),
				Width: currentWidth,
				Start: startIdx,
				End:   i,
			})
			currentLine.Reset()
			currentWidth = 0
			startIdx = i

			// Skip leading space on the new line
			if seg.Type == SpaceSegment {
				startIdx = i + 1
				continue
			}
		}

		currentLine.WriteString(seg.Text)
		currentWidth += seg.Width
	}

	if currentLine.Len() > 0 || startIdx < len(prepared.Segments) {
		lines = append(lines, LayoutLine{
			Text:  currentLine.String(),
			Width: currentWidth,
			Start: startIdx,
			End:   len(prepared.Segments),
		})
	}

	return &LayoutResult{
		Height:    len(lines) * lineHeight,
		LineCount: len(lines),
		Lines:     lines,
	}
}
