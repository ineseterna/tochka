package marketplace

import (
	"net/http"

	"github.com/gocraft/web"
)

func (c *Context) SellerMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	user, _ := FindUserByUsername(r.PathParams["store"])
	if user == nil {
		http.NotFound(w, r.Request)
		return
	}

	viewSeller := user.ViewUser(c.ViewUser.Language)
	c.ViewSeller = &viewSeller

	c.CanEdit = (c.ViewUser.Uuid == c.ViewSeller.Uuid) || c.ViewUser.IsAdmin || c.ViewUser.IsStaff
	if !c.CanEdit {
		http.NotFound(w, r.Request)
		return
	}

	next(w, r)
}

func (c *Context) SellerItemMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	c.CanEdit = c.ViewUser.IsAdmin || c.ViewUser.AllowedToSell

	if r.PathParams["item"] != "" && r.PathParams["item"] != "new" {
		item, _ := FindItemByUuid(r.PathParams["item"])
		if item == nil || (item.User.Uuid != c.ViewUser.Uuid && !c.CanEdit) {
			http.NotFound(w, r.Request)
			return
		}
		c.Item = *item
		viewItem := c.Item.ViewItem(c.ViewUser.Language)
		c.ViewItem = &viewItem
	}

	store, _ := FindUserByUsername(r.PathParams["store"])
	if store.Username != c.ViewUser.Username && !c.CanEdit {
		http.NotFound(w, r.Request)
		return
	}

	next(w, r)
}

func (c *Context) SellerItemPackageMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if r.PathParams["package"] != "" && r.PathParams["package"] != "new" {
		itemPackage, _ := FindPackageByUuid(r.PathParams["package"])
		if itemPackage != nil {
			if itemPackage.ItemUuid != c.Item.Uuid {
				http.NotFound(w, r.Request)
				return
			}
			c.Package = *itemPackage
			c.ViewPackage = itemPackage.ViewPackage()
		}
	}
	next(w, r)
}

func (c *Context) VendorMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	seller, _ := FindUserByUsername(r.PathParams["store"])
	if seller == nil || seller.Banned {
		http.NotFound(w, r.Request)
		return
	}

	warnings := FindActiveWarningsForUser(seller.Uuid)
	for _, w := range warnings {
		if !w.HasExpired() {
			seller.UserWarnings = append(seller.UserWarnings, w)
		}
	}

	seller.Items = FindItemsForSeller(seller.Uuid)
	c.ViewUserWarnings = warnings.ViewUserWarnings(c.ViewUser.Language)
	reviews, _ := FindRatingReviewsForVendor(seller.Uuid)
	seller.RatingReviews = reviews
	viewSeller := seller.ViewUser(c.ViewUser.Language)
	c.ViewSeller = &viewSeller
	c.ViewItems = Items(seller.Items).ViewItems(c.ViewUser.Language)
	c.CanEdit = (c.ViewUser.Uuid == c.ViewSeller.Uuid) || c.ViewUser.IsAdmin || c.ViewUser.IsStaff

	next(w, r)
}

func (c *Context) VendorItemMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	item, _ := FindItemByUuid(r.PathParams["item"])
	if item == nil || item.UserUuid != c.ViewSeller.Uuid {
		http.NotFound(w, r.Request)
		return
	}
	c.Item = *item
	viewItem := item.ViewItem(c.ViewUser.Language)
	c.ViewItem = &viewItem
	next(w, r)
}
