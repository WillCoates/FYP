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

	router.GET("/fields",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.FieldList(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/fields/view",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.Field(templateManager, logic), "web", logic),
			sessionManager))

	router.POST("/fields/view",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.FieldPostback(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/scripts",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.ScriptList(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/scripts/errors",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.ScriptErrors(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/scripts/add",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.ScriptAdd(templateManager, logic), "web", logic),
			sessionManager))

	router.POST("/scripts/add",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.ScriptAddPostback(templateManager, logic), "web", logic),
			sessionManager))

	router.GET("/scripts/edit",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.ScriptEdit(templateManager, logic), "web", logic),
			sessionManager))

	router.POST("/scripts/edit",
		framework.SessionMiddleware(
			framework.AuthMiddleware(controllers.ScriptEditPostback(templateManager, logic), "web", logic),
			sessionManager))

	return router
}
