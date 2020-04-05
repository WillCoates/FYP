package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
)

func SensorToken(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("SensorToken invoked without session")
			return
		}

		fmt.Fprintln(w, session.Values["token_sensor"])
	}
}
