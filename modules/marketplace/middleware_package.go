package marketplace

import (
	"net/http"

	"github.com/gocraft/web"
)

func (c *Context) PackageMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	itemPackage, _ := FindPackageByUuid(r.PathParams["package"])
	if itemPackage == nil {
		http.NotFound(w, r.Request)
		return
	}
	transaction := itemPackage.Transaction()
	if transaction == nil {
		http.NotFound(w, r.Request)
		return
	}
	if transaction.Buyer.Uuid != c.ViewUser.Uuid && transaction.Seller.Uuid != c.ViewUser.Uuid && !c.ViewUser.IsAdmin {
		http.NotFound(w, r.Request)
		return
	}
	c.Transaction = *transaction
	c.Package = *itemPackage
	c.ViewPackage = c.Package.ViewPackage()
	viewTransaction := transaction.ViewTransaction()
	c.ViewTransaction = &viewTransaction
	next(w, r)
}
