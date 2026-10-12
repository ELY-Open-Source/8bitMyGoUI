package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/chat"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var chatEntries = []Entry{
	{321, "SpeechBubble", "聊天气泡", "chat", "Their words on the left, yours on the right, tails pointing home.", demoSpeechBubble},
	{322, "ChatLog", "聊天记录", "chat", "Scrolls with the talk; at the end, it stays at the end.", demoChatLog},
	{323, "Composer", "输入栏", "chat", "Enter sends, Shift+Enter breaks the line. SEND ▶ works too.", demoComposer},
	{324, "TypingIndicator", "正在输入", "chat", "Three dots hop while someone types. Suspense!", demoTypingIndicator},
	{325, "ChatList", "会话列表", "chat", "Who said what last, and how much you missed.", demoChatList},
	{326, "ReadReceipt", "已读标记", "chat", "One check: sent. Two checks: seen. No excuses.", demoReadReceipt},
	{327, "Reactions", "表情回应", "chat", "Click a chip to add your face to the pile, again to take it back.", demoReactions},
	{328, "EmojiPicker", "表情选择", "chat", "Faces, the world and loot, eight to a row.", demoEmojiPicker},
	{329, "MentionInput", "提及输入", "chat", "Type @ and a few letters; Tab finishes the name.", demoMentionInput},
	{330, "DateSeparator", "日期分隔", "chat", "TODAY, YESTERDAY, then plain dates. Turns over at midnight.", demoDateSeparator},
	{331, "UnreadDivider", "未读分隔", "chat", "Everything under the red line is news.", demoUnreadDivider},
	{332, "SystemMessage", "系统消息", "chat", "Announcements from the game itself.", demoSystemMessage},
	{333, "ChannelList", "频道列表", "chat", "Bright channels with a dot have something new.", demoChannelList},
}

type chatDemo struct {
	draft, mention string
	msgs           []chat.Message
	conv, channel  int
	reactions      []chat.Reaction
	mails          []chat.Mail
	mail, vote     int
	news           int
	muted, camera  bool
	level          float64
	files          []string
	polls          []chat.PollOption
}

func chatState(c *ui.Context) *chatDemo {
	now := pixel.K(c).Now()
	return ui.Local(c.Root(), demoKey("chat"), func() chatDemo {
		return chatDemo{
			msgs: []chat.Message{
				{From: "Aria", Text: "Ready for the boss?", At: now.Add(-6 * time.Minute)},
				{From: "You", Text: "Potions packed. Lead the way!", At: now.Add(-5 * time.Minute), Mine: true},
				{From: "Brom", Text: "I brought snacks.", At: now.Add(-4 * time.Minute)},
			},
			reactions: []chat.Reaction{{Emoji: "love", Count: 3}, {Emoji: "laugh", Count: 1, Mine: true}, {Emoji: "fire", Count: 5}},
			mails:     mailbag(now),
			mail:      -1, vote: 1, camera: true, level: 0.6,
			files: []string{"world-map.png", "battle-theme.mp3", "notes.txt"},
			polls: []chat.PollOption{{Text: "Mage", Votes: 4}, {Text: "Knight", Votes: 6}, {Text: "Rogue", Votes: 3}, {Text: "Bard", Votes: 1}},
		}
	})
}

func demoSpeechBubble(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chatState(c)
	ui.Column(c).Width(k.D(220)).Gap(k.D(4)).Children(func() {
		chat.SpeechBubble(c, s.msgs[0], chat.SpeechBubbleOptions{})
		chat.SpeechBubble(c, s.msgs[1], chat.SpeechBubbleOptions{})
	})
}

func demoChatLog(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := chatState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		chat.ChatLog(c, "Party chat", chat.ChatLogOptions{Height: 120}, func() {
			chat.DateSeparator(c, k.Now(), chat.DateSeparatorOptions{Width: 190})
			chat.SystemMessage(c, "BROM HAS JOINED", chat.SystemMessageOptions{})
			for i, m := range s.msgs {
				if i == 2 {
					chat.UnreadDivider(c, chat.UnreadDividerOptions{Width: 190})
				}
				chat.SpeechBubble(c, m, chat.SpeechBubbleOptions{Width: 120})
			}
		})
		if small(c, "SAY HI", false) {
			s.msgs = append(s.msgs, chat.Message{From: "Mittens", Text: fmt.Sprintf("Mew #%d!", len(s.msgs)), At: pixel.K(c).Now()})
			log(fmt.Sprintf("%d messages", len(s.msgs)))
		}
	})
}

func demoComposer(c *ui.Context, log func(string)) {
	s := chatState(c)
	r := chat.Composer(c, "Message", &s.draft, chat.ComposerOptions{Placeholder: "Say something…"})
	if r.Changed {
		log(fmt.Sprintf("draft of %d letters", utf8.RuneCountInString(s.draft)))
	}
	// Counts, not words: logs never hold what was typed.
	if r.Sent != "" {
		log(fmt.Sprintf("sent %d letters", utf8.RuneCountInString(r.Sent)))
	}
}

func demoTypingIndicator(c *ui.Context, log func(string)) {
	chat.TypingIndicator(c, "Aria", chat.TypingIndicatorOptions{})
}

func demoChatList(c *ui.Context, log func(string)) {
	s := chatState(c)
	convs := []chat.Conversation{{Name: "Aria", Last: "Ready for the boss?", Unread: 2}, {Name: "Brom", Last: "I brought snacks."},
		{Name: "Guild Hall", Last: "Raid at eight. Bring torches.", Unread: 12}, {Name: "Mittens", Last: "Mew."}}
	r := chat.ChatList(c, "Chats", convs, &s.conv, chat.ChatListOptions{})
	if r.Changed {
		log("chat with " + convs[s.conv].Name)
	}
	if r.Submitted {
		log("open " + convs[s.conv].Name)
	}
}

func demoReadReceipt(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(10)).AlignItems(ui.Center).Children(func() {
		for _, r := range []struct {
			read bool
			text string
		}{{false, "sent"}, {true, "read"}} {
			ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
				chat.ReadReceipt(c, r.read, chat.ReadReceiptOptions{Scale: 2})
				pixel.Text(c, r.text, pixel.TextOptions{Slot: pixel.Faint})
			})
		}
	})
}

func demoReactions(c *ui.Context, log func(string)) {
	s := chatState(c)
	if r := chat.Reactions(c, "Reactions", s.reactions, chat.ReactionsOptions{}); r.Toggled >= 0 {
		re := s.reactions[r.Toggled]
		log(fmt.Sprintf("%s: %d, yours %v", re.Emoji, re.Count, re.Mine))
	}
}

func demoEmojiPicker(c *ui.Context, log func(string)) {
	if r := chat.EmojiPicker(c, "Emoji", chat.EmojiPickerOptions{}); r.Picked != "" {
		log("picked " + r.Picked)
	}
}

func demoMentionInput(c *ui.Context, log func(string)) {
	s := chatState(c)
	r := chat.MentionInput(c, "Mention", &s.mention, []string{"Aria", "Arthur", "Brom", "Mittens"}, chat.MentionInputOptions{Placeholder: "Try @ar"})
	if r.Changed {
		log(fmt.Sprintf("text of %d letters", utf8.RuneCountInString(s.mention)))
	}
	if r.Mentioned != "" {
		log("mentioned " + r.Mentioned)
	}
	if r.Submitted {
		log(fmt.Sprintf("sent %d letters", utf8.RuneCountInString(s.mention)))
	}
}

func demoDateSeparator(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	now := k.Now()
	ui.Column(c).Gap(k.D(4)).Children(func() {
		for _, d := range []int{0, -1, -40} {
			chat.DateSeparator(c, now.AddDate(0, 0, d), chat.DateSeparatorOptions{})
		}
	})
}

func demoUnreadDivider(c *ui.Context, log func(string)) {
	chat.UnreadDivider(c, chat.UnreadDividerOptions{})
}

func demoSystemMessage(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Width(k.D(200)).Gap(k.D(4)).Children(func() {
		chat.SystemMessage(c, "PLAYER 2 HAS JOINED", chat.SystemMessageOptions{})
		chat.SystemMessage(c, "The tavern is now open", chat.SystemMessageOptions{Slot: pixel.Good})
	})
}

func demoChannelList(c *ui.Context, log func(string)) {
	s := chatState(c)
	chans := []chat.Channel{{Name: "general", Unread: true}, {Name: "quests"}, {Name: "trading", Unread: true}, {Name: "off-topic"}}
	r := chat.ChannelList(c, "Channels", chans, &s.channel, chat.ChannelListOptions{})
	if r.Changed {
		log("# " + chans[s.channel].Name)
	}
	if r.Submitted {
		log("open # " + chans[s.channel].Name)
	}
}
