package pixel

import (
	"fmt"
	"strings"
)

// Sprite is hand-placed art in palette slots, in frames.
type Sprite struct {
	W, H   int
	frames [][]int8
}

const (
	cellClear int8 = -1
	cellInk   int8 = -2
)

// Art letters naming slots; '.' is clear, '#' the ink.
const artSlots = "kbflsenwaAigyrcm"

// MustSprite parses string art; blank lines split frames.
func MustSprite(art string) *Sprite {
	var frames [][]string
	var cur []string
	for _, line := range strings.Split(strings.TrimSpace(art), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if cur != nil {
				frames, cur = append(frames, cur), nil
			}
			continue
		}
		cur = append(cur, line)
	}
	if cur == nil {
		panic("pixel: empty sprite")
	}
	frames = append(frames, cur)
	sp := &Sprite{W: len(frames[0][0]), H: len(frames[0])}
	for _, f := range frames {
		if len(f) != sp.H {
			panic(fmt.Sprintf("pixel: sprite frames of %d and %d rows", sp.H, len(f)))
		}
		px := make([]int8, 0, sp.W*sp.H)
		for _, row := range f {
			if len(row) != sp.W {
				panic(fmt.Sprintf("pixel: sprite row %q is not %d wide", row, sp.W))
			}
			for i := range len(row) {
				px = append(px, cell(row[i]))
			}
		}
		sp.frames = append(sp.frames, px)
	}
	return sp
}

func cell(c byte) int8 {
	switch c {
	case '.':
		return cellClear
	case '#':
		return cellInk
	}
	if i := strings.IndexByte(artSlots, c); i >= 0 {
		return int8(i)
	}
	panic(fmt.Sprintf("pixel: sprite character %q", c))
}

func (s *Sprite) Frames() int { return len(s.frames) }

// Mirror returns s flipped left to right.
func (s *Sprite) Mirror() *Sprite {
	m := &Sprite{W: s.W, H: s.H}
	for _, px := range s.frames {
		out := make([]int8, len(px))
		for y := range s.H {
			for x := range s.W {
				out[y*s.W+x] = px[y*s.W+s.W-1-x]
			}
		}
		m.frames = append(m.frames, out)
	}
	return m
}

// Sprite draws frame of sp; # takes ink; frames wrap.
func (g Canvas) Sprite(x, y int, sp *Sprite, frame int, ink Slot, scale int) {
	g.sprite(x, y, sp, frame, ink, scale, false)
}

// Silhouette draws sp's shape all in one slot.
func (g Canvas) Silhouette(x, y int, sp *Sprite, frame int, s Slot, scale int) {
	g.sprite(x, y, sp, frame, s, scale, true)
}

func (g Canvas) sprite(x, y int, sp *Sprite, frame int, ink Slot, scale int, mono bool) {
	checkScale(scale)
	n := len(sp.frames)
	px := sp.frames[(frame%n+n)%n]
	for row := range sp.H {
		line := px[row*sp.W : (row+1)*sp.W]
		for col := 0; col < sp.W; {
			v := line[col]
			end := col + 1
			for end < sp.W && line[end] == v {
				end++
			}
			if v != cellClear {
				s := ink
				if v >= 0 && !mono {
					s = Slot(v)
				}
				g.Rect(x+col*scale, y+row*scale, (end-col)*scale, scale, s)
			}
			col = end
		}
	}
}
