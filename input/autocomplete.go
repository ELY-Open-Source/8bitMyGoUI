package input

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AutocompleteOptions struct {
	Placeholder string
	// Width is in dots; zero means 120.
	Width    int
	Disabled bool
}

type AutocompleteResult struct {
	Element *ui.Element
	// Accepted reports a completion taken by Tab or Enter.
	Changed, Accepted, Submitted bool
}

// Autocomplete ghosts the first completion; Tab or Enter takes it.
func Autocomplete(c *ui.Context, text string, value *string, suggestions []string, o AutocompleteOptions) AutocompleteResult {
	k := pixel.K(c)
	ed, e := field(c, value, pixel.EditorOptions{Label: text, Placeholder: o.Placeholder, Disabled: o.Disabled}, o.Width, nil, nil, nil)
	res := AutocompleteResult{Element: e, Changed: ed.Changed, Submitted: ed.Submitted}
	best := completion(suggestions, *value)
	tab, enter := best != "" && e.Shortcut(0, ui.KeyTab), best != "" && e.Shortcut(0, ui.KeyEnter)
	if tab || enter {
		*value, res.Changed, res.Accepted = best, true, true
		slog.Debug("8bit: autocomplete accepted", "label", text, "value", best)
		return res
	}
	if best != "" {
		ed.Input.Description(best)
		typed := *value
		ed.Element.DrawOver(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			rest := string([]rune(best)[len([]rune(typed)):])
			if x := pixel.Measure(typed); x+pixel.Measure(rest) <= g.W {
				g.Text(x, 0, rest, pixel.Faint, 1)
			}
		})
	}
	return res
}

// completion is the first suggestion extending v, ignoring case.
func completion(suggestions []string, v string) string {
	if v == "" {
		return ""
	}
	n := len([]rune(v))
	for _, s := range suggestions {
		if r := []rune(s); len(r) > n && strings.EqualFold(string(r[:n]), v) {
			return s
		}
	}
	return ""
}
