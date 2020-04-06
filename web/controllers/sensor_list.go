package controllers

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/WillCoates/FYP/common/protocol/sensors"
	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

		data.Sensors = make(chan sensorListSensorData)

		latestReadings, err := sensorService.GetLatestSensorReadings(ctx, &latestReadingsReq)
		if err != nil {
			state, ok := status.FromError(err)
			if !ok || state.Code() != codes.NotFound {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			close(data.Sensors)
		} else {
			// Read in background
			go func() {
				defer close(data.Sensors)
				for {
					reading, err := latestReadings.Recv()
					if err == io.EOF {
						break
					} else if err != nil {
						log.Println(err)
						return
					}

					var sensorData sensorListSensorData

					sensorData.Name = reading.UnitName
					sensorData.Reading = reading.Reading
					sensorData.MeasurementName = reading.Measurementname
					sensorData.MeasurementUnit = reading.Measurementunit
					sensorData.Sensor = reading.Sensor
					sensorData.Unit = reading.Unit
					sensorData.LastUpdated = formatTime(reading.Timestamp)

					data.Sensors <- sensorData
				}
			}()
		}

		// Render
		err = templateManager.Execute("sensor_list", w, data)
		if err != nil {
			log.Println("SensorList failed to render", err)
		}
	}
}
