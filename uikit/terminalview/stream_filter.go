package terminalview

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	terminalEscape        byte = 0x1b
	terminalBell          byte = 0x07
	maxOSCPayload              = 4096
	maxTerminalTitleRunes      = 128
)

type terminalStreamState uint8

const (
	terminalStreamText terminalStreamState = iota
	terminalStreamEscape
	terminalStreamCSI
	terminalStreamOSC
	terminalStreamOSCEscape
)

// terminalStreamFilter removes terminal metadata sequences that FLTK's native
// Fl_Terminal does not implement and would otherwise render as visible garbage.
// Display-oriented CSI sequences remain byte-for-byte intact. State survives
// arbitrary PTY chunk boundaries.
type terminalStreamFilter struct {
	state          terminalStreamState
	csi            []byte
	osc            []byte
	oscDiscard     bool
	titles         []string
	bells          int
	bracketedPaste bool
}

func (f *terminalStreamFilter) Reset() {
	f.state = terminalStreamText
	f.csi = f.csi[:0]
	f.osc = f.osc[:0]
	f.oscDiscard = false
	f.titles = f.titles[:0]
	f.bells = 0
	f.bracketedPaste = false
}

func (f *terminalStreamFilter) Filter(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	out := make([]byte, 0, len(data))
	for _, b := range data {
		switch f.state {
		case terminalStreamText:
			if b == terminalBell {
				f.bells++
				out = append(out, b)
			} else if b == terminalEscape {
				f.state = terminalStreamEscape
			} else {
				out = append(out, b)
			}
		case terminalStreamEscape:
			switch b {
			case ']':
				f.osc = f.osc[:0]
				f.oscDiscard = false
				f.state = terminalStreamOSC
			case '[':
				f.csi = append(f.csi[:0], terminalEscape, '[')
				f.state = terminalStreamCSI
			default:
				out = append(out, terminalEscape, b)
				f.state = terminalStreamText
			}
		case terminalStreamCSI:
			f.csi = append(f.csi, b)
			if b >= 0x40 && b <= 0x7e {
				f.trackPrivateMode(f.csi)
				if !isUnsupportedPrivateMode(f.csi) {
					out = append(out, f.csi...)
				}
				f.csi = f.csi[:0]
				f.state = terminalStreamText
			}
		case terminalStreamOSC:
			switch b {
			case terminalBell:
				f.finishOSC()
				f.state = terminalStreamText
			case terminalEscape:
				f.state = terminalStreamOSCEscape
			default:
				f.appendOSC(b)
			}
		case terminalStreamOSCEscape:
			switch b {
			case '\\', terminalBell:
				f.finishOSC()
				f.state = terminalStreamText
			case terminalEscape:
				// Remain here so repeated ESC bytes cannot leak OSC payload.
			default:
				f.appendOSC(terminalEscape)
				f.appendOSC(b)
				f.state = terminalStreamOSC
			}
		}
	}
	return out
}

func (f *terminalStreamFilter) appendOSC(b byte) {
	if f.oscDiscard {
		return
	}
	if len(f.osc) >= maxOSCPayload {
		f.osc = f.osc[:0]
		f.oscDiscard = true
		return
	}
	f.osc = append(f.osc, b)
}

func (f *terminalStreamFilter) finishOSC() {
	defer func() {
		f.osc = f.osc[:0]
		f.oscDiscard = false
	}()
	if f.oscDiscard {
		return
	}
	command, value, ok := bytes.Cut(f.osc, []byte{';'})
	if !ok || (string(command) != "0" && string(command) != "2") || !utf8.Valid(value) {
		return
	}
	title := strings.TrimSpace(string(value))
	if title == "" {
		return
	}
	for _, r := range title {
		if unicode.IsControl(r) {
			return
		}
	}
	runes := []rune(title)
	if len(runes) > maxTerminalTitleRunes {
		title = string(runes[:maxTerminalTitleRunes])
	}
	if len(f.titles) == 0 || f.titles[len(f.titles)-1] != title {
		f.titles = append(f.titles, title)
	}
}

func (f *terminalStreamFilter) takeTitles() []string {
	if len(f.titles) == 0 {
		return nil
	}
	titles := append([]string(nil), f.titles...)
	f.titles = f.titles[:0]
	return titles
}

func (f *terminalStreamFilter) takeBells() int {
	bells := f.bells
	f.bells = 0
	return bells
}

func (f *terminalStreamFilter) trackPrivateMode(sequence []byte) {
	if !privateModeContains(sequence, 2004) {
		return
	}
	switch sequence[len(sequence)-1] {
	case 'h':
		f.bracketedPaste = true
	case 'l':
		f.bracketedPaste = false
	}
}

func privateModeContains(sequence []byte, wanted int) bool {
	if len(sequence) < 5 || sequence[0] != terminalEscape || sequence[1] != '[' || sequence[2] != '?' {
		return false
	}
	value, hasDigit := 0, false
	for _, b := range sequence[3 : len(sequence)-1] {
		switch {
		case b >= '0' && b <= '9':
			value = value*10 + int(b-'0')
			hasDigit = true
		case b == ';':
			if hasDigit && value == wanted {
				return true
			}
			value, hasDigit = 0, false
		default:
			return false
		}
	}
	return hasDigit && value == wanted
}

func isUnsupportedPrivateMode(sequence []byte) bool {
	if len(sequence) < 4 || sequence[0] != terminalEscape || sequence[1] != '[' || sequence[2] != '?' {
		return false
	}
	final := sequence[len(sequence)-1]
	return final == 'h' || final == 'l'
}
