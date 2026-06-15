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
	"time"

	"github.com/0xYeah/fltk2go/uikit/view"
)

var ErrDisabled = errors.New("fltk2go automation debug server is disabled; set FLTK2GO_AUTOMATION_DEBUG=1")

type Server struct {
	httpServer *http.Server
	listener   net.Listener
}

type Config struct {
	Addr string
}

type SnapshotResponse struct {
	Nodes []view.AutomationNode `json:"nodes"`
}

type ActionResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
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
	s := &Server{listener: ln}
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
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "debug": true})
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, SnapshotResponse{Nodes: view.AutomationSnapshot()})
}

func (s *Server) handleClick(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ActionResponse{Error: err.Error()})
		return
	}
	if err := view.AutomationClick(req.ID); err != nil {
		writeJSON(w, http.StatusBadRequest, ActionResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true})
}

func (s *Server) handleSetText(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, ActionResponse{Error: err.Error()})
		return
	}
	if err := view.AutomationSetText(req.ID, req.Text); err != nil {
		writeJSON(w, http.StatusBadRequest, ActionResponse{Error: err.Error()})
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
		writeJSON(w, http.StatusOK, map[string]any{"name": "fltk2go-debug-automation", "transport": "streamable-http"})
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
		result, err := callTool(req.Params)
		if err != nil {
			res.Error = &rpcError{Code: -32000, Message: err.Error()}
		} else {
			res.Result = result
		}
	default:
		res.Error = &rpcError{Code: -32601, Message: "method not found"}
	}
	writeJSON(w, http.StatusOK, res)
}

func toolsList() []map[string]any {
	return []map[string]any{
		{"name": "fltk_snapshot", "description": "Return all registered FLTK/UIKit automation nodes and properties.", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
		{"name": "fltk_click", "description": "Invoke the debug click action for an automation id.", "inputSchema": map[string]any{"type": "object", "required": []string{"id"}, "properties": map[string]any{"id": map[string]any{"type": "string"}}}},
		{"name": "fltk_set_text", "description": "Set text for a text-capable automation id.", "inputSchema": map[string]any{"type": "object", "required": []string{"id", "text"}, "properties": map[string]any{"id": map[string]any{"type": "string"}, "text": map[string]any{"type": "string"}}}},
		{"name": "fltk_wait", "description": "Wait until an automation id is registered.", "inputSchema": map[string]any{"type": "object", "required": []string{"id"}, "properties": map[string]any{"id": map[string]any{"type": "string"}, "timeout_ms": map[string]any{"type": "integer"}}}},
	}
}

func callTool(raw json.RawMessage) (map[string]any, error) {
	var req struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	var payload any
	switch req.Name {
	case "fltk_snapshot":
		payload = SnapshotResponse{Nodes: view.AutomationSnapshot()}
	case "fltk_click":
		var args struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(req.Arguments, &args); err != nil {
			return nil, err
		}
		if err := view.AutomationClick(args.ID); err != nil {
			return nil, err
		}
		payload = ActionResponse{OK: true}
	case "fltk_set_text":
		var args struct{ ID, Text string }
		if err := json.Unmarshal(req.Arguments, &args); err != nil {
			return nil, err
		}
		if err := view.AutomationSetText(args.ID, args.Text); err != nil {
			return nil, err
		}
		payload = ActionResponse{OK: true}
	case "fltk_wait":
		var args struct {
			ID        string `json:"id"`
			TimeoutMS int    `json:"timeout_ms"`
		}
		if err := json.Unmarshal(req.Arguments, &args); err != nil {
			return nil, err
		}
		if args.TimeoutMS <= 0 {
			args.TimeoutMS = 1000
		}
		deadline := time.Now().Add(time.Duration(args.TimeoutMS) * time.Millisecond)
		for {
			if _, ok := view.AutomationLookup(args.ID); ok {
				payload = ActionResponse{OK: true}
				break
			}
			if time.Now().After(deadline) {
				return nil, view.ErrAutomationNodeNotFound
			}
			time.Sleep(25 * time.Millisecond)
		}
	default:
		return nil, fmt.Errorf("unknown tool %q", req.Name)
	}
	b, _ := json.Marshal(payload)
	return map[string]any{"content": []map[string]any{{"type": "text", "text": string(b)}}, "isError": false}, nil
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
