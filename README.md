# fltk2go

![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)
![FLTK](https://img.shields.io/badge/FLTK-1.4-blue)

**fltk2go** 是一个将 C++ FLTK 框架桥接到 Go 的现代 GUI 库，提供了类似 iOS/Mac Cocoa Framework (UIKit) 的高级声明式与命令式 API。
它不仅提供了原生跨平台的高性能，还带来了 `addSubview` 和基于闭包(Block)回调的现代开发体验。

*  Use FLTK like the Cocoa Framework for iOS/Mac.
*  中文：像iOS/Mac的Cocoa Framework一样用FLTK
* [FLTK Resource Doc](https://www.fltk.org/doc-1.4/index.html)

---

## 📖 实用文档入口

README 保持快速概览；更完整的运行、自动化、Agent、Playwright 集成说明放在 `docs/`：

- [整体运行说明 / Running Guide](docs/RUNNING.md) — 环境、运行 examples、测试、构建、常见问题。
- [Claude / Codex Automation Guide](docs/AUTOMATION_AGENT_GUIDE.md) — 通过 MCP/HTTP automation bridge 做语义化 GUI 调试。
- [Playwright Quick Integration](docs/PLAYWRIGHT_INTEGRATION.md) — 用 Playwright test runner 快速接入 FLTK2Go native GUI 自动化。
- [API Documentation](docs/API_DOCUMENTATION.md) — API 参考文档。

---

## 🏗 三层架构设计 (Architecture)

fltk2go 采用了清晰的三层架构设计，既保证了底层 C++ API 的高效执行，又为 Go 开发者提供了极其友好的高层封装：

```text
fltk2go/
├─ uikit/              # ⭐ 应用层：高级 UI 框架 (类 iOS Cocoa UIKit)
│  ├─ window/          # 窗口组件
│  ├─ button/          # 按钮组件
│  ├─ splitview/       # [新] 分割视图组件
│  ├─ loginview/       # [新] 登录视图组件
│  └─ navigationbar/   # [新] 导航栏组件
├─ runtime/            # ⚙️ 核心层：运行时核心
│  ├─ handle/          # unsafe.Pointer 生命周期管理
│  ├─ callback/        # Go ↔ C 回调表与事件分发
│  ├─ loop/            # UI loop / 渲染循环
│  └─ thread/          # 主线程约束安全保证
├─ fltk_bridge/        # 🌉 底层：C ABI / Cgo
└─ lib/                # C/C++ FLTK 静态链接库
```

- **底层 (fltk_bridge)**: 封装 C++ FLTK API 为纯 C ABI，并通过 cgo 暴露给 Go。
- **核心层 (runtime)**: 负责 Go ↔ C 之间的指针转换、生命周期控制、事件分发机制与主线程亲和性管理。
- **应用层 (uikit)**: 提供类似 Cocoa 风格的面向对象 API 组件，极大地简化了 UI 构建与事件绑定的复杂性。

---

## 🚀 快速上手 (Quick Start)

### 安装

```shell
go get github.com/0xYeah/fltk2go@latest
```

UIKit-style controls are available from the root `uikit` package as a facade, while existing subpackage imports remain supported.

```go
package main

import (
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit"
)

func main() {
	runtime.LockOSThread()

	win := uikit.NewUIWindow(420, 220, "UIKit Controls")
	root := win.RootView()

	slider := uikit.NewUISlider(&foundation.Rect{X: 24, Y: 40, Width: 360, Height: 28})
	slider.SetMinimumValue(0)
	slider.SetMaximumValue(100)
	slider.SetValue(50)
	root.AddSubview(slider)

	progress := uikit.NewUIProgressView(&foundation.Rect{X: 24, Y: 84, Width: 360, Height: 22})
	progress.SetMinimumValue(0)
	progress.SetMaximumValue(100)
	progress.SetProgress(50)
	root.AddSubview(progress)

	slider.OnValueChanged(func(value float64) {
		progress.SetProgress(value)
	})

	win.Show()
	fltk2go.Run()
}
```

Current UIKit wrappers include `UILabel`, `UIButton`, `Input`, `UITableView`,
`UISlider`, `UIProgressView`, `UISwitch`, `UIScrollView`, `UISplitView`,
`UIStackView`, and `UITextView`, plus root facade dialog helpers
`uikit.Message`, `uikit.Alert`, and `uikit.Choice`.

### Examples 调试结构

`examples` 采用“可 import 示例包 + 同目录独立入口”的结构，方便 IDE 断点调试：

- `examples/<name>/<name>.go`: 示例业务包，暴露 `BuildView(parent *view.UIView)`，主 launcher 和独立入口都会调用这里。
- `examples/<name>/main.go`: 带 `//go:build ignore` 的标准 `package main` 调试入口，可在 IDE 中直接运行或断点调试，不参与示例包编译。
- `examples/main.go`: 示例总入口，左侧选择示例，右侧内嵌预览，也可以启动对应独立入口。

```shell
cd examples
go run .
go run ./counter/main.go
go run ./slider_progress/main.go
```

### Debug automation / MCP JSON-RPC HTTP

fltk2go now includes an opt-in debug automation layer for native FLTK/UIKit apps. It is designed for agents and test runners that need stable control inspection/actions without fragile screen-coordinate clicking.

Debug-only behavior:

- The HTTP/MCP server starts only when `FLTK2GO_AUTOMATION_DEBUG=1` is set.
- Building with `-tags release` compiles a disabled stub, so release binaries cannot expose the debug server.
- Automation actions are dispatched onto the FLTK event loop via `Fl::awake()` and call Go handlers directly where possible instead of moving the physical mouse.
- Bind to `127.0.0.1` unless you are inside a trusted CI network. The debug server has no authentication and can trigger application actions.
- `Config.DirectActions` is only for unit tests without an FLTK event loop; real debug apps should leave it false.

```go
import "github.com/0xYeah/fltk2go/uikit/automation"

if automation.Enabled() {
	srv, err := automation.StartDebugServer(automation.Config{Addr: "127.0.0.1:8765"})
	if err != nil {
		panic(err)
	}
	defer srv.Close()
}
```

Add stable IDs to controls through their `UIView`:

```go
startButton.View().
	SetAutomationID("app.start").
	SetAutomationName("Start mapping").
	SetAutomationRole("button")
```

Use globally unique, stable IDs such as `screen.section.control`. Do not use localized button text, random values, or unstable row indexes as IDs. `UIButton` automatically exposes role/name and click actions after `OnTouchUpInside`; `Input` automatically exposes role/name and text get/set handlers.
`UILabel`, `UISlider`, and `UIProgressView` expose their current state through the snapshot `value` field so agents can assert UI changes without parsing pixels.

Run a debug build:

```shell
FLTK2GO_AUTOMATION_DEBUG=1 go run .
```

HTTP endpoints:

- `GET /debug/automation/snapshot` returns registered controls with id, role, name, label, text, bounds, visible/enabled state, and custom properties.
- `POST /debug/automation/click` with `{"id":"app.start"}` invokes the debug click action.
- `POST /debug/automation/set_text` with `{"id":"app.input","text":"hello"}` updates text-capable controls.
- `POST /mcp` exposes MCP-style JSON-RPC-over-HTTP tools: `fltk_snapshot`, `fltk_click`, `fltk_set_text`, and `fltk_wait`. It is a simple POST endpoint, not a full SSE streaming transport.

Agent workflow for Hermes/Codex/Claude Code:

```shell
# inspect current UI tree; nodes include actions like ["click"] or ["set_text"]
curl -s http://127.0.0.1:8765/debug/automation/snapshot

# fill a field and click a semantic button
curl -s -X POST http://127.0.0.1:8765/debug/automation/set_text \
  -H 'Content-Type: application/json' \
  -d '{"id":"login.username","text":"yerikokay"}'
curl -s -X POST http://127.0.0.1:8765/debug/automation/click \
  -H 'Content-Type: application/json' \
  -d '{"id":"login.submit"}'
```

MCP-style JSON-RPC examples:

```shell
curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"agent","version":"dev"}}}'

curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"fltk_snapshot","arguments":{}}}'

curl -s http://127.0.0.1:8765/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"fltk_click","arguments":{"id":"app.start"}}}'
```

Tool results include both `structuredContent` for agents and `content[0].text` containing the same JSON for simple clients. `fltk_wait` currently waits for automation ID registration only; follow actions with another `fltk_snapshot` to verify visual/application state.

The aggregate `examples/` launcher is automation-enabled when started with the same environment. It defaults to the Counter example and lets agents switch previews through the launcher list node:

```shell
cd examples
FLTK2GO_AUTOMATION_DEBUG=1 FLTK2GO_AUTOMATION_ADDR=127.0.0.1:8765 go run .

# Switch the right-hand preview to Input.
curl -s -X POST http://127.0.0.1:8765/debug/automation/set_text \
  -H 'Content-Type: application/json' \
  -d '{"id":"examples.launcher.list","text":"Input"}'

# Fill fields, click Update preview, then assert input.preview.value.
curl -s -X POST http://127.0.0.1:8765/debug/automation/set_text \
  -H 'Content-Type: application/json' \
  -d '{"id":"input.text","text":"hello"}'
curl -s -X POST http://127.0.0.1:8765/debug/automation/click \
  -H 'Content-Type: application/json' \
  -d '{"id":"input.update_preview"}'
```

Current launcher coverage includes stable IDs and state assertions for Counter, Input, and Slider & Progress. For example, `slider.max` updates `slider.volume.progress.value` and `slider.brightness.progress.value` to `100`.

For Claude Code / Codex integration patterns, prompts, MCP JSON-RPC examples, covered automation IDs, and troubleshooting, see [`docs/AUTOMATION_AGENT_GUIDE.md`](docs/AUTOMATION_AGENT_GUIDE.md).

### UIKit 代码片段示范

通过 `uikit` 快速创建一个现代化的窗口和组件：

```go
package main

import (
	"fmt"
	"runtime"

	"github.com/0xYeah/fltk2go"
	"github.com/0xYeah/fltk2go/foundation"
	"github.com/0xYeah/fltk2go/uikit/loginview"
	"github.com/0xYeah/fltk2go/uikit/window"
)

func main() {
	// ⚠️ 必须锁定主线程，防止 GUI goroutine 被调度到其他 OS 线程导致图形上下文异常
	runtime.LockOSThread()

	// 1. 创建窗口
	win := window.NewUIWindow(400, 300, "UIKit 快速上手")
	root := win.RootView()

	// 2. 创建高级组件 (LoginView)
	lv := loginview.NewLoginView(&foundation.Rect{X: 0, Y: 0, Width: 400, Height: 300})
	
	// 3. 闭包事件绑定
	lv.OnLoginClick(func(username, password string) {
		fmt.Printf("登录点击! 账号: %s, 密码: %s\n", username, password)
	})

	// 4. 添加到视图树
	root.AddSubview(lv)

	// 5. 显示并启动事件循环
	win.Show()
	fltk2go.Run()
}
```

---

## 📦 支持的组件库 (Components)

`fltk2go/uikit` 提供了一系列丰富的现代 UI 控件：

- **基础视图**: `UIView`, `UIWindow`
- **控制组件**: `UIButton`, `UILabel`, `UITextField`, `UICheckbox`, `UISwitch`, `UISlider`, `UIStepper`
- **数据视图**: `UITableView`, `UITreeView`, `UIImageView`, `UITextView`
- **布局视图**: `UIScrollView`, `UITabView`, `UILayout`
- **菜单导航**: `UIMenuBar`, `UIAlert`, `UIDropdown`
- **🆕 新增高级组件**:
  - **`UISplitView` (SplitView)**: 支持拖拽调整大小的分割视图组件，适合构建复杂多面板 IDE 类布局。
  - **`UILoginView` (LoginView)**: 开箱即用的登录面板组件，内置账号密码输入和登录交互封装。
  - **`UINavigationBar` (NavigationBar)**: 类 iOS 导航栏组件，支持左侧/右侧按钮、标题和导航堆栈概念。

---

## 📚 API 速查表 (API Cheat Sheet)

### 核心生命周期
| API | 描述 |
| --- | --- |
| `runtime.LockOSThread()` | **必需**。绑定主线程，防止 GUI 崩溃。 |
| `fltk2go.Run()` | 启动并阻塞执行 FLTK 主事件循环。 |
| `fltk2go.App(func)` | 封装了主线程锁定与初始化的快捷入口函数。 |

### UIKit 视图操作
| API | 描述 |
| --- | --- |
| `view.AddSubview(child)` | 将子视图添加到父视图中。 |
| `view.RemoveFromSuperview()` | 从父视图中移除当前视图。 |
| `view.SetFrame(rect)` | 重新设置视图的坐标和尺寸。 |

### 常用事件绑定
| API | 描述 |
| --- | --- |
| `button.OnTouchUpInside(func)` | 按钮点击事件回调。 |
| `textfield.OnTextChanged(func)` | 文本输入框内容变更回调。 |
| `slider.OnValueChanged(func)` | 滑块数值变化回调。 |
| `tableview.OnCellClick(func)` | 表格行点击回调。 |
| `loginview.OnLoginClick(func)` | 登录面板点击回调，直接获取账密。 |

---

## 🛠 构建 FLTK 依赖 (Build)

如果需要重新编译底层的 C++ FLTK 依赖库：

```shell
go run fltk_build/fltk_build.go fltk_build/manifest.go
```

## ⚠️ 注意事项 (Warning)

**防止 GUI goroutine 被调度到其他 OS 线程**

对于 Win32 / OpenGL / GDI+ / macOS Cocoa 等具有线程亲和性的系统，GUI API 必须在主线程执行。

```go
func main() {
    // 将当前 goroutine 绑定到当前操作系统线程。
    // 防止 goroutine 被调度到其他 OS 线程，导致 GUI / 图形上下文失效或异常。
    runtime.LockOSThread()
    // ... 其他代码
}
```

## 🐧 依赖项 (Linux Depends)

在 Linux 环境下编译需要安装以下开发库：

```shell
apt update
apt install -y \
  build-essential cmake pkg-config \
  libx11-dev libxext-dev libxinerama-dev libxcursor-dev libxrender-dev libxfixes-dev \
  libxft-dev \
  libgl1-mesa-dev libglu1-mesa-dev \
  mesa-common-dev
```
