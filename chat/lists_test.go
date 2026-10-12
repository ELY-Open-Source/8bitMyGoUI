package chat_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/chat"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestChatListBadgesAndChooses(t *testing.T) {
	wd := &world{clk: newClock()}
	convs := []chat.Conversation{{Name: "Aria", Last: "Ready?", Unread: 2}, {Name: "Brom", Last: "Snacks."}, {Name: "Guild", Last: "Raid at eight", Unread: 120}}
	sel, changed, submitted := 0, 0, 0
	tt := view(500, 300, wd, func(c *ui.Context) {
		r := chat.ChatList(c, "Chats", convs, &sel, chat.ChatListOptions{})
		if r.Changed {
			changed++
		}
		if r.Submitted {
			submitted++
		}
	})
	find(t, tt, "2")
	find(t, tt, "99+")
	click(t, tt, "Brom")
	if sel != 1 || changed != 1 {
		t.Fatalf("clicking Brom chose %d, changed %d", sel, changed)
	}
	// The chosen row is Face, readable in any palette.
	if r := find(t, tt, "Brom"); count(tt.Image(), int(r.X/2)+30, int(r.Y/2), 40, 8, pixel.Face) < 200 || count(tt.Image(), int(r.X/2)+30, int(r.Y/2), 40, 8, pixel.AccentLit) != 0 {
		t.Fatal("the chosen row is not Face")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if sel != 2 || submitted != 1 {
		t.Fatalf("Down and Enter: row %d, submitted %d", sel, submitted)
	}
	// Rows are named; typing a name chooses its row.
	find(t, tt, "Aria: Ready?, 2 UNREAD")
	tt.Key(0, ui.KeyB)
	if sel != 1 {
		t.Fatalf("typing b chose row %d", sel)
	}
	mustPanic(t, "row", func() {
		bad := 5
		view(300, 200, wd, func(c *ui.Context) { chat.ChatList(c, "Chats", convs, &bad, chat.ChatListOptions{}) })
	})
}

func TestChannelListMarksUnread(t *testing.T) {
	wd := &world{clk: newClock()}
	chans := []chat.Channel{{Name: "general", Unread: true}, {Name: "quests"}}
	// None is chosen, so no ▶ adds its Accent.
	sel := -1
	tt := view(400, 200, wd, func(c *ui.Context) { chat.ChannelList(c, "Channels", chans, &sel, chat.ChannelListOptions{}) })
	unread, read := find(t, tt, "general, NEW MESSAGES"), find(t, tt, "quests")
	// Unread channels get an Accent dot; read ones are Faint.
	img := tt.Image()
	if count(img, int(unread.X/2), int(unread.Y/2), int(unread.W/2), 16, pixel.Accent) == 0 || count(img, int(read.X/2), int(read.Y/2), int(read.W/2), 16, pixel.Accent) != 0 {
		t.Fatal("only unread channels carry a dot")
	}
	if count(img, int(read.X/2), int(read.Y/2), 40, 16, pixel.Faint) == 0 {
		t.Fatal("a read channel is Faint")
	}
}

func TestMemberListGroupsByPresence(t *testing.T) {
	wd := &world{clk: newClock()}
	tt := view(400, 300, wd, func(c *ui.Context) {
		chat.MemberList(c, "Members", []chat.Member{{Name: "Ari", Status: display.Offline}, {Name: "Bo", Status: display.Busy}, {Name: "Cy", Status: display.Online}}, chat.MemberListOptions{})
	})
	on, off := find(t, tt, "ONLINE · 2"), find(t, tt, "OFFLINE · 1")
	if bo, ari := find(t, tt, "Bo"), find(t, tt, "Ari"); !(on.Y < bo.Y && bo.Y < off.Y && off.Y < ari.Y) {
		t.Fatal("online members come first, each under its heading")
	}
	find(t, tt, "BUSY")
}

func TestUserCardActions(t *testing.T) {
	wd := &world{clk: newClock()}
	clicked := -1
	tt := view(400, 200, wd, func(c *ui.Context) {
		if r := chat.UserCard(c, "Aria", chat.UserCardOptions{Title: "Mage", Actions: []string{"MESSAGE", "INVITE"}}); r.Clicked >= 0 {
			clicked = r.Clicked
		}
	})
	find(t, tt, "Mage")
	click(t, tt, "INVITE")
	if clicked != 1 {
		t.Fatalf("INVITE reported %d", clicked)
	}
}

func TestMailboxOpensMail(t *testing.T) {
	wd := &world{clk: newClock()}
	mails := []chat.Mail{{From: "Guild", Subject: "Rank up", At: wd.clk.t}, {From: "Shop", Subject: "Sale", At: wd.clk.t, Read: true}}
	sel, opened := -1, -1
	tt := view(500, 300, wd, func(c *ui.Context) {
		if r := chat.Mailbox(c, "Inbox", mails, &sel, chat.MailboxOptions{}); r.Opened >= 0 {
			opened = r.Opened
		}
	})
	click(t, tt, "Guild: Rank up, NEW MESSAGES")
	tt.Key(0, ui.KeyEnter)
	if opened != 0 || !mails[0].Read {
		t.Fatalf("Enter opened %d, read %v", opened, mails[0].Read)
	}
}

func TestMailReaderAndLetter(t *testing.T) {
	wd := &world{clk: newClock()}
	m := chat.Mail{From: "Guild", Subject: "Rank up", Body: "Well done.", At: wd.clk.t}
	tt := view(500, 400, wd, func(c *ui.Context) {
		chat.MailReader(c, m, chat.MailReaderOptions{})
		chat.Letter(c, "Wear the scarf.", chat.LetterOptions{Signature: "Grandpa"})
	})
	for _, s := range []string{"FROM", "Guild", "SUBJECT", "Rank up", "DATE", "2026-09-02 10:30", "Well done."} {
		find(t, tt, s)
	}
	r := find(t, tt, "Wear the scarf. — Grandpa")
	// The parchment darkens to its rim.
	img := tt.Image()
	if count(img, int(r.X/2), int(r.Y/2)+4, 1, int(r.H/2)-8, pixel.Shade) == 0 || count(img, int(r.X/2)+10, int(r.Y/2)+30, 20, 1, pixel.Shade) != 0 {
		t.Fatal("only the letter's edge is dithered")
	}
	mustPanic(t, "narrow letter", func() {
		view(100, 100, wd, func(c *ui.Context) { chat.Letter(c, "x", chat.LetterOptions{Width: 30}) })
	})
}

func TestReadReceiptSentAndRead(t *testing.T) {
	wd := &world{clk: newClock()}
	tt := view(200, 100, wd, func(c *ui.Context) {
		chat.ReadReceipt(c, false, chat.ReadReceiptOptions{})
		chat.ReadReceipt(c, true, chat.ReadReceiptOptions{})
	})
	sent, read := find(t, tt, "SENT"), find(t, tt, "READ")
	img := tt.Image()
	if a, b := count(img, int(sent.X/2), int(sent.Y/2), 8, 5, pixel.Faint), count(img, int(read.X/2), int(read.Y/2), 8, 5, pixel.Info); b <= a || a == 0 {
		t.Fatalf("one check of %d dots, two of %d", a, b)
	}
}
