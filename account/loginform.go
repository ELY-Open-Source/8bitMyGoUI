package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type LoginFormOptions struct {
	// Failures counts wrong tries; each new one shakes the form.
	Failures int
	Disabled bool
}

type LoginFormResult struct {
	Element   *ui.Element
	Submitted bool
}

type loginKey struct{}

// login is the failures seen and whether the warning shows.
type login struct {
	seen   int
	failed bool
}

// shakes are the steps of a head shake, in dots.
var shakes = []int{-4, 4, -3, 3, -2, 2, -1, 1}

// LoginForm asks a name and password; the caller checks them.
func LoginForm(c *ui.Context, title string, name, password *string, o LoginFormOptions) LoginFormResult {
	k := pixel.K(c)
	if o.Failures < 0 {
		panic(fmt.Sprintf("account: LoginForm %q has %d failures", title, o.Failures))
	}
	box := ui.Box(c)
	lg := ui.Local(box, loginKey{}, func() login { return login{seen: o.Failures} })
	fresh := o.Failures > lg.seen
	if fresh {
		// A wrong try clears the password, naming no field.
		lg.seen, lg.failed, *password = o.Failures, true, ""
	}
	dx := 0
	if n := k.Shot(c, box, o.Failures, 12, len(shakes)); n < len(shakes) {
		dx = shakes[n]
	}
	var res LoginFormResult
	box.Padding(0, k.D(4-dx), 0, k.D(4+dx)).Children(func() {
		layout.Panel(c, layout.PanelOptions{Style: pixel.Raised, Padding: 6}, func() {
			text.ShadowText(c, title, text.ShadowTextOptions{})
			nm := input.TextInput(c, k.Msg("name"), name, input.TextInputOptions{Placeholder: k.Msg("name"), Width: 140, Disabled: o.Disabled})
			pw := input.PasswordInput(c, k.Msg("password"), password, input.PasswordInputOptions{Placeholder: k.Msg("password"), Width: 140, Disabled: o.Disabled})
			if fresh && !pw.Input.IsDisabled() {
				pw.Input.Focus()
			}
			if nm.Changed || pw.Changed {
				lg.failed = false
			}
			if lg.failed {
				pixel.Text(c, k.Msg("login_failed"), pixel.TextOptions{Slot: pixel.Bad})
			}
			ready := *name != "" && *password != ""
			clicked := input.Button(c, k.Msg("login"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled || !ready}).Clicked
			res.Submitted = ready && (clicked || pw.Submitted && !pw.Input.IsDisabled())
		}).Gap(k.D(5)).Role(ui.RoleGroup).Label(title)
	})
	res.Element = box
	// Never the name or password: only that someone tried.
	if res.Submitted {
		slog.Debug("8bit: login form submitted", "label", title)
	}
	return res
}
