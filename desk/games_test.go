package desk_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/desk"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestMinesweeperFirstClickIsSafe(t *testing.T) {
	won := 0
	tt := view(600, 500, newWorld(), func(c *ui.Context) {
		// Seven mines fill all but the click's 3×3: it wins.
		if desk.Minesweeper(c, "Mines", desk.MinesweeperOptions{Cols: 4, Rows: 4, Mines: 7, Seed: 3}).Won {
			won++
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if won != 1 {
		t.Fatalf("won %d times: %q", won, tt.Texts())
	}
	find(t, tt, "STAGE CLEAR!")
}

func TestMinesweeperFlagsTimesAndResets(t *testing.T) {
	wd := newWorld()
	off := false
	tt := view(600, 500, wd, func(c *ui.Context) {
		desk.Minesweeper(c, "Mines", desk.MinesweeperOptions{Mines: 3, Seed: 5, Disabled: off})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyF)
	find(t, tt, "002")
	r := find(t, tt, "MINEFIELD")
	// The flag flies in Bad on the first cell.
	if fx, fy, _, _ := dots(r); count(tt.Image(), fx, fy, 16, 16, pixel.Bad) == 0 {
		t.Fatal("no Bad flag on the first cell")
	}
	tt.RightClickAt(r.X+r.W-4, r.Y+r.H-4)
	find(t, tt, "001")
	// A flag lifts again; the first open starts the clock.
	tt.Key(0, ui.KeyF)
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyEnter)
	later(tt, wd, 5*time.Second)
	find(t, tt, "005")
	click(t, tt, "NEW GAME")
	find(t, tt, "003")
	find(t, tt, "000")
	off = true
	tt.Frame()
	tt.Key(0, ui.KeyF)
	find(t, tt, "003")
}

func TestSnakeRunsCrashesAndContinues(t *testing.T) {
	wd := newWorld()
	over, score := 0, -1
	tt := view(400, 300, wd, func(c *ui.Context) {
		r := desk.Snake(c, "Snake", desk.SnakeOptions{Cols: 6, Rows: 4, Seed: 9})
		score = r.Score
		if r.Over {
			over++
		}
	})
	x, y, w, h := dots(find(t, tt, "Snake"))
	frozen := append([]byte(nil), tt.Image().Pix...)
	later(tt, wd, time.Second)
	if string(frozen) != string(tt.Image().Pix) {
		t.Fatal("it moved before it started")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	// A turn back on its neck is ignored; it lives.
	tt.Key(0, ui.KeyLeft)
	later(tt, wd, 250*time.Millisecond)
	if over != 0 {
		t.Fatal("turning back on itself crashed it")
	}
	// Heading right from the fourth column, the wall comes fast.
	later(tt, wd, time.Second)
	if over != 1 || count(tt.Image(), x, y+16, w, h-16, pixel.Good) > 0 {
		t.Fatalf("over %d times, or the body is still Good", over)
	}
	later(tt, wd, time.Second)
	tt.Key(0, ui.KeyEnter)
	if over != 1 || score != 0 || count(tt.Image(), x, y+16, w, h-16, pixel.Good) == 0 {
		t.Fatalf("after CONTINUE: over %d, score %d", over, score)
	}
}

func TestSnakePausesWithoutFocus(t *testing.T) {
	wd := newWorld()
	tt := view(400, 300, wd, func(c *ui.Context) {
		desk.Snake(c, "Snake", desk.SnakeOptions{Seed: 9})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.SetFocused(false)
	tt.Frame()
	before := append([]byte(nil), tt.Image().Pix...)
	later(tt, wd, time.Second)
	if string(before) != string(tt.Image().Pix) {
		t.Fatal("it ran on without the keyboard")
	}
}

func TestDiskDefragSortsWhileRunning(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		running, done := false, 0
		tt := view(600, 300, wd, func(c *ui.Context) {
			if desk.DiskDefrag(c, "Disk", running, desk.DiskDefragOptions{Blocks: 64, Seed: 4}).Done {
				done++
			}
		})
		x, y, w, _ := dots(find(t, tt, "Disk"))
		// Inside the frame, Lit is a fragment or a read.
		frag := func() int { return count(tt.Image(), x+2, y+2, w-4, 12, pixel.Lit) }
		start := frag()
		later(tt, wd, 10*time.Second)
		if frag() != start || done != 0 {
			t.Fatalf("reduce %v: it sorted while stopped", reduce)
		}
		// A loader keeps time, so Reduce Motion runs it too.
		running = true
		tt.Frame()
		for range 20 {
			later(tt, wd, time.Second)
		}
		if done != 1 || frag() != 0 || start == 0 {
			t.Fatalf("reduce %v: done %d, %d fragmented dots of %d", reduce, done, frag(), start)
		}
	}
}
