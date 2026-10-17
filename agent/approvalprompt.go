package agent

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ApprovalPromptOptions struct {
	// Danger inks the command and YES in Bad.
	Danger bool
	// Width is in dots; zero means 240.
	Width    int
	Disabled bool
}

type ApprovalPromptResult struct {
	Element ui.Element
	// Allowed or Denied is the answer given this frame.
	Allowed, Denied bool
}

// ApprovalPrompt asks ALLOW SPELL? over a command: YES or NO.
func ApprovalPrompt(c *ui.Context, text, command string, o ApprovalPromptOptions) ApprovalPromptResult {
	k := pixel.K(c)
	if command == "" {
		panic("agent: ApprovalPrompt has no command")
	}
	w := pixel.Or("ApprovalPrompt Width", o.Width, 240)
	ask, warn := k.Msg("allow_spell"), k.Msg("danger")
	need := 28 + pixel.Measure(ask)
	if o.Danger {
		need += 20 + pixel.Measure(warn)
	}
	if w < need {
		panic(fmt.Sprintf("agent: ApprovalPrompt Width %d is under %d", w, need))
	}
	// The command shows whole to readers, three lines to eyes.
	lines, ink, yes := pixel.Wrap(command, w-20, 3), pixel.Ink, input.Primary
	if o.Danger {
		ink, yes = pixel.Bad, input.Danger
	}
	var res ApprovalPromptResult
	panel := ui.Column(c).Width(k.D(w)).Padding(k.D(6)).Gap(k.D(6)).Role(ui.RoleGroup).Label(text)
	if o.Disabled {
		panel.Description(k.Msg("disabled"))
	}
	res.Element = panel.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Double, pixel.Face)
	}).Children(func() {
		ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			sprite(c, k, icon.Robot, pixel.Ink)
			pixel.Text(c, ask, pixel.TextOptions{})
			if o.Danger {
				sprite(c, k, icon.Warning, pixel.Bad)
				pixel.Text(c, warn, pixel.TextOptions{Slot: pixel.Bad})
			}
		})
		ui.Box(c).Size(k.D(w-12), k.D(len(lines)*pixel.LineHeight+6)).Role(ui.RoleText).Label(command).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
			for i, l := range lines {
				g.Text(4, 3+i*pixel.LineHeight, l, ink, 1)
			}
		})
		ui.Row(c).Gap(k.D(6)).Children(func() {
			res.Allowed = input.Button(c, k.Msg("yes"), input.ButtonOptions{Variant: yes, Disabled: o.Disabled}).Clicked
			res.Denied = input.Button(c, k.Msg("no"), input.ButtonOptions{Disabled: o.Disabled}).Clicked
		})
	})
	if res.Allowed || res.Denied {
		slog.Debug("8bit: approval answered", "label", text, "allowed", res.Allowed)
	}
	return res
}

// sprite is a 12-dot icon, unnamed, in ink.
func sprite(c *ui.Context, k pixel.Kit, sp *pixel.Sprite, ink pixel.Slot) ui.Element {
	return ui.Box(c).Size(k.D(sp.W), k.D(sp.H)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Sprite(0, 0, sp, 0, ink, 1)
	})
}
