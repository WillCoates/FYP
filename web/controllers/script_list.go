package controllers

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type scriptListScriptData struct {
	ID           string
	Name         string
	LastModified string
	ErrorCount   int64
}

type scriptListData struct {
	Scripts chan scriptListScriptData
}

func ScriptList(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data scriptListData
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("SensorList invoked without session")
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
		var scriptsRequest scripting.GetScriptsRequest

		data.Scripts = make(chan scriptListScriptData)

		scripts, err := scriptService.GetScripts(ctx, &scriptsRequest)
		if err != nil {
			state, ok := status.FromError(err)
			if !ok || state.Code() != codes.NotFound {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			close(data.Scripts)
		} else {
			// Read in background
			go func() {
				defer close(data.Scripts)
				for {
					script, err := scripts.Recv()
					if err == io.EOF {
						break
					} else if err != nil {
						log.Println(err)
						return
					}

					var scriptData scriptListScriptData

					scriptData.ID = script.Id
					scriptData.Name = script.Name
					scriptData.ErrorCount = script.RecentErrorCount
					scriptData.LastModified = formatTime(script.LastModified)

					data.Scripts <- scriptData
				}
			}()
		}

		// Render
		err = templateManager.Execute("script_list", w, data)
		if err != nil {
			log.Println("SensorList failed to render", err)
		}
	}
}
