package gallery

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chat"
	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var mailEntries = []Entry{
	{334, "MemberList", "成员列表", "chat", "Who is around, who is away, who logged off for dinner.", demoMemberList},
	{335, "UserCard", "用户卡", "chat", "A face, a name, a title and things to do with them.", demoUserCard},
	{336, "Mailbox", "信箱", "chat", "Closed envelopes are new. Enter opens one.", demoMailbox},
	{337, "MailReader", "读信", "chat", "Who, what, when; then the letter itself.", demoMailReader},
	{338, "Letter", "信笺", "chat", "Parchment with worn, dithered edges and a signature.", demoLetter},
	{339, "NotificationBell", "通知铃", "chat", "Rings in steps when news arrives. Click to read it all.", demoNotificationBell},
	{340, "CallControls", "通话控制", "chat", "Mute, camera, and the big red button.", demoCallControls},
	{341, "VoiceMeter", "语音电平", "chat", "Green when you talk, red when you shout.", demoVoiceMeter},
	{342, "AttachmentChip", "附件", "chat", "A file's icon, name and size; × takes it off.", demoAttachmentChip},
	{343, "LinkPreview", "链接预览", "chat", "Title, site and a thumbnail to come. A click opens it.", demoLinkPreview},
	{344, "Poll", "投票", "chat", "Click to vote; click another to change your mind.", demoPoll},
	{345, "StickerPicker", "贴纸选择", "chat", "Big faces for big feelings.", demoStickerPicker},
}

// mailbag is the demo inbox.
func mailbag(now time.Time) []chat.Mail {
	return []chat.Mail{
		{From: "The Guild", Subject: "Your rank went up!", Body: "Congratulations, adventurer. Collect your new badge at the front desk.", At: now.Add(-time.Hour)},
		{From: "Aria", Subject: "Re: boss fight", Body: "Same time tomorrow? Bring more potions this time.", At: now.Add(-3 * time.Hour)},
		{From: "Shopkeeper", Subject: "Sale on shields", Body: "Two for one, today only.", At: now.Add(-26 * time.Hour), Read: true},
	}
}

func demoMemberList(c *ui.Context, log func(string)) {
	chat.MemberList(c, "Members", []chat.Member{{Name: "Aria", Status: display.Online}, {Name: "Brom", Status: display.Away},
		{Name: "Mittens", Status: display.Busy}, {Name: "Arthur", Status: display.Offline}, {Name: "Old Sage", Status: display.Offline}}, chat.MemberListOptions{})
}

func demoUserCard(c *ui.Context, log func(string)) {
	actions := []string{"MESSAGE", "INVITE"}
	if r := chat.UserCard(c, "Aria", chat.UserCardOptions{Title: "Level 12 Mage", Actions: actions}); r.Clicked >= 0 {
		log(actions[r.Clicked])
	}
}

func demoMailbox(c *ui.Context, log func(string)) {
	s := chatState(c)
	r := chat.Mailbox(c, "Inbox", s.mails, &s.mail, chat.MailboxOptions{})
	if r.Changed {
		log("mail " + s.mails[s.mail].Subject)
	}
	if r.Opened >= 0 {
		log("opened " + s.mails[r.Opened].Subject)
	}
}

func demoMailReader(c *ui.Context, log func(string)) {
	chat.MailReader(c, mailbag(pixel.K(c).Now())[0], chat.MailReaderOptions{})
}

func demoLetter(c *ui.Context, log func(string)) {
	chat.Letter(c, "Dear hero, the cave is colder than it looks. Wear the red scarf. Come home for supper.", chat.LetterOptions{Signature: "Grandpa"})
}

func demoNotificationBell(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chatState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		if chat.NotificationBell(c, "Notifications", s.news, chat.NotificationBellOptions{}).Clicked {
			s.news = 0
			log("all read")
		}
		if small(c, "NEWS!", false) {
			s.news++
			log(fmt.Sprintf("%d new", s.news))
		}
	})
}

func demoCallControls(c *ui.Context, log func(string)) {
	s := chatState(c)
	r := chat.CallControls(c, "Call", &s.muted, &s.camera, chat.CallControlsOptions{})
	if r.Changed {
		log(fmt.Sprintf("muted %v, camera %v", s.muted, s.camera))
	}
	if r.HungUp {
		log("hung up")
	}
}

func demoVoiceMeter(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chatState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		chat.VoiceMeter(c, "Voice", s.level, chat.VoiceMeterOptions{})
		if input.Slider(c, "Level", &s.level, input.SliderOptions{}).Changed {
			log(fmt.Sprintf("level %.2f", s.level))
		}
	})
}

func demoAttachmentChip(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chatState(c)
	sizes := map[string]int64{"world-map.png": 2_516_582, "battle-theme.mp3": 4_404_019, "notes.txt": 812}
	ui.Column(c).Gap(k.D(3)).AlignItems(ui.Start).Children(func() {
		// Keys keep each chip's state as earlier ones go.
		for _, f := range slices.Clone(s.files) {
			ui.Box(c).Key(f).Children(func() {
				if chat.AttachmentChip(c, f, sizes[f], chat.AttachmentChipOptions{}).Removed {
					s.files = slices.DeleteFunc(s.files, func(x string) bool { return x == f })
					log("removed " + f)
				}
			})
		}
		if len(s.files) == 0 && small(c, "RESTORE", false) {
			s.files = []string{"world-map.png", "battle-theme.mp3", "notes.txt"}
			log("restored")
		}
	})
}

func demoLinkPreview(c *ui.Context, log func(string)) {
	if chat.LinkPreview(c, "https://mygo.egoist.dev/docs/ui", "MyGo UI: native widgets, written in Go", "mygo.egoist.dev", chat.LinkPreviewOptions{}).Clicked {
		log("opened mygo.egoist.dev")
	}
}

func demoPoll(c *ui.Context, log func(string)) {
	s := chatState(c)
	if r := chat.Poll(c, "Best class for the dungeon?", s.polls, &s.vote, chat.PollOptions{}); r.Voted >= 0 {
		log("voted " + s.polls[r.Voted].Text)
	}
}

func demoStickerPicker(c *ui.Context, log func(string)) {
	if r := chat.StickerPicker(c, "Stickers", chat.StickerPickerOptions{}); r.Picked != "" {
		log("sticker " + r.Picked)
	}
}
