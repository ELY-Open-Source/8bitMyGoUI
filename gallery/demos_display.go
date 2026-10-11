package gallery

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var displayEntries = []Entry{
	{41, "Icon", "图标", "display", "A sprite in one ink slot at ×1 to ×4, named for screen readers or hidden.", demoIcon},
	{42, "Sprite", "精灵", "display", "A multi-slot sprite that plays its frames and can face the other way.", demoSprite},
	{43, "Picture", "像素图", "display", "Any image reduced to the palette with Bayer dither. Press P: it recolors.", demoPicture},
	{44, "Avatar", "头像", "display", "A face built from a name: skin, hair, eyes, mouth and background.", demoAvatar},
	{45, "AvatarGroup", "头像组", "display", "Overlapping faces and a chip for the rest.", demoAvatarGroup},
	{46, "Badge", "徽标", "display", "A count or a word in a notched pill; counts stop at 99+.", demoBadge},
	{47, "Tag", "标签", "display", "A chip with text, an icon and a remove box.", demoTag},
	{48, "StatusDot", "状态点", "display", "Online, away, busy or offline; a live online dot blinks.", demoStatusDot},
	{49, "Medal", "奖牌", "display", "Gold, silver and bronze on ribbons, rank on the face.", demoMedal},
	{50, "Trophy", "奖杯", "display", "A gold cup over an engraved plaque.", demoTrophy},
	{51, "Stat", "数据", "display", "A big value, its label and a change up or down.", demoStat},
	{52, "KeyValue", "键值", "display", "A label and a value joined by a dotted leader.", demoKeyValue},
	{53, "Portrait", "立绘框", "display", "A character in a frame over a nameplate.", demoPortrait},
	{54, "Ribbon", "绶带", "display", "A title banner with folded, notched tails.", demoRibbon},
	{55, "Stamp", "印章", "display", "A rubber stamp whose border wore unevenly.", demoStamp},
	{56, "Emoji", "表情", "display", "Sixteen faces placed dot by dot.", demoEmoji},
	{57, "Flag", "旗帜", "display", "A flag on a pole, waving in steps.", demoFlag},
	{58, "Gem", "宝石", "display", "A cut gem for each rarity, with a stepped sparkle.", demoGem},
	{59, "RarityFrame", "稀有度框", "display", "A border by rarity; a legendary one carries a running glint.", demoRarityFrame},
	{60, "Bubble", "气泡", "display", "A speech bubble with its tail on any side.", demoBubble},
	{61, "Swatch", "色板", "display", "The sixteen slots of the palette in use, by name.", demoSwatch},
	{62, "Placeholder", "占位图", "display", "A dithered box with a cross and its size, for art to come.", demoPlaceholder},
	{63, "Figure", "插图", "display", "Content over a caption in Faint.", demoFigure},
	{64, "QRCode", "二维码", "display", "A scannable QR code in the darkest and lightest slots.", demoQRCode},
	{65, "Signpost", "路牌", "display", "A post with up to three boards. Each board is a button.", demoSignpost},
}

// sunset is a smooth picture for Picture to reduce.
var sunset = func() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 96, 64))
	for y := range 64 {
		for x := range 96 {
			t := float64(y) / 63
			c := color.RGBA{uint8(250 - 120*t), uint8(170 - 120*t), uint8(90 + 90*t), 255}
			if dx, dy := float64(x-60), float64(y-34); dx*dx+dy*dy < 14*14 {
				c = color.RGBA{255, 220, 120, 255}
			}
			if float64(y) > 44+6*math.Sin(float64(x)/9) {
				c = color.RGBA{50, 40, 70, 255}
			}
			img.SetRGBA(x, y, c)
		}
	}
	return img
}()

func demoIcon(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.End).Children(func() {
		for sc := 1; sc <= 4; sc++ {
			display.Icon(c, icon.Heart, display.IconOptions{Slot: pixel.Bad, Scale: sc, Label: fmt.Sprintf("heart ×%d", sc)})
		}
	})
}

func demoSprite(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(12)).Children(func() {
		display.Sprite(c, slime, display.SpriteOptions{Scale: 3, FPS: 2, Label: "slime"})
		display.Sprite(c, slime, display.SpriteOptions{Scale: 3, FPS: 2, Mirror: true, Label: "slime, mirrored"})
	})
}

func demoPicture(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.End).Children(func() {
		display.Picture(c, sunset, display.PictureOptions{Label: "sunset"})
		display.Picture(c, sunset, display.PictureOptions{W: 48, H: 32, Label: "small sunset"})
	})
}

var party = []string{"Aria", "Brom", "Cass", "Dax", "Eve", "Finn", "Gus"}

func demoAvatar(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).Children(func() {
		for _, n := range party[:5] {
			display.Avatar(c, n, display.AvatarOptions{Scale: 2})
		}
		display.Avatar(c, "Sunny", display.AvatarOptions{Scale: 2, Picture: sunset})
	})
}

func demoAvatarGroup(c *ui.Context, log func(string)) {
	display.AvatarGroup(c, party, display.AvatarGroupOptions{Scale: 2})
}

func demoBadge(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		for _, n := range []int{3, 42, 128} {
			display.Badge(c, n, display.BadgeOptions{})
		}
		display.Badge(c, 0, display.BadgeOptions{Text: "NEW", Slot: pixel.Good})
		display.Badge(c, 0, display.BadgeOptions{Text: "HOT", Slot: pixel.Warn})
	})
}

func demoTag(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	tags := ui.Local(c.Root(), demoKey("tags"), func() []string { return []string{"fire", "ice", "thunder", "wind"} })
	ui.Column(c).Gap(k.D(6)).Children(func() {
		ui.Row(c).Wrap().Gap(k.D(4)).Children(func() {
			for i, t := range *tags {
				if display.Tag(c, t, display.TagOptions{Icon: icon.Sparkle, Removable: true}).Removed {
					*tags = append((*tags)[:i:i], (*tags)[i+1:]...)
					log("removed " + t)
				}
			}
			display.Tag(c, "locked", display.TagOptions{Icon: icon.Lock, Removable: true, Disabled: true})
		})
		if button(c, "RESET") {
			*tags = []string{"fire", "ice", "thunder", "wind"}
			log("tags reset")
		}
	})
}

func demoStatusDot(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		for _, s := range []display.Status{display.Online, display.Away, display.Busy, display.Offline} {
			ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
				display.StatusDot(c, s, display.StatusDotOptions{Live: true})
				label(c, []string{"Player 1", "Player 2", "Player 3", "Player 4"}[s], pixel.Ink)
			})
		}
	})
}

func demoMedal(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(10)).Children(func() {
		for _, r := range []display.Rank{display.Gold, display.Silver, display.Bronze} {
			display.Medal(c, r, display.MedalOptions{Scale: 2})
		}
	})
}

func demoTrophy(c *ui.Context, log func(string)) {
	display.Trophy(c, "Champion 2026", display.TrophyOptions{Scale: 2})
}

func demoStat(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(24)).Children(func() {
		display.Stat(c, "SCORE", "42,000", display.StatOptions{Delta: "12%"})
		display.Stat(c, "LIVES", "2", display.StatOptions{Delta: "1", Down: true})
	})
}

func demoKeyValue(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Width(k.D(160)).Gap(k.D(2)).Children(func() {
		for _, kv := range [][2]string{{"STR", "18"}, {"DEX", "12"}, {"INT", "7"}, {"LUCK", "99"}} {
			display.KeyValue(c, kv[0], kv[1], display.KeyValueOptions{})
		}
	})
}

func demoPortrait(c *ui.Context, log func(string)) {
	display.Portrait(c, "Sir Brom", display.PortraitOptions{Size: 48}, func() {
		display.Avatar(c, "Brom", display.AvatarOptions{Scale: 3})
	})
}

func demoRibbon(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		display.Ribbon(c, "STAGE 1-1", display.RibbonOptions{})
		display.Ribbon(c, "BONUS ROUND", display.RibbonOptions{Slot: pixel.Info})
	})
}

func demoStamp(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(10)).Children(func() {
		display.Stamp(c, "APPROVED", display.StampOptions{Slot: pixel.Good, Seed: 3})
		display.Stamp(c, "VOID", display.StampOptions{Seed: 9})
	})
}

func demoEmoji(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Wrap().Gap(k.D(6)).Width(k.D(240)).Children(func() {
		for _, name := range display.EmojiNames {
			display.Emoji(c, name, display.EmojiOptions{Scale: 2})
		}
	})
}

func demoFlag(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(12)).Children(func() {
		display.Flag(c, display.FlagOptions{Scale: 2, Label: "red flag"})
		display.Flag(c, display.FlagOptions{Slot: pixel.Info, Scale: 2, Label: "blue flag"})
	})
}

func demoGem(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).Children(func() {
		for r := display.Common; r <= display.Legendary; r++ {
			display.Gem(c, r, display.GemOptions{Scale: 2})
		}
	})
}

func demoRarityFrame(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	items := []string{"Stick", "Bow", "Rapier", "Dragon Lance", "Excalibur"}
	ui.Row(c).Wrap().Gap(k.D(6)).Children(func() {
		for r := display.Common; r <= display.Legendary; r++ {
			display.RarityFrame(c, r, display.RarityFrameOptions{}, func() { label(c, items[r], pixel.Ink) })
		}
	})
}

func demoBubble(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	lines := []string{"Hey!", "Look up.", "Psst.", "Over here."}
	ui.Row(c).Wrap().Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		for i, side := range []display.Side{display.Bottom, display.Top, display.Left, display.Right} {
			display.Bubble(c, display.BubbleOptions{Tail: side}, func() { label(c, lines[i], pixel.Ink) })
		}
	})
}

func demoSwatch(c *ui.Context, log func(string)) {
	display.Swatch(c, display.SwatchOptions{})
}

func demoPlaceholder(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).Children(func() {
		display.Placeholder(c, display.PlaceholderOptions{})
		display.Placeholder(c, display.PlaceholderOptions{W: 120, H: 40})
	})
}

func demoFigure(c *ui.Context, log func(string)) {
	display.Figure(c, "Fig. 1 — a slime, at rest", display.FigureOptions{}, func() {
		display.Sprite(c, slime, display.SpriteOptions{Scale: 4, Label: "slime"})
	})
}

func demoQRCode(c *ui.Context, log func(string)) {
	display.QRCode(c, "https://mygo.egoist.dev/docs/ui", display.QRCodeOptions{})
}

func demoSignpost(c *ui.Context, log func(string)) {
	boards := []string{"Town", "Dungeon", "Shop"}
	if r := display.Signpost(c, boards, display.SignpostOptions{}); r.Clicked >= 0 {
		log("to the " + boards[r.Clicked])
	}
}
