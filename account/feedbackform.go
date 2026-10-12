package account

import (
	"log/slog"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FeedbackFormOptions struct {
	// Width is in dots; zero means 200.
	Width    int
	Disabled bool
}

type FeedbackFormResult struct {
	Element *ui.Element
	Sent    bool
}

// FeedbackForm takes stars and words; SEND waits for a rating.
func FeedbackForm(c *ui.Context, text string, rating *int, message *string, o FeedbackFormOptions) FeedbackFormResult {
	k := pixel.K(c)
	w := pixel.Or("FeedbackForm Width", o.Width, 200)
	var res FeedbackFormResult
	res.Element = ui.Column(c).Width(k.D(w)).Gap(k.D(5)).Role(ui.RoleGroup).Label(text).Children(func() {
		input.Rating(c, k.Msg("rating"), rating, input.RatingOptions{Disabled: o.Disabled})
		input.TextArea(c, k.Msg("message"), message, input.TextAreaOptions{Width: w, Disabled: o.Disabled})
		res.Sent = input.Button(c, k.Msg("send"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled || *rating == 0}).Clicked
	})
	// The words stay private: the log counts them.
	if res.Sent {
		slog.Debug("8bit: feedback sent", "label", text, "rating", *rating, "runes", utf8.RuneCountInString(*message))
	}
	return res
}
