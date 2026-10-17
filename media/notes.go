package media

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// keyboard plays C4–C5: A–K white, WETYU black.
var keyboard = [13]struct {
	key    ui.Key
	letter string
}{
	{ui.KeyA, "A"}, {ui.KeyW, "W"}, {ui.KeyS, "S"}, {ui.KeyE, "E"}, {ui.KeyD, "D"}, {ui.KeyF, "F"}, {ui.KeyT, "T"},
	{ui.KeyG, "G"}, {ui.KeyY, "Y"}, {ui.KeyH, "H"}, {ui.KeyU, "U"}, {ui.KeyJ, "J"}, {ui.KeyK, "K"},
}

var names = [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

// note names the semitone above C4, like "C#4".
func note(semi int) string {
	return fmt.Sprintf("%s%d", names[semi%12], 4+semi/12)
}

// played are the semitones whose letters e heard this frame.
func played(e ui.Element) []int {
	var out []int
	for semi, kb := range keyboard {
		if pixel.Shortcut(e, 0, kb.key) {
			out = append(out, semi)
		}
	}
	return out
}
