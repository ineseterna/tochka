package marketplace

import (
	"github.com/gocraft/web"
)

func (c *Context) AdvertisingMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	ads, err := GetAdvertisings(1)
	if err != nil {
		next(w, r)
		return
	}
	c.Advertisings = ads
	next(w, r)
}
