package marketplace

import (
	"fmt"
	"github.com/gocraft/web"
	"net/http"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) EditAdvertisings(w web.ResponseWriter, r *web.Request) {
	ads, err := FindAdvertisingByVendor(c.ViewUser.Uuid)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	c.Advertisings = ads
	c.AdvertisingCost = MARKETPLACE_SETTINGS.AdvertisingCost
	c.Items = FindItemsForSeller(c.ViewUser.Uuid)
	c.ViewSeller = &c.ViewUser
	// c.ViewSeller = Seller{c.ViewUser.User}.ViewSeller(c.ViewUser.User.Language)
	c.USDBTCRate = GetCurrencyRate("BTC", "USD")
	util.RenderTemplate(w, "advertising/edit", c)
}

func (c *Context) AddAdvertisingsPOST(w web.ResponseWriter, r *web.Request) {
	count := 10000
	vendorUuid := c.ViewUser.Uuid
	comment := r.FormValue("text")
	itemUuid := r.FormValue("item")

	priceUSD := MARKETPLACE_SETTINGS.AdvertisingCost
	c.USDBTCRate = GetCurrencyRate("BTC", "USD")

	price := priceUSD / c.USDBTCRate

	userWallets := c.ViewUser.User.FindUserBitcoinWallets()
	if userWallets.Balance().Balance < price {
		c.Error = fmt.Sprintf("Please deposit %f BTC to your onsite wallet.", price)
		c.EditAdvertisings(w, r)
		return
	}

	addr, err := apis.GenerateBTCAddress("advertising")
	if err != nil {
		c.Error = err.Error()
		c.EditAdvertisings(w, r)
		return
	}

	_, err = userWallets.Send(addr, price)
	if err != nil {
		c.Error = err.Error()
		c.EditAdvertisings(w, r)
		return
	}

	err = CreateAdvertising(comment, count, vendorUuid, itemUuid)
	if err != nil {
		c.Error = err.Error()
		c.EditAdvertisings(w, r)
		return
	}

	http.Redirect(w, r.Request, "/seller/"+c.ViewSeller.Username+"/advertisings", 302)

}
