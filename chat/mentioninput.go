package chat

import (
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MentionInputOptions struct {
	Placeholder string
	// Width is in dots; zero means 160.
	Width    int
	Disabled bool
}

type MentionInputResult struct {
	Element                      *ui.Element
	Changed, Submitted, Accepted bool
	// Mentioned is the member just completed, else empty.
	Mentioned string
}

// MentionInput suggests members after @; Tab or Enter completes.
func MentionInput(c *ui.Context, text string, value *string, members []string, o MentionInputOptions) MentionInputResult {
	k := pixel.K(c)
	w := pixel.Or("MentionInput Width", o.Width, 160)
	// The word after the last @ picks matching members.
	var whole, names []string
	if at := strings.LastIndex(*value, "@"); at >= 0 && !strings.ContainsAny((*value)[at:], " \n") {
		typed := strings.ToLower((*value)[at+1:])
		for _, m := range members {
			if strings.HasPrefix(strings.ToLower(m), typed) {
				whole, names = append(whole, (*value)[:at+1]+m+" "), append(names, m)
			}
		}
	}
	var res MentionInputResult
	res.Element = ui.Column(c).Gap(k.D(1)).Children(func() {
		r := input.Autocomplete(c, text, value, whole, input.AutocompleteOptions{Placeholder: o.Placeholder, Width: w, Disabled: o.Disabled})
		res.Changed, res.Submitted, res.Accepted = r.Changed, r.Submitted, r.Accepted
		if r.Changed {
			slog.Debug("8bit: mention input changed", "label", text, "runes", utf8.RuneCountInString(*value))
		}
		for i, s := range whole {
			if r.Accepted && s == *value {
				res.Mentioned = names[i]
				slog.Debug("8bit: mention", "label", text, "member", names[i])
			}
		}
		if len(names) > 0 && !r.Accepted {
			pixel.Text(c, line("@"+strings.Join(names, "  @"), w), pixel.TextOptions{Slot: pixel.Faint})
		}
	})
	return res
}
