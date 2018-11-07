package marketplace

import (
	"github.com/gocraft/web"
)

func (c *Context) GlobalsMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	lang := c.ViewUser.Language
	if lang == "" {
		session := sessionStore.Load(r.Request)
		if len(r.URL.Query()["language"]) > 0 {
			lang = r.URL.Query()["language"][0]
			session.PutString(w, "language", lang)
		} else {
			slang, err := session.GetString("language")
			if err != nil {
				lang = slang
			}
		}
		c.Language = lang
		c.ViewUser.Language = lang
	}

	c.Localization = GetLocalization(lang)
	c.SiteName = MARKETPLACE_SETTINGS.SiteName
	c.SiteURL = MARKETPLACE_SETTINGS.SiteURL
	next(w, r)
}
