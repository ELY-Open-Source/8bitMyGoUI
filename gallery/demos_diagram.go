package gallery

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/diagram"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var diagramEntries = []Entry{
	{501, "Mermaid", "Mermaid 图", "diagram", "Mermaid source on the dot grid: a quest as a flowchart, a trade, a bag.", demoMermaid},
}

// mermaidSamples are a hero's day in Mermaid.
var mermaidSamples = []struct{ name, source string }{
	{"QUEST", "flowchart TD\n  Inn[Rest at the inn] --> Gate{Gate open?}\n  Gate -- no --> Key[Find the key]\n  Key --> Gate\n" +
		"  Gate -- yes --> Boss([Fight the boss])\n  style Boss fill:#c73a3a,stroke:#33283a,color:#fff"},
	{"TRADE", "sequenceDiagram\n  participant Hero\n  participant Shop\n  Hero->>Shop: 3 potions\n  Shop-->>Hero: 45 gold\n  Hero->>Shop: Pay 45\n  Shop-->>Hero: Thanks"},
	{"BAG", "pie title Bag by weight\n  \"Potions\" : 12\n  \"Armor\" : 30\n  \"Loot\" : 18"},
}

func demoMermaid(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	at := ui.Local(c.Root(), demoKey("mermaid"), func() int { return 0 })
	ui.Column(c).Gap(k.D(8)).Children(func() {
		names := make([]string, len(mermaidSamples))
		for i, s := range mermaidSamples {
			names[i] = s.name
		}
		if input.Segmented(c, "Diagram", names, at, input.SegmentedOptions{}).Changed {
			log("drew " + names[*at])
		}
		diagram.Mermaid(c, mermaidSamples[*at].source, diagram.MermaidOptions{Label: names[*at]})
	})
}
