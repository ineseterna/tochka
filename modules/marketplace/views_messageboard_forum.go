package marketplace

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ViewShowMessageboardImage(w web.ResponseWriter, r *web.Request) {
	size := "normal"
	if len(r.URL.Query()["size"]) > 0 {
		size = r.URL.Query()["size"][0]
	}
	util.ServeImage(r.PathParams["uuid"], size, w, r)
}

func (c *Context) ViewEditMessageboardThreadGET(w web.ResponseWriter, r *web.Request) {
	var editThread bool
	if r.PathParams["uuid"] != "" {
		thread, err := GetMessageboardThread(r.PathParams["uuid"])
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		viewThread := thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
		c.ViewThread = &viewThread
		editThread = true
	}
	c.CaptchaId = captcha.New()
	c.MessageboardSections = FindAllMessageboardSections()

	if editThread {
		util.RenderTemplate(w, "board/thread_edit", c)
	} else {
		util.RenderTemplate(w, "board/thread_new", c)
	}
}

func (c *Context) ViewEditMessageboardThreadPOST(w web.ResponseWriter, r *web.Request) {
	// vars
	var (
		thread      *Thread
		isNewThread bool
		err         error
	)

	// captcha
	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha")) || c.ViewUser.IsAdmin || c.ViewUser.IsStaff
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ViewEditMessageboardThreadGET(w, r)
		return
	}

	// new or existing thread
	if r.PathParams["uuid"] != "" {
		thread, err = GetMessageboardThread(r.PathParams["uuid"])
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
	} else {
		thread, err = CreateThread(
			"messageboard",
			"",
			r.FormValue("title"),
			r.FormValue("text"),
			c.ViewUser.User,
			nil,
			true,
		)
		if err != nil {
			c.Error = err.Error()
			c.ViewEditMessageboardThreadGET(w, r)
			return
		}
		isNewThread = true
	}

	// section
	secId, err := strconv.ParseInt(r.FormValue("section_id"), 10, 64)
	if err != nil {
		c.Error = err.Error()
		c.ViewEditMessageboardThreadGET(w, r)
		return
	}
	section, err := FindMessageboardSectionByID(int(secId))
	if err != nil {
		c.Error = err.Error()
		c.ViewEditMessageboardThreadGET(w, r)
		return
	}

	// set title, text and section
	thread.Title = r.FormValue("title")
	thread.Text = r.FormValue("text")
	thread.MessageboardSectionID = section.ID
	thread.Save()

	viewThread := thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	c.ViewThread = &viewThread
	err = thread.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewEditMessageboardThreadGET(w, r)
		return
	}

	// feed actions
	if isNewThread {
		CreateFeedItem(c.ViewUser.Uuid, "new_thread", "created new thread", thread.Uuid)
	}

	// redirect
	http.Redirect(w, r.Request, fmt.Sprintf("/board/?section=%d", thread.MessageboardSectionID), 302)
}

func (c *Context) ViewReplyToMessageboardThread(w web.ResponseWriter, r *web.Request) {
	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha")) || c.ViewUser.IsAdmin || c.ViewUser.IsAdmin
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ShowThread(w, r)
		return
	}
	thread, err := GetMessageboardThread(r.FormValue("thread_uuid"))
	if err != nil {
		c.Error = err.Error()
		c.ShowThread(w, r)
		return
	}
	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ShowThread(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ShowThread(w, r)
		return
	}

	CreateFeedItem(c.ViewUser.Uuid, "new_thread_reply", "replied in thread", message.Uuid)
	c.ShowThread(w, r)
}

func (c *Context) ViewListMessageboardSections(w web.ResponseWriter, r *web.Request) {
	c.ViewMessageboardSections = CacheGetAllMessageboardSections().AsNestedViewMessageboardSections(c.ViewUser.Language)
	c.MessageboardDailyStats = CacheGetMessageboardDailyStats()
	c.MessageboardSummaryStats = CacheGetMessageboardSummaryStats()
	c.ViewMessageboardThreads = CacheFindNewMessageboardThreadsTop5().ViewMessageboardThreads(c.ViewUser.Language)
	util.RenderTemplate(w, "board/board_sections", c)
}

func (c *Context) ViewListThreads(w web.ResponseWriter, r *web.Request) {

	// Param Parsing
	sectionID, _ := strconv.ParseInt(r.PathParams["section_id"], 10, 64)
	c.SelectedSectionID = int(sectionID)
	if len(r.URL.Query()["page"]) > 0 {
		strPage := r.URL.Query()["page"][0]
		page, err := strconv.ParseInt(strPage, 10, 32)
		if err != nil || page < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.Page = int(page) - 1
	}

	numberOfThreadsPerPage := 50
	messagebordThreads := FindMessageboardThreadsForUserUuid(c.SelectedSectionID, c.Page, numberOfThreadsPerPage, c.ViewUser.Uuid)
	c.ViewMessageboardThreads = MessageboardThreads(messagebordThreads).ViewMessageboardThreads(c.ViewUser.Language)

	// Section
	section, _ := FindMessageboardSectionByID(c.SelectedSectionID)
	c.ViewMessageboardSection = section.ViewMessageboardSection(c.ViewUser.Language)

	// Paging
	numberOfThreads := float64(CacheCountMessageboardThreads(c.SelectedSectionID))
	c.NumberOfPages = int(math.Ceil(numberOfThreads / float64(numberOfThreadsPerPage)))
	for i := 0; i < c.NumberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}
	c.Page += 1

	util.RenderTemplate(w, "board/threads", c)
}

func (c *Context) ShowThread(w web.ResponseWriter, r *web.Request) {

	if c.ViewUser.Uuid == "" {
		redirectUrl := "/auth/register"
		http.Redirect(w, r.Request, redirectUrl, 302)
		return
	}

	thread, err := GetMessageboardThread(r.PathParams["uuid"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	viewThread := thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	c.ViewThread = &viewThread
	c.NumberOfPages = int(math.Ceil(float64(len(c.ViewThread.Messages)) / 50.0))

	if len(r.URL.Query()["page"]) > 0 {
		strPage := r.URL.Query()["page"][0]
		page, err := strconv.ParseInt(strPage, 10, 32)
		if err != nil || page < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.Page = int(page) - 1
	}
	// paging
	for i := 0; i < c.NumberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}

	c.ViewThread.Messages = c.ViewThread.Messages[c.Page*50 : int(math.Min(float64(len(c.ViewThread.Messages)), float64(c.Page*50+50)))]
	c.Page = c.Page + 1
	c.CaptchaId = captcha.New()

	// c.ViewThreads = FindMessageboardThreads(c.SelectedSectionID).ViewThreads(c.ViewUser.Language, c.ViewUser.User)
	c.SelectedSection = c.ViewThread.Section

	UpdateThreadPerusalStatus(thread.Uuid, c.ViewUser.Uuid)

	util.RenderTemplate(w, "board/thread", c)
}
