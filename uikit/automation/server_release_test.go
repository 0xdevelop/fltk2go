//go:build release

package automation

import "testing"

func TestReleaseAutomationDisabled(t *testing.T) {
	t.Setenv("FLTK2GO_AUTOMATION_DEBUG", "1")
	if Enabled() {
		t.Fatal("Enabled() = true in release build, want false")
	}
	if srv, err := StartDebugServer(Config{Addr: "127.0.0.1:0"}); err != ErrDisabled || srv != nil {
		t.Fatalf("StartDebugServer() = (%v, %v), want nil ErrDisabled", srv, err)
	}
	var srv *Server
	if got := srv.Addr(); got != "" {
		t.Fatalf("nil Server Addr()=%q, want empty", got)
	}
	if err := srv.Close(); err != nil {
		t.Fatalf("nil Server Close()=%v, want nil", err)
	}
}
