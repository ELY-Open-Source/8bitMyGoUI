package feedback

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type ContinuePromptOptions struct {
	// Seconds counts down; zero means 10.
	Seconds  int
	Disabled bool
}

type ContinuePromptResult struct {
	Element           *ui.Element
	Yes, No, TimedOut bool
}

type promptKey struct{}

// prompt is when the count began, and if it ended.
type prompt struct {
	start time.Time
	over  bool
}

// ContinuePrompt asks CONTINUE? YES or NO, counting down.
func ContinuePrompt(c *ui.Context, o ContinuePromptOptions) ContinuePromptResult {
	k := pixel.K(c)
	n := pixel.Or("ContinuePrompt Seconds", o.Seconds, 10)
	col := ui.Column(c).Gap(k.D(6)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(k.Msg("continue"))
	st := ui.Local(col, promptKey{}, func() prompt { return prompt{start: k.Now()} })
	left := max(0, n-max(0, int(k.Now().Sub(st.start)/time.Second)))
	res := ContinuePromptResult{Element: col}
	if left == 0 && !st.over {
		st.over, res.TimedOut = true, true
	}
	if !st.over {
		k.Tick(c, 10)
	}
	col.Children(func() {
		text.ShadowText(c, k.Msg("continue"), text.ShadowTextOptions{Slot: pixel.Accent, Scale: 2})
		text.ShadowText(c, strconv.Itoa(left), text.ShadowTextOptions{Slot: pixel.Warn, Scale: 4})
		r := input.ButtonGroup(c, []string{k.Msg("yes"), k.Msg("no")}, input.ButtonGroupOptions{Label: k.Msg("continue"), Disabled: st.over || o.Disabled})
		res.Yes, res.No = r.Clicked == 0, r.Clicked == 1
	})
	if res.Yes || res.No {
		st.over = true
	}
	if res.Yes || res.No || res.TimedOut {
		slog.Debug("8bit: continue prompt", "yes", res.Yes, "no", res.No, "timed out", res.TimedOut)
	}
	return res
}
