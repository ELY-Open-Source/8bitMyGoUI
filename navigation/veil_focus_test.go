package navigation_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// veiled builds under a veil the test raises and lifts.
func veiled(inert *bool, build func(c *ui.Context)) *ui.Tester {
	return view(600, 400, newClock(), false, func(c *ui.Context) {
		s := pixel.Current(c)
		s.Inert = *inert
		pixel.Scope(c, s, func() { build(c) })
	})
}

func TestOpenMenusGiveFocusBackAfterAVeil(t *testing.T) {
	items := []navigation.MenuItem{{Label: "One"}, {Label: "Two"}}
	for name, build := range map[string]func(c *ui.Context){
		"dropdown": func(c *ui.Context) { navigation.DropdownMenu(c, "File", items, navigation.DropdownMenuOptions{}) },
		"menubar": func(c *ui.Context) {
			navigation.MenuBar(c, "Bar", []navigation.Menu{{Title: "Edit", Items: items}}, navigation.MenuBarOptions{})
		},
	} {
		inert := false
		tt := veiled(&inert, build)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyDown)
		if !tt.HasText("Two") {
			t.Fatalf("%s did not open: %q", name, tt.Texts())
		}
		inert = true
		tt.Frame()
		inert = false
		tt.Frame()
		// Only focus on the opener lets Down open it again.
		tt.Key(0, ui.KeyDown)
		if !tt.HasText("Two") {
			t.Errorf("%s did not get its focus back after the veil", name)
		}
	}
}
