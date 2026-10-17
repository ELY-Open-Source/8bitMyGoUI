package code

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CodeBlockOptions struct {
	Disabled bool
}

type CodeBlockResult struct {
	Element ui.Element
	Copied  bool
}

type copiedKey struct{}

// CodeBlock colors Go-like source by slot, numbers lines, copies.
func CodeBlock(c *ui.Context, text, src string, o CodeBlockOptions) CodeBlockResult {
	k := pixel.K(c)
	src = tidy(src)
	lines, slots := strings.Split(src, "\n"), paint(src)
	gutter := pixel.Measure(strconv.Itoa(len(lines))) + 6
	var res CodeBlockResult
	res.Element = layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 3}, func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			pixel.Text(c, text, pixel.TextOptions{Slot: pixel.Faint}).Grow(1)
			e := ui.Box(c)
			at := ui.Local(e, copiedKey{}, func() time.Time { return time.Time{} })
			// The button says COPIED! for two seconds.
			label := k.Msg("copy")
			if k.Since(c, *at, 4, 8) < 8 {
				label = k.Msg("copied")
			}
			e.Children(func() {
				if input.Button(c, label, input.ButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
					c.WriteClipboard(src)
					*at, res.Copied = k.Now(), true
					slog.Debug("8bit: code block copied", "label", text, "lines", len(lines))
				}
			})
		})
		off := 0
		for i, l := range lines {
			from := off
			off += len(l) + 1
			ui.Box(c).Size(k.D(gutter+pixel.Measure(l)), k.D(pixel.LineHeight)).Role(ui.RoleText).Label(l).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				n := strconv.Itoa(i + 1)
				g.Text(gutter-6-pixel.Measure(n), 0, n, pixel.Faint, 1)
				colored(g, gutter, 0, l, slots[from:from+len(l)])
			})
		}
	})
	return res
}
