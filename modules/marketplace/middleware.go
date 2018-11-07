package marketplace

import (
	"time"

	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) CurrencyMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.CurrencyRates = map[string]map[string]float64{}
	for _, fc := range FIAT_CURRENCIES {
		c.CurrencyRates[fc] = map[string]float64{}
		for _, cc := range CRYPTO_CURRENCIES {
			c.CurrencyRates[fc][cc] = GetCurrencyRate(cc, fc)
		}
	}
	next(w, r)
}

func (c *Context) LoggerMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	next(rw, req)
	startTime := time.Now()
	username := c.ViewUser.Username

	if username == "" {
		username = "anon"
	} else {
		username = "@" + username
	}

	util.Log.Info(
		"method:'%s' user:'%s' url':%s' agent:'%s' duration:'%d μs' status:'%d'\n",
		req.Method,
		username,
		req.URL.Path,
		req.UserAgent(),
		time.Since(startTime).Nanoseconds()/1000,
		rw.StatusCode(),
	)
}
