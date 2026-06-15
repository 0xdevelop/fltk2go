package fltk_bridge

import (
	"testing"
)

func TestSecretInput(t *testing.T) {
	win := NewWindow(200, 200, "Test SecretInput")
	win.Begin()
	secret := NewSecretInput(10, 10, 100, 30, "Password")
	win.End()

	if secret == nil {
		t.Fatal("NewSecretInput returned nil")
	}

	secret.SetValue("testpass")
	if secret.Value() != "testpass" {
		t.Errorf("Expected 'testpass', got '%s'", secret.Value())
	}
}
