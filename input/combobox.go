package input

import (
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ComboboxOptions struct {
	Placeholder string
	// Width is in dots; zero means 120.
	Width    int
	Disabled bool
}

type ComboboxResult struct {
	Element *ui.Element
	// Changed reports a new *value; Chosen, an option taken.
	Changed, Chosen bool
}

type comboKey struct{}

type combo struct {
	open, typed bool
	// lit is the lit option; ptr, the last pointer spot.
	lit int
	ptr [2]float32
	// top is the first row shown; wheel, unspent scroll.
	top   int
	wheel float32
}

// comboRows is how many options show at once.
const comboRows = 8

// Combobox lists matching options; Enter or a click takes one.
func Combobox(c *ui.Context, text string, value *string, options []string, o ComboboxOptions) ComboboxResult {
	k := pixel.K(c)
	ed, e := field(c, value, pixel.EditorOptions{Label: text, Placeholder: o.Placeholder, Disabled: o.Disabled}, o.Width, nil, nil, nil)
	ed.Input.Role(ui.RoleComboBox)
	st := ui.Local(e, comboKey{}, func() combo { return combo{} })
	res := ComboboxResult{Element: e, Changed: ed.Changed}
	take := func(opt string) {
		if *value != opt {
			*value, res.Changed = opt, true
		}
		st.open, st.typed, res.Chosen = false, false, true
	}
	if ed.Changed {
		st.open, st.typed, st.lit = true, true, 0
	}
	// Keys reach the wrapper; the editor's state gates them.
	on := !ed.Input.IsDisabled()
	if on && pixel.Shortcut(e, 0, ui.KeyDown) {
		if st.open {
			st.lit++
		} else {
			st.open, st.typed, st.lit = true, false, 0
		}
	}
	if st.open && on && pixel.Shortcut(e, 0, ui.KeyUp) {
		st.lit--
	}
	shown := options
	if st.typed {
		shown = matching(options, *value)
	}
	st.lit = max(0, min(st.lit, len(shown)-1))
	// The rows shown follow the lit option.
	st.top = max(0, min(st.top, st.lit, len(shown)-comboRows), st.lit-comboRows+1)
	if st.open && on && len(shown) > 0 && pixel.Shortcut(e, 0, ui.KeyEnter) {
		take(shown[st.lit])
	}
	if !on {
		st.open = false
	}
	show, held := st.open && len(shown) > 0, false
	panel := ui.PopoverBase(c, e, &show, func(panel *ui.Element) {
		panel.MinWidth(e.Bounds().W)
		panel.HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind != ui.InputScroll {
				return false
			}
			st.wheel += ev.DY
			turn(&st.wheel, k.D(pixel.LineHeight), func(d int) { st.lit = max(0, min(len(shown)-1, st.lit+d)) })
			return true
		})
		end := min(len(shown), st.top+comboRows)
		panel.DrawOver(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			if st.top > 0 {
				g.Sprite(g.W-8, 2, tinyUp, 0, pixel.Faint, 1)
			}
			if end < len(shown) {
				g.Sprite(g.W-8, g.H-5, tinyDown, 0, pixel.Faint, 1)
			}
		})
		// Only a pointer that moves lights the option under it.
		px, py, over := panel.PointerPosition()
		moved := over && [2]float32{px, py} != st.ptr
		st.ptr = [2]float32{px, py}
		menu(c, panel, func() {
			items := make([]*ui.Element, len(shown))
			for i := st.top; i < end; i++ {
				opt := shown[i]
				items[i] = ui.Row(c).Role(ui.RoleListItem).Label(opt)
				if items[i].Hovered() && moved {
					st.lit = i
				}
				if items[i].Pressed() {
					held = true
				}
				if items[i].Clicked() {
					take(opt)
				}
			}
			for i := st.top; i < end; i++ {
				row(k, items[i], shown[i], i == st.lit, shown[i] == *value)
			}
		})
	})
	if st.open && len(shown) > 0 && !show {
		st.open = false
	}
	// Losing focus closes it, unless pressing on the list.
	if panel != nil && panel.Pressed() {
		held = true
	}
	if panel != nil && st.open && !ed.Focused && !held {
		st.open = false
		// This frame built the list: hide it, rebuild without it.
		panel.Invisible()
		c.Invalidate()
	}
	if res.Chosen {
		slog.Debug("8bit: combobox chosen", "label", text, "runes", utf8.RuneCountInString(*value))
	}
	return res
}

// matching are options holding q, prefix matches first.
func matching(options []string, q string) []string {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return options
	}
	var starts, holds []string
	for _, opt := range options {
		switch l := strings.ToLower(opt); {
		case strings.HasPrefix(l, q):
			starts = append(starts, opt)
		case strings.Contains(l, q):
			holds = append(holds, opt)
		}
	}
	return append(starts, holds...)
}
