package input_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
)

func TestLogsCountWhatWasTyped(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(prev)
	name, tags, hero := "", []string{}, "Hero"
	tt := view(500, 400, nil, func(c *ui.Context) {
		input.NameEntry(c, "Initials", &name, input.NameEntryOptions{})
		input.TokenField(c, "Tags", &tags, input.TokenFieldOptions{})
		input.EditableText(c, "Name", &hero, input.EditableTextOptions{})
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	for range 4 {
		tt.Key(0, ui.KeyEnter)
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("zork")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Type("Quux")
	tt.Key(0, ui.KeyEnter)
	if name != "BBB" || len(tags) != 1 || !strings.HasSuffix(hero, "Quux") {
		t.Fatalf("typed %q, %q, %q", name, tags, hero)
	}
	for _, s := range []string{"BBB", "zork", "Quux"} {
		if strings.Contains(buf.String(), s) {
			t.Errorf("the log holds %q", s)
		}
	}
}

func TestAutocompleteLogsNoWords(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(prev)
	q := ""
	tt := view(400, 200, nil, func(c *ui.Context) {
		input.Autocomplete(c, "Find", &q, []string{"QuuxSecret"}, input.AutocompleteOptions{})
	})
	tt.Key(0, ui.KeyTab)
	tt.Type("Quu")
	tt.Key(0, ui.KeyTab)
	if q != "QuuxSecret" || strings.Contains(buf.String(), "Secret") {
		t.Fatalf("value %q; the log has the word: %v", q, strings.Contains(buf.String(), "Secret"))
	}
}

func TestComboboxLogsNoWords(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(prev)
	v := ""
	tt := view(400, 300, nil, func(c *ui.Context) { input.Combobox(c, "Find", &v, []string{"QuuxSecret"}, input.ComboboxOptions{}) })
	tt.Key(0, ui.KeyTab)
	tt.Type("QuuxSecret")
	tt.Key(0, ui.KeyEnter)
	if v != "QuuxSecret" || strings.Contains(buf.String(), "QuuxSecret") {
		t.Fatalf("value %q; the log has the word: %v", v, strings.Contains(buf.String(), "QuuxSecret"))
	}
}
