package marketplace

import (
	"fmt"
	"strings"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
)

func EventNewWarning(warning UserWarning) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.SiteURL
		text      = fmt.Sprintf("[@%s](%s/user/%s) has reported [@%s](%s/user/%s) for:\n> %s",
			warning.Reporter.Username, marketUrl, warning.Reporter.Username,
			warning.User.Username, marketUrl, warning.User.Username,
			strings.Replace(warning.Text, "\n", "\n > ", -1),
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookWarnings, text)
}

func EventWarningStatsUpdate(warning UserWarning, staff User) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.SiteURL
		text      = fmt.Sprintf("[@%s](%s/user/%s) has updated warning for [@%s](%s/user/%s) with severety **%s**",
			staff.Username, marketUrl, staff.Username,
			warning.User.Username, marketUrl, warning.User.Username,
			warning.Severety,
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookWarnings, text)
}
