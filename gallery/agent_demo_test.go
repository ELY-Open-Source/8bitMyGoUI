package gallery

import (
	"strings"
	"testing"
)

func TestAgentDemosLogTheirActions(t *testing.T) {
	// Clicks that only focus a component may log nothing.
	quiet := map[string]bool{"Prompt": true}
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/451-agentcard", []string{"READ A BIG FILE", "COMPACT"}}, {"/c/452-toolcall", []string{"CAST GREP"}},
		{"/c/453-planquest", []string{"COMPLETE STEP"}}, {"/c/454-thinkingorb", []string{"PONDER"}},
		{"/c/455-tokenmana", []string{"CAST (12K)", "DRINK ETHER"}}, {"/c/456-approvalprompt", []string{"YES", "NO"}},
		{"/c/457-agentlog", []string{"NEXT STEP"}}, {"/c/458-subagentparty", []string{"SUMMON", "NEXT TURN"}},
		{"/c/459-memorycrystals", []string{"Memories"}}, {"/c/462-modelselect", []string{"Models"}},
		{"/c/463-promptbox", []string{"Prompt", "type:ship it", "key:Enter"}}, {"/c/464-costcoins", []string{"RUN THE AGENT", "NEW MONTH"}},
		{"/c/465-checkpoints", []string{"RESTORE Tests green", "SAVE HERE"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			if k, ok := strings.CutPrefix(step, "key:"); ok {
				tt.Key(0, letters[k])
			} else if s, ok := strings.CutPrefix(step, "type:"); ok {
				tt.Type(s)
			} else {
				click(t, tt, step)
			}
			if len(a.logs[slug]) <= before && !quiet[step] {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
		// A prompt may hold secrets: the log counts, never quotes.
		if strings.Contains(strings.Join(a.logs[slug], " "), "ship it") {
			t.Errorf("%s logged the prompt: %q", slug, a.logs[slug])
		}
	}
}
