package tabletop

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type BingoCardOptions struct {
	Disabled bool
}

type BingoCardResult struct {
	Element ui.Element
	// Bingo is the first full line, this frame.
	Changed, Bingo bool
}

type bingoKey struct{}

// card is the cursor, and if a line was full.
type card struct {
	cur  int
	line bool
}

const bingoW, bingoH = 18, 16

// BingoCard marks a 5×5 card; a full line is BINGO!
func BingoCard(c *ui.Context, text string, numbers [25]int, marked *[25]bool, o BingoCardOptions) BingoCardResult {
	k := pixel.K(c)
	// Cells print three digits at most; the middle is FREE.
	for i, n := range numbers {
		if i != 12 && (n < 1 || n > 999) {
			panic(fmt.Sprintf("tabletop: BingoCard number %d at %d", n, i))
		}
	}
	col := ui.Column(c).Gap(k.D(2)).AlignItems(ui.Center).AlignSelf(ui.Start)
	st := ui.Local(col, bingoKey{}, func() card { return card{cur: 12} })
	// The middle is FREE: always marked.
	marked[12] = true
	res := BingoCardResult{}
	res.Element = col.Children(func() {
		ui.Box(c).Size(k.D(5*bingoW), k.D(7)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, l := range "BINGO" {
				g.Micro(i*bingoW+7, 1, string(l), pixel.Accent, 1)
			}
		})
		shell, base, ps, picked := pixel.Pad(c, text, 5, 25, bingoW, bingoH, &st.cur, o.Disabled)
		if picked >= 0 && picked != 12 {
			marked[picked], res.Changed = !marked[picked], true
		}
		say := k.Msg("free_space")
		if st.cur != 12 {
			say = strconv.Itoa(numbers[st.cur])
		}
		if marked[st.cur] {
			say += ", " + k.Msg("marked")
		}
		base.Description(fmt.Sprintf("%c %s", "BINGO"[st.cur%5], say))
		at, shown := st.cur, *marked
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range 25 {
				x, y := i%5*bingoW, i/5*bingoH
				// Cells share their edges; the last ones end inside.
				w, h := bingoW+1, bingoH+1
				if i%5 == 4 {
					w = bingoW
				}
				if i/5 == 4 {
					h = bingoH
				}
				g.Frame(x, y, w, h, pixel.Flat, pixel.Lit)
				if i == at && ps.Focus {
					g.Cursor(x+1, y+1, bingoW-1, bingoH-1)
				}
				if shown[i] && i != 12 {
					// A dauber's blot, dithered so the number reads.
					g.Disc(x+9, y+8, 6, pixel.Magic)
					g.Veil(x+3, y+2, 13, 13, pixel.Lit, 6)
				}
				// FREE is a star; it fits in any tongue.
				if i == 12 {
					g.Sprite(x+3, y+2, icon.Star, 0, pixel.Accent, 1)
					continue
				}
				s := strconv.Itoa(numbers[i])
				g.Micro(x+(bingoW+1-pixel.MicroWidth(s))/2, y+6, s, pixel.Ink, 1)
			}
			g.Finish(ps)
		})
		line := full(marked)
		res.Bingo, st.line = line && !st.line, line
		if line {
			pixel.Text(c, k.Msg("bingo"), pixel.TextOptions{Slot: pixel.Good, Scale: 2})
		}
	})
	if res.Changed || res.Bingo {
		slog.Debug("8bit: bingo card", "label", text, "cell", st.cur, "bingo", res.Bingo)
	}
	return res
}

// full reports a whole row, column or diagonal marked.
func full(m *[25]bool) bool {
	lines := [][5]int{{0, 6, 12, 18, 24}, {4, 8, 12, 16, 20}}
	for i := range 5 {
		lines = append(lines, [5]int{5 * i, 5*i + 1, 5*i + 2, 5*i + 3, 5*i + 4}, [5]int{i, i + 5, i + 10, i + 15, i + 20})
	}
	for _, l := range lines {
		if m[l[0]] && m[l[1]] && m[l[2]] && m[l[3]] && m[l[4]] {
			return true
		}
	}
	return false
}
