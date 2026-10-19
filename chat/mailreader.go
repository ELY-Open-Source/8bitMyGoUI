package chat

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type MailReaderOptions struct {
	// Width is in dots; zero means 220.
	Width int
}

// MailReader shows a mail's headers, a rule, its body.
func MailReader(c *ui.Context, m Mail, o MailReaderOptions) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("MailReader Width", o.Width, 220)
	return ui.Column(c).Width(k.D(w)).Gap(k.D(4)).Children(func() {
		head := [][2]string{{k.Msg("from"), m.From}, {k.Msg("subject"), m.Subject}, {k.Msg("date"), m.At.Format("2006-01-02 15:04")}}
		data.DescriptionList(c, head, data.DescriptionListOptions{TermWidth: 50})
		ui.Box(c).FillWidth().Height(k.D(3)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.HLine(0, 1, g.W, pixel.Edge)
		})
		para(c, k, m.Body, w, 0, pixel.Ink)
	})
}
