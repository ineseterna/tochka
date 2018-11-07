package marketplace

import (
	"github.com/dchest/captcha"
	"github.com/gocraft/web"
)

func (c *Context) ViewCaptchaImage(w web.ResponseWriter, r *web.Request) {
	captcha.WriteImage(w, r.PathParams["captcha_id"], 300, 100)
}
