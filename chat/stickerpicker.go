package chat

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type StickerPickerOptions struct {
	Disabled bool
}

type StickerPickerResult struct {
	Element ui.Element
	// Picked is the sticker chosen this frame, else empty.
	Picked string
}

type stickerKey struct{}

// StickerPicker picks a big face sticker from a grid.
func StickerPicker(c *ui.Context, text string, o StickerPickerOptions) StickerPickerResult {
	k := pixel.K(c)
	box := ui.Box(c)
	cur := ui.Local(box, stickerKey{}, func() int { return 0 })
	var res StickerPickerResult
	res.Element = box.Children(func() {
		if picked := grid(c, k, text, display.EmojiNames, 4, 3, cur, o.Disabled); picked != "" {
			res.Picked = picked
			slog.Debug("8bit: sticker picked", "label", text, "sticker", picked)
		}
	})
	return res
}
