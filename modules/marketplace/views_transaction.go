package marketplace

import (
	"math"
	"net/http"
	"strconv"

	btcqr "github.com/GeertJohan/go.btcqr"
	"github.com/dchest/captcha"
	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ShowTransaction(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	review, _ := FindRatingReviewByTransactionUuid(c.Transaction.Uuid)
	if review != nil {
		c.RatingReview = *review
	}
	if len(r.URL.Query()["section"]) > 0 {
		section := r.URL.Query()["section"][0]
		c.SelectedSection = section
	} else {
		c.SelectedSection = "payment"
	}
	util.RenderTemplate(w, "transaction/show", c)
}

func (c *Context) ShowTransactionPOST(w web.ResponseWriter, r *web.Request) {
	message, err := CreateMessage(r.FormValue("text"), c.Thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ShowTransaction(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ShowTransaction(w, r)
		return
	}

	c.TransactionMiddleware(w, r, c.ShowTransaction)
}

func (c *Context) UpdateTransaction(w web.ResponseWriter, r *web.Request) {
	transaction, _ := FindTransactionByUuid(r.PathParams["transaction"])
	if transaction == nil {
		http.NotFound(w, r.Request)
		return
	}
	transaction.UpdateTransactionStatus()
	viewTransaction := transaction.ViewTransaction()
	c.ViewTransaction = &viewTransaction
	c.ShowTransaction(w, r)
}

func (c *Context) ListCurrentTransactionStatuses(w web.ResponseWriter, r *web.Request) {
	pageSize := 20
	if len(r.URL.Query()["status"]) > 0 {
		c.SelectedStatus = r.URL.Query()["status"][0]
	}
	c.NumberOfTransactions = CountCurrentTransactionStatuses(c.ViewUser.Uuid, c.SelectedStatus, false)
	c.NumberOfPages = int(math.Ceil(float64(c.NumberOfTransactions) / float64(pageSize)))
	for i := 0; i < c.NumberOfPages; i++ { // paging
		c.Pages = append(c.Pages, i+1)
	}

	c.ViewCurrentTransactionStatuses = FindCurrentTransactionStatuses(
		c.ViewUser.Uuid, c.SelectedStatus, false, c.Page, pageSize).
		ViewCurrentTransactionStatuses(c.ViewUser.Language)

	util.RenderTemplate(w, "transaction/list", c)
}

func (c *Context) TransactionImage(w web.ResponseWriter, r *web.Request) {
	btcqr.DefaultConfig.Scheme = "guldencoin"

	var amount float64
	if c.Transaction.Type == "bitcoin" {
		amount = c.Transaction.BitcoinTransaction.Amount
	}
	if c.Transaction.Type == "bitcoin_cash" {
		amount = c.Transaction.BitcoinCashTransaction.Amount
	}

	req := &btcqr.Request{
		Address: c.Transaction.Uuid,
		Amount:  amount,
		Label:   c.Transaction.Description,
		Message: "",
	}
	code, err := req.GenerateQR()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	png := code.PNG()
	w.Header().Set("Content-type", "image/png")
	w.Write(png)
}

func (c *Context) CompleteTransactionPOST(w web.ResponseWriter, r *web.Request) {

	// Review
	review, _ := FindRatingReviewByTransactionUuid(c.Transaction.Uuid)

	// quality
	itemQuality, err := strconv.Atoi(r.FormValue("item_quality"))
	if err != nil || itemQuality < 0 || itemQuality > 5 {
		c.Error = "Wrong input for item quality"
		http.NotFound(w, r.Request)
		return
	}
	itemReview := r.FormValue("item_review")
	if len(itemReview) > 255 {
		itemReview = itemReview[0:255]
	}
	// package
	marketplaceQuality, err := strconv.Atoi(r.FormValue("marketplace_quality"))
	if err != nil || marketplaceQuality < 0 || marketplaceQuality > 5 {
		c.Error = "Wrong input for marketplace quality"
		http.NotFound(w, r.Request)
		return
	}
	marketplaceReview := r.FormValue("marketplace_review")
	if len(marketplaceReview) > 255 {
		marketplaceReview = marketplaceReview[0:255]
	}
	// seller
	sellerQuality, err := strconv.Atoi(r.FormValue("seller_quality"))
	if err != nil || sellerQuality < 0 || sellerQuality > 5 {
		c.Error = "Wrong input for seller quality"
		http.NotFound(w, r.Request)
		return
	}
	sellerReview := r.FormValue("seller_review")
	if len(sellerReview) > 255 {
		sellerReview = sellerReview[0:255]
	}

	if review == nil {
		review = &RatingReview{
			Uuid: util.GenerateUuid(),
		}

		CreateFeedItem(c.Transaction.BuyerUuid, "new_review", "added new review", review.Uuid)
	}

	pkg, _ := FindPackageByUuid(c.Transaction.PackageUuid)
	if pkg != nil {
		review.ItemUuid = pkg.ItemUuid
	}

	review.ItemReview = itemReview
	review.ItemScore = itemQuality
	review.MarketplaceReview = marketplaceReview
	review.MarketplaceScore = marketplaceQuality
	review.SellerReview = sellerReview
	review.SellerScore = sellerQuality
	review.TransactionUuid = c.Transaction.Uuid
	review.SellerUuid = c.Transaction.SellerUuid
	review.UserUuid = c.ViewUser.Uuid

	review.Save()
	if review != nil {
		c.RatingReview = *review
	}

	http.Redirect(w, r.Request, "/payments/"+c.Transaction.Uuid+"?section=review", 302)
}

func (c *Context) SetTransactionShippingStatus(w web.ResponseWriter, r *web.Request) {
	status := r.FormValue("shipping_status")
	if !(status == "DISPATCHED" || status == "SHIPPED") {
		http.NotFound(w, r.Request)
		return
	}
	c.Transaction.SetShippingStatus(status, "Shipping status changed to "+status, c.ViewUser.Uuid)
	redirectUrl := "/payments/" + c.Transaction.Uuid
	http.Redirect(w, r.Request, redirectUrl, 302)
}

func (c *Context) ReleaseTransaction(w web.ResponseWriter, r *web.Request) {
	if c.Transaction.SellerUuid != c.ViewUser.Uuid {
		err := c.Transaction.Release("User released transaction", c.ViewUser.Uuid)
		if err != nil {
			c.Transaction.SetTransactionStatus(
				c.Transaction.CurrentPaymentStatus(),
				c.Transaction.CurrentAmountPaid(),
				"Failed to release transaction",
				c.ViewUser.Uuid,
				nil,
			)
		}
	}
	http.Redirect(w, r.Request, "/payments/"+c.Transaction.Uuid, 302)
}

func (c *Context) CancelTransaction(w web.ResponseWriter, r *web.Request) {
	if c.Transaction.IsCompleted() && !c.Transaction.IsDispatched() && !c.Transaction.IsShipped() {
		err := c.Transaction.Cancel("User cancelled transaction", c.ViewUser.Uuid)
		if err != nil {
			c.Transaction.SetTransactionStatus(
				c.Transaction.CurrentPaymentStatus(),
				c.Transaction.CurrentAmountPaid(),
				"Failed to cancel transaction",
				c.ViewUser.Uuid,
				nil,
			)
		}
	}
	http.Redirect(w, r.Request, "/payments/"+c.Transaction.Uuid, 302)
}
