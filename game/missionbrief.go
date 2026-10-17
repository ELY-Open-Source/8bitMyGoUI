package game

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MissionBriefOptions struct {
	Disabled bool
}

type MissionBriefResult struct {
	Element ui.Element
	Start   bool
}

// MissionBrief sets the map beside the objectives and START.
func MissionBrief(c *ui.Context, q Quest, area []string, at Spot, marks []Marker, o MissionBriefOptions) MissionBriefResult {
	k := pixel.K(c)
	row := ui.Row(c).Gap(k.D(10)).AlignItems(ui.Start).Role(ui.RoleGroup).Label(q.Title)
	res := MissionBriefResult{Element: row}
	row.Children(func() {
		Minimap(c, q.Title, area, nil, at, marks, MinimapOptions{Cell: 5})
		ui.Column(c).Gap(k.D(8)).Children(func() {
			QuestTracker(c, q, QuestTrackerOptions{})
			res.Start = input.Button(c, k.Msg("start"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled}).Clicked
		})
	})
	if res.Start {
		slog.Debug("8bit: mission start", "quest", q.Title)
	}
	return res
}
