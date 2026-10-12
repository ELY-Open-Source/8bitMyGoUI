package desk

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ClipboardPanelOptions struct {
	// Width is in dots; zero means 200.
	Width    int
	Disabled bool
}

type ClipboardPanelResult struct {
	Element *ui.Element
	// Copied is the clip copied this frame, else -1.
	Copied int
}

type clipKey struct{}

// copied is the clip last copied and when.
type copied struct {
	clip string
	at   time.Time
}

// ClipboardPanel lists recent clips; a click copies one back.
func ClipboardPanel(c *ui.Context, text string, clips []string, o ClipboardPanelOptions) ClipboardPanelResult {
	k := pixel.K(c)
	seen := map[string]bool{}
	for i, s := range clips {
		if s == "" || seen[s] {
			panic(fmt.Sprintf("desk: ClipboardPanel clip %d is blank or twice", i))
		}
		seen[s] = true
	}
	w := pixel.Or("ClipboardPanel Width", o.Width, 200)
	col := ui.Column(c).Width(k.D(w)).Gap(k.D(2)).Role(ui.RoleList).Label(text)
	last := ui.Local(col, clipKey{}, func() copied { return copied{} })
	// COPIED! shows for two seconds on the clip copied.
	shown := ""
	if until := last.at.Add(2 * time.Second); k.Now().Before(until) {
		shown = last.clip
		k.At(c, until)
	}
	res := ClipboardPanelResult{Copied: -1}
	res.Element = col.Children(func() {
		if len(clips) == 0 {
			pixel.Text(c, k.Msg("empty"), pixel.TextOptions{Slot: pixel.Faint})
		}
		for i, clip := range clips {
			// Keyed by the clip: new ones come in on top.
			ui.Box(c).Key(clip).FillWidth().Children(func() {
				shell, base, st := pixel.Shell(c, func() *ui.Element {
					return ui.Box(c).Focusable().Role(ui.RoleButton).Label(clip).Description(k.Msg("copy")).Disabled(o.Disabled)
				})
				if !base.IsDisabled() && base.Clicked() {
					c.WriteClipboard(clip)
					res.Copied, *last = i, copied{clip, k.Now()}
				}
				done := clip == shown || res.Copied == i
				line := strings.Fields(clip)
				shell.FillWidth().Height(k.D(pixel.LineHeight + 6)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.Face)
					mark := ""
					if done {
						mark = k.Msg("copied")
						g.Text(g.W-4-pixel.Measure(mark), 3, mark, pixel.Good, 1)
					}
					g.Text(4, 3, pixel.Wrap(strings.Join(line, " "), g.W-12-pixel.Measure(mark), 1)[0], pixel.Ink, 1)
					g.Finish(st)
				})
			})
		}
	})
	if res.Copied >= 0 {
		slog.Debug("8bit: clip copied", "label", text, "clip", res.Copied)
	}
	return res
}
