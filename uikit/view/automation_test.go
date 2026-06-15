package view

import "testing"

func TestAutomationRegistryActionsAndSnapshot(t *testing.T) {
	var clicked bool
	v := (&UIView{}).
		SetAutomationID("demo.button").
		SetAutomationRole("button").
		SetAutomationName("Demo Button").
		SetAutomationProperty("semantic", "primary").
		OnAutomationClick(func() error {
			clicked = true
			return nil
		})
	defer v.SetAutomationID("")

	if err := AutomationClick("demo.button"); err != nil {
		t.Fatalf("AutomationClick error: %v", err)
	}
	if !clicked {
		t.Fatal("automation click handler was not invoked")
	}

	nodes := AutomationSnapshot()
	var found *AutomationNode
	for i := range nodes {
		if nodes[i].ID == "demo.button" {
			found = &nodes[i]
			break
		}
	}
	if found == nil {
		t.Fatal("registered automation node not found in snapshot")
	}
	if found.Role != "button" || found.Name != "Demo Button" || found.Properties["semantic"] != "primary" {
		t.Fatalf("unexpected snapshot: %#v", *found)
	}
}

func TestAutomationSetText(t *testing.T) {
	var text string
	v := (&UIView{}).SetAutomationID("demo.input").SetAutomationTextHandlers(func(s string) error {
		text = s
		return nil
	}, func() (string, bool) { return text, true })
	defer v.SetAutomationID("")

	if err := AutomationSetText("demo.input", "hello"); err != nil {
		t.Fatalf("AutomationSetText error: %v", err)
	}
	if text != "hello" {
		t.Fatalf("text = %q, want hello", text)
	}
	if got := v.AutomationSnapshot().Text; got != "hello" {
		t.Fatalf("snapshot text = %q, want hello", got)
	}
}
