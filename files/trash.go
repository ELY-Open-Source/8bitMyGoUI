package files

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type TrashOptions struct {
	Disabled bool
}

type TrashResult struct {
	Element ui.Element
	Emptied bool
}

var (
	binEmpty = pixel.MustSprite(`
		................
		......eeee......
		..eeeeelleeeee..
		..elllllllllle..
		..eeeeeeeeeeee..
		...essssssssse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...essssssssse..
		....eeeeeeeee...
		................`)
	// Crumpled paper peeks out from under the lid.
	binFull = pixel.MustSprite(`
		......ee........
		.....elle.ee....
		..eeeeleeeelee..
		..elllllllllle..
		..eeeeeeeeeeee..
		...essssssssse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...eslslslslse..
		...essssssssse..
		....eeeeeeeee...
		................`)
)

// Trash is a bin with its count; EMPTY clears it.
func Trash(c *ui.Context, text string, count int, o TrashOptions) TrashResult {
	k := pixel.K(c)
	if count < 0 {
		panic(fmt.Sprintf("files: Trash count %d", count))
	}
	art, say := binEmpty, k.Msg("empty")
	if count > 0 {
		art, say = binFull, fmt.Sprintf(k.Msg("in_bin"), count)
	}
	var res TrashResult
	res.Element = ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text).Children(func() {
		ui.Box(c).Size(k.D(32), k.D(32)).Draw(func(p *ui.Painter, r ui.Rect) {
			k.Canvas(p, r).Sprite(0, 0, art, 0, pixel.Ink, 2)
		})
		ui.Column(c).Gap(k.D(3)).Children(func() {
			pixel.Text(c, say, pixel.TextOptions{})
			res.Emptied = input.Button(c, k.Msg("empty_trash"), input.ButtonOptions{Size: input.Small, Variant: input.Danger, Disabled: o.Disabled || count == 0}).Clicked
		})
	})
	if res.Emptied {
		slog.Debug("8bit: trash emptied", "label", text, "count", count)
	}
	return res
}
