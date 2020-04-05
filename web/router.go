package main

import (
	"net/http"

	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/controllers"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
)

func CreateRouter(logic *business.Logic, sessionManager *framework.SessionManager, templateManager *framework.TemplateManager) http.Handler {
	router := httprouter.New()

	router.GET("/counter",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.Counter(), "web", logic),
			sessionManager))

	router.GET("/",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.SensorList(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/sensors",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.SensorList(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/sensor",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.Sensor(templateManager, logic), "web", logic),
			sessionManager))

	router.POST("/sensor",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.SensorPostback(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/sensors/add",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.SensorAdd(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/sensors/token",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.SensorToken(templateManager, logic), "sensor", logic),
			sessionManager))

	return router
}
