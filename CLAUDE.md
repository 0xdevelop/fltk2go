# CLAUDE.md

## Project Characteristics
### Project Vision
- iOS OC style GUI with rich controls, using `addSubview` with block callbacks.

### Technical Approach
- Utilizes C++/C/Go stack with CGo for integration.
- Ensures runtime thread safety for robust execution.

## Development & Build Guidelines
### Examples (案例库)
- 所有的分项案例都可以独立运行（如 `go run examples/button/main.go`）。
- `examples/main.go` 是所有分项案例的聚合体入口，运行它可以进入并分步演示所有子项。
- 新增组件时，必须在 `examples/` 目录下添加独立示例，并将其注册到聚合体入口中。

### Build Artifacts (构建产物)
- Debug 过程中 `build` 出来的临时二进制可执行文件（如 `*.exe`、macOS 下的二进制文件），**用完即刻删除**，禁止提交到版本库，保持工作区干净。
- 任何临时文件、实验目录、构建探针都必须放在项目根目录的 `tmp/` 下；禁止使用 `/tmp`、`/private/tmp` 或系统临时目录。

## AI Workflow Guidelines
### AI Assistant Workflow
- **Autonomous Task Execution**: The AI assistant will autonomously execute tasks as per defined guidelines.
- **Token Management**: Efficient management of tokens for seamless operations.
- **Phase-Based Commits**: Commits will be categorized based on phases with an accompanying changelog to track progress and changes effectively.
