package actions

import (
	_ "embed"

	"github.com/chromedp/chromedp"
)

//go:embed js/auto_click.js
var AutoClickJS string

func AutoClick() chromedp.EvaluateAction {
	return chromedp.Evaluate(AutoClickJS, nil)
}

//go:embed js/all_onclick_value.js
var AllOnClickValueJS string

func AllOnClickValue() chromedp.EvaluateAction {
	return chromedp.Evaluate(AllOnClickValueJS, nil)
}
