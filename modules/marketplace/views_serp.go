package marketplace

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) parseItemsQuery(w web.ResponseWriter, r *web.Request) {

	if len(r.URL.Query()["shipping-to"]) > 0 {
		c.ShippingTo = r.URL.Query()["shipping-to"][0]
	}
	if len(r.URL.Query()["shipping-from"]) > 0 {
		c.ShippingFrom = r.URL.Query()["shipping-from"][0]
	}

	if len(r.URL.Query()["category"]) > 0 {
		cid, err := strconv.ParseInt(r.URL.Query()["category"][0], 10, 64)
		c.CategoryID = int(cid)
		if err == nil {
			ic, err := FindCategoryByID(c.CategoryID)
			if err == nil {
				if ic.ParentCategory() != nil {
					if ic.ParentCategory().ParentCategory() != nil {
						c.Category = ic.ParentCategory().ParentCategory().NameEn
					} else {
						c.Category = ic.ParentCategory().NameEn
					}
					c.SubCategory = ic.NameEn
				} else {
					c.Category = ic.NameEn
				}
			}
		}
	}

	if len(r.URL.Query()["city-id"]) > 0 {
		cityId, err := strconv.ParseInt(r.URL.Query()["city-id"][0], 10, 32)
		if err != nil || cityId < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.CityID = int(cityId)

		city, err := FindCityByID(c.CityID)
		if err == nil {
			c.City = city.NameEn
		} else {
			loc := GetLocalization(c.ViewUser.Language)
			c.City = loc.Items.AllCountries
		}
	} else {
		loc := GetLocalization(c.ViewUser.Language)
		c.City = loc.Items.AllCountries
	}

	if len(r.URL.Query()["page"]) > 0 {
		strPage := r.URL.Query()["page"][0]
		page, err := strconv.ParseInt(strPage, 10, 32)
		if err != nil || page < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.Page = int(page) - 1
	}

	if len(r.URL.Query()["query"]) > 0 {
		c.Query = r.URL.Query()["query"][0]
	}

	if len(r.URL.Query()["sortby"]) > 0 {
		c.SortBy = r.URL.Query()["sortby"][0]
	} else {
		c.SortBy = "popularity"
	}

	if len(r.URL.Query()["account"]) > 0 {
		c.Account = r.URL.Query()["account"][0]
	} else {
		c.Account = "all"
	}

	if len(r.URL.Query()["shipping-from"]) > 0 {
		c.ShippingFrom = r.URL.Query()["shipping-from"][0]
	}

	if len(r.URL.Query()["shipping-to"]) > 0 {
		c.ShippingTo = r.URL.Query()["shipping-to"][0]
	}

	c.SelectedPackageType = r.PathParams["package_type"]
}

func (c *Context) listAvailableItems(w web.ResponseWriter, r *web.Request) {
	c.parseItemsQuery(w, r)

	userUuid := ""
	if c.ViewSeller != nil {
		userUuid = c.ViewSeller.Uuid
	}

	ais := FindAvailableItems(userUuid)

	filteredAvailableItems := ais.
		Filter(
			c.CategoryID,
			c.CityID,
			c.SelectedPackageType,
			c.Query,
			c.ShippingTo,
			c.ShippingFrom,
			c.Account,
			userUuid,
		).
		Sort("popularity")

	menuFilteredItems := ais.Filter(
		c.CategoryID,
		0,
		c.SelectedPackageType,
		c.Query,
		c.ShippingTo,
		c.ShippingFrom,
		c.Account,
		userUuid,
	)
	numberOfAvailableItems := len(filteredAvailableItems)

	start := c.Page * MARKETPLACE_SETTINGS.ItemsPerPage
	if 0 > start {
		start = 0
	}

	finish := start + MARKETPLACE_SETTINGS.ItemsPerPage
	if numberOfAvailableItems < finish {
		finish = numberOfAvailableItems
	}

	pagedAvailableItems := filteredAvailableItems[start:finish]

	c.NumberOfPages = int(math.Ceil(float64(numberOfAvailableItems) / float64(MARKETPLACE_SETTINGS.ItemsPerPage)))
	c.ViewAvailableItems = pagedAvailableItems.ViewAvailableItems(c.ViewUser.Language, c.ViewUser.Currency)
	c.ItemCategories = CacheGetCategories(c.SelectedPackageType, c.ShippingTo, c.ShippingFrom, userUuid, c.CityID)

	c.ShippingToList = filteredAvailableItems.ShippingToList()
	c.ShippingFromList = filteredAvailableItems.ShippingFromList()
	c.GeoCities = menuFilteredItems.DropCitiesList()

	// paging
	for i := 0; i < c.NumberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}
	c.Page += 1
}

func (c *Context) ListAvailableItems(w web.ResponseWriter, r *web.Request) {
	c.listAvailableItems(w, r)
	util.RenderTemplate(w, "serp/list_items", c)
}

func (c *Context) ListAvailableVendors(w web.ResponseWriter, r *web.Request) {

	c.parseItemsQuery(w, r)

	ais := FindAvailableItems("")

	filteredAvailableItems := ais.
		Filter(c.CategoryID, c.CityID, c.SelectedPackageType, c.Query, c.ShippingTo, c.ShippingFrom, c.Account, "")

	vendors := filteredAvailableItems.VendorList().Sort(c.SortBy)

	menuFilteredItems := ais.Filter(c.CategoryID, 0, c.SelectedPackageType, c.Query, c.ShippingTo, c.ShippingTo, c.Account, "")
	numberOfAvailableItems := len(vendors)

	start := c.Page * MARKETPLACE_SETTINGS.ItemsPerPage
	if 0 > start {
		start = 0
	}

	finish := start + MARKETPLACE_SETTINGS.ItemsPerPage
	if numberOfAvailableItems < finish {
		finish = numberOfAvailableItems
	}

	pagedVendors := vendors[start:finish]

	c.NumberOfPages = int(math.Ceil(float64(numberOfAvailableItems) / float64(MARKETPLACE_SETTINGS.ItemsPerPage)))
	c.ViewVendors = pagedVendors.ViewVendors(c.ViewUser.Language)
	c.ItemCategories = CacheGetCategories(c.SelectedPackageType, c.ShippingTo, c.ShippingFrom, "", c.CityID)

	c.ShippingToList = filteredAvailableItems.ShippingToList()
	c.ShippingFromList = filteredAvailableItems.ShippingFromList()
	c.GeoCities = menuFilteredItems.DropCitiesList()

	// paging
	for i := 0; i < c.NumberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}
	c.Page += 1

	// rendering
	util.RenderTemplate(w, "serp/list_vendors", c)
}
