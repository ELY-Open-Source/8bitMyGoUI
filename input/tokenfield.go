package input

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type TokenFieldOptions struct {
	Placeholder string
	// Width is in dots; zero means 200.
	Width    int
	Disabled bool
}

type TokenFieldResult struct {
	Element ui.Element
	Changed bool
}

type tokenKey struct{}

// editorSlot keys the editor so added tags keep its focus.
type editorSlot struct{}

// TokenField shows tokens as tags; Enter or comma adds.
func TokenField(c *ui.Context, text string, tokens *[]string, o TokenFieldOptions) TokenFieldResult {
	k := pixel.K(c)
	e := ui.Row(c).Wrap().Width(k.D(pixel.Or("TokenField Width", o.Width, 200))).Padding(k.D(3), k.D(4)).Gap(k.D(3)).AlignItems(ui.Center)
	e.Role(ui.RoleGroup).Label(text).Description(strings.Join(*tokens, ", "))
	draft := ui.Local(e, tokenKey{}, func() string { return "" })
	// The row stays enabled, so off is checked here.
	off := o.Disabled || pixel.Current(c).Inert
	res := TokenFieldResult{Element: e}
	add := func(s string) {
		if s = strings.TrimSpace(s); s != "" && !slices.Contains(*tokens, s) {
			*tokens, res.Changed = append(*tokens, s), true
		}
	}
	// A comma typed since the last frame cuts the draft.
	if parts := strings.Split(*draft, ","); len(parts) > 1 && !off {
		for _, p := range parts[:len(parts)-1] {
			add(p)
		}
		*draft = parts[len(parts)-1]
	}
	if !off && pixel.Shortcut(e, 0, ui.KeyEnter) {
		add(*draft)
		*draft = ""
	}
	if !off && *draft == "" && len(*tokens) > 0 && pixel.Shortcut(e, 0, ui.KeyBackspace) {
		*tokens, res.Changed = (*tokens)[:len(*tokens)-1], true
	}
	var ed pixel.Edit
	e.Children(func() {
		for i := 0; i < len(*tokens); i++ {
			tok, removed := (*tokens)[i], false
			ui.Box(c).Key(tok).Children(func() {
				removed = display.Tag(c, tok, display.TagOptions{Removable: true, Disabled: o.Disabled}).Removed
			})
			if removed {
				*tokens, res.Changed = slices.Delete(*tokens, i, i+1), true
				i--
			}
		}
		ui.Box(c).Key(editorSlot{}).Grow(1).MinWidth(k.D(40)).Children(func() {
			ed = pixel.Editor(c, draft, pixel.EditorOptions{Label: text, Placeholder: o.Placeholder, Disabled: o.Disabled})
		})
	})
	focus := ed.Focused
	e.Draw(func(p *ui.Painter, r ui.Rect) { well(k.Canvas(p, r), pixel.Edge, focus) })
	if res.Changed {
		slog.Debug("8bit: tokens", "label", text, "tokens", len(*tokens))
	}
	return res
}
