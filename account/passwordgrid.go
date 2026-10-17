package account

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// symbols are what an old cartridge password is spelled with.
const symbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!?-."

type PasswordGridOptions struct {
	// Length is the password's symbols, 4–24; zero means 12.
	Length   int
	Disabled bool
}

type PasswordGridResult struct {
	Element ui.Element
	// Submitted is END on a full password.
	Changed, Submitted bool
}

type gridKey struct{}

// PasswordGrid spells *password by picking symbols from a grid.
func PasswordGrid(c *ui.Context, text string, password *string, o PasswordGridOptions) PasswordGridResult {
	k := pixel.K(c)
	n := pixel.Or("PasswordGrid Length", o.Length, 12)
	if n < 4 || n > 24 || len(*password) > n || strings.Trim(*password, symbols) != "" {
		panic(fmt.Sprintf("account: PasswordGrid of %d symbols holds %d bad or extra", n, len(*password)))
	}
	old := *password
	var res PasswordGridResult
	res.Element = ui.Column(c).Gap(k.D(5)).Role(ui.RoleGroup).Label(text)
	cur := ui.Local(res.Element, gridKey{}, func() int { return 0 })
	res.Element.Children(func() {
		typed := *password
		ui.Box(c).Size(k.D(n*9+(n-1)/4*4), k.D(13)).Role(ui.RoleText).Label(fmt.Sprintf("%s %d/%d", text, len(typed), n)).Description(typed).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			// Blanks in fours, as the old screens printed them.
			for i := range n {
				x, line := i*9+i/4*4, pixel.Faint
				if i < len(typed) {
					g.Micro(x+2, 4, typed[i:i+1], pixel.Ink, 1)
				}
				if i == len(typed) {
					line = pixel.Accent
				}
				g.HLine(x+1, 11, 6, line)
			}
		})
		shell, base, st, picked := pixel.Pad(c, k.Msg("symbols"), 8, len(symbols), 11, 11, cur, o.Disabled)
		if picked >= 0 && len(*password) < n {
			*password += symbols[picked : picked+1]
		}
		if !base.IsDisabled() && pixel.Shortcut(base, 0, ui.KeyBackspace) && *password != "" {
			*password = (*password)[:len(*password)-1]
		}
		base.Description(symbols[*cur : *cur+1])
		at := *cur
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range len(symbols) {
				x, y, ink := i%8*11, i/8*11, pixel.Ink
				if i == at {
					g.Rect(x, y, 11, 11, pixel.Accent)
					ink = pixel.AccentInk
				}
				g.Micro(x+4, y+3, symbols[i:i+1], ink, 1)
			}
			g.Finish(st)
		})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if input.Button(c, k.Msg("back"), input.ButtonOptions{Size: input.Small, Disabled: o.Disabled || *password == ""}).Clicked {
				*password = (*password)[:len(*password)-1]
			}
			res.Submitted = input.Button(c, k.Msg("end"), input.ButtonOptions{Size: input.Small, Variant: input.Primary, Disabled: o.Disabled || len(*password) < n}).Clicked
		})
	})
	res.Changed = *password != old
	// Lengths only: the password stays off the log.
	if res.Changed || res.Submitted {
		slog.Debug("8bit: password grid", "label", text, "length", len(*password), "submitted", res.Submitted)
	}
	return res
}
