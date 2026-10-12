package desk_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/desk"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
)

func TestStartMenuGetsFocusBackAfterAVeil(t *testing.T) {
	wd := newWorld()
	programs := []desk.Program{{Name: "Paint", Icon: icon.File}, {Name: "Notes", Icon: icon.File}}
	tt := view(400, 300, wd, func(c *ui.Context) { desk.StartMenu(c, "Start", programs, desk.StartMenuOptions{}) })
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	if !tt.HasText("Notes") {
		t.Fatalf("the menu did not open: %q", tt.Texts())
	}
	wd.inert = true
	tt.Frame()
	wd.inert = false
	tt.Frame()
	tt.Key(0, ui.KeyUp)
	if !tt.HasText("Notes") {
		t.Fatal("START did not get its focus back after the veil")
	}
}
