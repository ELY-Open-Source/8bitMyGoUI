package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

var textEntries = []Entry{
	{21, "Heading", "标题", "text", "Titles at ×3, ×2 and ×1, each with a 1-dot shadow.", demoHeading},
	{22, "Paragraph", "段落", "text", "Wrapped text with a line gap, alignment and a line limit that ends in an ellipsis.", demoParagraph},
	{23, "ShadowText", "投影字", "text", "Text with a 1-dot drop shadow in Shade.", demoShadowText},
	{24, "OutlineText", "描边字", "text", "Every glyph dot ringed in Edge, readable on any background.", demoOutlineText},
	{25, "GradientText", "渐变字", "text", "Glyph dots dithered from one slot at the top to another at the bottom.", demoGradientText},
	{26, "WaveText", "波浪字", "text", "Characters bob on a sine wave, stepped at 10 fps.", demoWaveText},
	{27, "RainbowText", "彩虹字", "text", "Each character takes the next slot of a cycling band: palette cycling.", demoRainbowText},
	{28, "Typewriter", "打字机", "text", "Dialogue that types itself out. Click or press Enter to show the rest.", demoTypewriter},
	{29, "Marquee", "跑马灯", "text", "Text that scrolls a dot a step through a window. Point at it to pause.", demoMarquee},
	{30, "BlinkText", "闪烁字", "text", "Text that blinks, and stays put under Reduce Motion.", demoBlinkText},
	{31, "ShakeText", "颤抖字", "text", "Characters jitter a dot in a seeded pattern, for fear and earthquakes.", demoShakeText},
	{32, "Kbd", "键帽", "text", "A keycap with depth; pressed, it sinks.", demoKbd},
	{33, "Code", "行内代码", "text", "An inline code chip on Well with an Edge border.", demoCode},
	{34, "Quote", "引用", "text", "A block quote: an Accent bar, a pixel quote mark, the words and a cite.", demoQuote},
	{35, "Link", "链接", "text", "Underlined text that opens a page on click or Enter.", demoLink},
	{36, "SevenSegment", "七段数码", "text", "Digits in seven segments, the unlit ones left as ghosts.", demoSevenSegment},
	{37, "Odometer", "里程表", "text", "A number whose digit wheels roll to each new value.", demoOdometer},
	{38, "MicroText", "微缩字", "text", "Upper-case text in the 3×5 micro font, for tiny labels.", demoMicroText},
	{39, "Highlight", "高亮", "text", "Matches of a query marked with a dithered highlighter.", demoHighlight},
	{40, "ScrambleText", "解码字", "text", "Text that settles from random glyphs, left to right.", demoScrambleText},
}

func demoHeading(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		text.Heading(c, "LEVEL 1", text.HeadingOptions{Level: 1})
		text.Heading(c, "Level 2", text.HeadingOptions{Level: 2, Slot: pixel.Accent})
		text.Heading(c, "Level 3 第三级", text.HeadingOptions{Level: 3})
	})
}

func demoParagraph(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	words := "The cartridge slid into the slot with a click. The screen blinked twice, then the music began, and a small knight stepped out of the castle gate."
	ui.Column(c).Gap(k.D(10)).Width(k.D(220)).Children(func() {
		text.Paragraph(c, words, text.ParagraphOptions{Gap: 2})
		text.Paragraph(c, words, text.ParagraphOptions{Align: text.Center, MaxLines: 2, Slot: pixel.Faint})
	})
}

func demoShadowText(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		text.ShadowText(c, "SCORE 004200", text.ShadowTextOptions{Slot: pixel.Accent, Scale: 2})
		text.ShadowText(c, "Hi-score 最高分", text.ShadowTextOptions{})
	})
}

func demoOutlineText(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).Padding(k.D(6)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Dither(0, 0, g.W, g.H, pixel.Info, pixel.Good, 8)
	}).Children(func() {
		text.OutlineText(c, "STAGE CLEAR", text.OutlineTextOptions{Scale: 2})
	})
}

func demoGradientText(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		text.GradientText(c, "SUNSET", text.GradientTextOptions{Top: pixel.Warn, Bottom: pixel.Bad, Scale: 3})
		text.GradientText(c, "Crystal 水晶", text.GradientTextOptions{Top: pixel.Info, Bottom: pixel.Magic, Scale: 2})
	})
}

func demoWaveText(c *ui.Context, log func(string)) {
	text.WaveText(c, "~ WELCOME TO THE BEACH ~", text.WaveTextOptions{Slot: pixel.Info, Scale: 2})
}

func demoRainbowText(c *ui.Context, log func(string)) {
	text.RainbowText(c, "SUPER PIXEL BROS", text.RainbowTextOptions{Scale: 2})
}

func demoTypewriter(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	lines := []string{
		"Old man: It's dangerous to go alone. Take this cartridge.",
		"Old man: Blow on it first. Trust me.",
		"Old man: That's all I know. Off you go.",
	}
	at := ui.Local(c.Root(), demoKey("typewriter"), func() int { return 0 })
	ui.Column(c).Gap(k.D(6)).Width(k.D(240)).Children(func() {
		r := text.Typewriter(c, lines[*at], text.TypewriterOptions{CPS: 24})
		if r.Done {
			log(fmt.Sprintf("line %d done", *at+1))
		}
		if input.Button(c, "NEXT", input.ButtonOptions{}).Clicked {
			*at = (*at + 1) % len(lines)
			log("next line")
		}
	})
}

func demoMarquee(c *ui.Context, log func(string)) {
	text.Marquee(c, "NEW HIGH SCORE! ★ INSERT COIN ★ PRESS START ★", text.MarqueeOptions{Width: 160, Slot: pixel.Accent})
}

func demoBlinkText(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(12)).Children(func() {
		text.BlinkText(c, "PRESS START", text.BlinkTextOptions{Scale: 2})
		text.BlinkText(c, "LOW HP", text.BlinkTextOptions{Slot: pixel.Bad, FPS: 12})
	})
}

func demoShakeText(c *ui.Context, log func(string)) {
	text.ShakeText(c, "IT'S RIGHT BEHIND YOU", text.ShakeTextOptions{Slot: pixel.Bad, Scale: 2, Seed: 7})
}

func demoKbd(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(4)).AlignItems(ui.End).Children(func() {
		for _, key := range []string{"↑", "↑", "↓", "↓", "←", "→", "←", "→", "B", "A"} {
			text.Kbd(c, key, text.KbdOptions{})
		}
		text.Kbd(c, "START", text.KbdOptions{Pressed: true})
	})
}

func demoCode(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
		pixel.Text(c, "Call", pixel.TextOptions{})
		text.Code(c, "pixel.Use(c, s)", text.CodeOptions{})
		pixel.Text(c, "once per frame.", pixel.TextOptions{})
	})
}

func demoQuote(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).Width(k.D(240)).Children(func() {
		text.Quote(c, "It's dangerous to go alone! Take this.", text.QuoteOptions{Cite: "An old man in a cave"})
	})
}

func demoLink(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(12)).Children(func() {
		if text.Link(c, "MyGo docs", "https://mygo.egoist.dev/docs/ui", text.LinkOptions{}).Clicked {
			log("opened MyGo docs")
		}
		text.Link(c, "Lost link", "https://example.com", text.LinkOptions{Disabled: true})
	})
}

func demoSevenSegment(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		text.SevenSegment(c, "12:59", text.SevenSegmentOptions{Scale: 2})
		text.SevenSegment(c, "-0042", text.SevenSegmentOptions{Slot: pixel.Good})
	})
}

func demoOdometer(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	score := ui.Local(c.Root(), demoKey("odometer"), func() int { return 1990 })
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		text.Odometer(c, *score, text.OdometerOptions{Digits: 6})
		if input.Button(c, "+15", input.ButtonOptions{}).Clicked {
			*score += 15
			log(fmt.Sprintf("score %d", *score))
		}
		if input.Button(c, "+500", input.ButtonOptions{}).Clicked {
			*score += 500
			log(fmt.Sprintf("score %d", *score))
		}
	})
}

func demoMicroText(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		text.MicroText(c, "hp 32/40 mp 12/18", text.MicroTextOptions{})
		text.MicroText(c, "lv 99", text.MicroTextOptions{Slot: pixel.Accent, Scale: 3})
	})
}

func demoHighlight(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		text.Highlight(c, "Slime, King Slime, Metal Slime", "slime", text.HighlightOptions{})
		text.Highlight(c, "Potion, Hi-Potion, Ether", "potion", text.HighlightOptions{Mark: pixel.Info})
	})
}

func demoScrambleText(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	words := []string{"ACCESS GRANTED", "DECRYPTING...", "WELCOME, PLAYER 1"}
	at := ui.Local(c.Root(), demoKey("scramble"), func() int { return 0 })
	ui.Column(c).Gap(k.D(6)).Children(func() {
		text.ScrambleText(c, words[*at], text.ScrambleTextOptions{Slot: pixel.Good, Scale: 2})
		if input.Button(c, "AGAIN", input.ButtonOptions{}).Clicked {
			*at = (*at + 1) % len(words)
			log("scramble " + words[*at])
		}
	})
}
