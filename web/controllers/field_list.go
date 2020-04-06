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

type fieldListData struct {
	Fields chan sensors.Field
}

func FieldList(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data fieldListData
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("FieldList invoked without session")
			return
		}

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

		data.Fields = make(chan sensors.Field)

		fields, err := sensorService.GetFields(ctx, &fieldRequest)
		if err != nil {
			state, ok := status.FromError(err)
			if !ok || state.Code() != codes.NotFound {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			close(data.Fields)
		} else {
			// Read in background
			go func() {
				defer close(data.Fields)
				for {
					field, err := fields.Recv()
					if err == io.EOF {
						break
					} else if err != nil {
						log.Println(err)
						return
					}

					data.Fields <- *field
				}
			}()
		}

		// Render
		err = templateManager.Execute("field_list", w, data)
		if err != nil {
			log.Println("FieldList failed to render", err)
		}
	}
}
