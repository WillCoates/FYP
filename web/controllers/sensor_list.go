package controllers

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/WillCoates/FYP/common/protocol/sensors"
	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type sensorListSensorData struct {
	Name            string
	MeasurementName string
	Reading         string
	MeasurementUnit string
	LastUpdated     string
	Unit            string
	Sensor          string
}

type sensorListData struct {
	Sensors chan sensorListSensorData
}

const timeFormat = "15:04:05 Mon 02 Jan 2006"

func SensorList(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data sensorListData
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("SensorList invoked without session")
			return
		}

		// Query parameters
		includeHiddenStr := r.URL.Query().Get("include_hidden")
		includeHidden := includeHiddenStr == "on"

		// Sensor service connection
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

		// Get sensors
		var latestReadingsReq sensors.GetSensorReadingsRequest

		latestReadingsReq.IgnoreHidden = !includeHidden

		latestReadings, err := sensorService.GetLatestSensorReadings(ctx, &latestReadingsReq)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		data.Sensors = make(chan sensorListSensorData)

		// Read in background
		go func() {
			defer close(data.Sensors)
			for {
				reading, err := latestReadings.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Println(err)
				}

				var sensorData sensorListSensorData
				lastUpdated := time.Unix(reading.Timestamp, 0)

				sensorData.Name = reading.UnitName
				sensorData.Reading = reading.Reading
				sensorData.MeasurementName = reading.Measurementname
				sensorData.MeasurementUnit = reading.Measurementunit
				sensorData.Sensor = reading.Sensor
				sensorData.Unit = reading.Unit
				sensorData.LastUpdated = lastUpdated.Local().Format(timeFormat)

				data.Sensors <- sensorData
			}
		}()

		// Render
		err = templateManager.Execute("sensor_list", w, data)
		if err != nil {
			log.Println("SensorList failed to render", err)
		}
	}
}
