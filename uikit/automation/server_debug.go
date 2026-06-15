//go:build !release

package automation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/0xYeah/fltk2go/fltk_bridge"
	"github.com/0xYeah/fltk2go/uikit/view"
)

var ErrDisabled = errors.New("fltk2go automation debug server is disabled; set FLTK2GO_AUTOMATION_DEBUG=1")

type Server struct {
	httpServer    *http.Server
	listener      net.Listener
	directActions bool
	actionTimeout time.Duration
}

type Config struct {
	Addr string

	// DirectActions runs automation handlers on the HTTP goroutine. It is intended
	// for unit tests that do not have an FLTK event loop. Production debug sessions
	// should leave this false so actions are dispatched through Fl::awake().
	DirectActions bool

	// ActionTimeout bounds UI-thread dispatch. Defaults to 2 seconds.
	ActionTimeout time.Duration
}

type SnapshotResponse struct {
	Nodes []view.AutomationNode `json:"nodes"`
}

type AutomationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

type ActionResponse struct {
	OK    bool             `json:"ok"`
	Error *AutomationError `json:"error,omitempty"`
}

func Enabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("FLTK2GO_AUTOMATION_DEBUG")))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func StartDebugServer(cfg Config) (*Server, error) {
	if !Enabled() {
		return nil, ErrDisabled
	}
	addr := strings.TrimSpace(cfg.Addr)
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	timeout := cfg.ActionTimeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	s := &Server{listener: ln, directActions: cfg.DirectActions, actionTimeout: timeout}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/debug/automation/snapshot", s.handleSnapshot)
	mux.HandleFunc("/debug/automation/click", s.handleClick)
	mux.HandleFunc("/debug/automation/set_text", s.handleSetText)
	mux.HandleFunc("/mcp", s.handleMCP)
	s.httpServer = &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = s.httpServer.Serve(ln) }()
	return s, nil
}

func (s *Server) Addr() string {
	if s == nil || s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

func (s *Server) Close() error {
	if s == nil || s.httpServer == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "debug": true})
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	nodes, err := runOnUIThread(r.Context(), s, func() ([]view.AutomationNode, error) {
		return view.AutomationSnapshot(), nil
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ActionResponse{OK: false, Error: automationError(err, "")})
		return
	}
	writeJSON(w, http.StatusOK, SnapshotResponse{Nodes: nodes})
}

func (s *Server) handleClick(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	var req struct {
		ID string `json:"id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ActionResponse{OK: false, Error: automationError(err, req.ID)})
		return
	}
	if _, err := runOnUIThread(r.Context(), s, func() (struct{}, error) {
		return struct{}{}, view.AutomationClick(req.ID)
	}); err != nil {
		writeJSON(w, http.StatusBadRequest, ActionResponse{OK: false, Error: automationError(err, req.ID)})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true})
}

func (s *Server) handleSetText(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	var req struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ActionResponse{OK: false, Error: automationError(err, req.ID)})
		return
	}
	if _, err := runOnUIThread(r.Context(), s, func() (struct{}, error) {
		return struct{}{}, view.AutomationSetText(req.ID, req.Text)
	}); err != nil {
		writeJSON(w, http.StatusBadRequest, ActionResponse{OK: false, Error: automationError(err, req.ID)})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true})
}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id,omitempty"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *Server) handleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, http.StatusOK, map[string]any{"name": "fltk2go-debug-automation", "transport": "mcp-json-rpc-http", "streaming": false})
		return
	}
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	var req rpcRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, rpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: -32700, Message: err.Error()}})
		return
	}
	if req.ID == nil && strings.HasPrefix(req.Method, "notifications/") {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	res := rpcResponse{JSONRPC: "2.0", ID: req.ID}
	switch req.Method {
	case "initialize":
		res.Result = map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "fltk2go-debug-automation", "version": "0.1.0"},
		}
	case "tools/list":
		res.Result = map[string]any{"tools": toolsList()}
	case "tools/call":
		res.Result = s.callTool(r.Context(), req.Params)
	default:
		res.Error = &rpcError{Code: -32601, Message: "method not found"}
	}
	writeJSON(w, http.StatusOK, res)
}

func toolsList() []map[string]any {
	return []map[string]any{
		{
			"name":        "fltk_snapshot",
			"description": "Return all registered FLTK/UIKit automation nodes and properties.",
			"inputSchema": map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{}},
		},
		{
			"name":        "fltk_click",
			"description": "Invoke the debug click action for an automation id.",
			"inputSchema": objectSchema([]string{"id"}, map[string]any{
				"id": map[string]any{"type": "string", "description": "Stable automation id, e.g. app.login.submit"},
			}),
		},
		{
			"name":        "fltk_set_text",
			"description": "Set text for a text-capable automation id.",
			"inputSchema": objectSchema([]string{"id", "text"}, map[string]any{
				"id":   map[string]any{"type": "string", "description": "Stable automation id, e.g. login.username"},
				"text": map[string]any{"type": "string", "description": "Text to put into the target control"},
			}),
		},
		{
			"name":        "fltk_wait",
			"description": "Wait until an automation id is registered. This waits for registry presence, not full visual stability.",
			"inputSchema": objectSchema([]string{"id"}, map[string]any{
				"id":         map[string]any{"type": "string", "description": "Stable automation id to wait for"},
				"timeout_ms": map[string]any{"type": "integer", "description": "Maximum wait time in milliseconds", "default": 1000, "minimum": 1, "maximum": 30000},
			}),
		},
	}
}

func objectSchema(required []string, props map[string]any) map[string]any {
	return map[string]any{"type": "object", "required": required, "additionalProperties": false, "properties": props}
}

func (s *Server) callTool(ctx context.Context, raw json.RawMessage) map[string]any {
	var req struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return toolError(err, "")
	}
	if len(req.Arguments) == 0 {
		req.Arguments = json.RawMessage(`{}`)
	}

	switch req.Name {
	case "fltk_snapshot":
		nodes, err := runOnUIThread(ctx, s, func() ([]view.AutomationNode, error) {
			return view.AutomationSnapshot(), nil
		})
		if err != nil {
			return toolError(err, "")
		}
		return toolOK(SnapshotResponse{Nodes: nodes})
	case "fltk_click":
		var args struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(req.Arguments, &args); err != nil {
			return toolError(err, "")
		}
		if _, err := runOnUIThread(ctx, s, func() (struct{}, error) {
			return struct{}{}, view.AutomationClick(args.ID)
		}); err != nil {
			return toolError(err, args.ID)
		}
		return toolOK(ActionResponse{OK: true})
	case "fltk_set_text":
		var args struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		}
		if err := json.Unmarshal(req.Arguments, &args); err != nil {
			return toolError(err, "")
		}
		if _, err := runOnUIThread(ctx, s, func() (struct{}, error) {
			return struct{}{}, view.AutomationSetText(args.ID, args.Text)
		}); err != nil {
			return toolError(err, args.ID)
		}
		return toolOK(ActionResponse{OK: true})
	case "fltk_wait":
		var args struct {
			ID        string `json:"id"`
			TimeoutMS int    `json:"timeout_ms"`
		}
		if err := json.Unmarshal(req.Arguments, &args); err != nil {
			return toolError(err, "")
		}
		if args.TimeoutMS <= 0 {
			args.TimeoutMS = 1000
		}
		if args.TimeoutMS > 30000 {
			args.TimeoutMS = 30000
		}
		deadline := time.Now().Add(time.Duration(args.TimeoutMS) * time.Millisecond)
		for {
			if _, ok := view.AutomationLookup(args.ID); ok {
				return toolOK(ActionResponse{OK: true})
			}
			if time.Now().After(deadline) {
				return toolError(view.ErrAutomationNodeNotFound, args.ID)
			}
			select {
			case <-ctx.Done():
				return toolError(ctx.Err(), args.ID)
			case <-time.After(25 * time.Millisecond):
			}
		}
	default:
		return toolError(fmt.Errorf("unknown tool %q", req.Name), "")
	}
}

func toolOK(payload any) map[string]any {
	b, _ := json.Marshal(payload)
	return map[string]any{
		"content":           []map[string]any{{"type": "text", "text": string(b)}},
		"structuredContent": payload,
		"isError":           false,
	}
}

func toolError(err error, id string) map[string]any {
	payload := ActionResponse{OK: false, Error: automationError(err, id)}
	b, _ := json.Marshal(payload)
	return map[string]any{
		"content":           []map[string]any{{"type": "text", "text": string(b)}},
		"structuredContent": payload,
		"isError":           true,
	}
}

func runOnUIThread[T any](ctx context.Context, s *Server, fn func() (T, error)) (T, error) {
	if s == nil || s.directActions {
		return fn()
	}
	if ctx == nil {
		ctx = context.Background()
	}

	type result struct {
		value T
		err   error
	}
	var canceled atomic.Bool
	done := make(chan result, 1)
	if ok := fltk_bridge.Awake(func() {
		if canceled.Load() {
			return
		}
		value, err := fn()
		done <- result{value: value, err: err}
	}); !ok {
		var zero T
		return zero, errors.New("failed to dispatch automation action to FLTK event loop")
	}

	timer := time.NewTimer(s.actionTimeout)
	defer timer.Stop()
	select {
	case res := <-done:
		return res.value, res.err
	case <-ctx.Done():
		canceled.Store(true)
		var zero T
		return zero, ctx.Err()
	case <-timer.C:
		canceled.Store(true)
		var zero T
		return zero, errors.New("timed out waiting for FLTK event loop automation action")
	}
}

func automationError(err error, id string) *AutomationError {
	if err == nil {
		return nil
	}
	code := "automation_error"
	switch {
	case errors.Is(err, view.ErrAutomationIDRequired):
		code = "id_required"
	case errors.Is(err, view.ErrAutomationNodeNotFound):
		code = "node_not_found"
	case errors.Is(err, view.ErrAutomationActionUnsupported):
		code = "action_unsupported"
	case errors.Is(err, context.Canceled):
		code = "request_canceled"
	case errors.Is(err, context.DeadlineExceeded):
		code = "request_deadline_exceeded"
	case strings.Contains(err.Error(), "timed out"):
		code = "ui_dispatch_timeout"
	case strings.Contains(err.Error(), "unknown tool"):
		code = "unknown_tool"
	}
	return &AutomationError{Code: code, Message: err.Error(), ID: id}
}

func allowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	writeJSON(w, http.StatusMethodNotAllowed, ActionResponse{OK: false, Error: &AutomationError{Code: "method_not_allowed", Message: "method not allowed"}})
	return false
}

func decodeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body required")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
