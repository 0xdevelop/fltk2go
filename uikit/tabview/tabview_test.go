package tabview_test

import (
	"strings"
	"testing"

	"github.com/0xdevelop/fltk2go/fltk_bridge"
	"github.com/0xdevelop/fltk2go/foundation"
	"github.com/0xdevelop/fltk2go/uikit/tabview"
	"github.com/0xdevelop/fltk2go/uikit/view"
)

// mockViewable 实现一个简单的 Viewable
type mockViewable struct {
	v view.UIView
}

func newMockViewable() *mockViewable {
	m := &mockViewable{}
	raw := fltk_bridge.NewGroup(0, 0, 100, 100, "")
	m.v.BindRaw(raw)
	return m
}

func (m *mockViewable) View() *view.UIView {
	return &m.v
}

func TestUITabView_Creation(t *testing.T) {
	tv := tabview.NewUITabView(&foundation.Rect{X: 10, Y: 10, Width: 400, Height: 300})
	if tv == nil {
		t.Fatal("Expected UITabView to be created")
	}
	if tv.Raw() == nil {
		t.Fatal("Expected Raw() to return underlying Group widget")
	}
	if tv.View() == nil {
		t.Fatal("Expected View() to return UIView")
	}
}

func TestUITabView_AddTab(t *testing.T) {
	tv := tabview.NewUITabView(nil)

	child1 := newMockViewable()
	child2 := newMockViewable()

	tv.AddTab("Tab 1", child1)
	tv.AddTab("Tab 2", child2)

	// Since we added two tabs, activeIndex should be 0 by default
	// and there should be multiple children in the tab bar and content area
	if tv.Raw().ChildCount() < 2 {
		t.Fatalf("Expected UITabView to have children")
	}
}

func TestUITabView_SelectTab(t *testing.T) {
	tv := tabview.NewUITabView(nil)

	child1 := newMockViewable()
	child2 := newMockViewable()

	tv.AddTab("Tab 1", child1)
	tv.AddTab("Tab 2", child2)

	called := false
	tv.OnTabChanged(func(idx int) {
		if idx == 1 {
			called = true
		}
	})

	tv.SelectTab(1)
	if !called {
		t.Fatal("Expected OnTabChanged callback to be triggered for index 1")
	}
	if tv.ActiveIndex() != 1 || tv.Count() != 2 {
		t.Fatalf("unexpected tab state: active=%d count=%d", tv.ActiveIndex(), tv.Count())
	}
}

func TestUITabViewSelectAdjacentWrapsAndNotifies(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.AddTabWithID("one", "One", nil)
	tv.AddTabWithID("two", "Two", nil)
	tv.AddTabWithID("three", "Three", nil)

	var changed []int
	tv.OnTabChanged(func(index int) { changed = append(changed, index) })

	if !tv.SelectPrevious() || tv.ActiveIndex() != 2 {
		t.Fatalf("previous from first did not wrap: active=%d", tv.ActiveIndex())
	}
	if !tv.SelectNext() || tv.ActiveIndex() != 0 {
		t.Fatalf("next from last did not wrap: active=%d", tv.ActiveIndex())
	}
	if got, want := changed, []int{2, 0}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("selection callbacks = %#v, want %#v", got, want)
	}
}

func TestUITabViewSelectAdjacentRequiresMultipleTabs(t *testing.T) {
	var nilTabs *tabview.UITabView
	if nilTabs.SelectNext() || nilTabs.SelectPrevious() {
		t.Fatal("nil tab view reported a selection change")
	}
	tv := tabview.NewUITabView(nil)
	if tv.SelectNext() || tv.SelectPrevious() {
		t.Fatal("empty tab view reported a selection change")
	}
	tv.AddTabWithID("only", "Only", nil)
	if tv.SelectNext() || tv.SelectPrevious() || tv.ActiveIndex() != 0 {
		t.Fatalf("single tab should remain stable: active=%d", tv.ActiveIndex())
	}
}

func TestUITabViewMoveTabPreservesActiveIdentityAndAutomationOrder(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.SetAutomationID("sessions")
	tv.AddTabWithID("one", "One", nil)
	tv.AddTabWithID("two", "Two", nil)
	tv.AddTabWithID("three", "Three", nil)
	tv.SelectTab(1)

	changed := 0
	tv.OnTabChanged(func(int) { changed++ })
	if !tv.MoveTab(1, 0) {
		t.Fatal("MoveTab(1, 0) failed")
	}
	if got := []string{tv.TabID(0), tv.TabID(1), tv.TabID(2)}; got[0] != "two" || got[1] != "one" || got[2] != "three" {
		t.Fatalf("tab order after move = %#v", got)
	}
	if tv.ActiveIndex() != 0 || tv.TabID(tv.ActiveIndex()) != "two" {
		t.Fatalf("active identity drifted: index=%d id=%q", tv.ActiveIndex(), tv.TabID(tv.ActiveIndex()))
	}
	if changed != 0 {
		t.Fatalf("reordering active identity emitted %d selection callbacks", changed)
	}
	if node, ok := view.AutomationLookup("sessions.tab.two"); !ok || node.AutomationSnapshot().Properties["index"] != "0" {
		t.Fatalf("moved tab automation index is stale: ok=%t node=%#v", ok, node)
	}
	if node, ok := view.AutomationLookup("sessions.tab.one"); !ok || node.AutomationSnapshot().Properties["index"] != "1" {
		t.Fatalf("displaced tab automation index is stale: ok=%t node=%#v", ok, node)
	}
}

func TestUITabViewMoveTabRejectsInvalidOrNoopMoves(t *testing.T) {
	var nilTabs *tabview.UITabView
	if nilTabs.MoveTab(0, 1) {
		t.Fatal("nil tab view accepted a move")
	}
	tv := tabview.NewUITabView(nil)
	tv.AddTabWithID("one", "One", nil)
	tv.AddTabWithID("two", "Two", nil)
	if tv.MoveTab(-1, 0) || tv.MoveTab(0, 2) || tv.MoveTab(0, 0) {
		t.Fatal("invalid or no-op move was accepted")
	}
	if tv.TabID(0) != "one" || tv.TabID(1) != "two" || tv.ActiveIndex() != 0 {
		t.Fatal("rejected move changed tab state")
	}
}

func TestUITabViewStableIDsAndDynamicRemoval(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.SetAutomationID("sessions")
	tv.AddTabWithID("local", "Local", newMockViewable())
	tv.AddTabWithID("server", "Server", newMockViewable())
	tv.AddTabWithID("logs", "Logs", newMockViewable())

	if got := tv.IndexOfID("server"); got != 1 {
		t.Fatalf("IndexOfID(server) = %d, want 1", got)
	}
	if _, ok := view.AutomationLookup("sessions.tab.server"); !ok {
		t.Fatal("stable semantic tab id was not registered")
	}

	tv.SelectTab(1)
	if !tv.RemoveTab(0) {
		t.Fatal("RemoveTab(0) failed")
	}
	if tv.ActiveIndex() != 0 || tv.TabID(0) != "server" {
		t.Fatalf("active tab identity changed after removing preceding tab: active=%d id=%q", tv.ActiveIndex(), tv.TabID(0))
	}
	if _, ok := view.AutomationLookup("sessions.tab.local"); ok {
		t.Fatal("removed tab remained in automation registry")
	}

	if !tv.SetTabTitle(0, "Server · 运行中") {
		t.Fatal("SetTabTitle failed")
	}
	if node, ok := view.AutomationLookup("sessions.tab.server"); !ok || node.AutomationName() != "Server · 运行中" {
		t.Fatal("updated title was not published semantically")
	}
}

func TestUITabViewRemovalSelectsNearestTabAndNotifies(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.AddTabWithID("one", "One", nil)
	tv.AddTabWithID("two", "Two", nil)
	tv.AddTabWithID("three", "Three", nil)
	tv.SelectTab(2)

	var changed []int
	tv.OnTabChanged(func(index int) { changed = append(changed, index) })
	if !tv.RemoveTab(2) {
		t.Fatal("RemoveTab(active) failed")
	}
	if tv.ActiveIndex() != 1 || tv.TabID(1) != "two" {
		t.Fatalf("expected nearest remaining tab, active=%d id=%q", tv.ActiveIndex(), tv.TabID(tv.ActiveIndex()))
	}
	if len(changed) != 1 || changed[0] != 1 {
		t.Fatalf("unexpected selection callbacks: %#v", changed)
	}

	tv.RemoveTab(1)
	tv.RemoveTab(0)
	if tv.ActiveIndex() != -1 || tv.Count() != 0 {
		t.Fatalf("empty tab view state is invalid: active=%d count=%d", tv.ActiveIndex(), tv.Count())
	}
}

func TestUITabViewCloseRequestUsesCurrentStableTabIndex(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.SetAutomationID("sessions")
	tv.AddTabWithID("one", "One", nil)
	tv.AddTabWithID("two", "Two", nil)
	tv.SetTabsClosable(true)

	var requested []int
	tv.OnTabCloseRequested(func(index int) { requested = append(requested, index) })
	if err := view.AutomationClick("sessions.tab.two.close"); err != nil {
		t.Fatalf("close automation failed: %v", err)
	}
	if len(requested) != 1 || requested[0] != 1 {
		t.Fatalf("close requests = %#v, want [1]", requested)
	}
	if tv.Count() != 2 {
		t.Fatal("close request must not remove a tab before the owner accepts it")
	}

	if !tv.RemoveTab(0) {
		t.Fatal("failed to remove preceding tab")
	}
	if err := view.AutomationClick("sessions.tab.two.close"); err != nil {
		t.Fatalf("close automation after reindex failed: %v", err)
	}
	if len(requested) != 2 || requested[1] != 0 {
		t.Fatalf("reindexed close requests = %#v, want [1 0]", requested)
	}
}

func TestUITabViewRequestTabCloseHonorsCurrentOwnerPolicy(t *testing.T) {
	var nilTabs *tabview.UITabView
	if nilTabs.RequestTabClose(0) {
		t.Fatal("nil tab view accepted a close request")
	}

	tv := tabview.NewUITabView(nil)
	tv.AddTabWithID("pinned", "Pinned", nil)
	tv.AddTabWithID("ordinary", "Ordinary", nil)
	tv.SetTabsClosable(true)
	if _, ok := tv.SetTabPinned(0, true); !ok {
		t.Fatal("failed to pin protected tab")
	}

	var requested []string
	tv.OnTabCloseRequested(func(index int) { requested = append(requested, tv.TabID(index)) })
	if tv.RequestTabClose(0) {
		t.Fatal("pinned tab accepted a close request")
	}
	if !tv.RequestTabClose(1) {
		t.Fatal("ordinary tab close request was rejected")
	}
	if got := strings.Join(requested, ","); got != "ordinary" {
		t.Fatalf("close requests = %q, want ordinary", got)
	}

	if !tv.SetTabClosable(1, false) || tv.RequestTabClose(1) {
		t.Fatal("per-tab close veto was not enforced")
	}
	tv.SetTabsClosable(false)
	if tv.RequestTabClose(1) {
		t.Fatal("global close veto was not enforced")
	}
}

func TestUITabViewCloseAffordanceVisibilityAndLifecycle(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.SetAutomationID("sessions")
	tv.AddTabWithID("local", "Local", nil)

	if err := view.AutomationClick("sessions.tab.local.close"); err != view.ErrAutomationNodeUnavailable {
		t.Fatalf("hidden close action error = %v, want unavailable", err)
	}
	tv.SetTabsClosable(true)
	if node, ok := view.AutomationLookup("sessions.tab.local.close"); !ok || node.AutomationName() != "Close Local" {
		t.Fatalf("close affordance semantic state missing: ok=%t node=%#v", ok, node)
	}
	tv.SetTabTitle(0, "Local · running")
	if node, _ := view.AutomationLookup("sessions.tab.local.close"); node.AutomationName() != "Close Local · running" {
		t.Fatalf("close affordance title was stale: %q", node.AutomationName())
	}
	if !tv.RemoveTab(0) {
		t.Fatal("RemoveTab failed")
	}
	if _, ok := view.AutomationLookup("sessions.tab.local.close"); ok {
		t.Fatal("removed close affordance remained in automation registry")
	}
}

func TestUITabViewPerTabClosableStateSurvivesReorder(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.SetAutomationID("sessions")
	tv.AddTabWithID("protected", "Protected", nil)
	tv.AddTabWithID("ordinary", "Ordinary", nil)
	tv.SetTabsClosable(true)

	if !tv.SetTabClosable(0, false) || tv.TabClosable(0) {
		t.Fatal("failed to disable the protected tab close affordance")
	}
	if err := view.AutomationClick("sessions.tab.protected.close"); err != view.ErrAutomationNodeUnavailable {
		t.Fatalf("protected close action error = %v, want unavailable", err)
	}

	var requested []string
	tv.OnTabCloseRequested(func(index int) { requested = append(requested, tv.TabID(index)) })
	if err := view.AutomationClick("sessions.tab.ordinary.close"); err != nil {
		t.Fatalf("ordinary close action failed: %v", err)
	}
	if !tv.MoveTab(0, 1) || tv.TabClosable(1) {
		t.Fatal("per-tab close policy did not follow stable identity across reorder")
	}
	if err := view.AutomationClick("sessions.tab.protected.close"); err != view.ErrAutomationNodeUnavailable {
		t.Fatalf("moved protected close action error = %v, want unavailable", err)
	}
	if len(requested) != 1 || requested[0] != "ordinary" {
		t.Fatalf("close requests = %#v, want only ordinary", requested)
	}
	if !tv.SetTabClosable(1, true) || !tv.TabClosable(1) {
		t.Fatal("failed to restore moved tab close affordance")
	}
}

func TestUITabViewPinnedTabsFormLeadingPartition(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.SetAutomationID("sessions")
	for _, id := range []string{"one", "two", "three", "four"} {
		tv.AddTabWithID(id, id, nil)
	}
	tv.SelectTab(0)
	tv.SetTabsClosable(true)

	index, ok := tv.SetTabPinned(2, true)
	if !ok || index != 0 {
		t.Fatalf("pin three = index %d, ok %t; want 0, true", index, ok)
	}
	index, ok = tv.SetTabPinned(2, true)
	if !ok || index != 1 {
		t.Fatalf("pin two = index %d, ok %t; want 1, true", index, ok)
	}
	if got := []string{tv.TabID(0), tv.TabID(1), tv.TabID(2), tv.TabID(3)}; strings.Join(got, ",") != "three,two,one,four" {
		t.Fatalf("pinned order = %v", got)
	}
	if tv.PinnedCount() != 2 || !tv.TabPinned(0) || !tv.TabPinned(1) || tv.TabPinned(2) {
		t.Fatalf("invalid pinned partition: count=%d", tv.PinnedCount())
	}
	if tv.TabID(tv.ActiveIndex()) != "one" {
		t.Fatalf("active identity drifted after pinning: %q", tv.TabID(tv.ActiveIndex()))
	}
	if node, ok := view.AutomationLookup("sessions.tab.three"); !ok || node.AutomationSnapshot().Properties["pinned"] != "true" {
		t.Fatalf("pinned semantic state missing: ok=%t node=%#v", ok, node)
	}
	if err := view.AutomationClick("sessions.tab.three.close"); err != view.ErrAutomationNodeUnavailable {
		t.Fatalf("pinned close action error = %v, want unavailable", err)
	}

	index, ok = tv.SetTabPinned(0, false)
	if !ok || index != 1 || tv.TabID(1) != "three" || tv.TabPinned(1) {
		t.Fatalf("unpin did not move to ordinary boundary: index=%d ok=%t", index, ok)
	}
	if node, ok := view.AutomationLookup("sessions.tab.three.close"); !ok || !node.AutomationSnapshot().Visible {
		t.Fatalf("unpin did not restore close affordance: ok=%t node=%#v", ok, node)
	}
}

func TestUITabViewMoveTabRejectsPinnedPartitionCrossing(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	for _, id := range []string{"pinned", "ordinary-one", "ordinary-two"} {
		tv.AddTabWithID(id, id, nil)
	}
	if _, ok := tv.SetTabPinned(0, true); !ok {
		t.Fatal("failed to pin first tab")
	}
	if tv.MoveTab(0, 1) || tv.MoveTab(1, 0) {
		t.Fatal("MoveTab crossed the pinned partition")
	}
	if !tv.MoveTab(1, 2) {
		t.Fatal("ordinary tabs could not reorder within their partition")
	}
	if got := []string{tv.TabID(0), tv.TabID(1), tv.TabID(2)}; strings.Join(got, ",") != "pinned,ordinary-two,ordinary-one" {
		t.Fatalf("unexpected order after guarded moves: %v", got)
	}
}

func TestUITabViewContextMenuRequestUsesCurrentStableIdentity(t *testing.T) {
	tv := tabview.NewUITabView(nil)
	tv.AddTabWithID("one", "One", nil)
	tv.AddTabWithID("two", "Two", nil)
	tv.AddTabWithID("three", "Three", nil)
	tv.SelectTab(1)

	var requested []tabview.TabContextMenuState
	tv.OnTabContextMenuRequested(func(state tabview.TabContextMenuState) {
		requested = append(requested, state)
	})
	if !tv.RequestTabContextMenu(2) {
		t.Fatal("valid tab context-menu request was rejected")
	}
	if len(requested) != 1 || requested[0].ID != "three" || requested[0].Title != "Three" || requested[0].Index != 2 || requested[0].Selected {
		t.Fatalf("unexpected context-menu state: %#v", requested)
	}

	if !tv.MoveTab(2, 0) || !tv.RequestTabContextMenu(0) {
		t.Fatal("moved tab context-menu request failed")
	}
	if got := requested[1]; got.ID != "three" || got.Index != 0 || got.Selected {
		t.Fatalf("context-menu identity was stale after reorder: %#v", got)
	}
	if tv.RequestTabContextMenu(-1) || tv.RequestTabContextMenu(3) {
		t.Fatal("invalid tab context-menu request was accepted")
	}
}

func TestUITabViewContextMenuRequiresOwnerHandler(t *testing.T) {
	var nilTabs *tabview.UITabView
	if nilTabs.RequestTabContextMenu(0) {
		t.Fatal("nil tab view accepted a context-menu request")
	}
	tv := tabview.NewUITabView(nil)
	tv.AddTabWithID("one", "One", nil)
	if tv.RequestTabContextMenu(0) {
		t.Fatal("context-menu request without an owner handler was consumed")
	}
}
