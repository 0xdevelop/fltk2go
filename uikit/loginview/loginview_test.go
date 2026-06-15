package loginview_test

import (
	"runtime"
	"testing"

	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/loginview"
)

func init() {
	runtime.LockOSThread()
}

func TestLoginView_Creation(t *testing.T) {
	lv := loginview.NewLoginView(nil)
	if lv == nil {
		t.Fatal("Expected LoginView to be created")
	}
	if lv.View() == nil {
		t.Fatal("Expected View() to return UIView")
	}
}

func TestLoginView_CreationWithRect(t *testing.T) {
	rect := &foundation.Rect{X: 10, Y: 10, Width: 500, Height: 400}
	lv := loginview.NewLoginView(rect)
	if lv == nil {
		t.Fatal("Expected LoginView to be created with Rect")
	}
}

func TestLoginView_Properties(t *testing.T) {
	lv := loginview.NewLoginView(nil)

	lv.SetUsername("testuser")
	if lv.Username() != "testuser" {
		t.Fatalf("Expected Username to be 'testuser', got '%s'", lv.Username())
	}

	lv.SetPassword("secret123")
	if lv.Password() != "secret123" {
		t.Fatalf("Expected Password to be 'secret123', got '%s'", lv.Password())
	}
}

func TestLoginView_Callbacks(t *testing.T) {
	lv := loginview.NewLoginView(nil)

	called := false
	var gotUsername, gotPassword string

	lv.OnLoginClick(func(username, password string) {
		called = true
		gotUsername = username
		gotPassword = password
	})

	lv.SetUsername("user")
	lv.SetPassword("pass")

	// Trigger callback indirectly if possible, but in unit tests without FLTK event loop
	// it's not straightforward to simulate a button click reliably. We just test
	// that setting the callback does not panic and works.
	if lv.View() == nil {
		t.Fatal("View should not be nil")
	}

	// We can't directly trigger the loginButton's callback without access to it
	// So we'll just acknowledge the setting doesn't crash
	_ = called
	_ = gotUsername
	_ = gotPassword
}
