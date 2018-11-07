package marketplace

import (
	"net/http"

	"github.com/gocraft/web"
)

func (c *Context) TransactionStatsMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.ViewUser.Uuid != "" {
		c.NumberOfTransactions = CountCurrentTransactionStatuses(c.ViewUser.Uuid, "", true)
	}
	next(w, r)
}

func (c *Context) TransactionMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	transaction, _ := FindTransactionByUuid(r.PathParams["transaction"])
	if transaction == nil {
		http.NotFound(w, r.Request)
		return
	}

	if !c.ViewUser.IsAdmin && !c.ViewUser.IsStaff && (transaction.BuyerUuid != c.ViewUser.Uuid && transaction.SellerUuid != c.ViewUser.Uuid) {
		http.NotFound(w, r.Request)
		return
	}

	pckg, _ := FindPackageByUuid(transaction.PackageUuid)
	if pckg != nil {
		c.Package = *pckg
		c.ViewPackage = pckg.ViewPackage()
	}
	c.Transaction = *transaction
	viewTransaction := transaction.ViewTransaction()
	c.ViewTransaction = &viewTransaction

	thread, err := GetTransactionThread(*transaction, "")
	if err == nil {
		c.Thread = *thread
		viewThread := thread.ViewThread(c.ViewUser.User.Language, c.ViewUser.User)
		c.ViewThread = &viewThread
	}

	next(w, r)
}
