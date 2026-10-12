package desk

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DiskDefragOptions struct {
	// Blocks on the disk; zero means 160.
	Blocks int
	// Seed scatters the files; zero takes the clock.
	Seed uint64
}

type DiskDefragResult struct {
	Element *ui.Element
	// Done is the frame the last block falls in place.
	Done bool
}

type defragKey struct{}

// disk is the first layout, its swaps and moves made.
type disk struct {
	from        []int
	plan        [][2]int
	moves, last int
	done        bool
}

// files color the six files on the disk.
var files = []pixel.Slot{pixel.Info, pixel.Good, pixel.Magic, pixel.Accent, pixel.Warn, pixel.Bad}

const dblock, dcols = 6, 32

// DiskDefrag moves scattered blocks into order while running.
func DiskDefrag(c *ui.Context, text string, running bool, o DiskDefragOptions) DiskDefragResult {
	k := pixel.K(c)
	n := pixel.Or("DiskDefrag Blocks", o.Blocks, 160)
	e := ui.Box(c).Role(ui.RoleProgress).Label(text)
	d := ui.Local(e, defragKey{}, func() disk {
		seed := o.Seed
		if seed == 0 {
			seed = uint64(k.Now().UnixNano())
		}
		from := scatter(n, seed)
		return disk{from: from, plan: swaps(from), last: -1}
	})
	if len(d.from) != n {
		panic(fmt.Sprintf("desk: DiskDefrag went from %d blocks to %d", len(d.from), n))
	}
	// A loader keeps time: it runs on under Reduce Motion.
	if running && d.moves < len(d.plan) {
		tick := k.Tick(c, 15)
		// A clock stepping back moves nothing; a long gap, four.
		if d.last >= 0 {
			d.moves = min(len(d.plan), d.moves+max(0, min(tick-d.last, 4)))
		}
		d.last = tick
	} else {
		d.last = -1
	}
	res := DiskDefragResult{Element: e, Done: d.moves == len(d.plan) && !d.done}
	d.done = d.moves == len(d.plan)
	cells := slices.Clone(d.from)
	for _, s := range d.plan[:d.moves] {
		cells[s[0]], cells[s[1]] = cells[s[1]], cells[s[0]]
	}
	want := slices.Clone(d.from)
	slices.SortFunc(want, func(a, b int) int { return order(a) - order(b) })
	read, write := -1, -1
	if running && d.moves > 0 && !d.done {
		read, write = d.plan[d.moves-1][0], d.plan[d.moves-1][1]
	}
	pct := 100
	if len(d.plan) > 0 {
		pct = 100 * d.moves / len(d.plan)
	}
	status := fmt.Sprintf(k.Msg("percent_done"), pct)
	e.Description(status)
	if res.Done {
		slog.Debug("8bit: defrag done", "label", text, "moves", d.moves)
	}
	rows := (n + dcols - 1) / dcols
	w, h := dcols*dblock+4, rows*dblock+4
	e.Size(k.D(w), k.D(h+pixel.LineHeight+2)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, w, h, pixel.Sunken, pixel.Well)
		for i, f := range cells {
			x, y := 2+i%dcols*dblock, 2+i/dcols*dblock
			switch {
			case i == read:
				g.Rect(x, y, dblock-1, dblock-1, pixel.Lit)
			case i == write:
				g.Rect(x, y, dblock-1, dblock-1, pixel.Ink)
			case f < 0:
				g.Dot(x+2, y+2, pixel.Shade)
			default:
				g.Rect(x, y, dblock-1, dblock-1, files[f])
				// A block not yet in its place is still fragmented.
				if f != want[i] {
					g.Checker(x, y, dblock-1, dblock-1, pixel.Lit)
				}
			}
		}
		g.Text(0, h+2, status, pixel.Ink, 1)
	})
	return res
}

// order sorts files first, by number, and free blocks last.
func order(f int) int {
	if f < 0 {
		return len(files)
	}
	return f
}

// scatter puts six files on two thirds of n blocks.
func scatter(n int, seed uint64) []int {
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	cells := make([]int, n)
	for i := range cells {
		cells[i] = -1
	}
	used := n * 2 / 3
	for j, at := range rng.Perm(n)[:used] {
		cells[at] = j * len(files) / max(1, used)
	}
	return cells
}

// swaps sorts the disk a block at a time.
func swaps(cells []int) [][2]int {
	cur, want := slices.Clone(cells), slices.Clone(cells)
	slices.SortFunc(want, func(a, b int) int { return order(a) - order(b) })
	var plan [][2]int
	for p := range cur {
		if cur[p] == want[p] {
			continue
		}
		for q := p + 1; q < len(cur); q++ {
			if cur[q] == want[p] {
				cur[p], cur[q] = cur[q], cur[p]
				plan = append(plan, [2]int{q, p})
				break
			}
		}
	}
	return plan
}
