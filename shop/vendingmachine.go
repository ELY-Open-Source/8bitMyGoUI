package shop

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Snack is a slot in the machine; nameless means empty.
type Snack struct {
	Name string
	// Sprite is at most 16 dots square.
	Sprite *pixel.Sprite
	Price  int
	Stock  int
}

type VendingMachineOptions struct {
	Disabled bool
}

type VendingMachineResult struct {
	Element ui.Element
	// Vended is the slot dropped this frame, else -1.
	Vended int
}

type vendKey struct{}

// vend is the code typed, the line shown, the drops.
type vend struct {
	code, msg  string
	cur, drops int
	last       int
}

var keys = []string{"A", "B", "C", "D", "1", "2", "3", "4", "CLR", "OK"}

const vcell = 30

// VendingMachine sells slots A1 to D4 by code; credit pays.
func VendingMachine(c *ui.Context, text string, slots [16]Snack, credit *int, o VendingMachineOptions) VendingMachineResult {
	k := pixel.K(c)
	amount("VendingMachine credit", *credit)
	said := make([]string, 0, 16)
	for i, s := range slots {
		if s.Name == "" {
			continue
		}
		amount("VendingMachine price", s.Price)
		if s.Sprite == nil || s.Sprite.W > 16 || s.Sprite.H > 16 || s.Stock < 0 {
			panic(fmt.Sprintf("shop: VendingMachine slot %d %q needs a small sprite and stock", i, s.Name))
		}
		said = append(said, code(i)+" "+s.Name+" "+strconv.Itoa(s.Price))
	}
	row := ui.Row(c).Gap(k.D(8)).AlignItems(ui.Start).Role(ui.RoleGroup).Label(text)
	v := ui.Local(row, vendKey{}, func() vend { return vend{last: -1} })
	res := VendingMachineResult{Vended: -1}
	res.Element = row.Children(func() {
		glass := ui.Box(c).Size(k.D(4*vcell+4), k.D(4*vcell+28)).Role(ui.RoleImage).Label(k.Msg("snacks")).Description(strings.Join(said, ", "))
		// The bought snack bounces into the tray in steps.
		bounce, last := k.Shot(c, glass, v.drops, 8, 4), v.last
		glass.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
			for i, s := range slots {
				x, y := 2+i%4*vcell, 2+i/4*vcell
				g.Micro(x+2, y+2, code(i), pixel.Faint, 1)
				if s.Name == "" {
					continue
				}
				g.Sprite(x+(vcell-s.Sprite.W)/2, y+8, s.Sprite, 0, pixel.Ink, 1)
				p := strconv.Itoa(s.Price)
				g.Micro(x+vcell-2-pixel.MicroWidth(p), y+vcell-7, p, pixel.Warn, 1)
				if s.Stock == 0 {
					g.Veil(x+1, y+1, vcell-2, vcell-2, pixel.Shade, 8)
				}
			}
			ty := 4*vcell + 6
			g.Rect(4, ty, g.W-8, 18, pixel.Shade)
			if last >= 0 && slots[last].Sprite != nil {
				sp := slots[last].Sprite
				dy := []int{-6, 0, -2, 0}[min(bounce, 3)]
				g.Sprite((g.W-sp.W)/2, ty+18-sp.H-1+dy, sp, 0, pixel.Ink, 1)
			}
		})
		ui.Column(c).Gap(k.D(4)).Children(func() {
			line := v.msg
			if line == "" {
				line = v.code
			}
			ui.Box(c).Size(k.D(4*22), k.D(pixel.LineHeight+4)).Role(ui.RoleStatus).Label(line).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
				g.Text(4, 2, pixel.Wrap(line, g.W-8, 1)[0], pixel.Good, 1)
			})
			shell, base, st, picked := pixel.Pad(c, k.Msg("keypad"), 4, len(keys), 22, 18, &v.cur, o.Disabled)
			if picked >= 0 {
				res.Vended = press(k, v, keys[picked], slots, credit)
			}
			say := keys[v.cur]
			switch say {
			case "CLR":
				say = k.Msg("clear")
			case "OK":
				say = k.Msg("ok")
			}
			base.Description(say)
			cur := v.cur
			shell.Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i, key := range keys {
					x, y := i%4*22, i/4*18
					g.Frame(x+1, y+1, 20, 16, pixel.Raised, pixel.Face)
					if i == cur && st.Focus {
						g.Cursor(x, y, 22, 18)
					}
					// CLR and OK are a cross and a check.
					switch key {
					case "CLR":
						g.Sprite(x+5, y+3, icon.Close, 0, pixel.Bad, 1)
					case "OK":
						g.Sprite(x+5, y+3, icon.Check, 0, pixel.Good, 1)
					default:
						g.Micro(x+(22-pixel.MicroWidth(key))/2, y+7, key, pixel.Ink, 1)
					}
				}
				g.Finish(st)
			})
			money(c, k, Gold, *credit, pixel.Ink)
		})
	})
	if res.Vended >= 0 {
		slog.Debug("8bit: vended", "label", text, "slot", code(res.Vended))
	}
	return res
}

// code names slot i: A1 to D4.
func code(i int) string {
	return string(rune('A'+i/4)) + strconv.Itoa(i%4+1)
}

// press takes a key; OK on a full code vends.
func press(k pixel.Kit, v *vend, key string, slots [16]Snack, credit *int) int {
	switch {
	case len(key) == 1 && key[0] >= 'A' && key[0] <= 'D':
		v.code, v.msg = key, ""
	case len(key) == 1 && len(v.code) == 1:
		v.code += key
	case key == "CLR":
		v.code, v.msg = "", ""
	case key == "OK" && len(v.code) == 2:
		i := int(v.code[0]-'A')*4 + int(v.code[1]-'1')
		s := slots[i]
		v.code = ""
		switch {
		case s.Name == "" || s.Stock == 0:
			v.msg = k.Msg("sold_out")
		case *credit < s.Price:
			v.msg = k.Msg("insert_coin")
		default:
			*credit -= s.Price
			v.msg, v.drops, v.last = k.Msg("enjoy"), v.drops+1, i
			return i
		}
	case key == "OK":
		v.msg = k.Msg("code_prompt")
	}
	return -1
}
