package desk

import (
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CalculatorOptions struct {
	Disabled bool
}

type CalculatorResult struct {
	Element *ui.Element
	// Changed is a key pressed this frame.
	Changed bool
}

type calcKey struct{}

// calc is the entry, total, waiting sign and cursor.
type calc struct {
	entry, op string
	acc       float64
	over      bool
	cur       int
}

var calcKeys = []string{"7", "8", "9", "÷", "4", "5", "6", "×", "1", "2", "3", "-", "C", "0", "=", "+"}

// signs are keyboard keys for the operator keys.
var signs = []struct {
	mods ui.Modifiers
	key  ui.Key
	sign string
}{{0, ui.KeyMinus, "-"}, {0, ui.KeySlash, "÷"}, {0, ui.KeyEqual, "="}, {ui.Shift, ui.KeyEqual, "+"}, {ui.Shift, ui.Key8, "×"}, {0, ui.KeyC, "C"}}

const keyW, keyH = 24, 20

// Calculator works: digits, signs, = and C; ÷0 ends it.
func Calculator(c *ui.Context, text string, o CalculatorOptions) CalculatorResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(4))
	cl := ui.Local(col, calcKey{}, func() calc { return calc{} })
	res := CalculatorResult{}
	res.Element = col.Children(func() {
		screen := ui.Box(c).Size(k.D(4*keyW), k.D(pixel.LineHeight+6)).Role(ui.RoleStatus)
		shell, base, st, picked := pixel.Pad(c, text, 4, len(calcKeys), keyW, keyH, &cl.cur, o.Disabled)
		// Digits and signs typed on the keyboard press their keys.
		if picked < 0 && !base.IsDisabled() && !pixel.Typing(c) {
			for i := range 10 {
				if pixel.Shortcut(base, 0, ui.Key0+ui.Key(i)) {
					picked = indexOf(strconv.Itoa(i))
				}
			}
			for _, s := range signs {
				if pixel.Shortcut(base, s.mods, s.key) {
					picked = indexOf(s.sign)
				}
			}
		}
		if picked >= 0 {
			cl.press(calcKeys[picked])
			res.Changed = true
			slog.Debug("8bit: calculator", "label", text, "key", calcKeys[picked])
		}
		shown, ink := cl.display(), pixel.Ink
		if cl.over {
			shown, ink = k.Msg("game_over"), pixel.Bad
		}
		screen.Label(shown).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
			s := pixel.Wrap(shown, g.W-8, 1)[0]
			g.Text(g.W-4-pixel.Measure(s), 3, s, ink, 1)
		})
		base.Description(calcKeys[cl.cur])
		at, over := cl.cur, cl.over
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, key := range calcKeys {
				x, y := i%4*keyW, i/4*keyH
				fill, ink := pixel.Face, pixel.Ink
				// After GAME OVER, C glows: it is the CONTINUE.
				if key == "C" && over {
					fill, ink = pixel.Accent, pixel.AccentInk
				}
				g.Frame(x+1, y+1, keyW-2, keyH-2, pixel.Raised, fill)
				if i == at && st.Focus {
					g.Cursor(x, y, keyW, keyH)
					ink = pixel.Ink
				}
				g.Text(x+(keyW-pixel.Measure(key))/2, y+2, key, ink, 1)
			}
			g.Finish(st)
		})
	})
	return res
}

// indexOf is a key's place on the pad.
func indexOf(key string) int {
	for i, k := range calcKeys {
		if k == key {
			return i
		}
	}
	panic("desk: no calculator key " + key)
}

// press takes a key: digits build, signs fold in.
func (cl *calc) press(key string) {
	switch key {
	case "C":
		*cl = calc{cur: cl.cur}
	case "=":
		cl.fold()
		cl.op = ""
	case "+", "-", "×", "÷":
		cl.fold()
		cl.op = key
	default:
		switch {
		case cl.entry == "0":
			cl.entry = key
		case len(cl.entry) < 12:
			cl.entry += key
		}
	}
}

// fold applies the waiting sign to the entry.
func (cl *calc) fold() {
	if cl.entry == "" {
		return
	}
	v := 0.0
	for _, r := range cl.entry {
		v = v*10 + float64(r-'0')
	}
	cl.entry = ""
	switch cl.op {
	case "":
		cl.acc = v
	case "+":
		cl.acc += v
	case "-":
		cl.acc -= v
	case "×":
		cl.acc *= v
	case "÷":
		cl.acc /= v
	}
	// ÷0 and past twelve digits are GAME OVER.
	if math.IsNaN(cl.acc) || math.Abs(cl.acc) >= 1e12 {
		cl.over = true
	}
}

// display is the entry while typing, else the total.
func (cl *calc) display() string {
	if cl.entry != "" {
		return cl.entry
	}
	v := cl.acc
	whole := len(strconv.Itoa(int(math.Abs(v))))
	s := strconv.FormatFloat(v, 'f', max(0, 11-whole), 64)
	if strings.Contains(s, ".") {
		s = strings.TrimSuffix(strings.TrimRight(s, "0"), ".")
	}
	if s == "-0" {
		return "0"
	}
	return s
}
