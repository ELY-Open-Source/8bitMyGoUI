package code

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Tool is what a press on the canvas does.
type Tool int

const (
	Pen Tool = iota
	Fill
	Erase
	Pick
)

type PixelEditorOptions struct {
	// Zoom is dots per sprite dot, 2–16; zero means 8.
	Zoom     int
	Disabled bool
}

type PixelEditorResult struct {
	Element *ui.Element
	// Changed is the art; Copied is COPY.
	Changed, Copied bool
}

type brushKey struct{}

// brush is tool, slot, cursor, last dab, grid, art size.
type brush struct {
	tool, slot, cur, last, w, h int
	grid                        bool
}

// PixelEditor paints *art, one frame of string art.
func PixelEditor(c *ui.Context, text string, art *string, o PixelEditorOptions) PixelEditorResult {
	k := pixel.K(c)
	z := pixel.Or("PixelEditor Zoom", o.Zoom, 8)
	sp := pixel.MustSprite(*art)
	if z < 2 || z > 16 || sp.Frames() != 1 {
		panic(fmt.Sprintf("code: PixelEditor of %d frames at zoom %d", sp.Frames(), z))
	}
	w, h := sp.W, sp.H
	// cells hold a slot, or -1 where clear.
	cells := make([]int, w*h)
	for i := range cells {
		cells[i] = -1
		if s, ok := sp.At(0, i%w, i/w, pixel.Ink); ok {
			cells[i] = int(s)
		}
	}
	var res PixelEditorResult
	res.Element = ui.Row(c).Gap(k.D(6)).AlignItems(ui.Start)
	b := ui.Local(res.Element, brushKey{}, func() brush { return brush{slot: int(pixel.Ink), last: -1, grid: true} })
	// New art drops a dab measured on the old.
	if b.w != w || b.h != h {
		b.w, b.h, b.last, b.cur = w, h, -1, min(b.cur, w*h-1)
	}
	apply := func(i int) bool {
		switch Tool(b.tool) {
		case Pen, Erase:
			v := b.slot
			if Tool(b.tool) == Erase {
				v = -1
			}
			if cells[i] == v {
				return false
			}
			cells[i] = v
			return true
		case Fill:
			return flood(cells, w, h, i, b.slot)
		}
		// The picker takes a slot, then hands back the pen.
		if cells[i] >= 0 {
			b.slot, b.tool = cells[i], int(Pen)
			slog.Debug("8bit: pixel editor picked", "label", text, "slot", pixel.Slot(b.slot))
		}
		return false
	}
	res.Element.Children(func() {
		shell, base, st, picked := pixel.Pad(c, text, w, w*h, z, z, &b.cur, o.Disabled)
		if !base.IsDisabled() {
			// Pen and eraser lay dots along a drag.
			tool := Tool(b.tool)
			if x, y, _ := base.PointerPosition(); base.Pressed() && (tool == Pen || tool == Erase) {
				cx, cy := int(x/k.Dot())/z, int(y/k.Dot())/z
				if i := cy*w + cx; cx >= 0 && cy >= 0 && cx < w && cy < h && i != b.last {
					// A fast drag still lays every dot between.
					from := i
					if b.last >= 0 {
						from = b.last
					}
					steps(from%w, from/w, cx, cy, func(j int) { res.Changed = apply(j) || res.Changed }, w)
					b.last = i
				}
			} else {
				b.last = -1
			}
			if picked >= 0 {
				res.Changed = apply(picked) || res.Changed
			}
		}
		v := cells[b.cur]
		say := "clear"
		if v >= 0 {
			say = pixel.Slot(v).String()
		}
		base.Description(fmt.Sprintf("%d, %d: %s", b.cur%w, b.cur/w, say))
		canvas(shell, k, st, cells, w, z, b.cur, b.grid)
		tools(c, k, text, b, art, o.Disabled, &res)
	})
	if res.Changed {
		*art = pixel.SpriteOf(w, h, func(x, y int) (pixel.Slot, bool) {
			v := cells[y*w+x]
			return pixel.Slot(max(v, 0)), v >= 0
		}).Art()
		slog.Debug("8bit: pixel editor painted", "label", text, "tool", b.tool, "slot", pixel.Slot(b.slot))
	}
	return res
}

// canvas paints cells, checkers where clear, grid and cursor.
func canvas(shell *ui.Element, k pixel.Kit, st pixel.State, cells []int, w, z, cur int, grid bool) {
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, v := range cells {
			x, y := i%w*z, i/w*z
			if v < 0 {
				g.Rect(x, y, z, z, pixel.Lit)
				g.Rect(x, y, z/2, z/2, pixel.Face)
				g.Rect(x+z/2, y+z/2, z-z/2, z-z/2, pixel.Face)
				continue
			}
			g.Rect(x, y, z, z, pixel.Slot(v))
		}
		if grid {
			for x := z; x < g.W; x += z {
				g.VLine(x, 0, g.H, pixel.Shade)
			}
			for y := z; y < g.H; y += z {
				g.HLine(0, y, g.W, pixel.Shade)
			}
		}
		if st.Focus {
			x, y := cur%w*z, cur/w*z
			g.Rect(x, y, z, 1, pixel.Accent)
			g.Rect(x, y+z-1, z, 1, pixel.Accent)
			g.Rect(x, y, 1, z, pixel.Accent)
			g.Rect(x+z-1, y, 1, z, pixel.Accent)
		}
		g.Finish(st)
	})
}

// tools are the tool switch, slot pad, GRID and COPY.
func tools(c *ui.Context, k pixel.Kit, text string, b *brush, art *string, disabled bool, res *PixelEditorResult) {
	ui.Column(c).Gap(k.D(4)).Children(func() {
		names := []string{k.Msg("pen"), k.Msg("fill"), k.Msg("erase"), k.Msg("pick")}
		if input.Segmented(c, k.Msg("tool"), names, &b.tool, input.SegmentedOptions{Disabled: disabled}).Changed {
			slog.Debug("8bit: pixel editor tool", "label", text, "tool", names[b.tool])
		}
		old := b.slot
		// The pad's cursor is the slot laid down.
		shell, base, st, _ := pixel.Pad(c, k.Msg("color"), 8, int(pixel.Slots), 10, 10, &b.slot, disabled)
		base.Description(pixel.Slot(b.slot).String())
		if b.slot != old {
			slog.Debug("8bit: pixel editor slot", "label", text, "slot", pixel.Slot(b.slot))
		}
		at := b.slot
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for s := range int(pixel.Slots) {
				x, y := s%8*10, s/8*10
				if s == at {
					g.Rect(x, y, 10, 10, pixel.Accent)
				}
				g.Frame(x+1, y+1, 8, 8, pixel.Flat, pixel.Slot(s))
			}
			g.Finish(st)
		})
		if input.Toggle(c, k.Msg("grid"), &b.grid, input.ToggleOptions{Disabled: disabled}).Changed {
			slog.Debug("8bit: pixel editor grid", "label", text, "on", b.grid)
		}
		if input.Button(c, k.Msg("copy"), input.ButtonOptions{Disabled: disabled}).Clicked {
			c.WriteClipboard(*art)
			res.Copied = true
			slog.Debug("8bit: pixel editor copied", "label", text)
		}
	})
}

// flood lays slot over the like cells around i.
func flood(cells []int, w, h, i, slot int) bool {
	from := cells[i]
	if from == slot {
		return false
	}
	stack := []int{i}
	for len(stack) > 0 {
		j := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cells[j] != from {
			continue
		}
		cells[j] = slot
		x, y := j%w, j/w
		for _, n := range [4][2]int{{x - 1, y}, {x + 1, y}, {x, y - 1}, {x, y + 1}} {
			if n[0] >= 0 && n[1] >= 0 && n[0] < w && n[1] < h {
				stack = append(stack, n[1]*w+n[0])
			}
		}
	}
	return true
}

// steps calls at for each cell on a line.
func steps(x0, y0, x1, y1 int, at func(i int), w int) {
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := sign(x1-x0), sign(y1-y0)
	for e := dx + dy; ; {
		at(y0*w + x0)
		if x0 == x1 && y0 == y1 {
			return
		}
		if e2 := 2 * e; e2 >= dy {
			e, x0 = e+dy, x0+sx
		} else {
			e, y0 = e+dx, y0+sy
		}
	}
}

func abs(v int) int { return max(v, -v) }

func sign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}
