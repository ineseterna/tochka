package marketplace

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gocraft/web"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ViewStaffListSupportTickets(w web.ResponseWriter, r *web.Request) {
	var (
		err          error
		pageSize     int = 50
		selectedPage int = 0
	)

	c.SelectedSection = "new-open"
	if len(r.URL.Query()["section"]) > 0 {
		c.SelectedSection = r.URL.Query()["section"][0]
	}

	if len(r.URL.Query()["page"]) > 0 {
		selectedPageStr := r.URL.Query()["page"][0]
		page, err := strconv.Atoi(selectedPageStr)
		if err == nil {
			selectedPage = page - 1
		}
	}

	numberOfPages := int(math.Ceil(float64(CountSupportTicketsByStatus(c.SelectedSection)) / float64(pageSize)))
	for i := 0; i < numberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}

	tickets, err := FindSupportTicketsByStatus(c.SelectedSection, c.SelectedPage, pageSize)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	c.ViewSupportTickets = tickets.ViewSupportTickets(c.ViewUser.Language)
	c.SelectedPage = selectedPage + 1
	util.RenderTemplate(w, "staff/users_support_tickets", c)
}

func (c *Context) ViewStaffUserFinance(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	viewSeller := user.ViewUser(c.ViewUser.Language)
	c.ViewSeller = &viewSeller

	c.ViewSeller.BitcoinWallets = c.ViewSeller.FindUserBitcoinWallets()
	c.ViewSeller.EthereumWallets = c.ViewSeller.FindUserEthereumWallets()
	c.ViewSeller.BitcoinCashWallets = c.ViewSeller.FindUserBitcoinCashWallets()

	for _, w := range c.ViewSeller.BitcoinWallets {
		w.UpdateBalance(false)
	}

	for _, w := range c.ViewSeller.EthereumWallets {
		w.UpdateBalance(false)
	}

	for _, w := range c.ViewSeller.BitcoinCashWallets {
		w.UpdateBalance(false)
	}

	c.ViewSeller.BitcoinBalance = c.ViewSeller.BitcoinWallets.Balance()
	c.ViewSeller.EthereumBalance.Balance = c.ViewSeller.EthereumWallets.Balance().Balance
	c.ViewSeller.BitcoinCashBalance = c.ViewSeller.BitcoinCashWallets.Balance()

	c.ViewSeller.BitcoinWallet = c.ViewSeller.BitcoinWallets[0]
	c.ViewSeller.EthereumWallet = c.ViewSeller.EthereumWallets[0]
	c.ViewSeller.BitcoinCashWallet = c.ViewSeller.BitcoinCashWallets[0]

	c.UserSettingsHistory = SettingsChangeHistoryByUser(user.Uuid)

	util.RenderTemplate(w, "staff/users_user_finance", c)
}

func (c *Context) ViewStaffUserTickets(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	viewSeller := user.ViewUser(c.ViewUser.Language)
	c.ViewSeller = &viewSeller

	tickets, err := FindSupportTicketsForUser(*user)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewSupportTickets = tickets.ViewSupportTickets(c.ViewUser.Language)
	util.RenderTemplate(w, "staff/users_user_tickets", c)
}

func (c *Context) ViewStaffUserPayments(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	viewSeller := user.ViewUser(c.ViewUser.Language)
	c.ViewSeller = &viewSeller

	c.ViewCurrentTransactionStatuses = FindCurrentTransactionStatuses(
		user.Uuid, c.SelectedStatus, false, 0, 100).
		ViewCurrentTransactionStatuses(c.ViewUser.Language)

	util.RenderTemplate(w, "staff/users_user_payments", c)
}

func (c *Context) ViewStaffUserAdminActions(w web.ResponseWriter, r *web.Request) {

	user, err := FindUserByUsername(r.PathParams["username"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	viewSeller := user.ViewUser(c.ViewUser.Language)
	c.ViewSeller = &viewSeller
	util.RenderTemplate(w, "staff/users_user_admin_actions", c)
}
