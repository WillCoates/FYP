package controllers

import (
	"log"
	"net/http"

	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
)

type sensorAddData struct {
	UserID string
}

func SensorAdd(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data sensorAddData
		token, ok := framework.GetToken(r.Context())
		if !ok {
			log.Println("SensorAdd invoked without token")
			return
		}

		data.UserID = token.Payload.Subject

		err := templateManager.Execute("sensor_add", w, data)
		if err != nil {
			log.Println("SensorAdd failed to execute template", err)
		}
	}
}
