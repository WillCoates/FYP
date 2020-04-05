package controllers

import (
	"net/http"

	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
)

type dashboardData struct {
	AverageHumidity framework.Graph
}

func Dashboard(templateManager *framework.TemplateManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data dashboardData
		data.AverageHumidity.ID = "avgHumidity"
		data.AverageHumidity.Title = "Average Humidity"
		templateManager.Execute("dashboard", w, data)
	}
}
