package marketplace

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"
	"github.com/wcharczuk/go-chart"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Staff List
*/

func (c *Context) ViewStaff(w web.ResponseWriter, r *web.Request) {
	http.Redirect(w, r.Request, "/staff/users", 302)
}

/*
	Staff
*/

func (c *Context) ViewStaffListItems(w web.ResponseWriter, r *web.Request) {

	var (
		err      error
		page     int = 1
		pageSize int = 200
	)

	if len(r.URL.Query()["page"]) > 0 {
		selectedPageStr := r.URL.Query()["page"][0]
		page, err = strconv.Atoi(selectedPageStr)
	}

	if err == nil {
		c.SelectedPage = page
	}

	// numberOfItems := CountUnreviewedItems()
	// numberOfPages := int(math.Ceil(float64(numberOfItems) / float64(pageSize)))
	// for i := 0; i < numberOfPages; i++ {
	// 	c.Pages = append(c.Pages, i+1)
	// }

	unreviewdItems := FindUnreviewedItems(page-1, pageSize)
	c.ViewItems = unreviewdItems.ViewItems(c.ViewUser.Language)

	util.RenderTemplate(w, "staff/items", c)
}

/*
	Disputes
*/

func (c *Context) ViewStaffListDisputes(w web.ResponseWriter, r *web.Request) {
	// transaction type
	if len(r.URL.Query()["status"]) > 0 {
		c.SelectedStatus = r.URL.Query()["status"][0]
	}
	// pages
	pageSize := 20
	if len(r.URL.Query()["page"]) > 0 {
		strPage := r.URL.Query()["page"][0]
		page, err := strconv.ParseInt(strPage, 10, 32)
		if err != nil || page < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.Page = int(page)
	} else {
		c.Page = 1
	}

	c.SelectedStatus = ""
	if len(r.URL.Query()["status"]) > 0 {
		c.SelectedStatus = r.URL.Query()["status"][0]
	}

	c.NumberOfTransactions = CountDisputedTransactions(c.ViewUser.Uuid, c.SelectedStatus)
	c.NumberOfPages = int(math.Ceil(float64(c.NumberOfTransactions) / float64(pageSize)))
	for i := 0; i < c.NumberOfPages; i++ { // paging
		c.Pages = append(c.Pages, i+1)
	}

	transactions := Transactions(GetDisputedTransactionsPaged(pageSize, c.Page-1, c.ViewUser.Uuid, c.SelectedStatus))
	c.ViewTransactions = transactions.ViewTransactions()
	util.RenderTemplate(w, "staff/disputes", c)
}

/*
	Staff
*/

func (c *Context) ViewStaffListStaff(w web.ResponseWriter, r *web.Request) {
	users, err := FindStaffMemebers()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewExtendedUsers = ExtendedUsers(users).ViewExtendedUsers(c.Language)
	util.RenderTemplate(w, "staff/staff", c)
}

/*
	Vendors
*/

func (c *Context) ViewStaffListVendors(w web.ResponseWriter, r *web.Request) {
	users, err := FindValidVerificationRequests()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewVendors = users.ViewVendors(c.ViewUser.Language)
	util.RenderTemplate(w, "staff/vendors", c)
}

func (c *Context) ViewStaffVendorVerificationThreadGET(w web.ResponseWriter, r *web.Request) {
	user, err := c.getUserForTrustPage(r)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	viewSeller := user.ViewUser(c.ViewUser.Language)
	c.ViewSeller = &viewSeller

	thread, err := GetVendorVerificationThread(user, true)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	c.CaptchaId = captcha.New()
	viewThread := thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	c.ViewThread = &viewThread
	util.RenderTemplate(w, "staff/vendors_verification_thread", c)
}

func (c *Context) ViewStaffVendorVerificationThreadPOST(w web.ResponseWriter, r *web.Request) {

	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha")) || c.ViewUser.IsStaff || c.ViewUser.IsAdmin
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ViewVerificationThreadGET(w, r)
		return
	}

	user, err := c.getUserForTrustPage(r)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	thread, err := GetVendorVerificationThread(user, false)
	if err != nil {
		c.Error = err.Error()
		c.ViewVerificationThreadGET(w, r)
		return
	}
	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewVerificationThreadGET(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewVerificationThreadGET(w, r)
		return
	}

	EventNewTrustedVendorThreadPost(*c.ViewUser.User, user, *message)
	c.ViewStaffVendorVerificationThreadGET(w, r)
}

/*
	Warnings
*/

func (c *Context) ViewStaffListWarnings(w web.ResponseWriter, r *web.Request) {
	warnings := FindAllActiveWarnings()
	c.ViewUserWarnings = warnings.ViewUserWarnings(c.ViewUser.Language)
	util.RenderTemplate(w, "staff/warnings", c)
}

/*
	Statistics
*/

func getStats() []StatsItem {
	dt, _ := time.Parse(time.RFC3339, "2016-12-26T00:00:00+00:00")
	return CacheGetMarketplaceStats(dt)
}

func (c *Context) ViewStaffStats(w web.ResponseWriter, r *web.Request) {
	// User Statistics
	c.NumberOfUsers = CountUsers(nil)
	// Vendor Statistics
	c.NumberOfVendors = CountVendors(nil)
	c.NumberOfFreeVendors = CountFreeVendors(nil)
	c.NumberOfGoldVendors = CountGoldVendors(nil)
	c.NumberOfSilverVendors = CountSilverVendors(nil)
	c.NumberOfBronzeVendors = CountBronzeVendors(nil)
	// Item Statistrics
	c.NumberOfItems = CountItems()
	c.StatsItems = getStats()

	util.RenderTemplate(w, "staff/stats", c)
}

func (c *Context) ViewStatsNumberOfUsersGraph(w web.ResponseWriter, r *web.Request) {

	statsItems := getStats()

	xValues := []time.Time{}
	yValues := []float64{}

	for _, si := range statsItems {
		yValues = append(yValues, float64(si.NumberOfUsers))
		xValues = append(xValues, si.Date)
	}

	continuosSeries := chart.TimeSeries{
		XValues: xValues,
		YValues: yValues,
	}

	graph := chart.Chart{
		Series: []chart.Series{continuosSeries},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}

func (c *Context) ViewStatsNumberOfVendorsGraph(w web.ResponseWriter, r *web.Request) {

	statsItems := getStats()

	xValues := []time.Time{}
	yValues := []float64{}

	for _, si := range statsItems {
		yValues = append(yValues, float64(si.NumberOfVendors))
		xValues = append(xValues, si.Date)
	}

	continuosSeries := chart.TimeSeries{
		XValues: xValues,
		YValues: yValues,
	}

	graph := chart.Chart{
		Series: []chart.Series{continuosSeries},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}

func (c *Context) ViewStatsBTCTradeAmountGraph(w web.ResponseWriter, r *web.Request) {

	statsItems := getStats()

	xValues := []time.Time{}
	yValues := []float64{}

	for _, si := range statsItems {
		yValues = append(yValues, float64(si.BTCTradeAmount))
		xValues = append(xValues, si.Date)
	}

	continuosSeries := chart.TimeSeries{
		XValues: xValues,
		YValues: yValues,
	}

	graph := chart.Chart{
		Series: []chart.Series{continuosSeries},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}

func (c *Context) ViewStatsETHTradeAmountGraph(w web.ResponseWriter, r *web.Request) {

	statsItems := getStats()

	xValues := []time.Time{}
	yValues := []float64{}

	for _, si := range statsItems {
		yValues = append(yValues, float64(si.ETHTradeAmount))
		xValues = append(xValues, si.Date)
	}

	continuosSeries := chart.TimeSeries{
		XValues: xValues,
		YValues: yValues,
	}

	graph := chart.Chart{
		Series: []chart.Series{continuosSeries},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}

func (c *Context) ViewStatsBCHTradeAmountGraph(w web.ResponseWriter, r *web.Request) {

	statsItems := getStats()

	xValues := []time.Time{}
	yValues := []float64{}

	for _, si := range statsItems {
		yValues = append(yValues, float64(si.BCHTradeAmount))
		xValues = append(xValues, si.Date)
	}

	continuosSeries := chart.TimeSeries{
		XValues: xValues,
		YValues: yValues,
	}

	graph := chart.Chart{
		Series: []chart.Series{continuosSeries},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}

func (c *Context) ViewStatsNumberOfTransactionsGraph(w web.ResponseWriter, r *web.Request) {

	statsItems := getStats()

	xValues := []time.Time{}
	yValues := []float64{}

	for _, si := range statsItems {
		yValues = append(yValues,
			float64(si.NumberOfBTCTransactionsCreated)+
				float64(si.NumberOfBCHTransactionsCreated)+
				float64(si.NumberOfETHTransactionsCreated))
		xValues = append(xValues, si.Date)
	}

	continuosSeries := chart.TimeSeries{
		XValues: xValues,
		YValues: yValues,
	}

	graph := chart.Chart{
		Series: []chart.Series{continuosSeries},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}

/*
	News
*/

func (c *Context) ViewStaffEditNewsGET(w web.ResponseWriter, r *web.Request) {
	thread, err := GetNewsThread(c.ViewUser.Language)
	if err != nil {
		panic(err)
		http.NotFound(w, r.Request)
		return
	}
	c.SelectedSection = "news"
	c.CaptchaId = captcha.New()
	viewThread := thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	c.ViewThread = &viewThread
	util.RenderTemplate(w, "staff/news", c)
}

func (c *Context) ViewStaffEditNewsPOST(w web.ResponseWriter, r *web.Request) {

	thread, err := GetNewsThread(c.ViewUser.Language)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha")) || c.ViewUser.IsAdmin || c.ViewUser.IsAdmin
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ViewStaffEditNewsGET(w, r)
		return
	}

	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewStaffEditNewsGET(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewStaffEditNewsGET(w, r)
		return
	}

	c.ViewStaffEditNewsGET(w, r)
}

/*
	Categories
*/

func (c *Context) ViewStaffCategories(w web.ResponseWriter, r *web.Request) {
	c.ItemCategories = FindAllCategories()
	util.RenderTemplate(w, "staff/categories_list", c)
}

func (c *Context) ViewStaffCategoriesEdit(w web.ResponseWriter, r *web.Request) {

	if r.PathParams["id"] != "new" {
		catId, err := strconv.ParseInt(r.PathParams["id"], 10, 64)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		category, err := FindCategoryByID(int(catId))
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.ItemCategory = *category
	}

	categories := FindAllCategories()

	cat := Category{
		Name: "",
		ID:   "0",
	}
	c.Categories = append(c.Categories, cat)

	translateCat := func(ic ItemCategory, lvl int) Category {
		cat := Category{
			ID: fmt.Sprintf("%d", ic.ID),
		}

		switch c.ViewUser.Language {
		case "ru":
			cat.Name = ic.NameRu
		case "de":
			cat.Name = ic.NameDe
		case "es":
			cat.Name = ic.NameEs
		case "fr":
			cat.Name = ic.NameFr
		case "rs":
			cat.Name = ic.NameRs
		case "tr":
			cat.Name = ic.NameTr
		default:
			cat.Name = ic.NameEn
		}

		for i := 0; i < lvl-1; i++ {
			cat.Name = "-" + cat.Name
		}

		return cat
	}

	for _, cat1 := range categories {
		c.Categories = append(c.Categories, translateCat(cat1, 1))
		for _, cat2 := range cat1.Subcategories {
			c.Categories = append(c.Categories, translateCat(cat2, 2))
			for _, cat3 := range cat2.Subcategories {
				c.Categories = append(c.Categories, translateCat(cat3, 3))
			}
		}
	}

	c.Category = fmt.Sprintf("%d", c.ItemCategory.ParentID)

	util.RenderTemplate(w, "staff/categories_edit", c)
}

func (c *Context) ViewStaffCategoriesEditPOST(w web.ResponseWriter, r *web.Request) {

	var category ItemCategory

	if r.PathParams["id"] != "new" {
		catId, err := strconv.ParseInt(r.PathParams["id"], 10, 64)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		cat, err := FindCategoryByID(int(catId))
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		category = *cat
	}

	err := r.ParseForm()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	if r.FormValue("parent_id") != "" {
		parId, err := strconv.ParseUint(r.FormValue("parent_id"), 10, 64)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		if parId != 0 {
			parentCat, err := FindCategoryByID(int(parId))
			if err != nil {
				http.NotFound(w, r.Request)
				return
			}
			category.ParentID = parentCat.ID
		} else {
			category.ParentID = 0
		}
	}

	category.Icon = r.FormValue("icon")
	category.NameEn = r.FormValue("name_en")
	category.NameRu = r.FormValue("name_ru")
	category.NameDe = r.FormValue("name_de")
	category.NameEs = r.FormValue("name_es")
	category.NameFr = r.FormValue("name_fr")
	category.NameRs = r.FormValue("name_rs")
	category.NameTr = r.FormValue("name_tr")
	category.Save()

	http.Redirect(w, r.Request, fmt.Sprintf("/staff/item_categories/"), 302)
}

/*
	Advertisings
*/

func (c *Context) ViewStaffAdvertisings(w web.ResponseWriter, r *web.Request) {
	var err error
	c.Advertisings, err = FindAllAdvertising()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	c.AdvertisingCost = MARKETPLACE_SETTINGS.AdvertisingCost
	util.RenderTemplate(w, "staff/advertising", c)
}

func (c *Context) ViewStaffCategoriesDelete(w web.ResponseWriter, r *web.Request) {
	catId, err := strconv.ParseUint(r.PathParams["id"], 10, 64)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	category, err := FindCategoryByID(int(catId))
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	err = category.Remove()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	http.Redirect(w, r.Request, fmt.Sprintf("/staff/item_categories/"), 302)
}

func (c *Context) ViewStaffListReportedMessages(w web.ResponseWriter, r *web.Request) {
	c.ViewMessages = FindReportedMessages().ViewMessages(c.ViewUser.Language)
	util.RenderTemplate(w, "staff/reported_messages", c)
}

func (c *Context) ViewStaffShowReportedMessage(w web.ResponseWriter, r *web.Request) {
	message, err := FindMessageByUuid(r.PathParams["message_uuid"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	if !message.IsReported {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
	util.RenderTemplate(w, "staff/reported_message", c)
}
