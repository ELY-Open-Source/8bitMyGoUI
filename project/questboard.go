package project

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type QuestBoardOptions struct {
	// Width is each column's in dots; zero means 96.
	Width    int
	Disabled bool
}

type QuestBoardResult struct {
	Element ui.Element
	// Moved reports a card that changed column or place.
	Moved bool
}

type boardKey struct{}

// board is the cursor; held while a press carries it.
type board struct {
	col, row    int
	press, held bool
}

// Cards hold two lines of title.
const headH, cardH, colGap = pixel.LineHeight + 4, 2*pixel.LineHeight + 6, 4

var stages = [3]string{"todo", "doing", "done"}

// QuestBoard keeps cards in TODO, DOING and DONE.
func QuestBoard(c *ui.Context, text string, cols *[3][]string, o QuestBoardOptions) QuestBoardResult {
	k := pixel.K(c)
	cw, most := pixel.Or("QuestBoard Width", o.Width, 96), 1
	for _, col := range cols {
		if slices.Contains(col, "") {
			panic("project: QuestBoard has a card with no title")
		}
		most = max(most, len(col))
	}
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	b := ui.Local(base, boardKey{}, func() board { return board{} })
	fit := func() {
		b.col = max(0, min(2, b.col))
		b.row = max(0, min(len(cols[b.col])-1, b.row))
	}
	fit()
	res := QuestBoardResult{Element: shell}
	// move carries the cursor's card to a column and place.
	move := func(col, row int) {
		if col < 0 || col > 2 || len(cols[b.col]) == 0 {
			return
		}
		limit := len(cols[col])
		if col == b.col {
			limit--
		}
		if row = max(0, min(limit, row)); col == b.col && row == b.row {
			return
		}
		title := cols[b.col][b.row]
		cols[b.col] = slices.Delete(cols[b.col], b.row, b.row+1)
		cols[col] = slices.Insert(cols[col], row, title)
		b.col, b.row, res.Moved = col, row, true
	}
	pressed := false
	if !base.IsDisabled() {
		for _, m := range []struct {
			key    ui.Key
			dc, dr int
		}{{ui.KeyUp, 0, -1}, {ui.KeyDown, 0, 1}, {ui.KeyLeft, -1, 0}, {ui.KeyRight, 1, 0}} {
			if pixel.Shortcut(base, 0, m.key) {
				b.col, b.row = b.col+m.dc, b.row+m.dr
				fit()
			}
			if pixel.Shortcut(base, ui.Alt, m.key) {
				move(b.col+m.dc, b.row+m.dr)
			}
		}
		_, _, pressed = base.Dragged()
	}
	if pressed {
		x, y, _ := base.PointerPosition()
		col, dy := int(x/k.Dot())/(cw+colGap), int(y/k.Dot())-headH
		row := dy / cardH
		switch {
		case !b.press:
			// A new press takes the card under it, if any.
			b.held = col >= 0 && col < 3 && dy >= 0 && row < len(cols[col])
			if b.held {
				b.col, b.row = col, row
			}
		case b.held:
			move(col, max(0, dy)/cardH)
		}
	}
	b.held, b.press = b.held && pressed, pressed
	names := [3]string{k.Msg(stages[0]), k.Msg(stages[1]), k.Msg(stages[2])}
	if n := len(cols[b.col]); n > 0 {
		base.Description(fmt.Sprintf("%s %d/%d, %s", names[b.col], b.row+1, n, cols[b.col][b.row]))
	} else {
		base.Description(names[b.col] + ", " + k.Msg("empty"))
	}
	shown := [3][]string{slices.Clone(cols[0]), slices.Clone(cols[1]), slices.Clone(cols[2])}
	at, held := [2]int{b.col, b.row}, b.held
	shell.Size(k.D(3*cw+2*colGap), k.D(headH+most*cardH+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for ci, col := range shown {
			x := ci * (cw + colGap)
			g.Frame(x, 0, cw, g.H, pixel.Sunken, pixel.Well)
			g.Text(x+4, 2, names[ci]+" "+strconv.Itoa(len(col)), pixel.Faint, 1)
			g.HLine(x+2, headH-1, cw-4, pixel.Shade)
			for ri, t := range col {
				cx, cy := x+3, headH+ri*cardH+1
				on := ci == at[0] && ri == at[1]
				fill := pixel.Face
				// A carried card lifts off its shadow.
				if on && held {
					g.Rect(cx+1, cy+1, cw-6, cardH-2, pixel.Shade)
					cx, cy, fill = cx-1, cy-1, pixel.Lit
				}
				g.Frame(cx, cy, cw-6, cardH-2, pixel.Raised, fill)
				ink, room := pixel.Ink, cw-14
				if ci == 2 {
					ink, room = pixel.Faint, room-14
					g.Sprite(cx+cw-20, cy+3, icon.Check, 0, pixel.Good, 1)
				}
				for i, l := range pixel.Wrap(t, room, 2) {
					g.Text(cx+4, cy+2+i*pixel.LineHeight, l, ink, 1)
				}
				if on && (st.Focus || held) {
					g.Rect(cx-1, cy-1, cw-4, 1, pixel.Accent)
					g.Rect(cx-1, cy+cardH-2, cw-4, 1, pixel.Accent)
					g.Rect(cx-1, cy-1, 1, cardH, pixel.Accent)
					g.Rect(cx+cw-6, cy-1, 1, cardH, pixel.Accent)
				}
			}
		}
		g.Finish(st)
	})
	if res.Moved {
		slog.Debug("8bit: quest board", "label", text, "column", b.col, "row", b.row)
	}
	return res
}
