package dialog

import "testing"

func TestChoiceRejectsUnsupportedOptionCountsWithoutPanicking(t *testing.T) {
	if got := Choice("message"); got != -1 {
		t.Fatalf("Choice with no options = %d, want -1", got)
	}
	if got := Choice("message", "A", "B", "C"); got != -1 {
		t.Fatalf("Choice with three options = %d, want -1", got)
	}
}
