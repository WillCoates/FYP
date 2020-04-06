package controllers

import (
	"context"
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

const scriptTemplateSrc = `if data["sensor"] == "TC" then
  local reading = tonumber(data["value"])
  if reading > 25 then
    publish("greenhouse/heater", "off")
  elseif reading < 15 then
  	publish("greenhouse/heater", "on")
  end
end`

func ScriptAdd(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data scriptEditData

		// Template script
		data.Script.Source = scriptTemplateSrc
		data.Script.Details = new(scripting.ScriptDetails)
		data.Script.Details.Name = "Untitled Script"
		data.Script.Details.Subscriptions = []string{"sensor/#"}
		data.Errors = make(chan scriptEditDataError)
		close(data.Errors)

		// Render
		err := templateManager.Execute("script_add", w, data)
		if err != nil {
			log.Println("ScriptAdd failed to render", err)
		}
	}
}

func ScriptAddPostback(templateManager *framework.TemplateManager, logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		session, ok := framework.GetSession(r.Context())
		if !ok {
			log.Println("ScriptAdd invoked without session")
			return
		}

		md := metadata.Pairs("Authorization", session.Values["token_web"])
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		grpcClient, err := grpc.Dial(logic.Config["script_service"], grpc.WithInsecure())
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		defer grpcClient.Close()

		scriptingService := scripting.NewScriptingServiceClient(grpcClient)

		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		switch r.Form.Get("submit") {
		case "Add":
			var editedScriptInfo scripting.Script
			editedScriptInfo.Source = r.Form.Get("src")
			editedScriptInfo.Details = new(scripting.ScriptDetails)
			editedScriptInfo.Details.Id = r.URL.Query().Get("id")
			editedScriptInfo.Details.Name = r.Form.Get("name")
			// Remove carriage returns to ensure string is UNIX formatted
			editedScriptInfo.Details.Subscriptions = strings.Split(strings.Replace(r.Form.Get("subs"), "\r", "", -1), "\n")
			addedScript, err := scriptingService.AddScript(ctx, &editedScriptInfo)
			if err != nil {
				log.Println("ScriptAdd failed to add script", err)
				w.WriteHeader(500)
				return
			}
			http.Redirect(w, r, "/scripts/edit?id="+addedScript.Details.Id, http.StatusSeeOther)
			return
		}
		w.WriteHeader(500)
	}
}
