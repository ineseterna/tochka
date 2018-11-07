package marketplace

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ViewAPILogin(user User, w web.ResponseWriter, r *web.Request) {
	var err error
	c.APISession, err = CreateAPISession(user)
	if err != nil {
		c.Error = err.Error()
		util.APIResponse(w, r, c)
		return
	}

	now := time.Now()
	user.LastLoginDate = &now
	user.Save()

	EventUserLoggedIn(user)
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPILoginRegisterGET(w web.ResponseWriter, r *web.Request) {
	if c.ViewUser.Uuid != "" {
		http.NotFound(w, r.Request)
		return
	}
	c.CaptchaId = captcha.New()
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPILoginPOST(w web.ResponseWriter, r *web.Request) {
	if r.FormValue("decryptedmessage") == "" {
		var (
			isCaptchaValid    = captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
			user, _           = FindUserByUsername(r.FormValue("username"))
			isLoginSuccessful = isCaptchaValid && (user != nil) && user.CheckPassphrase(r.FormValue("passphrase"))
		)
		if !isCaptchaValid {
			c.Error = "Invalid captcha"
			c.ViewAPILoginRegisterGET(w, r)
			return
		}
		if user == nil || !isLoginSuccessful {
			c.Error = "Failed to authenticate"
			c.ViewAPILoginRegisterGET(w, r)
			return
		}
		if user.TwoFactorAuthentication {
			session, _ := CreateAPISession(*user)
			c.APISession = session
			c.APISession.SecondFactorSecretText = util.GenerateUuid()
			c.APISession.Save()

			c.SecretText, _ = util.EncryptText(c.APISession.SecondFactorSecretText, user.Pgp)
			util.APIResponse(w, r, c)
		} else {
			c.ViewAPILogin(*user, w, r)
		}
	} else {
		var (
			secretText       = c.APISession.SecondFactorSecretText
			decryptedmessage = strings.Trim(r.FormValue("decryptedmessage"), "\n ")
		)
		if decryptedmessage == secretText {
			c.ViewAPILogin(c.APISession.User, w, r)
			return
		} else {
			c.Error = "Could not authenticate"
			c.ViewAPILoginRegisterGET(w, r)
			return
		}
	}
}

func (c *Context) ViewAPIRegisterPOST(w web.ResponseWriter, r *web.Request) {
	if c.ViewUser.Uuid != "" {
		http.NotFound(w, r.Request)
		return
	}
	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ViewAPILoginRegisterGET(w, r)
		return
	}
	if r.FormValue("passphrase") != r.FormValue("passphrase_2") {
		c.Error = "Passphrases do not match"
		c.ViewAPILoginRegisterGET(w, r)
		return
	}

	user, err := CreateUser(r.FormValue("username"), r.FormValue("passphrase"))
	if err != nil {
		c.Error = err.Error()
		c.ViewAPILoginRegisterGET(w, r)
		return
	}

	c.APISession, err = CreateAPISession(*user)
	if err != nil {
		c.Error = err.Error()
		c.ViewAPILoginRegisterGET(w, r)
		return
	}

	EventUserRegistred(*user)
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPISERP(w web.ResponseWriter, r *web.Request) {
	c.listAvailableItems(w, r)
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPIShowItem(w web.ResponseWriter, r *web.Request) {
	c.viewShowItem(w, r)
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPIShowStore(w web.ResponseWriter, r *web.Request) {
	c.viewShowStore(w, r)
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPIBookPackage(w web.ResponseWriter, r *web.Request) {
	transactionCount := 0
	for _, t := range FindTransactionsForBuyer(c.ViewUser.Uuid) {
		if t.CurrentPaymentStatus() == "PENDING" {
			transactionCount += 1
		}
	}
	if transactionCount > 10 {
		c.Error = "You have more than 10 active reservations"
		util.APIResponse(w, r, c)
		return
	}
	if c.ViewItem.User.Uuid == c.ViewUser.Uuid || c.ViewUser.IsSeller {
		c.Error = "You can't purchase your own items"
		util.APIResponse(w, r, c)
		return
	}
	shippingId, _ := strconv.ParseInt(r.FormValue("shipping_id"), 10, 64)
	shippingOption, _ := FindShippingOptionById(uint(shippingId))
	quantity, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		quantity = int64(1)
	}
	groups := c.ViewItem.PackagesWithoutReservation().GroupsTable()
	itemPackage, err := groups.GetPackageByHash(r.PathParams["hash"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	transaction, err := CreateTransactionForPackage(
		*itemPackage,
		*c.ViewUser.User,
		r.FormValue("type"),
		int(quantity),
		shippingOption,
		r.FormValue("shipping_address"),
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	transaction.FundFromUserWallets(*c.ViewUser.User)
	viewTransaction := transaction.ViewTransaction()
	c.ViewTransaction = &viewTransaction
	util.APIResponse(w, r, c)
	return
}

func (c *Context) ViewAPIListCurrentTransactionStatuses(w web.ResponseWriter, r *web.Request) {
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
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPIShowTransactionGET(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	review, _ := FindRatingReviewByTransactionUuid(c.Transaction.Uuid)
	if review != nil {
		c.RatingReview = *review
	}
	util.APIResponse(w, r, c)
}

func (c *Context) ViewAPIShowTransactionPOST(w web.ResponseWriter, r *web.Request) {
	message, err := CreateMessage(r.FormValue("text"), c.Thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewAPIShowTransactionGET(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewAPIShowTransactionGET(w, r)
		return
	}

	c.TransactionMiddleware(w, r, c.ViewAPIShowTransactionGET)
}

func (c *Context) ViewAPIWallet(w web.ResponseWriter, r *web.Request) {
	util.APIResponse(w, r, c)
}
