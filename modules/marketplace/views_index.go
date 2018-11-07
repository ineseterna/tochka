package marketplace

import (
	"github.com/gocraft/web"
	"net/http"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) Index(w web.ResponseWriter, r *web.Request) {
	if c.ViewUser.Uuid == "" {
		util.RenderTemplate(w, "index", c)
	} else {
		redirectUrl := "/marketplace"
		http.Redirect(w, r.Request, redirectUrl, 302)
	}
}
