# AGENTS.md

8bitMyGoUI: pixel-art components in the manner of 16-bit consoles, for MyGo native UI. Go module `git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI`.
Progress truth: `TASKS.md`. Requirements: `PRD.md` and `docs/prd/components-*.md`.

## Stack

| Item | Value |
| --- | --- |
| Go | 1.27.1 |
| MyGo | `github.com/egoist/mygo v0.2.10`, pinned in go.mod |
| Font | Ark Pixel Font 12px proportional zh_Hans v2026.09.25 (OFL-1.1, `pixel/OFL.txt`). Rebuild: `go run ./tools/fontpack -in ark-pixel-12px-proportional-zh_hans.bdf` (BDF sha256 `535f7024…8c86`) |
| QR | `rsc.io/qr v0.2.0` (BSD-3), encoder behind `display.QRCode` |
| Web | `cmd/web` built with `GOOS=js GOARCH=wasm`; the browser runs the gallery under `ui.Tester` and blits frames to a canvas |
| Hosting | Cloudflare Workers static assets from `web/dist`; code on Gitea `git.zacharyzhang.com` |

Read MyGo source before using an API: `$(go env GOMODCACHE)/github.com/egoist/mygo@v0.2.10` (`ui/*.go`, `docs/ui/*.md`). Never guess.

## Layout

```text
pixel/             engine: settings, palette, canvas, dither, font, sprites, clock, editor, messages
icon/              hand-placed sprites
<category>/        one component per file, package = folder name
gallery/           registry, demos (demos_<category>.go), app view, invariant tests
cmd/gallery/       native window
cmd/web/           js/wasm host
web/               shell, glue JS, build.sh, deploy.sh, wrangler.jsonc, _headers, Playwright tests; web/dist is build output
tools/fontpack/    BDF to font12.bin
internal/archtest/ import and API rules
```

## Tiers

Imports point down only. Same-tier imports are forbidden. `internal/archtest` enforces it.

0 pixel · 1 icon · 2 layout, text, display · 3 input, feedback · 4 navigation, data, overlay, datetime, chart, fx · 5 game, chat, media, code, account, files, project, shop, agent, desk, tabletop · 6 gallery · 7 cmd

## Style rules

- Draw on the dot grid through `pixel.Canvas`. One dot is `Dot` DIPs.
- Name palette slots, never colors. Every pixel must be a palette color.
- Only axis-aligned rectangles reach the screen: no radius, path, shadow, gradient, opacity, scaled bitmap.
- Tone comes from Bayer dither between two slots.
- Animation steps through `Canvas.Step`, `Canvas.Busy` or `Kit.Tick`. No easing.
- Text through the pixel font only.

## MyGo rules

- Use Bases, containers, `List`, `GridView`, `Outline`, `Router`, `Overlay`, `DialogBase`, `PopoverBase`, `FocusGroup`, `Local`, `HandleInput`, `TextCaret`, `Shortcut`, clipboard. Full list: PRD §2.4.
- Never call styled widgets, `Background`, `Border*`, `Radius`, `Shadow`, `Gradient`, `Opacity`, `Stripes`, `Tooltip`, `Menu`, `ContextMenu`, `Transition`, `Animate*`, `Loop`, `c.Toast*`, or painter calls other than `Fill` and `Clip`. Only `pixel` calls `After`.
- A disabled element never paints. Interactive components draw on an enabled shell; the base inside fills it, carries behavior, label and `Disabled`, and turns MyGo's focus ring off.

## Component rules

- Signatures and results: PRD §5.1. Zero options give the default look. Sizes in options are dots.
- `k := pixel.K(c)` once per component; layout sizes through `k.D(n)`.
- Containers take content as a final `func()`; interactive results are `NameResult`.
- A Clicked, Shortcut, Changed or drag during the build makes MyGo build the frame again, so what they change shows at once. State changed any other way after its element was built (`Pressed`, `Edit.Changed`) is read at paint time.
- Never call `Selected`: MyGo paints it in the accent color. Radios are `CheckboxBase` with `RoleRadio`.
- Clip only boxes that sit on whole dots; split widths with whole-dot shares, not equal `Grow`. `Grow` starts from zero size, so never use it on a content-sized child.
- Row builders that MyGo runs during layout (`List`, `GridView`, `Outline`) go through `pixel.Lazy(c, row)`.
- Single-key shortcuts only while `pixel.Typing(c)` is false; MyGo routes shortcuts before a focused field.
- Copy shown by a component comes from `k.Msg`; business copy comes from the caller.
- Fail fast: panic on programmer errors such as bad options. Never swallow an error.
- Composites log state changes: `slog.Debug("8bit: <component> <event>", …)`. Never log secrets.
- No cross-frame `*ui.Element` caching. No goroutines touching elements.
- Files under 500 lines. Comments: one line, ten words at most, only when needed. No dead code, no versioned APIs.

## Gallery entry

Each component adds one `Entry` in `gallery/demos_<category>.go`: ID, Name, Zh, Category, Summary, Demo. The demo reports every event through `log`. The code panel shows the demo function's own source.

## Tests

- `go vet ./... && go test ./...` before reporting done.
- Behavior tests drive `ui.Tester` and assert values.
- `gallery` checks every entry under every palette: only palette colors, whole dots. `PIXEL_SHOTS=<dir> go test ./gallery -run TestInvariants` writes PNGs for visual review.
- Web: `npm install && npm test` in `web/` builds `web/dist` and runs Playwright on a local server; `BASE_URL=<url> npx playwright test` checks a deployed site. Tests read `window.__8bit` (`ready`, `stats()`, `texts()`, `page()`, `find(label)`, `log`).

## Deploy

- `npm run deploy` in `web/` (clean tree only): builds, uploads a Worker version, runs Playwright on its preview URL, promotes it to 100%, runs Playwright on https://8bitmygoui.zhangyanghaha0407.workers.dev.
- Rollback: `npx wrangler rollback <version-id>` in `web/`; `npx wrangler deployments list` shows versions.
- Scripts and wasm sit under `a/`, named by content hash, and `_headers` makes them immutable; `index.html` revalidates. Never cache a path whose content can change.

## Review

Each TASKS.md item is reviewed in a Herdr pane split to the right, by `codex --yolo -m gpt-6.1-sol -c model_reasoning_effort=high -c service_tier=fast`. At most three rounds. The reviewer reports defects against PRD and this file, ends with `VERDICT: PASS` or `VERDICT: FAIL`, and does not reject work for missing enhancements.

## Secrets

The Gitea token lives only in the environment (`GITEA_TOKEN`). Never write it to a file, a commit or a log.
