package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// glyph is a 12×12 emoji and its # dots' slot.
type glyph struct {
	sp  *pixel.Sprite
	ink pixel.Slot
}

// sets are the emoji by category, keyed by message.
var sets = []struct {
	key   string
	names []string
}{
	{"faces", display.EmojiNames},
	{"nature", []string{"sun", "moon", "cloud", "snow", "bolt", "drop", "leaf", "flower", "fire", "star", "sparkle", "heart"}},
	{"items", []string{"coin", "gem", "potion", "sword", "shield", "chest", "crown", "trophy", "skull", "gift", "key", "scroll"}},
}

var things = map[string]glyph{
	"sun": {icon.Sun, pixel.Warn}, "moon": {icon.Moon, pixel.Warn}, "cloud": {icon.Cloud, pixel.Faint},
	"snow": {icon.Snowflake, pixel.Info}, "bolt": {icon.Bolt, pixel.Warn}, "drop": {icon.Drop, pixel.Info},
	"leaf": {icon.Leaf, pixel.Good}, "flower": {icon.Flower, pixel.Accent}, "fire": {icon.Flame, pixel.Bad},
	"star": {icon.Star, pixel.Warn}, "sparkle": {icon.Sparkle, pixel.Warn}, "heart": {icon.Heart, pixel.Bad},
	"coin": {icon.Coin, pixel.Warn}, "gem": {icon.Gem, pixel.Info}, "potion": {icon.Potion, pixel.Magic},
	"sword": {icon.Sword, pixel.Edge}, "shield": {icon.Shield, pixel.Info}, "chest": {icon.Chest, pixel.Warn},
	"crown": {icon.Crown, pixel.Warn}, "trophy": {icon.Trophy, pixel.Warn}, "skull": {icon.Skull, pixel.Ink},
	"gift": {icon.Gift, pixel.Accent}, "key": {icon.Key, pixel.Warn}, "scroll": {icon.Scroll, pixel.Faint},
}

// emoji is the face or thing named name.
func emoji(name string) glyph {
	if g, ok := things[name]; ok {
		return g
	}
	for _, n := range display.EmojiNames {
		if n == name {
			return glyph{display.EmojiSprite(name), pixel.Ink}
		}
	}
	panic(fmt.Sprintf("chat: no emoji %q", name))
}
