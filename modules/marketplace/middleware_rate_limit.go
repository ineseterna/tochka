package marketplace

import (
	"time"

	"github.com/gocraft/web"
)

var userToTockerLimiter map[string](chan time.Time)

func (c *Context) RateLimitMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.ViewUser.Username != "" {
		<-getUserTicker(c.ViewUser.Username)
		next(w, r)
	} else {
		next(w, r)
	}
}

func TickerLimiter(rps, burst int) (c chan time.Time, cancel func()) {
	// create the buffered channel and prefill it
	c = make(chan time.Time, burst)
	for i := 0; i < burst; i++ {
		c <- time.Now()
	}

	// create a ticker with the interval 1/rps
	t := time.NewTicker(time.Second / time.Duration(rps))

	// add to the channel with each tick
	go func() {
		for t := range t.C {
			select {
			case c <- t: // add the tick to channel
			default: // channel already full, drop the tick
			}
		}
		close(c) // close channel when the ticker is stopped
	}()

	return c, t.Stop
}

func getUserTicker(username string) chan time.Time {
	if l, ok := userToTockerLimiter[username]; ok {
		return l
	} else {
		l, _ := TickerLimiter(16, 1)
		userToTockerLimiter[username] = l
		return l
	}
}

func init() {
	userToTockerLimiter = map[string](chan time.Time){}
}
