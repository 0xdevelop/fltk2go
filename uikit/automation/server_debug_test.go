//go:build !release

package automation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/0xdevelop/fltk2go/uikit/view"
)

func TestStartDebugServerRequiresEnv(t *testing.T) {
	t.Setenv("FLTK2GO_AUTOMATION_DEBUG", "")
	if _, err := StartDebugServer(Config{}); err != ErrDisabled {
		t.Fatalf("StartDebugServer error = %v, want ErrDisabled", err)
	}
}

func TestEnabledEnvParsing(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  bool
	}{
		{"1", true}, {"true", true}, {"TRUE", true}, {" yes ", true}, {"on", true},
		{"", false}, {"0", false}, {"false", false}, {"no", false}, {"whatever", false},
	} {
		t.Setenv("FLTK2GO_AUTOMATION_DEBUG", tc.value)
		if got := Enabled(); got != tc.want {
			t.Fatalf("Enabled(%q) = %v, want %v", tc.value, got, tc.want)
		}
	}
}

func TestDebugServerHTTPAndMCP(t *testing.T) {
	t.Setenv("FLTK2GO_AUTOMATION_DEBUG", "1")
	clicked := false
	text := ""
	buttonView := (&view.UIView{}).SetAutomationID("demo.button").SetAutomationRole("button").OnAutomationClick(func() error {
		clicked = true
		return nil
	})
	defer buttonView.SetAutomationID("")
	inputView := (&view.UIView{}).SetAutomationID("demo.input").SetAutomationRole("textbox").SetAutomationTextHandlers(func(s string) error {
		text = s
		return nil
	}, func() (string, bool) { return text, true })
	defer inputView.SetAutomationID("")

	srv, err := StartDebugServer(Config{Addr: "127.0.0.1:0", DirectActions: true})
	if err != nil {
		t.Fatalf("StartDebugServer: %v", err)
	}
	defer srv.Close()

	base := "http://" + srv.Addr()
	var snapshot SnapshotResponse
	getJSON(t, base+"/debug/automation/snapshot", &snapshot)
	if len(snapshot.Nodes) == 0 {
		t.Fatal("snapshot returned no nodes")
	}

	var action ActionResponse
	postJSON(t, base+"/debug/automation/click", map[string]string{"id": "demo.button"}, &action)
	if !action.OK || !clicked {
		t.Fatalf("click action=%#v clicked=%v", action, clicked)
	}

	postJSON(t, base+"/debug/automation/set_text", map[string]string{"id": "demo.input", "text": "hello"}, &action)
	if !action.OK || text != "hello" {
		t.Fatalf("set_text action=%#v text=%q", action, text)
	}

	rpc := rpcCall(t, base, "tools/list", map[string]any{})
	result := rpc["result"].(map[string]any)
	if len(result["tools"].([]any)) != 4 {
		t.Fatalf("tools/list returned %#v", result)
	}

	rpc = rpcCall(t, base, "tools/call", map[string]any{"name": "fltk_snapshot", "arguments": map[string]any{}})
	assertToolOK(t, rpc)
	rpc = rpcCall(t, base, "tools/call", map[string]any{"name": "fltk_click", "arguments": map[string]any{"id": "demo.button"}})
	assertToolOK(t, rpc)
	rpc = rpcCall(t, base, "tools/call", map[string]any{"name": "fltk_set_text", "arguments": map[string]any{"id": "demo.input", "text": "world"}})
	assertToolOK(t, rpc)
	if text != "world" {
		t.Fatalf("mcp set_text text=%q, want world", text)
	}
	rpc = rpcCall(t, base, "tools/call", map[string]any{"name": "fltk_wait", "arguments": map[string]any{"id": "demo.input", "timeout_ms": 50}})
	assertToolOK(t, rpc)
}

func TestDebugServerErrorsAndMethodChecks(t *testing.T) {
	t.Setenv("FLTK2GO_AUTOMATION_DEBUG", "1")
	srv, err := StartDebugServer(Config{Addr: "127.0.0.1:0", DirectActions: true})
	if err != nil {
		t.Fatalf("StartDebugServer: %v", err)
	}
	defer srv.Close()
	base := "http://" + srv.Addr()

	resp, err := http.Post(base+"/debug/automation/snapshot", "application/json", bytes.NewReader([]byte(`{}`)))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("snapshot POST status=%d, want 405", resp.StatusCode)
	}

	resp, err = http.Post(base+"/debug/automation/click", "application/json", bytes.NewReader([]byte(`{"id":"missing"}`)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("missing click status=%d, want 400", resp.StatusCode)
	}
	var action ActionResponse
	if err := json.NewDecoder(resp.Body).Decode(&action); err != nil {
		t.Fatal(err)
	}
	if action.OK || action.Error == nil || action.Error.Code != "node_not_found" {
		t.Fatalf("unexpected missing click response %#v", action)
	}

	rpc := rpcCall(t, base, "tools/call", map[string]any{"name": "fltk_click", "arguments": map[string]any{"id": "missing"}})
	result := rpc["result"].(map[string]any)
	if result["isError"] != true {
		t.Fatalf("tool error result=%#v", result)
	}
}

func TestDebugServerDispatchTimeoutWithoutLoop(t *testing.T) {
	t.Setenv("FLTK2GO_AUTOMATION_DEBUG", "1")
	srv, err := StartDebugServer(Config{Addr: "127.0.0.1:0", ActionTimeout: 5 * time.Millisecond})
	if err != nil {
		t.Fatalf("StartDebugServer: %v", err)
	}
	defer srv.Close()

	resp, err := http.Get("http://" + srv.Addr() + "/debug/automation/snapshot")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("snapshot without loop status=%d, want 500", resp.StatusCode)
	}
}

func getJSON(t *testing.T, url string, out any) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status=%d", url, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatal(err)
	}
}

func postJSON(t *testing.T, url string, in any, out any) {
	t.Helper()
	buf, _ := json.Marshal(in)
	resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST %s status=%d", url, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatal(err)
	}
}

func rpcCall(t *testing.T, base, method string, params any) map[string]any {
	t.Helper()
	rpc := map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": params}
	buf, _ := json.Marshal(rpc)
	resp, err := http.Post(base+"/mcp", "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("mcp %s: %v", method, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("mcp %s status=%d", method, resp.StatusCode)
	}
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out["error"] != nil {
		t.Fatalf("mcp %s error=%#v", method, out["error"])
	}
	return out
}

func assertToolOK(t *testing.T, rpc map[string]any) {
	t.Helper()
	result := rpc["result"].(map[string]any)
	if result["isError"] != false {
		t.Fatalf("tool result error=%#v", result)
	}
	if _, ok := result["structuredContent"].(map[string]any); !ok {
		t.Fatalf("missing structuredContent in %#v", result)
	}
}
