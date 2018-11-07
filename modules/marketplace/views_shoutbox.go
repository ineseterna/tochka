package marketplace

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ViewShoutboxGET(w web.ResponseWriter, r *web.Request) {

	c.SelectedSection = "news"
	if len(r.URL.Query()["section"]) > 0 {
		c.SelectedSection = r.URL.Query()["section"][0]
	}

	switch c.SelectedSection {
	case "shoutbox":
		thread, err := GetShoutboxThread("en")
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		viewThread := thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
		c.ViewThread = &viewThread
		if len(c.ViewThread.Messages) > 30 {
			c.ViewThread.Messages = c.ViewThread.Messages[0:30]
		}
	case "news":
		lang := "en"
		thread, err := GetNewsThread(lang)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		viewThread := thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
		c.ViewThread = &viewThread
	}

	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "shoutbox", c)
}

func (c *Context) ViewShoutboxPOST(w web.ResponseWriter, r *web.Request) {
	thread, err := GetShoutboxThread("en")
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha")) || c.ViewUser.IsAdmin || c.ViewUser.IsAdmin
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ViewShoutboxGET(w, r)
		return
	}

	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewShoutboxGET(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewShoutboxGET(w, r)
		return
	}

	EventNewShoutboxPost(*c.ViewUser.User, *message)
	c.ViewShoutboxGET(w, r)
}
