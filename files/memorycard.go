package files

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Save is a game's save: a title and its icon.
type Save struct {
	Title string
	// Icon is at most 16 dots square.
	Icon *pixel.Sprite
	// Slot inks the icon's # dots; zero means Ink.
	Slot pixel.Slot
}

type MemoryCardOptions struct {
	Disabled bool
}

type MemoryCardResult struct {
	Element ui.Element
	Changed bool
	// Deleted is the block whose DELETE was pressed, else -1.
	Deleted int
}

const blocks, block = 15, 24

// MemoryCard lays fifteen blocks; saves fill them in order.
func MemoryCard(c *ui.Context, text string, saves []Save, chosen *int, o MemoryCardOptions) MemoryCardResult {
	k := pixel.K(c)
	if len(saves) > blocks || *chosen < 0 || *chosen >= blocks {
		panic(fmt.Sprintf("files: MemoryCard has %d saves, block %d chosen", len(saves), *chosen))
	}
	for _, s := range saves {
		if s.Title == "" || s.Icon == nil || s.Icon.W > 16 || s.Icon.H > 16 {
			panic(fmt.Sprintf("files: MemoryCard save %q needs an icon of 16 dots at most", s.Title))
		}
	}
	old := *chosen
	res := MemoryCardResult{Deleted: -1}
	res.Element = ui.Column(c).Width(k.D(5 * (block + 2))).Gap(k.D(4)).Children(func() {
		pixel.Text(c, fmt.Sprintf(k.Msg("blocks_used"), len(saves)), pixel.TextOptions{Slot: pixel.Faint})
		// Arrows choose at once; the cursor is the block.
		shell, base, st, _ := pixel.Pad(c, text, 5, blocks, block+2, block+2, chosen, o.Disabled)
		say := k.Msg("free")
		if *chosen < len(saves) {
			say = saves[*chosen].Title
		}
		base.Description(fmt.Sprintf("%d, %s", *chosen+1, say))
		at := *chosen
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range blocks {
				x, y := i%5*(block+2)+1, i/5*(block+2)+1
				g.Frame(x, y, block, block, pixel.Sunken, pixel.Well)
				if i < len(saves) {
					s := saves[i]
					g.Sprite(x+(block-s.Icon.W)/2, y+(block-s.Icon.H)/2, s.Icon, 0, s.Slot, 1)
				} else {
					g.Checker(x+2, y+2, block-4, block-4, pixel.Shade)
				}
				if i == at {
					ring(g, x-1, y-1, block+2, block+2)
				}
			}
			g.Finish(st)
		})
		ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			pixel.Text(c, say, pixel.TextOptions{Wrap: true, MaxLines: 1}).Grow(1)
			free := *chosen >= len(saves)
			if input.Button(c, k.Msg("delete"), input.ButtonOptions{Size: input.Small, Variant: input.Danger, Disabled: o.Disabled || free, Label: k.Msg("delete") + " " + say}).Clicked {
				res.Deleted = *chosen
			}
		})
	})
	res.Changed = *chosen != old
	if res.Changed || res.Deleted >= 0 {
		slog.Debug("8bit: memory card", "label", text, "block", *chosen, "deleted", res.Deleted)
	}
	return res
}

// ring draws a two-dot Accent ring inside the box.
func ring(g pixel.Canvas, x, y, w, h int) {
	for d := range 2 {
		g.Rect(x+d, y+d, w-2*d, 1, pixel.Accent)
		g.Rect(x+d, y+h-1-d, w-2*d, 1, pixel.Accent)
		g.Rect(x+d, y+d, 1, h-2*d, pixel.Accent)
		g.Rect(x+w-1-d, y+d, 1, h-2*d, pixel.Accent)
	}
}
