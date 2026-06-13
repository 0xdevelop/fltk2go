package slider

import "testing"

func TestUISliderZeroValueIsNilSafe(t *testing.T) {
	var s *UISlider
	if s.View() != nil {
		t.Fatal("nil slider View() should return nil")
	}
	if s.Raw() != nil {
		t.Fatal("nil slider Raw() should return nil")
	}
	s.SetMinimum(0)
	s.SetMaximum(100)
	s.SetStep(1)
	s.SetValue(50)
	s.SetType(0)
	s.OnValueChanged(func(float64) {})
	if got := s.Value(); got != 0 {
		t.Fatalf("nil slider Value() = %v, want 0", got)
	}
}
