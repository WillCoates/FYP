package controllers

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type scriptErrorsErrorData struct {
	ScriptID  string
	Name      string
	Timestamp string
	Message   string
}

type scriptErrorsData struct {
	Errors chan scriptErrorsErrorData
}

func ScriptErrors(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data scriptErrorsData
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("SensorErrors invoked without session")
			return
		}

		// Script service connection
		md := metadata.Pairs("Authorization", session.Values["token_web"])
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		grpcClient, err := grpc.Dial(logic.Config["script_service"], grpc.WithInsecure())
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		defer grpcClient.Close()

		scriptService := scripting.NewScriptingServiceClient(grpcClient)

		// Get sensors
		var errorsRequest scripting.GetScriptErrorsRequest
		errorsRequest.Since = time.Now().Unix() - 24*60*60
		errorsRequest.Limit = 500

		data.Errors = make(chan scriptErrorsErrorData)

		errors, err := scriptService.GetScriptErrors(ctx, &errorsRequest)
		if err != nil {
			state, ok := status.FromError(err)
			if !ok || state.Code() != codes.NotFound {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			close(data.Errors)
		} else {
			// Read in background
			go func() {
				defer close(data.Errors)
				for {
					scriptError, err := errors.Recv()
					if err == io.EOF {
						break
					} else if err != nil {
						log.Println(err)
						return
					}

					var errorData scriptErrorsErrorData

					errorData.ScriptID = scriptError.Script.Id
					errorData.Name = scriptError.Script.Name
					errorData.Message = scriptError.Message
					errorData.Timestamp = formatTime(scriptError.Timestamp)

					data.Errors <- errorData
				}
			}()
		}

		// Render
		err = templateManager.Execute("script_errors", w, data)
		if err != nil {
			log.Println("SensorList failed to render", err)
		}
	}
}
