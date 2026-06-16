# Playwright Quick Integration for FLTK2Go Automation

Playwright is usually used for browser automation, but its test runner and HTTP client are also a practical way to drive FLTK2Go native GUI apps through the debug automation bridge.

Use this approach when you want:

- repeatable end-to-end tests around native FLTK/UIKit examples;
- CI-friendly assertions against JSON snapshots instead of screenshots/OCR;
- one familiar test runner for web, API, and native GUI debug automation;
- traces, retries, reporters, and parallel test orchestration from Playwright.

This guide uses the FLTK2Go automation bridge exposed by:

```text
GET  /healthz
GET  /debug/automation/snapshot
POST /debug/automation/click
POST /debug/automation/set_text
POST /mcp
```

The bridge is debug-only and must be enabled with `FLTK2GO_AUTOMATION_DEBUG=1`.

## Important model

Playwright is not clicking the native window directly here. Instead:

1. Playwright starts or connects to the FLTK2Go app.
2. Playwright calls the app's local automation HTTP/MCP endpoint.
3. The app dispatches actions onto the FLTK event loop.
4. Playwright reads the next JSON snapshot and asserts state.

That means tests are stable across screen sizes, window positions, language rendering, and headless/non-headless environments as long as the FLTK app can run.

## Install Playwright test runner

From the repository root, keep local tooling under `tmp/` when possible:

```shell
mkdir -p tmp/playwright
cd tmp/playwright
npm init -y
npm install -D @playwright/test
```

You do not need browser binaries for pure HTTP/API tests. If you also want browser screenshots or a browser UI, run:

```shell
npx playwright install chromium
```

## Start the examples app manually

Terminal 1:

```shell
cd examples
FLTK2GO_AUTOMATION_DEBUG=1 \
FLTK2GO_AUTOMATION_ADDR=127.0.0.1:8765 \
GOCACHE=../tmp/go-cache \
go run .
```

Terminal 2:

```shell
curl -s http://127.0.0.1:8765/healthz
```

Expected:

```json
{"debug":true,"ok":true}
```

## Minimal Playwright API test

Create `tmp/playwright/fltk2go-automation.spec.ts`:

```ts
import { expect, test, request } from '@playwright/test';

const BASE = process.env.FLTK2GO_AUTOMATION_URL ?? 'http://127.0.0.1:8765';

type AutomationNode = {
  id?: string;
  role?: string;
  name?: string;
  label?: string;
  text?: string;
  value?: string;
  actions?: string[];
  enabled: boolean;
  visible: boolean;
  bounds: { x: number; y: number; width: number; height: number };
  children?: AutomationNode[];
};

async function snapshot(api: ReturnType<typeof request.newContext> extends Promise<infer T> ? T : never) {
  const res = await api.get(`${BASE}/debug/automation/snapshot`);
  expect(res.ok()).toBeTruthy();
  return (await res.json()) as { nodes: AutomationNode[] };
}

function findNode(nodes: AutomationNode[], id: string): AutomationNode {
  const node = nodes.find((n) => n.id === id);
  if (!node) throw new Error(`missing automation node: ${id}`);
  return node;
}

test('counter increments through semantic automation', async () => {
  const api = await request.newContext();

  const health = await api.get(`${BASE}/healthz`);
  expect(await health.json()).toEqual({ debug: true, ok: true });

  let tree = await snapshot(api);
  expect(findNode(tree.nodes, 'counter.increment').actions).toContain('click');
  expect(findNode(tree.nodes, 'counter.title').value).toBe('Clicked 0 count');

  const click = await api.post(`${BASE}/debug/automation/click`, {
    data: { id: 'counter.increment' },
  });
  expect(click.ok()).toBeTruthy();

  tree = await snapshot(api);
  expect(findNode(tree.nodes, 'counter.title').value).toBe('Clicked 1 count');

  await api.dispose();
});
```

Run it:

```shell
cd tmp/playwright
FLTK2GO_AUTOMATION_URL=http://127.0.0.1:8765 npx playwright test fltk2go-automation.spec.ts
```

## MCP-style Playwright helper

If you want to exercise the MCP-style `/mcp` endpoint instead of raw HTTP endpoints, use a small helper:

```ts
import { expect, request, test } from '@playwright/test';

const BASE = process.env.FLTK2GO_AUTOMATION_URL ?? 'http://127.0.0.1:8765';

async function mcpTool(api: any, name: string, args: Record<string, unknown>) {
  const res = await api.post(`${BASE}/mcp`, {
    data: {
      jsonrpc: '2.0',
      id: Date.now(),
      method: 'tools/call',
      params: { name, arguments: args },
    },
  });
  expect(res.ok()).toBeTruthy();
  const body = await res.json();
  expect(body.error).toBeFalsy();
  return body.result;
}

test('switch to Input and verify preview through MCP tools', async () => {
  const api = await request.newContext();

  let result = await mcpTool(api, 'fltk_set_text', {
    id: 'examples.launcher.list',
    text: 'Input',
  });
  expect(result.isError).toBe(false);

  for (const [id, text] of [
    ['input.text', 'hello'],
    ['input.integer', '42'],
    ['input.float', '3.14'],
    ['input.password', 'pw'],
    ['input.note', 'note'],
  ]) {
    result = await mcpTool(api, 'fltk_set_text', { id, text });
    expect(result.isError).toBe(false);
  }

  result = await mcpTool(api, 'fltk_click', { id: 'input.update_preview' });
  expect(result.isError).toBe(false);

  result = await mcpTool(api, 'fltk_snapshot', {});
  const nodes = result.structuredContent.nodes;
  const preview = nodes.find((n: any) => n.id === 'input.preview');
  expect(preview.value).toContain('Text: hello');
  expect(preview.value).toContain('Integer: 42');

  await api.dispose();
});
```

## Start and stop the app from Playwright

For local development it is often simpler to start the app manually. For CI or a one-command test, Playwright can spawn the app:

```ts
import { spawn, type ChildProcessWithoutNullStreams } from 'node:child_process';
import { request, test as base, expect } from '@playwright/test';

const PORT = process.env.FLTK2GO_AUTOMATION_PORT ?? '8765';
const BASE = `http://127.0.0.1:${PORT}`;

let app: ChildProcessWithoutNullStreams | undefined;

async function waitForHealth(timeoutMs = 15000) {
  const api = await request.newContext();
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    try {
      const res = await api.get(`${BASE}/healthz`, { timeout: 1000 });
      if (res.ok()) {
        await api.dispose();
        return;
      }
    } catch {
      // app not ready yet
    }
    await new Promise((resolve) => setTimeout(resolve, 250));
  }
  await api.dispose();
  throw new Error(`FLTK2Go app did not become healthy at ${BASE}`);
}

base.beforeAll(async () => {
  app = spawn('go', ['run', '.'], {
    cwd: '../../examples',
    env: {
      ...process.env,
      FLTK2GO_AUTOMATION_DEBUG: '1',
      FLTK2GO_AUTOMATION_ADDR: `127.0.0.1:${PORT}`,
      GOCACHE: '../tmp/go-cache',
    },
  });
  app.stdout.on('data', (chunk) => process.stdout.write(chunk));
  app.stderr.on('data', (chunk) => process.stderr.write(chunk));
  await waitForHealth();
});

base.afterAll(async () => {
  if (app && !app.killed) app.kill('SIGTERM');
});

base('app is healthy', async () => {
  const api = await request.newContext();
  const res = await api.get(`${BASE}/healthz`);
  expect(res.ok()).toBeTruthy();
  await api.dispose();
});
```

Adjust `cwd` for your test directory layout. Prefer project-local caches such as `tmp/go-cache` so generated artifacts stay out of the repository.

## Recommended test cases

A compact smoke suite should cover:

1. **Counter**
   - snapshot has `counter.increment` with `click` action;
   - click it;
   - assert `counter.title.value` changes.

2. **Input**
   - set `examples.launcher.list` to `Input`;
   - fill `input.text`, `input.integer`, `input.float`, `input.password`, `input.note`;
   - click `input.update_preview`;
   - assert `input.preview.value` contains the entered values;
   - click `input.clear` and assert reset state.

3. **Slider & Progress**
   - set `examples.launcher.list` to `Slider & Progress`;
   - click `slider.max`;
   - assert `slider.volume.progress.value === "100"`;
   - assert `slider.brightness.progress.value === "100"`;
   - click `slider.reset` and assert values are `0`.

4. **Error path**
   - set `examples.launcher.list` to a missing title;
   - expect HTTP `400` or MCP `isError: true` with a structured error.

## CI notes

- The native FLTK app still needs a display server on Linux. Use your CI's desktop runner or an Xvfb setup.
- Keep the debug server on loopback (`127.0.0.1`).
- Use unique ports per parallel worker, for example `8765 + workerIndex`.
- Always shut down the spawned app in `afterAll`.
- Do not run these tests against `-tags release` binaries; the automation server is intentionally disabled there.

## Troubleshooting

### Snapshot returns `ui_dispatch_timeout`

The FLTK event loop is not processing `Fl::awake()` callbacks. Confirm that the app:

- called `runtime.LockOSThread()`;
- called `fltk_bridge.Lock()` before using the automation server from HTTP goroutines;
- started the automation server before entering `fltk2go.Run()`;
- has not blocked the main event loop with long synchronous work.

### Playwright gets `ECONNREFUSED`

The app is not listening yet or is using a different port. Poll `/healthz` before running assertions and keep `FLTK2GO_AUTOMATION_ADDR` / `FLTK2GO_AUTOMATION_URL` aligned.

### `node_not_found`

The desired example may not be selected. Re-snapshot and check `examples.launcher.list.value`, then select the correct preview with `fltk_set_text` / `/debug/automation/set_text`.

### `action_unsupported`

Read the node's `actions` array. Labels and progress bars expose `value`; buttons expose `click`; inputs expose `set_text`.

### Tests pass locally but fail in CI

Check for missing display server, stale app processes holding the port, or parallel workers sharing one port. Use per-worker ports and print app stdout/stderr from the spawned process.
