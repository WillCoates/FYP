package controllers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type scriptEditDataError struct {
	Timestamp string
	Message   string
}

type scriptEditData struct {
	Errors chan scriptEditDataError
	Script scripting.Script
}

func ScriptEdit(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data scriptEditData
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("ScriptEdit invoked without session")
			return
		}

		// Get parameters
		scriptID := r.URL.Query().Get("id")

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

		// Get script
		var scriptRequest scripting.GetScriptRequest

		scriptRequest.Id = scriptID

		script, err := scriptService.GetScript(ctx, &scriptRequest)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		data.Script = *script

		// Get errors
		var errorRequest scripting.GetScriptErrorsRequest
		errorRequest.Id = []string{script.Details.Id}
		errorRequest.Since = script.Details.LastModified
		errorRequest.Limit = 100

		errors, err := scriptService.GetScriptErrors(ctx, &errorRequest)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		// Process errors in background
		data.Errors = make(chan scriptEditDataError)

		go func() {
			defer close(data.Errors)
			for {
				scriptError, err := errors.Recv()
				if err == io.EOF {
					return
				} else if err != nil {
					log.Println("ScriptEdit error while processing errors", err)
					continue
				}
				var errorData scriptEditDataError
				errorData.Message = scriptError.Message
				errorData.Timestamp = formatTime(scriptError.Timestamp)
				data.Errors <- errorData
			}
		}()

		// Render
		err = templateManager.Execute("script_edit", w, data)
		if err != nil {
			log.Println("ScriptEdit failed to render", err)
		}
	}
}

func ScriptEditPostback(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	pageRender := ScriptEdit(templateManager, logic)
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("ScriptEdit invoked without session")
			return
		}

		md := metadata.Pairs("Authorization", session.Values["token_web"])
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		grpcClient, err := grpc.Dial(logic.Config["script_service"], grpc.WithInsecure())
		if err != nil {
			log.Println("ScriptEdit failed to connect GRPC", err)
			w.WriteHeader(500)
			return
		}
		defer grpcClient.Close()

		scriptingService := scripting.NewScriptingServiceClient(grpcClient)

		err = r.ParseForm()
		if err != nil {
			log.Println("ScriptEdit failed ot parse form", err)
			w.WriteHeader(500)
			return
		}
		switch r.Form.Get("submit") {
		case "Edit":
			var editedScriptInfo scripting.Script
			editedScriptInfo.Source = r.Form.Get("src")
			editedScriptInfo.Details = new(scripting.ScriptDetails)
			editedScriptInfo.Details.Id = r.URL.Query().Get("id")
			editedScriptInfo.Details.Name = r.Form.Get("name")

			// Remove carriage returns to ensure string is UNIX formatted
			subs := strings.Split(strings.Replace(r.Form.Get("subs"), "\r", "", -1), "\n")
			// Ensure last line isn't blank
			if subs[len(subs)-1] == "" {
				subs = subs[:len(subs)-1]
			}
			editedScriptInfo.Details.Subscriptions = subs

			_, err = scriptingService.UpdateScript(ctx, &editedScriptInfo)
			if err != nil {
				log.Println("ScriptEdit failed to edit script", err)
				w.WriteHeader(500)
				return
			}
		case "Delete":
			var deletingScriptInfo scripting.Script
			deletingScriptInfo.Details = new(scripting.ScriptDetails)
			deletingScriptInfo.Details.Id = r.URL.Query().Get("id")

			_, err = scriptingService.DeleteScript(ctx, &deletingScriptInfo)
			if err != nil {
				log.Println("ScriptEdit failed to delete script", err)
				w.WriteHeader(500)
				return
			}

			http.Redirect(w, r, "/scripts", http.StatusSeeOther)
			return
		}
		pageRender(w, r, nil)
	}
}
