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
	"google.golang.org/grpc/metadata"
)

type fieldData struct {
	Field   sensors.Field
	Sensors chan sensorListSensorData
}

func Field(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data fieldData
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("FieldList invoked without session")
			return
		}

		fieldName := r.URL.Query().Get("name")

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

		// Get fields
		var fieldRequest sensors.GetFieldsRequest
		fieldRequest.Name = []string{fieldName}

		fields, err := sensorService.GetFields(ctx, &fieldRequest)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		field, err := fields.Recv()
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		data.Field = *field

		// Get sensors
		var sensorsRequest sensors.GetSensorReadingsRequest
		sensorsRequest.Site = fieldRequest.Name

		sensors, err := sensorService.GetLatestSensorReadings(ctx, &sensorsRequest)
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
				sensorInfo, err := sensors.Recv()
				if err == io.EOF {
					break
				}
				var sensor sensorListSensorData
				sensor.LastUpdated = formatTime(sensorInfo.Timestamp)
				sensor.MeasurementName = sensorInfo.Measurementname
				sensor.MeasurementUnit = sensorInfo.Measurementunit
				sensor.Name = sensorInfo.UnitName
				sensor.Reading = sensorInfo.Reading
				sensor.Sensor = sensorInfo.Sensor
				sensor.Unit = sensorInfo.Unit
				data.Sensors <- sensor
			}
		}()

		// Render
		err = templateManager.Execute("field", w, data)
		if err != nil {
			log.Println("FieldList failed to render", err)
		}
	}
}

func FieldPostback(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("FieldList invoked without session")
			return
		}

		fieldName := r.URL.Query().Get("name")

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

		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		switch r.Form.Get("submit") {
		case "Update":
			var updatedField sensors.Field
			updatedField.Name = fieldName
			updatedField.Crop = r.Form.Get("crop")
			sensorService.UpdateField(ctx, &updatedField)
			http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)

		case "Delete":
			var deletedField sensors.Field
			deletedField.Name = fieldName
			sensorService.DeleteField(ctx, &deletedField)
			http.Redirect(w, r, "/fields", http.StatusSeeOther)
		}

		w.WriteHeader(500)
	}
}
