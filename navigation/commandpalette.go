package navigation

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CommandPaletteOptions struct {
	Disabled bool
}

type CommandPaletteResult struct {
	Element *ui.Element
	// Ran is the command run this frame, or -1.
	Ran int
}

type paletteKey struct{}

// palette is the query, lit match and first shown row.
type palette struct {
	query    string
	lit, top int
}

// paletteRows is how many matches show at once.
const paletteRows = 8

// CommandPalette filters commands as you type; Cmd+K opens it.
func CommandPalette(c *ui.Context, text string, commands []MenuItem, open *bool, o CommandPaletteOptions) CommandPaletteResult {
	k := pixel.K(c)
	off := o.Disabled || pixel.Current(c).Inert
	if !off && c.Shortcut(ui.Cmd, ui.KeyK) {
		*open = !*open
	}
	if off {
		*open = false
	}
	home := ui.Box(c)
	st := ui.Local(home, paletteKey{}, func() palette { return palette{} })
	res := CommandPaletteResult{Element: home, Ran: -1}
	if !*open {
		*st = palette{}
		return res
	}
	shown := filter(commands, st.query)
	run := func(i int) {
		if !commands[i].Disabled {
			res.Ran, *open = i, false
		}
	}
	res.Element = ui.DialogBase(c, open, func(back, panel *ui.Element) {
		back.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Checker(0, 0, g.W, g.H, pixel.Back)
		})
		panel.Width(k.D(240)).Padding(k.D(6)).Gap(k.D(4)).Label(text).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
		}).Children(func() {
			f := input.TextInput(c, text, &st.query, input.TextInputOptions{Placeholder: "> …", Width: 228})
			f.Input.AutoFocus()
			if f.Changed {
				st.lit = 0
			}
			if f.Element.Shortcut(0, ui.KeyDown) {
				st.lit++
			}
			if f.Element.Shortcut(0, ui.KeyUp) {
				st.lit--
			}
			st.lit = max(0, min(st.lit, len(shown)-1))
			st.top = max(0, min(st.top, st.lit, len(shown)-paletteRows), st.lit-paletteRows+1)
			// A Shortcut, not Submitted: the frame builds again closed.
			if f.Element.Shortcut(0, ui.KeyEnter) && len(shown) > 0 {
				run(shown[st.lit])
			}
			ui.Column(c).AlignItems(ui.Stretch).Children(func() {
				for j := st.top; j < min(len(shown), st.top+paletteRows); j++ {
					it, lit := commands[shown[j]], j == st.lit
					row := ui.Row(c).Role(ui.RoleListItem).Label(it.Label).Height(k.D(pixel.LineHeight))
					if it.Disabled {
						row.Description(k.Msg("disabled"))
					}
					if row.Clicked() {
						run(shown[j])
					}
					row.Draw(func(p *ui.Painter, r ui.Rect) { menuRow(k.Canvas(p, r), it, lit) })
				}
			})
		})
	})
	if res.Ran >= 0 {
		slog.Debug("8bit: command palette", "label", text, "ran", commands[res.Ran].Label)
	}
	return res
}

// filter indexes commands holding q, prefix matches first.
func filter(commands []MenuItem, q string) []int {
	q = strings.ToLower(strings.TrimSpace(q))
	var starts, holds []int
	for i, it := range commands {
		switch l := strings.ToLower(it.Label); {
		case it.Separator:
		case strings.HasPrefix(l, q):
			starts = append(starts, i)
		case strings.Contains(l, q):
			holds = append(holds, i)
		}
	}
	return append(starts, holds...)
}
