# AGENTS.md

8bitMyGoUI: pixel-art components in the manner of 16-bit consoles, for MyGo native UI. Go module `git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI`.
Progress truth: `TASKS.md`. Requirements: `PRD.md` and `docs/prd/components-*.md`.

## Stack

| Item | Value |
| --- | --- |
| Go | 1.27.1 |
| MyGo | `github.com/egoist/mygo v0.2.10`, pinned in go.mod |
| Font | Ark Pixel Font 12px proportional zh_Hans (OFL-1.1), packed into `pixel/font12.bin` by `tools/fontpack` |
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
web/               shell, glue JS, build.sh, Playwright tests; web/dist is build output
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
- Never call styled widgets, `Background`, `Border*`, `Radius`, `Shadow`, `Gradient`, `Opacity`, `Stripes`, `Tooltip`, `Menu`, `ContextMenu`, `Transition`, `Animate*`, `Loop`, `c.Toast*`, or painter drawing other than `Fill`. Only `pixel` calls `After`.
- A disabled element never paints. Interactive components draw on an enabled shell; the base inside fills it, carries behavior, label and `Disabled`, and turns MyGo's focus ring off.

## Component rules

- Signatures and results: PRD §5.1. Zero options give the default look. Sizes in options are dots.
- `k := pixel.K(c)` once per component; layout sizes through `k.D(n)`.
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
- Web: `web/build.sh`, then `npx playwright test` in `web/`.

## Review

Each TASKS.md item is reviewed in a Herdr pane split to the right, by `codex --yolo -m gpt-6.1-sol -c model_reasoning_effort=high -c service_tier=fast`. At most three rounds. The reviewer reports defects against PRD and this file, ends with `VERDICT: PASS` or `VERDICT: FAIL`, and does not reject work for missing enhancements.

## Secrets

The Gitea token lives only in the environment (`GITEA_TOKEN`). Never write it to a file, a commit or a log.
