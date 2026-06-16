# Claude / Codex Automation Debugging Guide

This guide shows how Claude Code, Codex, or any other coding agent can debug FLTK2Go apps through the opt-in automation HTTP/MCP bridge instead of relying on screenshots, OCR, or fragile coordinate clicks.

The automation bridge is debug-only. It is enabled by `FLTK2GO_AUTOMATION_DEBUG=1` and disabled in `-tags release` builds.

## What the bridge exposes

When an app starts `automation.StartDebugServer`, agents can call:

| Endpoint / tool | Purpose |
| --- | --- |
| `GET /healthz` | Check that the debug bridge is running. |
| `GET /debug/automation/snapshot` | Read the current automation tree. |
| `POST /debug/automation/click` | Invoke a semantic click by automation id. |
| `POST /debug/automation/set_text` | Set text, or select a launcher item that exposes text selection. |
| `POST /mcp` + `fltk_snapshot` | MCP-style snapshot tool. |
| `POST /mcp` + `fltk_click` | MCP-style click tool. |
| `POST /mcp` + `fltk_set_text` | MCP-style text/selection tool. |
| `POST /mcp` + `fltk_wait` | Wait until an automation id is registered. |

Snapshots contain stable, agent-friendly fields:

```json
{
  "id": "counter.increment",
  "role": "button",
  "name": "Increment counter",
  "label": "点我 +1",
  "actions": ["click"],
  "enabled": true,
  "visible": true,
  "bounds": {"x": 20, "y": 80, "width": 160, "height": 44}
}
```

Stateful controls can expose `value`, for example labels, sliders, and progress bars:

```json
{
  "id": "slider.volume.progress",
  "role": "progressbar",
  "name": "Volume progress",
  "value": "100",
  "enabled": true,
  "visible": true
}
```

## Start the examples app for agent debugging

From the repository root:

```shell
cd examples
FLTK2GO_AUTOMATION_DEBUG=1 \
FLTK2GO_AUTOMATION_ADDR=127.0.0.1:8765 \
GOCACHE=../tmp/go-cache \
go run .
```

Use a different port when running multiple debug sessions:

```shell
FLTK2GO_AUTOMATION_ADDR=127.0.0.1:9876 go run .
```

Keep the server bound to `127.0.0.1` unless you are inside a trusted CI network. The bridge has no authentication and can trigger application actions.

## Raw HTTP workflow

### 1. Health check

```shell
curl -s http://127.0.0.1:8765/healthz
```

Expected:

```json
{"debug":true,"ok":true}
```

### 2. Read the UI tree

```shell
curl -s http://127.0.0.1:8765/debug/automation/snapshot
```

### 3. Click a button

```shell
curl -s -X POST http://127.0.0.1:8765/debug/automation/click \
  -H 'Content-Type: application/json' \
  -d '{"id":"counter.increment"}'
```

### 4. Set text or select an example

The examples launcher exposes `examples.launcher.list` as a text-selectable list node. Use it to switch the right-hand preview:

```shell
curl -s -X POST http://127.0.0.1:8765/debug/automation/set_text \
  -H 'Content-Type: application/json' \
  -d '{"id":"examples.launcher.list","text":"Input"}'
```

Then fill the Input example:

```shell
curl -s -X POST http://127.0.0.1:8765/debug/automation/set_text \
  -H 'Content-Type: application/json' \
  -d '{"id":"input.text","text":"hello"}'

curl -s -X POST http://127.0.0.1:8765/debug/automation/click \
  -H 'Content-Type: application/json' \
  -d '{"id":"input.update_preview"}'
```

Read `input.preview.value` from the next snapshot to assert the result.

## MCP-style JSON-RPC workflow

The `/mcp` endpoint is a simple JSON-RPC-over-HTTP bridge with MCP-style tool payloads. It is not a full SSE streaming transport.

### Initialize

```shell
curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"agent","version":"dev"}}}'
```

### List tools

```shell
curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
```

### Snapshot

```shell
curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"fltk_snapshot","arguments":{}}}'
```

### Wait for a node

```shell
curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"fltk_wait","arguments":{"id":"counter.increment","timeout_ms":5000}}}'
```

### Click

```shell
curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"fltk_click","arguments":{"id":"counter.increment"}}}'
```

### Set text / select

```shell
curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"fltk_set_text","arguments":{"id":"examples.launcher.list","text":"Slider & Progress"}}}'
```

Tool results include both human-readable text and structured data:

```json
{
  "content": [{"type":"text","text":"{\"ok\":true}"}],
  "structuredContent": {"ok": true},
  "isError": false
}
```

On tool-level failures, JSON-RPC still succeeds but `isError` is true:

```json
{
  "structuredContent": {
    "ok": false,
    "error": {
      "code": "node_not_found",
      "message": "automation node not found",
      "id": "missing.id"
    }
  },
  "isError": true
}
```

Agents should treat `isError: true` as a recoverable UI action failure: re-snapshot, wait, or choose another automation id.

## Claude Code usage pattern

Claude Code can drive the bridge from shell commands in the repository. A reliable prompt pattern is:

```text
Run the FLTK2Go examples app with FLTK2GO_AUTOMATION_DEBUG=1 on 127.0.0.1:8765.
Use curl against /mcp or /debug/automation/*, not screen coordinates.
First call fltk_snapshot, choose controls by stable id, perform actions, then call fltk_snapshot again to verify state fields like value/actions/enabled/visible.
Do not rely on OCR unless the automation tree lacks the needed state.
```

Suggested Claude Code loop:

1. Start the app in a background shell.
2. Poll `/healthz` until it returns `ok`.
3. Call `fltk_snapshot`.
4. Choose a node by `id`, `role`, and `actions`.
5. Call `fltk_wait` before interacting with dynamic previews.
6. Call `fltk_click` or `fltk_set_text`.
7. Call `fltk_snapshot` again and assert `value`, `text`, or node presence.
8. Stop the app process after the debugging run.

Example verification task:

```text
Use the automation bridge to verify that clicking counter.increment changes counter.title.value from Clicked 0 count to Clicked 1 count. Then switch examples.launcher.list to Input, fill input.text with hello, click input.update_preview, and verify input.preview.value contains Text: hello.
```

## Codex usage pattern

Codex can use the same bridge through shell commands or a small Python helper. Prefer JSON parsing over manual string matching.

Minimal Python helper:

```python
import json
import urllib.request

BASE = "http://127.0.0.1:8765"

def call(method, path, payload=None):
    data = json.dumps(payload).encode() if payload is not None else None
    req = urllib.request.Request(
        BASE + path,
        data=data,
        method=method,
        headers={"Content-Type": "application/json"},
    )
    with urllib.request.urlopen(req, timeout=5) as resp:
        return json.load(resp)

def mcp_tool(name, arguments):
    return call("POST", "/mcp", {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/call",
        "params": {"name": name, "arguments": arguments},
    })["result"]

def snapshot_nodes():
    result = mcp_tool("fltk_snapshot", {})
    return result["structuredContent"]["nodes"]

def find_node(node_id):
    for node in snapshot_nodes():
        if node.get("id") == node_id:
            return node
    raise AssertionError(f"missing node: {node_id}")

mcp_tool("fltk_click", {"id": "counter.increment"})
assert "Clicked" in find_node("counter.title")["value"]
```

Codex prompts should explicitly request semantic automation:

```text
Use the FLTK2Go automation bridge at http://127.0.0.1:8765. Do not click by coordinates. Parse JSON snapshots and assert state through node.value, node.text, node.actions, node.enabled, and node.visible.
```

## Example IDs currently covered

### Launcher

| ID | Role | Actions / state |
| --- | --- | --- |
| `examples.launcher.list` | `listbox` | `set_text`, `text`, `value` |
| `examples.launcher.run_selected` | `button` | `click` |
| `examples.preview` | `region` | `children` |

### Counter

| ID | Role | Actions / state |
| --- | --- | --- |
| `counter.title` | `text` | `value` |
| `counter.increment` | `button` | `click` |

### Input

| ID | Role | Actions / state |
| --- | --- | --- |
| `input.text` | `textbox` | `set_text`, `text` |
| `input.integer` | `textbox` | `set_text`, `text` |
| `input.float` | `textbox` | `set_text`, `text` |
| `input.password` | `textbox` | `set_text`, `text` |
| `input.note` | `textbox` | `set_text`, `text` |
| `input.preview` | `text` | `value` |
| `input.update_preview` | `button` | `click` |
| `input.clear` | `button` | `click` |

### Slider & Progress

Use either `Slider & Progress` or the HTML title `Slider &amp; Progress` when selecting the launcher item.

| ID | Role | Actions / state |
| --- | --- | --- |
| `slider.volume` | `slider` | `value` |
| `slider.volume.progress` | `progressbar` | `value` |
| `slider.volume.label` | `text` | `value` |
| `slider.brightness` | `slider` | `value` |
| `slider.brightness.progress` | `progressbar` | `value` |
| `slider.brightness.label` | `text` | `value` |
| `slider.reset` | `button` | `click` |
| `slider.set_half` | `button` | `click` |
| `slider.max` | `button` | `click` |

## Troubleshooting

### `/debug/automation/snapshot` times out

If snapshot returns `ui_dispatch_timeout`, the app may not be processing the FLTK event loop. Make sure the app called `runtime.LockOSThread()`, started the debug server before `fltk2go.Run()`, and called `fltk_bridge.Lock()` before using `Fl::awake()` from HTTP goroutines.

### `node_not_found`

The selected preview may not be loaded, or a previous preview was replaced. Call `fltk_snapshot`, inspect `examples.launcher.list.value`, then use `fltk_set_text` on `examples.launcher.list` and `fltk_wait` for the target id.

### `action_unsupported`

The node exists but does not expose the requested action. Check the node's `actions` array and use a supported operation. For example, labels expose `value` but are not clickable.

### Port already in use

Use another loopback port:

```shell
FLTK2GO_AUTOMATION_ADDR=127.0.0.1:9877 go run .
```

### Release builds

`go build -tags release` compiles a disabled automation stub. Do not expect `/healthz` or `/mcp` to exist in release binaries.

## Safety notes

- Keep the server on `127.0.0.1` by default.
- Do not expose the bridge to untrusted networks.
- The bridge is designed for development and CI debugging, not production control.
- Always verify actions with a follow-up snapshot; never assume a click succeeded just because the HTTP request returned 200.
