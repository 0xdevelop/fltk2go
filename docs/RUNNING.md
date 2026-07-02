# FLTK2Go Running Guide

This document is a practical runbook for building, running, debugging, and validating FLTK2Go apps and examples.

## Requirements

- Go 1.18+.
- A working desktop/display environment for native FLTK windows.
- The repository's bundled FLTK bridge sources and generated static libraries.

On Linux CI or a headless VM, provide a display server such as the CI desktop session or Xvfb before launching GUI examples.

## Quick run

From the repository root:

```shell
go test ./...
cd examples
go run .
```

Run a single example entry:

```shell
cd examples
go run ./counter/main.go
go run ./input/main.go
go run ./slider_progress/main.go
```

The aggregate launcher (`examples/main.go`) shows a left-side example list, a middle description panel, and a right-side live preview. Independent `examples/<name>/main.go` files are marked with `//go:build ignore` so they can be run directly for IDE debugging without becoming part of the example package build.

## Minimal application skeleton

```go
package main

import (
    "runtime"

    "github.com/0xdevelop/fltk2go"
    "github.com/0xdevelop/fltk2go/foundation"
    "github.com/0xdevelop/fltk2go/uikit"
)

func main() {
    runtime.LockOSThread()

    win := uikit.NewUIWindow(480, 320, "FLTK2Go App")
    root := win.RootView()

    btn := uikit.NewUIButton(&foundation.Rect{X: 24, Y: 24, Width: 180, Height: 44}, "Click me")
    btn.OnTouchUpInside(func() {
        btn.SetTitle("Clicked")
    })
    root.AddSubview(btn)

    win.Show()
    fltk2go.Run()
}
```

The important lifecycle rules are:

1. Call `runtime.LockOSThread()` in `main` before creating GUI objects.
2. Build all initial UI on the main thread.
3. Call `win.Show()`.
4. Enter the FLTK event loop with `fltk2go.Run()`.

## Running with debug automation

The automation bridge is opt-in and debug-only:

```shell
cd examples
FLTK2GO_AUTOMATION_DEBUG=1 \
FLTK2GO_AUTOMATION_ADDR=127.0.0.1:8765 \
GOCACHE=../tmp/go-cache \
go run .
```

Check it:

```shell
curl -s http://127.0.0.1:8765/healthz
curl -s http://127.0.0.1:8765/debug/automation/snapshot
```

Use this for semantic automation with Claude Code, Codex, Playwright, or CI scripts. See:

- [Claude / Codex Automation Debugging Guide](AUTOMATION_AGENT_GUIDE.md)
- [Playwright Quick Integration](PLAYWRIGHT_INTEGRATION.md)

## Build and test commands

Recommended verification from the repository root:

```shell
GOCACHE=./tmp/go-cache go test ./...
(cd examples && GOCACHE=../tmp/go-cache go test ./...)
(cd examples && GOCACHE=../tmp/go-cache go build ./...)
go vet ./...
git diff --check
```

Debug automation release guard:

```shell
go test -tags release ./uikit/automation
```

Race check for automation packages:

```shell
go test -race ./uikit/view ./uikit/automation
```

## Rebuilding bundled FLTK dependencies

If you need to rebuild the underlying C++ FLTK dependency artifacts:

```shell
go run fltk_build/fltk_build.go fltk_build/manifest.go
```

Only do this when updating bridge/dependency internals. Most application and UIKit work does not require rebuilding FLTK itself.

## Development workflow

1. Add or update UIKit components under `uikit/`.
2. Add a focused example under `examples/<name>/`.
3. Register the example in `examples/main.go` when it should appear in the aggregate launcher.
4. Add stable automation IDs for controls that agents or tests should operate.
5. Run package tests and examples build.
6. Keep generated caches and temporary artifacts under project `tmp/`.

## Troubleshooting

### GUI window does not show

- Confirm a display server is available.
- Confirm `runtime.LockOSThread()` is called before GUI creation.
- Confirm `win.Show()` is called before `fltk2go.Run()`.

### Automation snapshot times out

- Confirm `FLTK2GO_AUTOMATION_DEBUG=1` is set.
- Confirm the app is running and `/healthz` responds.
- Confirm the app's main event loop is not blocked.
- Confirm the debug app calls `fltk_bridge.Lock()` before serving automation requests from background goroutines.

### Port already in use

Choose a different loopback port:

```shell
FLTK2GO_AUTOMATION_ADDR=127.0.0.1:9876 go run .
```

### Release binary has no automation server

This is expected. `-tags release` compiles the disabled automation stub.

## Safety

The debug automation bridge has no authentication and can trigger app actions. Keep it on `127.0.0.1` unless you are in a trusted CI environment.
