package textfield_test

import (
	"runtime"
	"testing"

	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/uikit/textfield"
)

func init() {
	runtime.LockOSThread()
}

func TestUITextField_Creation(t *testing.T) {
	tf := textfield.NewUITextField(0, 0, 100, 30, "Placeholder")
	if tf == nil {
		t.Fatal("Expected UITextField to be created")
	}
	if tf.Raw() == nil {
		t.Fatal("Expected Raw() to not be nil")
	}
	if tf.View() == nil {
		t.Fatal("Expected View() to not be nil")
	}
}

func TestUITextField_Properties(t *testing.T) {
	tf := textfield.NewUITextField(0, 0, 100, 30, "OldPlaceholder")

	if tf.Placeholder() != "OldPlaceholder" {
		t.Fatalf("Expected Placeholder to be 'OldPlaceholder', got '%s'", tf.Placeholder())
	}

	tf.SetPlaceholder("NewPlaceholder")
	if tf.Placeholder() != "NewPlaceholder" {
		t.Fatalf("Expected Placeholder to be 'NewPlaceholder', got '%s'", tf.Placeholder())
	}

	tf.SetText("HelloWorld")
	if tf.Text() != "HelloWorld" {
		t.Fatalf("Expected Text to be 'HelloWorld', got '%s'", tf.Text())
	}

	tf.SetEnabled(false)
	if tf.IsEnabled() {
		t.Fatal("Expected IsEnabled to be false")
	}

	tf.SetEnabled(true)
	if !tf.IsEnabled() {
		t.Fatal("Expected IsEnabled to be true")
	}

	// Test styles
	tf.SetFontSize(14)
	tf.SetFont(fltk_bridge.HELVETICA)
	tf.SetTextColor(0xFF000000)
	tf.SetBackgroundColor(0x00FF0000)
	tf.SetCornerRadius()
	tf.SetBorderStyle(fltk_bridge.FLAT_BOX)
}

func TestUITextField_OnChange(t *testing.T) {
	tf := textfield.NewUITextField(0, 0, 100, 30, "")
	called := false
	tf.OnChange(func() {
		called = true
	})

	// Just verify setting the callback doesn't panic.
	if tf.Raw() == nil {
		t.Fatal("Raw should not be nil")
	}
	_ = called // Ignore unused variable
}

func TestUISecretTextField_Creation(t *testing.T) {
	stf := textfield.NewUISecretTextField(0, 0, 100, 30, "Password")
	if stf == nil {
		t.Fatal("Expected UISecretTextField to be created")
	}
	if stf.Raw() == nil {
		t.Fatal("Expected Raw() to not be nil")
	}
	if stf.View() == nil {
		t.Fatal("Expected View() to not be nil")
	}
}

func TestUISecretTextField_Properties(t *testing.T) {
	stf := textfield.NewUISecretTextField(0, 0, 100, 30, "Secret")

	stf.SetText("MySecret")
	if stf.Text() != "MySecret" {
		t.Fatalf("Expected Text to be 'MySecret', got '%s'", stf.Text())
	}

	stf.SetPlaceholder("NewSecret")
	if stf.Placeholder() != "NewSecret" {
		t.Fatalf("Expected Placeholder to be 'NewSecret', got '%s'", stf.Placeholder())
	}

	stf.SetEnabled(false)
	if stf.IsEnabled() {
		t.Fatal("Expected IsEnabled to be false")
	}

	stf.SetEnabled(true)
	if !stf.IsEnabled() {
		t.Fatal("Expected IsEnabled to be true")
	}

	// Test styles
	stf.SetFontSize(14)
	stf.SetFont(fltk_bridge.HELVETICA)
	stf.SetTextColor(0xFF000000)
	stf.SetBackgroundColor(0x00FF0000)
	stf.SetCornerRadius()
	stf.SetBorderStyle(fltk_bridge.FLAT_BOX)
}

func TestUISecretTextField_OnChange(t *testing.T) {
	stf := textfield.NewUISecretTextField(0, 0, 100, 30, "")
	called := false
	stf.OnChange(func() {
		called = true
	})

	// Just verify setting the callback doesn't panic.
	if stf.Raw() == nil {
		t.Fatal("Raw should not be nil")
	}
	_ = called // Ignore unused variable
}
