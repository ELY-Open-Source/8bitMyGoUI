package gallery

import (
	"strings"
	"testing"
)

func TestDeskDemosLogTheirActions(t *testing.T) {
	// Focusing, opening a menu or moving a light logs nothing.
	quiet := map[string]bool{"key:Tab": true, "Start": true, "key:Up": true, "FILE": true, "OPEN:": true}
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/466-desktop", []string{"MY COMPUTER"}}, {"/c/473-minesweeper", []string{"NEW GAME"}},
		{"/c/479-snake", []string{"Snake pit"}}, {"/c/481-rundialog", []string{"OPEN:", "type:ls"}},
		{"/c/467-desktopicon", []string{"MY COMPUTER", "key:Enter"}}, {"/c/468-taskbar", []string{"ONLINE"}},
		{"/c/469-startmenu", []string{"Start", "key:Up", "key:Enter"}}, {"/c/470-appwindow", []string{"CLOSE", "REOPEN"}},
		{"/c/471-calculator", []string{"Pocket calculator"}}, {"/c/472-notepad", []string{"FILE", "SAVE"}},
		{"/c/474-controlpanel", []string{"Control panel"}}, {"/c/475-trayicons", []string{"ONLINE"}},
		{"/c/477-bootscreen", []string{"REBOOT"}}, {"/c/478-screensaver", []string{"START SAVER", "key:A"}},
		{"/c/480-dock", []string{"MAIL"}}, {"/c/481-rundialog", []string{"OK", "CANCEL"}},
		{"/c/482-diskdefrag", []string{"START", "NEW DISK"}}, {"/c/483-wallpaperpicker", []string{"Wallpaper"}},
		{"/c/484-clipboardpanel", []string{"git push gitea main"}}, {"/c/485-lockscreen", []string{"LOCK", "key:A"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			if k, ok := strings.CutPrefix(step, "key:"); ok {
				tt.Key(0, letters[k])
				tt.Frame()
			} else if s, ok := strings.CutPrefix(step, "type:"); ok {
				tt.Type(s)
			} else {
				click(t, tt, step)
			}
			if len(a.logs[slug]) <= before && !quiet[step] {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}
