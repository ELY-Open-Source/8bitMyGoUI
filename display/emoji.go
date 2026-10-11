package display

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var faceBase = []string{
	"...eeeeee...",
	"..eyyyyyye..",
	".eyyyyyyyye.",
	"eyyyyyyyyyye",
	"eyyyyyyyyyye",
	"eyyyyyyyyyye",
	"eyyyyyyyyyye",
	"eyyyyyyyyyye",
	"eyyyyyyyyyye",
	".eyyyyyyyye.",
	"..eyyyyyye..",
	"...eeeeee...",
}

// EmojiNames lists the faces in order.
var EmojiNames = []string{"smile", "grin", "laugh", "wink", "love", "cool", "wow", "sad",
	"cry", "angry", "sleep", "sick", "dead", "think", "meh", "blush"}

// features are each face's rows over the base.
var features = map[string]map[int]string{
	"smile": {4: "...k....k...", 5: "...k....k...", 7: "..k......k..", 8: "...kkkkkk..."},
	"grin":  {4: "...k....k...", 5: "...k....k...", 7: "..kkkkkkkk..", 8: "..kllllllk..", 9: "...kkkkkk..."},
	"laugh": {4: "...k....k...", 5: "..k.k..k.k..", 7: "..kkkkkkkk..", 8: "..krrrrrrk..", 9: "...kkkkkk..."},
	"wink":  {4: "........k...", 5: "..kkk...k...", 7: "..k......k..", 8: "...kkkkkk..."},
	"love":  {3: "..r.r..r.r..", 4: "..rrr..rrr..", 5: "...r....r...", 7: "..k......k..", 8: "...kkkkkk..."},
	"cool":  {4: ".kkkkkkkkkk.", 5: "..kkk..kkk..", 7: "..k.........", 8: "...kkkkk...."},
	"wow":   {4: "..kk....kk..", 5: "..kk....kk..", 7: ".....kk.....", 8: "....k..k....", 9: ".....kk....."},
	"sad":   {4: "...k....k...", 5: "...k....k...", 8: "...kkkkkk...", 9: "..k......k.."},
	"cry":   {4: "..kkk..kkk..", 5: "..c......c..", 6: "..c......c..", 8: "...kkkkkk...", 9: "..k......k.."},
	"angry": {3: "..k......k..", 4: "...kk..kk...", 5: "...k....k...", 8: "...kkkkkk...", 9: "..k......k.."},
	"sleep": {0: ".........kkk", 1: "..........k.", 2: ".........kkk", 5: "..kkk..kkk..", 8: "....kkkk...."},
	"sick":  {4: "..kkk..kkk..", 5: "...k....k...", 7: "..g.g.g.g...", 8: "...g.g.g.g..", 9: "...kk..kk..."},
	"dead":  {3: "..k.k..k.k..", 4: "...k....k...", 5: "..k.k..k.k..", 8: "...kkkkkk..."},
	"think": {3: "..kkk.......", 4: "...k....k...", 5: "...k....k...", 8: ".....kkk...."},
	"meh":   {4: "...k....k...", 5: "...k....k...", 8: "...kkkkkk..."},
	"blush": {4: "...k....k...", 5: "...k....k...", 6: ".AA......AA.", 7: "..k......k..", 8: "...kkkkkk..."},
}

var emoji = func() map[string]*pixel.Sprite {
	out := map[string]*pixel.Sprite{}
	for _, name := range EmojiNames {
		rows := append([]string(nil), faceBase...)
		for i, over := range features[name] {
			row := []byte(rows[i])
			for j := range over {
				if over[j] != '.' {
					row[j] = over[j]
				}
			}
			rows[i] = string(row)
		}
		out[name] = pixel.MustSprite(strings.Join(rows, "\n"))
	}
	return out
}()

type EmojiOptions struct {
	Scale int
}

// Emoji draws one of sixteen hand-placed faces by name.
func Emoji(c *ui.Context, name string, o EmojiOptions) *ui.Element {
	sp, ok := emoji[name]
	if !ok {
		panic(fmt.Sprintf("display: no emoji %q", name))
	}
	return Icon(c, sp, IconOptions{Scale: o.Scale, Label: name})
}
