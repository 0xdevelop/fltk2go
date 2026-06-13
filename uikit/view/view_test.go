package view

import "testing"

func TestUIViewNilSafety(t *testing.T) {
	var v *UIView
	if v.Raw() != nil {
		t.Fatal("nil view Raw() should return nil")
	}
	if v.Superview() != nil {
		t.Fatal("nil view Superview() should return nil")
	}
	v.AddSubview(nil)
	v.RemoveFromSuperview()
}

func TestUIViewLifecycleWithoutHostIsNoop(t *testing.T) {
	v := &UIView{}
	if v.Raw() != nil {
		t.Fatal("zero UIView Raw() should return nil")
	}
	if v.Superview() != nil {
		t.Fatal("zero UIView Superview() should return nil")
	}
	v.AddSubview(nil)
	v.RemoveFromSuperview()
}
