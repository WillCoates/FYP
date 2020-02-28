package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
)

func Dashboard(templateManager *framework.TemplateManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		_, ok := framework.GetSession(r.Context())
		if !ok {
			w.WriteHeader(500)
			fmt.Fprint(w, "<pre>500 Internal Server Error\nNo session provided</pre>")
			log.Println("Dashboard: No session found, missing middleware?")
			return
		}
		templateManager.Execute("dashboard", w, nil)
	}
}
