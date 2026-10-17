package files

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Access is what a link lets its holder do.
type Access int

const (
	CanView Access = iota
	CanComment
	CanEdit
)

type SharePanelOptions struct {
	// Width is in dots; zero means 220.
	Width    int
	Disabled bool
}

type SharePanelResult struct {
	Element         ui.Element
	Copied, Changed bool
}

type shareKey struct{}

// SharePanel shows a link to copy and its access.
func SharePanel(c *ui.Context, text, link string, access *Access, o SharePanelOptions) SharePanelResult {
	k := pixel.K(c)
	if link == "" || *access < CanView || *access > CanEdit {
		panic(fmt.Sprintf("files: SharePanel link of %d bytes, access %d", len(link), *access))
	}
	var res SharePanelResult
	res.Element = ui.Column(c).Width(k.D(pixel.Or("SharePanel Width", o.Width, 220))).Gap(k.D(4)).Role(ui.RoleGroup).Label(text).Children(func() {
		ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			// The link is cut to fit; readers hear it whole.
			field := ui.Box(c).Grow(1).Height(k.D(pixel.LineHeight + 4)).Role(ui.RoleText).Label(k.Msg("link")).Description(link)
			lw, fits := pixel.WrapWidth(c, field)
			shown := pixel.Wrap(link, max(0, lw-6), 1)[0]
			field.Draw(func(p *ui.Painter, r ui.Rect) {
				fits(r)
				g := k.Canvas(p, r)
				g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
				g.Clip(2, 0, g.W-4, g.H, func() { g.Text(3, 2, shown, pixel.Ink, 1) })
			})
			e := ui.Box(c)
			at := ui.Local(e, shareKey{}, func() time.Time { return time.Time{} })
			// The button says COPIED! for two seconds.
			label := k.Msg("copy")
			if k.Since(c, *at, 4, 8) < 8 {
				label = k.Msg("copied")
			}
			e.Children(func() {
				if input.Button(c, label, input.ButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
					c.WriteClipboard(link)
					*at, res.Copied = k.Now(), true
				}
			})
		})
		i := int(*access)
		if input.Segmented(c, k.Msg("access"), strings.Fields(k.Msg("access_levels")), &i, input.SegmentedOptions{Disabled: o.Disabled}).Changed {
			*access, res.Changed = Access(i), true
		}
	})
	// A share link is a key; it never reaches logs.
	if res.Copied || res.Changed {
		slog.Debug("8bit: share panel", "label", text, "copied", res.Copied, "access", *access)
	}
	return res
}
