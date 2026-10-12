package desk

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type MinesweeperOptions struct {
	// Cols, Rows, Mines: zeros mean 9, 9 and 10.
	Cols, Rows, Mines int
	// Seed lays the mines; zero takes the clock.
	Seed     uint64
	Disabled bool
}

type MinesweeperResult struct {
	Element *ui.Element
	// Won or Lost is the game ending this frame.
	Won, Lost bool
	// Flagged is a flag set or lifted; Restarted, the face.
	Flagged, Restarted bool
}

type mineKey struct{}

// field is a game: mines, opened and flagged cells, clock.
type field struct {
	mine, open, flag []bool
	laid, over, won  bool
	boom, cur        int
	start, end       time.Time
}

const mcell = 16

// Minesweeper: Enter or click opens; F or right click flags.
func Minesweeper(c *ui.Context, label string, o MinesweeperOptions) MinesweeperResult {
	k := pixel.K(c)
	cols, rows, n := pixel.Or("Minesweeper Cols", o.Cols, 9), pixel.Or("Minesweeper Rows", o.Rows, 9), pixel.Or("Minesweeper Mines", o.Mines, 10)
	if n >= cols*rows {
		panic(fmt.Sprintf("desk: Minesweeper %d mines in %d cells", n, cols*rows))
	}
	col := ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(label)
	f := ui.Local(col, mineKey{}, func() field { return newField(cols * rows) })
	if len(f.mine) != cols*rows {
		*f = newField(cols * rows)
	}
	res := MinesweeperResult{}
	res.Element = col.Children(func() {
		ui.Row(c).Gap(k.D(12)).AlignItems(ui.Center).Children(func() {
			flags := 0
			for _, on := range f.flag {
				if on {
					flags++
				}
			}
			text.SevenSegment(c, fmt.Sprintf("%03d", max(-99, n-flags)), text.SevenSegmentOptions{})
			face := smile
			switch {
			case f.over && f.won:
				face = cool
			case f.over:
				face = dead
			}
			if input.IconButton(c, face, k.Msg("new_game"), input.IconButtonOptions{Disabled: o.Disabled}).Clicked {
				*f, res.Restarted = newField(cols*rows), true
				slog.Debug("8bit: minesweeper reset", "label", label)
			}
			text.SevenSegment(c, fmt.Sprintf("%03d", min(999, f.seconds(c, k))), text.SevenSegmentOptions{})
		})
		shell, base, st, picked := pixel.Pad(c, k.Msg("minefield"), cols, cols*rows, mcell, mcell, &f.cur, o.Disabled)
		if !base.IsDisabled() && !f.over {
			if !pixel.Typing(c) && base.Shortcut(0, ui.KeyF) {
				res.Flagged = f.toggle(f.cur)
			}
			if base.RightClicked() {
				if x, y, _ := base.PointerPosition(); x >= 0 && y >= 0 {
					if cx, cy := int(x/k.D(mcell)), int(y/k.D(mcell)); cx < cols && cy < rows {
						f.cur = cy*cols + cx
						res.Flagged = f.toggle(f.cur)
					}
				}
			}
			if picked >= 0 && !f.flag[picked] {
				if !f.laid {
					f.lay(cols, rows, n, picked, o.Seed, k.Now())
				}
				f.reveal(cols, rows, picked, k.Now())
				res.Won, res.Lost = f.over && f.won, f.over && !f.won
			}
		}
		base.Description(f.say(k, cols, f.cur))
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range f.mine {
				f.drawCell(g, cols, rows, i%cols*mcell, i/cols*mcell, i)
			}
			if st.Focus {
				x, y := f.cur%cols*mcell, f.cur/cols*mcell
				g.Rect(x, y, mcell, 1, pixel.Accent)
				g.Rect(x, y+mcell-1, mcell, 1, pixel.Accent)
				g.Rect(x, y, 1, mcell, pixel.Accent)
				g.Rect(x+mcell-1, y, 1, mcell, pixel.Accent)
			}
			g.Finish(st)
		})
		switch {
		case f.over && f.won:
			pixel.Text(c, k.Msg("stage_clear"), pixel.TextOptions{Slot: pixel.Good})
		case f.over:
			// GAME OVER, and the face is the CONTINUE.
			pixel.Text(c, k.Msg("game_over"), pixel.TextOptions{Slot: pixel.Bad})
		}
	})
	if res.Flagged {
		slog.Debug("8bit: minesweeper flag", "label", label, "cell", f.cur, "on", f.flag[f.cur])
	}
	if res.Won || res.Lost {
		slog.Debug("8bit: minesweeper over", "label", label, "won", res.Won)
	}
	return res
}

func newField(cells int) field {
	return field{mine: make([]bool, cells), open: make([]bool, cells), flag: make([]bool, cells), boom: -1}
}

// seconds since the first click, frozen at the end.
func (f *field) seconds(c *ui.Context, k pixel.Kit) int {
	switch {
	case !f.laid:
		return 0
	case f.over:
		return max(0, int(f.end.Sub(f.start)/time.Second))
	}
	// A clock stepping back holds the count at zero.
	s := max(0, int(k.Now().Sub(f.start)/time.Second))
	k.At(c, f.start.Add(time.Duration(s+1)*time.Second))
	return s
}

// toggle flags or unflags a shut cell, reporting a change.
func (f *field) toggle(i int) bool {
	if f.open[i] {
		return false
	}
	f.flag[i] = !f.flag[i]
	return true
}

// lay hides mines, none at the first click or beside.
func (f *field) lay(cols, rows, n, first int, seed uint64, now time.Time) {
	if seed == 0 {
		seed = uint64(now.UnixNano())
	}
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	var spots []int
	for i := range f.mine {
		if abs(i%cols-first%cols) > 1 || abs(i/cols-first/cols) > 1 {
			spots = append(spots, i)
		}
	}
	// Too few spots left? Then only the click is safe.
	if len(spots) < n {
		spots = spots[:0]
		for i := range f.mine {
			if i != first {
				spots = append(spots, i)
			}
		}
	}
	rng.Shuffle(len(spots), func(a, b int) { spots[a], spots[b] = spots[b], spots[a] })
	for _, i := range spots[:n] {
		f.mine[i] = true
	}
	f.laid, f.start = true, now
}

// near counts the mines around a cell.
func (f *field) near(cols, rows, i int) int {
	n := 0
	for _, j := range around(cols, rows, i) {
		if f.mine[j] {
			n++
		}
	}
	return n
}

// around lists a cell's neighbors on the board.
func around(cols, rows, i int) []int {
	var out []int
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			x, y := i%cols+dx, i/cols+dy
			if (dx != 0 || dy != 0) && x >= 0 && x < cols && y >= 0 && y < rows {
				out = append(out, y*cols+x)
			}
		}
	}
	return out
}

// reveal opens a cell; empties spread, a mine ends it.
func (f *field) reveal(cols, rows, i int, now time.Time) {
	if f.mine[i] {
		f.over, f.boom, f.end = true, i, now
		return
	}
	todo := []int{i}
	for len(todo) > 0 {
		j := todo[len(todo)-1]
		todo = todo[:len(todo)-1]
		if f.open[j] || f.flag[j] {
			continue
		}
		f.open[j] = true
		if f.near(cols, rows, j) == 0 {
			todo = append(todo, around(cols, rows, j)...)
		}
	}
	for j := range f.mine {
		if !f.mine[j] && !f.open[j] {
			return
		}
	}
	f.over, f.won, f.end = true, true, now
	copy(f.flag, f.mine)
}

// say reads a cell out: where, and what shows.
func (f *field) say(k pixel.Kit, cols, i int) string {
	at := fmt.Sprintf("%d,%d ", i%cols+1, i/cols+1)
	switch {
	case f.open[i]:
		return at + fmt.Sprint(f.near(cols, len(f.mine)/cols, i))
	case f.flag[i]:
		return at + k.Msg("flag")
	case f.over && f.mine[i]:
		return at + k.Msg("mine")
	}
	return at + k.Msg("unknown")
}

// numbers ink the mine counts, one to eight.
var numbers = [...]pixel.Slot{pixel.Info, pixel.Good, pixel.Bad, pixel.Magic, pixel.Warn, pixel.Accent, pixel.Ink, pixel.Faint}

// drawCell draws a cell: raised, flagged, counted or mined.
func (f *field) drawCell(g pixel.Canvas, cols, rows, x, y, i int) {
	switch {
	case f.open[i]:
		g.Rect(x, y, mcell, mcell, pixel.Well)
		g.Rect(x, y, mcell, 1, pixel.Shade)
		g.Rect(x, y, 1, mcell, pixel.Shade)
		if n := f.near(cols, rows, i); n > 0 {
			g.Micro(x+7, y+6, fmt.Sprint(n), numbers[n-1], 1)
		}
	case f.over && f.mine[i] && !f.flag[i]:
		back := pixel.Well
		if i == f.boom {
			back = pixel.Bad
		}
		g.Rect(x, y, mcell, mcell, back)
		g.Sprite(x+2, y+2, icon.Bomb, 0, pixel.Warn, 1)
	default:
		g.Frame(x, y, mcell, mcell, pixel.Raised, pixel.Face)
		if f.flag[i] {
			g.Sprite(x+2, y+2, icon.Flag, 0, pixel.Bad, 1)
		}
	}
}

var (
	smile = pixel.MustSprite(`
		...eeeeee...
		..eyyyyyye..
		.eyyyyyyyye.
		eyyykyykyyye
		eyyykyykyyye
		eyyyyyyyyyye
		eykyyyyyykye
		eyykyyyykyye
		eyyykkkkyyye
		.eyyyyyyyye.
		..eyyyyyye..
		...eeeeee...`)
	dead = pixel.MustSprite(`
		...eeeeee...
		..eyyyyyye..
		.eyyyyyyyye.
		eykykyykykye
		eyykyyyykyye
		eykykyykykye
		eyyyyyyyyyye
		eyyykkkkyyye
		eyykyyyykyye
		.eyyyyyyyye.
		..eyyyyyye..
		...eeeeee...`)
	cool = pixel.MustSprite(`
		...eeeeee...
		..eyyyyyye..
		.eyyyyyyyye.
		ekkkkkkkkkke
		eykkkyykkkye
		eyykyyyykyye
		eyyyyyyyyyye
		eykyyyyyykye
		eyykkkkkkyye
		.eyyyyyyyye.
		..eyyyyyye..
		...eeeeee...`)
)
