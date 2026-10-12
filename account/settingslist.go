package account

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Setting is a label beside the control it builds.
type Setting struct {
	Label   string
	Control func()
}

// SettingsSection is a titled run of settings.
type SettingsSection struct {
	Title    string
	Settings []Setting
}

type SettingsListOptions struct {
	// Width is in dots; zero means 240.
	Width int
}

// SettingsList sets controls beside labels under section titles.
func SettingsList(c *ui.Context, text string, sections []SettingsSection, o SettingsListOptions) *ui.Element {
	k := pixel.K(c)
	return ui.Column(c).Width(k.D(pixel.Or("SettingsList Width", o.Width, 240))).Gap(k.D(2)).Role(ui.RoleGroup).Label(text).Children(func() {
		for i, s := range sections {
			top := 4
			if i == 0 {
				top = 0
			}
			ui.Box(c).FillWidth().Height(k.D(pixel.LineHeight+3)).Margin(k.D(top), 0, 0, 0).Role(ui.RoleText).Label(s.Title).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Text(0, 0, s.Title, pixel.Accent, 1)
				g.HLine(0, pixel.LineHeight+1, g.W, pixel.Accent)
			})
			for j, st := range s.Settings {
				last := j == len(s.Settings)-1
				ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(k.D(6)).Padding(k.D(2), 0).Children(func() {
					// The control says its own name; the label only shows.
					mark(c, k, st.Label, pixel.Ink, 0).Grow(1)
					st.Control()
				}).DrawOver(func(p *ui.Painter, r ui.Rect) {
					// A dotted rule between rows, none under the last.
					if !last {
						g := k.Canvas(p, r)
						for x := 0; x < g.W; x += 2 {
							g.Dot(x, g.H-1, pixel.Shade)
						}
					}
				})
			}
		}
	})
}
