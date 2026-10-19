package chat

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type EmojiPickerOptions struct {
	Disabled bool
}

type EmojiPickerResult struct {
	Element ui.Element
	// Picked is the emoji chosen this frame, else empty.
	Picked string
}

type pickerKey struct{}

// pickerState is the open category and the cursor in it.
type pickerState struct{ set, cur int }

// EmojiPicker picks an emoji from one category's grid.
func EmojiPicker(c *ui.Context, text string, o EmojiPickerOptions) EmojiPickerResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(3))
	st := ui.Local(col, pickerKey{}, func() pickerState { return pickerState{} })
	var res EmojiPickerResult
	res.Element = col.Children(func() {
		tabs := make([]string, len(sets))
		for i, s := range sets {
			tabs[i] = k.Msg(s.key)
		}
		if input.Segmented(c, k.Msg("category"), tabs, &st.set, input.SegmentedOptions{Disabled: o.Disabled}).Changed {
			st.cur = 0
		}
		names := sets[st.set].names
		if picked := grid(c, k, text, names, 8, 1, &st.cur, o.Disabled); picked != "" {
			res.Picked = picked
			slog.Debug("8bit: emoji picked", "label", text, "emoji", picked)
		}
	})
	return res
}

// grid lays out names as a Pad of emoji.
func grid(c *ui.Context, k pixel.Kit, text string, names []string, cols, sc int, cur *int, disabled bool) string {
	cell := 12*sc + 4
	shell, base, st, picked := pixel.Pad(c, text, cols, len(names), cell, cell, cur, disabled)
	base.Description(names[*cur])
	at := *cur
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, n := range names {
			x, y := i%cols*cell, i/cols*cell
			if i == at && st.Focus {
				g.Cursor(x, y, cell, cell)
			}
			art := emoji(n)
			g.Sprite(x+2, y+2, art.sp, 0, art.ink, sc)
		}
		g.Finish(st)
	})
	if picked < 0 {
		return ""
	}
	return names[picked]
}
