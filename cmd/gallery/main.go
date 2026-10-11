// Command gallery opens the 8bitMyGoUI gallery in a native window.
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/gallery"
)

func main() {
	app := gallery.New("/")
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "8bitMyGoUI",
			Width:     1280,
			Height:    800,
			MinWidth:  640,
			MinHeight: 400,
			Content:   ui.View(app.View),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
