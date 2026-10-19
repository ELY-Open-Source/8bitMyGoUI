package desk

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SnakeOptions struct {
	// Cols and Rows are zeros for 16 by 12.
	Cols, Rows int
	// Seed places the food; zero takes the clock.
	Seed     uint64
	Disabled bool
}

type SnakeResult struct {
	Element ui.Element
	// Score is food eaten; Started and Over are this frame's.
	Score         int
	Started, Over bool
}

type snakeKey struct{}

// pit is a game: body head first, heading, food, score.
type pit struct {
	body          [][2]int
	dir, next     [2]int
	food          [2]int
	score, last   int
	playing, dead bool
	rng           *rand.Rand
}

const scell = 8

// Snake: arrows steer, food grows it; crashes are GAME OVER.
func Snake(c *ui.Context, text string, o SnakeOptions) SnakeResult {
	k := pixel.K(c)
	cols, rows := pixel.Or("Snake Cols", o.Cols, 16), pixel.Or("Snake Rows", o.Rows, 12)
	if cols < 6 || rows < 4 {
		panic(fmt.Sprintf("desk: Snake %d×%d is too small", cols, rows))
	}
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	s := ui.Local(shell, snakeKey{}, func() pit { return pit{} })
	res := SnakeResult{Element: shell}
	if s.rng == nil || len(s.body) == 0 {
		s.reset(cols, rows, o.Seed, uint64(k.Now().UnixNano()))
	}
	if !base.IsDisabled() {
		for _, m := range []struct {
			key ui.Key
			d   [2]int
		}{{ui.KeyUp, [2]int{0, -1}}, {ui.KeyDown, [2]int{0, 1}}, {ui.KeyLeft, [2]int{-1, 0}}, {ui.KeyRight, [2]int{1, 0}}} {
			// A turn back onto its own neck is ignored.
			if pixel.Shortcut(base, 0, m.key) && m.d != [2]int{-s.dir[0], -s.dir[1]} && !s.dead {
				res.Started = res.Started || !s.playing
				s.next, s.playing = m.d, true
			}
		}
		// Enter, Space or a click starts it, or CONTINUEs.
		enter, space := pixel.Shortcut(base, 0, ui.KeyEnter), pixel.Shortcut(base, 0, ui.KeySpace)
		if enter || space || base.Clicked() {
			if s.dead {
				s.reset(cols, rows, o.Seed, uint64(k.Now().UnixNano()))
			}
			res.Started, s.playing = res.Started || !s.playing, true
		}
	}
	// It runs 8 steps a second, only while focused.
	if s.playing && !s.dead && base.Focused() && !base.IsDisabled() {
		tick := k.Tick(c, 8)
		if s.last >= 0 {
			for range min(tick-s.last, 3) {
				if s.step(cols, rows) {
					res.Over = true
					break
				}
			}
		}
		s.last = tick
	} else {
		s.last = -1
	}
	res.Score = s.score
	if res.Started {
		slog.Debug("8bit: snake start", "label", text)
	}
	if res.Over {
		slog.Debug("8bit: snake over", "label", text, "score", s.score)
	}
	word := ""
	switch {
	case s.dead:
		word = k.Msg("game_over")
	case !s.playing:
		word = k.Msg("press_start")
	case !base.Focused():
		word = k.Msg("paused")
	}
	score := fmt.Sprintf("%s %d", k.Msg("score"), s.score)
	base.Description(score + ", " + word)
	body, food, dead := slices.Clone(s.body), s.food, s.dead
	w, h := cols*scell+4, rows*scell+4
	shell.Size(k.D(w), k.D(h+pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Text(0, 0, pixel.Wrap(score, w, 1)[0], pixel.Ink, 1)
		g.Frame(0, pixel.LineHeight, w, h, pixel.Sunken, pixel.Well)
		at := func(c [2]int) (int, int) { return 2 + c[0]*scell, pixel.LineHeight + 2 + c[1]*scell }
		fx, fy := at(food)
		g.Rect(fx+2, fy+2, scell-4, scell-4, pixel.Bad)
		g.Dot(fx+4, fy+1, pixel.Good)
		for i, b := range body {
			x, y := at(b)
			ink := pixel.Good
			if dead {
				ink = pixel.Faint
			}
			g.Rect(x+1, y+1, scell-2, scell-2, ink)
			if i == 0 {
				g.Rect(x+2, y+2, 2, 2, pixel.Ink)
				g.Rect(x+5, y+2, 2, 2, pixel.Ink)
			}
		}
		// The word sits on top, clear of the snake.
		if word != "" {
			word := pixel.Wrap(word, w-8, 1)[0]
			ww := pixel.Measure(word)
			g.Rect((w-ww)/2-3, pixel.LineHeight+3, ww+6, pixel.LineHeight+2, pixel.Face)
			g.Text((w-ww)/2, pixel.LineHeight+4, word, pixel.Ink, 1)
		}
		g.Finish(st)
	})
	return res
}

// reset lays a three-long snake mid-left, heading right.
func (s *pit) reset(cols, rows int, seed, clock uint64) {
	if seed == 0 {
		seed = clock
	}
	y := rows / 2
	*s = pit{body: [][2]int{{3, y}, {2, y}, {1, y}}, dir: [2]int{1, 0}, next: [2]int{1, 0}, last: -1,
		rng: rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))}
	s.feed(cols, rows)
}

// feed drops food on a free cell.
func (s *pit) feed(cols, rows int) {
	var free [][2]int
	for y := range rows {
		for x := range cols {
			if !slices.Contains(s.body, [2]int{x, y}) {
				free = append(free, [2]int{x, y})
			}
		}
	}
	if len(free) > 0 {
		s.food = free[s.rng.IntN(len(free))]
	}
}

// step moves the head on; it reports a crash.
func (s *pit) step(cols, rows int) bool {
	s.dir = s.next
	head := [2]int{s.body[0][0] + s.dir[0], s.body[0][1] + s.dir[1]}
	grow := head == s.food
	rest := s.body
	if !grow {
		rest = s.body[:len(s.body)-1]
	}
	if head[0] < 0 || head[0] >= cols || head[1] < 0 || head[1] >= rows || slices.Contains(rest, head) {
		s.dead = true
		return true
	}
	s.body = append([][2]int{head}, rest...)
	if grow {
		s.score++
		s.feed(cols, rows)
	}
	return false
}
