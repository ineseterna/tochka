package main

import (
	"fmt"
	"net/http"

	"github.com/gocraft/web"
	// "github.com/gorilla/csrf"
	"github.com/jasonlvhit/gocron"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/marketplace"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/settings"
)

var (
	APP_SETTINGS = settings.GetSettings()
)

func runCrons() {
	if !APP_SETTINGS.Debug {
		marketplace.StartTransactionsCron()
		marketplace.StartWalletsCron()
		marketplace.StartStatsCron()
		marketplace.StartSERPCron()
	}

	marketplace.StartCurrencyCron()

	<-gocron.Start()

}

func runWebserver() {
	// Root router
	rootRouter := web.New(marketplace.Context{})
	rootRouter.Middleware(web.StaticMiddleware("public"))
	// Marketplace router
	marketplace.ConfigureRouter(rootRouter.Subrouter(marketplace.Context{}, "/"))
	// Start the server
	address := fmt.Sprintf("%s:%s", APP_SETTINGS.Host, APP_SETTINGS.Port)
	// csrfProtection := csrf.Protect([]byte(settings.CSRFEncryption), csrf.FieldName("csrf_token"))
	println("Running server on " + address)
	http.ListenAndServe(address, rootRouter)
}

func runServer() {
	go runCrons()
	runWebserver()

}
