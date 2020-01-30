package main

import (
	"net/http"

	"github.com/WillCoates/FYP/auth/business"
	"github.com/WillCoates/FYP/auth/routes"
	"github.com/julienschmidt/httprouter"
)

func CreateRouter(logic *business.Logic) http.Handler {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/auth", routes.GetLogin(logic))
	router.POST("/auth", routes.PostLogin(logic))
	router.GET("/token", routes.GetToken(logic))
	router.POST("/token", routes.GetToken(logic))
	return router
}
