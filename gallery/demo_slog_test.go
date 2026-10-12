package gallery

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDemoEventsKeepTypedTextOutOfLogs(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(prev)
	// A wrong vault code: four digits, then WRONG CODE.
	a, tt := open("/c/094-pininput")
	click(t, tt, "Vault code")
	tt.Type("9876")
	shown := strings.Join(a.logs["094-pininput"], "\n")
	if !strings.Contains(shown, "WRONG CODE") || !strings.Contains(buf.String(), "8bit: demo") {
		t.Fatalf("no demo event logged: %q", shown)
	}
	if strings.Contains(buf.String(), "9876") || strings.Contains(shown, "9876") {
		t.Fatalf("the code reached a log: %q", shown)
	}
}

func TestDemoPanelsShowWhatLogsDoNot(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(prev)
	a, tt := open("/c/099-editabletext")
	click(t, tt, "Hero name")
	tt.Key(0, ui.KeyEnter)
	tt.Type("Zork")
	tt.Key(0, ui.KeyEnter)
	shown := strings.Join(a.logs["099-editabletext"], "\n")
	if !strings.Contains(shown, "Zork") || strings.Contains(buf.String(), "Zork") {
		t.Fatalf("panel %q; the log has the name: %v", shown, strings.Contains(buf.String(), "Zork"))
	}
}
