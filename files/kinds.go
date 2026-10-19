package files

import (
	"fmt"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Kind is what a file holds; it picks the icon.
type Kind int

const (
	Folder Kind = iota
	Text
	Image
	Audio
	Code
	Archive
)

type kindArt struct {
	msg string
	art *pixel.Sprite
}

// kinds name each Kind's message and 12×12 page (S7).
var kinds = []kindArt{
	{"folder", pixel.MustSprite(`
			............
			............
			.eeee.......
			.eyyye......
			.eyyyyeeeee.
			.eyyyyyyyye.
			.eeeeeeeeee.
			.elllyyyyye.
			.eyyyyyyyye.
			.eyyyyyyyye.
			.eeeeeeeeee.
			............`)},
	{"text_file", pixel.MustSprite(`
			.eeeeeee....
			.elllllee...
			.elnnnlese..
			.elllllesse.
			.ellllleeee.
			.elllllllle.
			.elnnnnnnle.
			.elllllllle.
			.elnnnnnnle.
			.elllllllle.
			.elnnnnllle.
			.eeeeeeeeee.`)},
	{"image_file", pixel.MustSprite(`
			.eeeeeee....
			.elllllee...
			.elllllese..
			.elllllesse.
			.ellllleeee.
			.elllllllle.
			.elccccyyle.
			.elccgcccle.
			.elcgggcgle.
			.elggggggle.
			.elllllllle.
			.eeeeeeeeee.`)},
	{"audio_file", pixel.MustSprite(`
			.eeeeeee....
			.elllllee...
			.elllllese..
			.elllllesse.
			.ellllleeee.
			.ellllmmlle.
			.ellllmlmle.
			.ellllmllle.
			.ellmmmllle.
			.elmmmmllle.
			.ellmmlllle.
			.eeeeeeeeee.`)},
	{"code_file", pixel.MustSprite(`
			.eeeeeee....
			.elllllee...
			.elllllese..
			.elllllesse.
			.ellllleeee.
			.elllllllle.
			.ellcllclle.
			.elcllllcle.
			.ecllllllce.
			.elcllllcle.
			.ellcllclle.
			.eeeeeeeeee.`)},
	{"archive_file", pixel.MustSprite(`
			.eeeeeee....
			.ellselee...
			.elleslese..
			.ellselesse.
			.ellesleeee.
			.ellselllle.
			.elllllllle.
			.elyyyyllle.
			.elyeeyllle.
			.elyyyyllle.
			.elllllllle.
			.eeeeeeeeee.`)},
}

// look returns kind's message and art; unknown kinds panic.
func look(kind Kind) kindArt {
	if kind < 0 || int(kind) >= len(kinds) {
		panic(fmt.Sprintf("files: unknown kind %d", kind))
	}
	return kinds[kind]
}
