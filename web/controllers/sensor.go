package controllers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/WillCoates/FYP/common/protocol/sensors"
	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type sensorData struct {
	Graph      framework.Graph
	Sensor     *sensors.SensorInfo
	Since      string
	SinceNames map[int64]string
	Unit       string
	SensorName string
}

func Sensor(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	sinceValues := map[string]int64{
		"1 hour ago":  3600,
		"1 day ago":   24 * 60 * 60,
		"1 week ago":  7 * 24 * 60 * 60,
		"1 month ago": 30 * 24 * 60 * 60,
		"All time":    0,
	}
	sinceNames := make(map[int64]string)
	for name, since := range sinceValues {
		sinceNames[since] = name
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data sensorData
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("Sensor invoked without session")
			return
		}

		// Since persistence
		var since int64 = 24 * 60 * 60
		sinceStr, ok := session.Values["sensor_values_since"]
		if ok {
			since, _ = strconv.ParseInt(sinceStr, 10, 64)
		}
		sinceQuery := r.URL.Query().Get("since")

		if sinceQuery != "" {
			since, ok = sinceValues[sinceQuery]
			if ok {
				session.Values["sensor_values_since"] = strconv.FormatInt(since, 10)
			}
		}

		data.Since = sinceNames[since]
		data.SinceNames = sinceNames

		// Query parameters
		data.SensorName = r.URL.Query().Get("sensor")
		data.Unit = r.URL.Query().Get("unit")

		// Sensor service connection
		md := metadata.Pairs("Authorization", session.Values["token_web"])
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		grpcClient, err := grpc.Dial(logic.Config["sensor_service"], grpc.WithInsecure())
		if err != nil {
			log.Println("Sensor failed to connect to GRPC", err)
			w.WriteHeader(500)
			return
		}
		defer grpcClient.Close()

		sensorService := sensors.NewSensorsServiceClient(grpcClient)

		// Retrieve sensor information
		var sensorReq sensors.GetSensorsRequest
		sensorReq.Sensor = []string{data.SensorName}
		sensorReq.Unit = []string{data.Unit}
		sensorReq.IncludeHidden = true

		sensorsResult, err := sensorService.GetSensors(ctx, &sensorReq)
		if err != nil {
			log.Println("Sensor failed to get sensor info", err)
			w.WriteHeader(500)
			return
		}

		data.Sensor, err = sensorsResult.Recv()
		if err != nil {
			log.Println("Sensor failed to read sensor info", err)
			w.WriteHeader(500)
			return
		}

		// Retrieve sensor readings
		var readingReq sensors.GetSensorReadingsRequest
		readingReq.Sensor = []string{r.URL.Query().Get("sensor")}
		readingReq.Unit = []string{r.URL.Query().Get("unit")}
		if since != 0 {
			readingReq.Since = time.Now().Unix() - since
		}

		readings, err := sensorService.GetSensorReadings(ctx, &readingReq)
		if err != nil {
			log.Println("Sensor failed to get readings", err)
			w.WriteHeader(500)
			return
		}

		for {
			reading, err := readings.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println("Sensor failed to read readings", err)
				w.WriteHeader(500)
				return
			}

			val, err := strconv.ParseFloat(reading.Reading, 64)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			data.Graph.Data = append(data.Graph.Data, framework.GraphPoint{reading.Timestamp, val})
			data.Graph.Title = reading.Measurementname
		}
		data.Graph.ID = "readings"

		// Render
		err = templateManager.Execute("sensor", w, data)
		if err != nil {
			log.Println("Sensor failed to render", err)
		}
	}
}

func SensorPostback(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	pageRender := Sensor(templateManager, logic)
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("Sensor invoked without session")
			return
		}

		md := metadata.Pairs("Authorization", session.Values["token_web"])
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		grpcClient, err := grpc.Dial(logic.Config["sensor_service"], grpc.WithInsecure())
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		defer grpcClient.Close()

		sensorService := sensors.NewSensorsServiceClient(grpcClient)

		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		switch r.Form.Get("submit") {
		case "Edit":
			var editedSensorInfo sensors.SensorInfo
			editedSensorInfo.Hidden = r.Form.Get("hidden") == "on"
			editedSensorInfo.Measurementname = r.Form.Get("measurementname")
			editedSensorInfo.Measurementunit = r.Form.Get("measurementunit")
			editedSensorInfo.Name = r.Form.Get("name")
			editedSensorInfo.Sensor = r.URL.Query().Get("sensor")
			editedSensorInfo.Site = r.Form.Get("site")
			editedSensorInfo.Unit = r.URL.Query().Get("unit")

			editedSensorInfo.Latitude, err = strconv.ParseFloat(r.Form.Get("latitude"), 64)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}

			editedSensorInfo.Longitude, err = strconv.ParseFloat(r.Form.Get("longitude"), 64)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}

			_, err = sensorService.UpdateSensor(ctx, &editedSensorInfo)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}

		case "Delete":
			var deletedSensorInfo sensors.SensorInfo
			deletedSensorInfo.Sensor = r.URL.Query().Get("sensor")
			deletedSensorInfo.Unit = r.URL.Query().Get("unit")

			sensorService.DeleteSensor(ctx, &deletedSensorInfo)

			http.Redirect(w, r, "/sensors", http.StatusSeeOther)
			return
		}
		pageRender(w, r, nil)
	}
}
