package chat

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Mail is a letter; Read opens its envelope.
type Mail struct {
	From, Subject, Body string
	At                  time.Time
	Read                bool
}

type MailboxOptions struct {
	// Width and Height are dots; zeros mean 220 by 120.
	Width, Height int
	Disabled      bool
}

type MailboxResult struct {
	Element *ui.Element
	Changed bool
	// Opened is the mail opened this frame, else -1.
	Opened int
}

// Mailbox lists envelopes; Enter or a double click opens.
func Mailbox(c *ui.Context, text string, mails []Mail, selected *int, o MailboxOptions) MailboxResult {
	k := pixel.K(c)
	w := pixel.Or("Mailbox Width", o.Width, 220)
	room := w - 14 - 16 - 30
	name := func(i int) string {
		m := mails[i]
		if !m.Read {
			return m.From + ": " + m.Subject + ", " + k.Msg("new_messages")
		}
		return m.From + ": " + m.Subject
	}
	r := list(c, k, text, len(mails), selected, w, pixel.Or("Mailbox Height", o.Height, 120), o.Disabled, name, func(i int) {
		m := mails[i]
		env, ink, from := icon.MailOpen, pixel.Faint, pixel.Faint
		if !m.Read {
			env, ink, from = icon.Mail, pixel.Accent, pixel.Ink
		}
		ui.Box(c).FillWidth().Height(k.D(2 * pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Sprite(0, 4, env, 0, ink, 1)
			g.Text(16, 0, line(m.From, room), from, 1)
			g.Text(16, pixel.LineHeight, line(m.Subject, room+28), pixel.Faint, 1)
			t := clock(m.At)
			g.Micro(g.W-pixel.MicroWidth(t)-2, 3, t, pixel.Faint, 1)
		})
	})
	res := MailboxResult{Element: r.Element, Changed: r.Changed, Opened: -1}
	if r.Submitted && *selected >= 0 {
		mails[*selected].Read, res.Opened = true, *selected
		slog.Debug("8bit: mailbox opened", "label", text, "mail", *selected)
	}
	return res
}
