package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/WillCoates/FYP/auth/business"
	"github.com/WillCoates/FYP/auth/model"
	"github.com/WillCoates/FYP/common/auth"
	"github.com/julienschmidt/httprouter"
)

// https://github.com/rabbitmq/rabbitmq-auth-backend-http

func RabbitUser(logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprint(w, "deny")
			log.Println("RabbitUser failed to parse form data")
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		// @ is never in a token, treat as email
		var token *auth.Token
		if strings.ContainsRune(username, '@') {
			token, err = logic.Authenticate(context.Background(), username, password, "special", 1)
		} else {
			token, err = logic.DecodeTokenStr(username)
		}

		if err != nil {
			fmt.Fprint(w, "deny")
			log.Println("RabbitUser failed to decode token", err)
			return
		}

		user, err := logic.GetUser(context.Background(), token)
		if err != nil {
			fmt.Fprint(w, "deny")
			log.Println("RabbitUser failed to get user from token")
			return
		}

		fmt.Fprint(w, "allow")

		if user.SpecialPerms != nil {
			for _, perm := range user.SpecialPerms {
				if strings.HasPrefix(perm, "rabbit-") {
					fmt.Fprint(w, "", strings.TrimPrefix(perm, "rabbit-"))
				}
			}
		}
	}
}

func RabbitVhost(logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// We don't use vhosts, so just allow for simplicity
		fmt.Fprint(w, "allow")
	}
}

func RabbitResource(logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprint(w, "deny", err)
			return
		}
		username := r.FormValue("username")
		resource := r.FormValue("resource")
		name := r.FormValue("name")
		permission := r.FormValue("permission")

		log.Println("RabbitResource", username, resource, name, permission)

		var user *model.User

		if strings.ContainsRune(username, '@') {
			user, err = logic.GetUserByEmail(context.Background(), username)
		} else {
			token, err := logic.DecodeTokenStr(username)
			if err != nil {
				fmt.Fprint(w, "deny", err)
				return
			}
			user, err = logic.GetUser(context.Background(), token)
		}

		if err != nil {
			fmt.Fprint(w, "deny", err)
			return
		}

		if user.SpecialPerms != nil {
			for _, perm := range user.SpecialPerms {
				if perm == "rabbit-administrator" {
					fmt.Fprint(w, "allow")
					return
				}
			}
		}

		if permission != "read" && permission != "write" && permission != "configure" {
			fmt.Fprint(w, "deny")
			return
		}

		switch resource {
		case "topic":
			if name == "amq.topic" {
				fmt.Fprint(w, "allow")
			} else {
				fmt.Fprint(w, "deny")
			}

		case "exchange":
			if name == "amq.topic" {
				fmt.Fprint(w, "allow")
			} else {
				fmt.Fprint(w, "deny")
			}

		case "queue":
			if strings.HasPrefix(name, "mqtt-subscription") {
				fmt.Fprint(w, "allow")
			} else {
				fmt.Fprint(w, "deny")
			}

		default:
			fmt.Fprint(w, "deny")
		}

	}
}

func RabbitTopic(logic *business.Logic) httprouter.Handle {
	/*
		topicRegex, err := regexp.Compile("^([a-zA-Z0-9]+)\\.([a-zA-Z0-9]+)(\\.([a-zA-Z0-9]+))?(\\.[a-zA-Z0-9]+)*$")
		if err != nil {
			log.Println("Failed to compile regex")
			log.Fatalln(err)
		}
	*/

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := r.ParseForm()

		if err != nil {
			fmt.Fprint(w, "deny")
			return
		}

		var username = r.FormValue("username")
		var name = r.FormValue("name")
		var permission = r.FormValue("permission")
		var routingKey = r.FormValue("routing_key")

		log.Println("RabbitTopic", username, name, permission, routingKey)

		var user *model.User

		if strings.ContainsRune(username, '@') {
			user, err = logic.GetUserByEmail(context.Background(), username)
		} else {
			token, err := logic.DecodeTokenStr(username)
			if err != nil {
				fmt.Fprint(w, "deny", err)
				return
			}
			user, err = logic.GetUser(context.Background(), token)
		}

		if user.SpecialPerms != nil {
			for _, perm := range user.SpecialPerms {
				if perm == "rabbit-administrator" {
					fmt.Fprint(w, "allow")
					return
				}
			}
		}

		if name != "amq.topic" {
			fmt.Fprint(w, "deny")
			return
		}

		/*
			match := topicRegex.FindStringSubmatch(routingKey)
			if len(match) < 2 {
				fmt.Fprintf(w, "deny")
				return
			}
			topicUser := match[1]
			//topicArea := match[2]
		*/
		index := strings.Index(routingKey, ".")
		var topicUser string

		if index == -1 {
			topicUser = routingKey
		} else {
			topicUser = routingKey[:index]
		}

		if topicUser != user.ID.Hex() {
			fmt.Fprint(w, "deny")
			return
		}

		fmt.Fprint(w, "allow")

		/* TODO
		var requiredPermission string

		switch {
		case topicArea == "sensor" && permission == "read":
			requiredPermission = "readSensor"
		case topicArea == "sensor" && permission == "write":
			requiredPermission = "writeSensor"
		case topicArea != "sensor" && permission == "read":
			requiredPermission = "readMq"
		case topicArea != "sensor" && permission == "write":
			requiredPermission = "writeMq"
		default:
			fmt.Fprintln(w, "deny")
			return
		}

		result := "deny"
		permissions, err := logic.GetTokenPermissions(context.Background(), token)

		for perm := range permissions {
			if perm == requiredPermission {
				result = "allow"
			}
		}

		fmt.Fprint(w, result)
		*/
	}
}
