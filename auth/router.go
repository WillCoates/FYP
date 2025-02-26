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
	router.GET("/keys", routes.GetKeys(logic))
	router.GET("/keys.sig", routes.GetKeysSignature(logic))
	router.POST("/rabbit/user", routes.RabbitUser(logic))
	router.POST("/rabbit/vhost", routes.RabbitVhost(logic))
	router.POST("/rabbit/resource", routes.RabbitResource(logic))
	router.POST("/rabbit/topic", routes.RabbitTopic(logic))
	return router
}
