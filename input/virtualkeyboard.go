package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type VirtualKeyboardOptions struct {
	Disabled bool
}

type VirtualKeyboardResult struct {
	Element            ui.Element
	Changed, Submitted bool
}

type keyboardKey struct{}

// caps are ten keys a row; \b erases, \n enters.
var caps = []rune("1234567890QWERTYUIOPASDFGHJKL\bZXCVBNM- \n")

var (
	backKey = pixel.MustSprite(`
		..#####
		.#....#
		#..#.##
		.#....#
		..#####
	`)
	enterKey = pixel.MustSprite(`
		.....#
		.#...#
		##...#
		######
		##....
		.#....
	`)
)

// VirtualKeyboard types into *value from keys on screen.
func VirtualKeyboard(c *ui.Context, text string, value *string, o VirtualKeyboardOptions) VirtualKeyboardResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(4))
	cur := ui.Local(col, keyboardKey{}, func() int { return 0 })
	res := VirtualKeyboardResult{Element: col}
	col.Children(func() {
		ui.Box(c).Size(k.D(140), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			well(g, pixel.Edge, false)
			rs := []rune(*value)
			for len(rs) > 0 && pixel.Measure(string(rs))+10 > g.W {
				rs = rs[1:]
			}
			g.Text(3, 2, string(rs), pixel.Ink, 1)
			g.Rect(4+pixel.Measure(string(rs)), 4, 5, 12, pixel.Accent)
		})
		shell, base, st, picked := pixel.Pad(c, text, 10, len(caps), 14, pixel.LineHeight, cur, o.Disabled)
		if picked >= 0 {
			switch r := caps[picked]; r {
			case '\b':
				if rs := []rune(*value); len(rs) > 0 {
					*value, res.Changed = string(rs[:len(rs)-1]), true
				}
			case '\n':
				res.Submitted = true
			default:
				*value, res.Changed = *value+string(r), true
			}
		}
		base.Description(*value)
		focused := base.Focused()
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, cp := range caps {
				x, y := i%10*14, i/10*pixel.LineHeight
				fill := pixel.Face
				if focused && i == *cur {
					fill = pixel.Lit
				}
				g.Frame(x, y, 13, pixel.LineHeight-1, pixel.Raised, fill)
				switch cp {
				case '\b':
					g.Sprite(x+3, y+5, backKey, 0, pixel.Ink, 1)
				case '\n':
					g.Sprite(x+4, y+4, enterKey, 0, pixel.Accent, 1)
				case ' ':
					g.Rect(x+3, y+10, 7, 1, pixel.Ink)
					g.Rect(x+3, y+8, 1, 2, pixel.Ink)
					g.Rect(x+9, y+8, 1, 2, pixel.Ink)
				default:
					l := string(cp)
					g.Text(x+(13-pixel.Measure(l))/2, y, l, pixel.Ink, 1)
				}
			}
			g.Finish(st)
		})
	})
	if res.Submitted {
		slog.Debug("8bit: keyboard entered", "label", text)
	}
	return res
}
