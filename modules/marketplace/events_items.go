package marketplace

import (
	"fmt"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
)

func EventNewItem(item Item) {
	go apis.PostMattermostRawEvent(MARKETPLACE_SETTINGS.MattermostIncomingHookItems, createItemMattermostEvent(item))
}

func createItemMattermostEvent(item Item) apis.MattermostEvent {

	var (
		user = item.User
	)

	return apis.MattermostEvent{
		Attachments: []apis.MattermostEventAttachment{
			apis.MattermostEventAttachment{
				Fallback:   "Item: " + item.Uuid,
				AuthorName: user.Username,
				AuthorLink: fmt.Sprintf("%s/user/%s", MARKETPLACE_SETTINGS.SiteURL, user.Username),
				Fields: []apis.MattermostEventField{
					{
						Title: "Item Name",
						Value: fmt.Sprintf("[%s](%s/payments/%s/item/%s)", item.Name, MARKETPLACE_SETTINGS.SiteURL, user.Username, item.Uuid),
						Short: true,
					},
					{
						Title: "Category",
						Value: item.ItemCategory.NameEn,
						Short: true,
					},
					{
						Title: "Description",
						Value: item.Description,
						Short: false,
					},
				},
			},
		},
	}
}
