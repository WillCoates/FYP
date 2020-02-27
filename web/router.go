package main

import (
	"net/http"

	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/controllers"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
)

func CreateRouter(logic *business.Logic, sessionManager *framework.SessionManager) http.Handler {
	router := httprouter.New()

	router.GET("/counter",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.Counter(), "web", logic),
			sessionManager),
	)

	return router
}
