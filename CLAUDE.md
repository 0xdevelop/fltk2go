# fltk2go — Claude 工作上下文

## 项目定位

为 FLTK C++ GUI 库提供 Go 绑定，**对外 API 模仿 iOS/Cocoa UIKit 风格**。
目标：让 Go 开发者用 `UIWindow` / `UIView` / `UILabel` / `UIButton` 等熟悉的命名写原生桌面 GUI。

---

## 架构：两层严格分离

```
fltk_bridge/          ← 第 1 层：CGO 裸绑定，1:1 映射 FLTK C++ API
uikit/                ← 第 2 层：UIKit 风格高层封装，对外暴露的真正 API
  window/  label/  button/  input/  splitview/  tableview/  view/
  colors/  screen/
examples/             ← 各功能独立演示，go run examples/xxx/main.go 单文件运行
```

**原则：example 和业务代码只应调用 `uikit/`，不直接大量调用 `fltk_bridge/`。**
仅在 uikit 层封装不足时才在 example 里临时用 `fltk_bridge` 直接调用，并同步补齐 uikit 封装。

---

## uikit 组件规范

### 必须实现的接口

每个 uikit 组件都必须实现 `view.Viewable`：

```go
type Viewable interface {
    View() *UIView
}
```

### 组件结构模板

```go
type UIFoo struct {
    v   view.UIView          // 持有 host（Window/Group）和 raw（底层 widget）
    raw *fltk_bridge.XxxWidget
}

func NewUIFoo(r *foundation.Rect, ...) *UIFoo {
    raw := fltk_bridge.NewXxx(r.X, r.Y, r.Width, r.Height, ...)
    f := &UIFoo{raw: raw}
    f.v.BindRaw(raw)          // 绑定底层 widget 到 UIView
    return f
}

func (f *UIFoo) View() *view.UIView { return &f.v }
func (f *UIFoo) Raw() *fltk_bridge.XxxWidget { return f.raw }
```

### 命名约定

| 层级 | 前缀/规范 | 示例 |
|------|-----------|------|
| uikit 类型 | `UI` 前缀 | `UIWindow`, `UILabel`, `UIButton` |
| uikit 构造函数 | `NewUI` 或 `New`（无歧义时） | `NewUILabel`, `tableview.New` |
| uikit 回调 | iOS 语义 | `OnTouchUpInside`, `SetDataSource`, `SetDelegate` |
| fltk_bridge 类型 | 与 FLTK 一致 | `Button`, `Slider`, `Tabs` |

---

## FLTK 自动挂载机制（关键！）

FLTK 用"当前 Group"隐式父容器机制：

1. `NewWindowWithPosition(...)` 在 C++ 构造函数里调用 `begin()`，Window 成为当前 Group
2. 此后任何 `fltk_bridge.NewXxx(...)` 调用都会**自动 add 到当前 Group**
3. `NewTabs/NewGroup/NewFlex` 等 Group 子类的构造函数也会调用 `begin()`，切换当前 Group
4. `group.End()` 恢复上一个 Group（父容器）为当前 Group

```
NewWindow()          → Window 为当前 Group
  NewLabel()         → 自动加入 Window ✓
  NewTabs()          → 自动加入 Window；Tabs 成为当前 Group
    NewGroup("Tab1") → 自动加入 Tabs；Group1 成为当前 Group
      NewChoice()    → 自动加入 Group1 ✓
      NewBox()       → 自动加入 Group1 ✓
    group1.End()     → Tabs 重新成为当前 Group
    NewGroup("Tab2") → 自动加入 Tabs；Group2 成为当前 Group
      ...
    group2.End()
  tabs.End()         → Window 重新成为当前 Group
```

**`root.AddSubview(widget)`** 内部调用 `rawWin.Add(widget.raw)`。
由于 widget 通常已经被自动挂载，这个调用是幂等的（FLTK 不会重复添加）。
对于 Window 直属子组件，`AddSubview` 是显式、安全的写法。
对于 Tabs/Group 内部组件，**不能调用 `root.AddSubview`**（host 为 nil，调用无效）。

---

## 颜色系统

### FLTK 颜色格式（底层）

```go
// 格式：0xRRGGBB00（末字节为 0 表示直接 RGB 颜色）
const BLUE uint = 0x42A5F500   // Material Blue 500
```

传给 fltk_bridge：`fltk_bridge.Color(rgb_uint)`

### uikit/colors 包（高层）

```go
// 支持多种输入格式
colors.ColorWithRGB(0x42A5F5)         // uint
colors.ColorWithRGB("#42A5F5")        // hex string
colors.ColorWithRGB(66, 165, 245)     // r,g,b
colors.ColorWithRGB("blue")           // 命名色

// 预定义语义色
colors.Blue / colors.Red / colors.Background / colors.Selection
```

---

## 已实现的 uikit 组件

| 组件 | 包 | 底层 | 状态 |
|------|----|------|------|
| UIWindow | `uikit/window` | `fltk_bridge.Window` | ✅ |
| UIView | `uikit/view` | `fltk_bridge.Widget` interface | ✅ |
| UILabel | `uikit/label` | `fltk_bridge.Box` | ✅ |
| UIButton | `uikit/button` | `fltk_bridge.Button`（含 Checkbox/Radio/Toggle） | ✅ |
| Input | `uikit/input` | `fltk_bridge.Input/IntInput/FloatInput` | ✅ |
| SplitView | `uikit/splitview` | `fltk_bridge.Flex` | ✅ |
| TableView | `uikit/tableview` | `fltk_bridge.Table`（自定义 BridgeTable） | ✅ |
| Color | `uikit/colors` | `fltk_bridge.Color` | ✅ |

**未封装的 fltk_bridge 组件（待 uikit 化）**：
`Slider`, `ValueSlider`, `Progress`, `Spinner`, `Tabs`, `Choice`, `MenuBar`, `MenuButton`, `Tree`, `Browser`, `Chart`, `Scroll`, `Grid`

---

## TableView DataSource / Delegate 模式

```go
// DataSource（必须）
type DataSource interface {
    NumberOfRows(tv *TableView) int
    CellForRow(tv *TableView, row int) *TableViewCell
}

// Delegate（可选）
type Delegate interface {
    DidSelectRow(tv *TableView, row int)
    RowHeight(tv *TableView, row int) int  // 返回 0 用默认高度
}
```

Cell 复用：`tv.Dequeue("reuseID")` → 配置 → 返回；`tv.Enqueue(cell)` 回收。

---

## 顶层入口

```go
// 方式 1（推荐）：fltk2go.App 封装了 LockOSThread + Run
fltk2go.App(func() {
    win := window.NewUIWindow(800, 600, "Title")
    // ...build UI...
    win.Show()
})

// 方式 2：手动（用于 example）
runtime.LockOSThread()
// ...build UI...
win.Show()
fltk2go.Run()
```

**`runtime.LockOSThread()` 是必须的**：FLTK（OpenGL/Win32/GDI）要求 GUI 在固定 OS 线程上运行。

---

## Example 约定

- 每个文件独立的 `package main` + `func main()`，**单独运行**：`go run examples/xxx_example.go`
- 不能 `go build ./examples/`（多 `func main` 冲突）
- 颜色常量在每个文件内独立定义（各自运行，不共享包级别）
- 文件命名：`<功能>_example.go`

---

## Cocoa / UIKit 风格演进方向

当前阶段：绝对坐标定位（x, y, w, h）。

**下一步规划**（按优先级）：

1. **UIViewController** — 管理 View 生命周期（viewDidLoad / viewWillAppear 等）
2. **Auto Layout（Constraint）** — `NSLayoutConstraint` 风格，替代硬编码坐标
   - `view.topAnchor.constraint(equalTo: parent.topAnchor, constant: 20)`
   - 内部可用 FLTK Flex 或自研约束求解器实现
3. **UINavigationController** — 页面栈，push/pop
4. **事件/通知机制** — 类似 `NotificationCenter`
5. **响应链（Responder Chain）** — 键盘事件向上冒泡

**添加新 Cocoa 模式时**：先在此文件更新设计，再写实现。
