package marketplace

import (
	"fmt"
	"strings"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
)

func EventNewTrustedVendorRequest(vendor User) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.SiteURL
		text      = fmt.Sprintf("[@%s](%s/user/%s) has requested for a trusted vendor status",
			vendor.Username, marketUrl, vendor.Username,
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookTrustedVendors, text)
}

func EventNewTrustedVendorThreadPost(user User, vendor User, message Message) {
	var (
		marketUrl = MARKETPLACE_SETTINGS.SiteURL
		text      = fmt.Sprintf("[@%s](%s/user/%s) has posted in vendor verification thread [@%s](%s/staff/vendors/%s):\n> %s",
			user.Username, marketUrl, user.Username,
			vendor.Username, marketUrl, vendor.Username,
			strings.Replace(message.Text, "\n", "\n > ", -1), //------------|
		)
	)
	go apis.PostMattermostEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookTrustedVendors, text)
}
