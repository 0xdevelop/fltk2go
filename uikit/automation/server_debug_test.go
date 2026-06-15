//go:build !release

package automation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/0xYeah/fltk2go/uikit/view"
)

func TestStartDebugServerRequiresEnv(t *testing.T) {
	t.Setenv("FLTK2GO_AUTOMATION_DEBUG", "")
	if _, err := StartDebugServer(Config{}); err != ErrDisabled {
		t.Fatalf("StartDebugServer error = %v, want ErrDisabled", err)
	}
}

func TestDebugServerHTTPAndMCP(t *testing.T) {
	t.Setenv("FLTK2GO_AUTOMATION_DEBUG", "1")
	clicked := false
	v := (&view.UIView{}).SetAutomationID("demo.button").SetAutomationRole("button").OnAutomationClick(func() error {
		clicked = true
		return nil
	})
	defer v.SetAutomationID("")

	srv, err := StartDebugServer(Config{Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatalf("StartDebugServer: %v", err)
	}
	defer srv.Close()

	base := "http://" + srv.Addr()
	resp, err := http.Get(base + "/debug/automation/snapshot")
	if err != nil {
		t.Fatalf("snapshot http get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("snapshot status = %d", resp.StatusCode)
	}

	body := bytes.NewBufferString(`{"id":"demo.button"}`)
	resp, err = http.Post(base+"/debug/automation/click", "application/json", body)
	if err != nil {
		t.Fatalf("click http post: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || !clicked {
		t.Fatalf("click status=%d clicked=%v", resp.StatusCode, clicked)
	}

	rpc := map[string]any{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": map[string]any{}}
	buf, _ := json.Marshal(rpc)
	resp, err = http.Post(base+"/mcp", "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("mcp tools/list: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("mcp status = %d", resp.StatusCode)
	}
}
