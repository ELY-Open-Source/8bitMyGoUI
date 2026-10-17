package desk

import (
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type RunDialogOptions struct {
	// Width is in dots; zero means 240.
	Width    int
	Disabled bool
}

type RunDialogResult struct {
	Element ui.Element
	// Canceled is CANCEL, the close box or Escape.
	Changed, Submitted, Canceled bool
}

// RunDialog asks for a program: OK runs, CANCEL quits.
func RunDialog(c *ui.Context, text string, value *string, o RunDialogOptions) RunDialogResult {
	k := pixel.K(c)
	w := pixel.Or("RunDialog Width", o.Width, 240)
	open := k.Msg("open_label")
	win := ui.Column(c).Width(k.D(w)).Padding(k.D(3)).Gap(k.D(4)).Role(ui.RoleGroup).Label(text)
	if o.Disabled {
		win.Description(k.Msg("disabled"))
	}
	active := win.FocusWithin()
	var res RunDialogResult
	win.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.Face)
	}).Children(func() {
		ui.Row(c).FillWidth().Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Grow(1).Height(k.D(pixel.LineHeight + 2)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
				barLook(k.Canvas(p, r), k.Msg("run"), active)
			})
			res.Canceled = input.IconButton(c, icon.Close, k.Msg("close"), input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked
		})
		ui.Row(c).Gap(k.D(6)).Padding(0, k.D(4)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(k.D(24), k.D(24)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Sprite(0, 0, icon.Terminal, 0, pixel.Ink, 2)
			})
			ui.Box(c).Width(k.D(w - 44)).Children(func() {
				pixel.Text(c, k.Msg("run_prompt"), pixel.TextOptions{Wrap: true, MaxLines: 2})
			})
		})
		var f input.TextInputResult
		ui.Row(c).Gap(k.D(4)).Padding(0, k.D(4)).AlignItems(ui.Center).Children(func() {
			// The field carries the name; this is for eyes.
			ui.Box(c).Size(k.D(pixel.Measure(open)), k.D(pixel.LineHeight)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Text(0, 0, open, pixel.Ink, 1)
			})
			f = input.TextInput(c, open, value, input.TextInputOptions{Width: w - 16 - pixel.Measure(open), Disabled: o.Disabled})
			res.Changed = f.Changed
		})
		blank := strings.TrimSpace(*value) == ""
		ui.Row(c).Gap(k.D(6)).Padding(0, k.D(4), k.D(4), 0).AlignSelf(ui.End).Children(func() {
			ok := input.Button(c, k.Msg("ok"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled || blank}).Clicked
			res.Submitted = !blank && !f.Input.IsDisabled() && (ok || f.Submitted)
			res.Canceled = input.Button(c, k.Msg("cancel"), input.ButtonOptions{Disabled: o.Disabled}).Clicked || res.Canceled
		})
	})
	// Escape quits from anywhere inside, as CANCEL does.
	if !o.Disabled && !pixel.Current(c).Inert && pixel.Shortcut(win, 0, ui.KeyEscape) {
		res.Canceled = true
	}
	res.Element = win
	if res.Changed {
		slog.Debug("8bit: run dialog typed", "label", text, "runes", utf8.RuneCountInString(*value))
	}
	if res.Submitted || res.Canceled {
		slog.Debug("8bit: run dialog", "label", text, "runes", utf8.RuneCountInString(*value), "submitted", res.Submitted)
	}
	return res
}
