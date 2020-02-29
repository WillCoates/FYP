package routes

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/WillCoates/FYP/auth/business"
	"github.com/WillCoates/FYP/auth/model"
	"github.com/julienschmidt/httprouter"
)

func generateLoginTemplate() *template.Template {
	template, err := template.New("loginTemplate").ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}
	return template
}

var loginTemplate *template.Template = generateLoginTemplate()

type loginData struct {
	LoginEmail    string
	RegisterEmail string
	RegisterName  string
	LoginError    string
	RegisterError string
}

func renderlogin(w io.Writer, data *loginData) error {
	return loginTemplate.ExecuteTemplate(w, "login", data)
}

func GetLogin(logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var data loginData
		err := renderlogin(w, &data)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}

func PostLogin(logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := r.ParseForm()
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		query := r.URL.Query()
		responseType := query.Get("response_type")
		clientID := query.Get("client_id")
		redirectURI := query.Get("redirect_uri")
		state := query.Get("state")
		challengeText := query.Get("code_challenge")
		challengeAlgo := query.Get("code_challenge_method")
		var challenge *model.Challenge

		if challengeText != "" {
			challenge, err = model.NewChallenge(challengeText, challengeAlgo)
			if err != nil {
				w.Write([]byte("Failed to decode challenge"))
				return
			}
		}

		if responseType != "code" {
			w.Write([]byte("Invalid response type"))
			return
		}

		target, err := url.Parse(redirectURI)

		if err != nil {
			w.Write([]byte("Failed to parse redirect"))
			return
		}

		var data loginData
		switch r.FormValue("submit") {
		case "Login":
			data.LoginEmail = r.FormValue("email")
			password := r.FormValue("password")

			token, err := logic.Authenticate(context.Background(), data.LoginEmail, password, clientID, 0)

			if err != nil {
				data.LoginError = err.Error()
			} else {
				code, err := logic.CreateAuthCode(context.Background(), token, redirectURI, challenge)
				if err != nil {
					log.Println("Failed to generate auth code")
					log.Println(err)
					data.LoginError = "Failed to process request, please try again in a few minutes."
				} else {
					query := target.Query()
					query.Add("code", code)
					query.Add("state", state)
					target.RawQuery = query.Encode()
					w.Header().Add("Location", target.String())
					w.WriteHeader(302)
					return
				}
			}

		case "Register":
			data.RegisterEmail = r.FormValue("email")
			data.RegisterName = r.FormValue("name")
			password := r.FormValue("password")

			err := logic.Register(context.Background(), data.RegisterEmail, data.RegisterName, password)
			if err != nil {
				data.RegisterError = err.Error()
			} else {
				token, err := logic.Authenticate(context.Background(), data.RegisterEmail, password, clientID, 0)

				if err != nil {
					data.RegisterError = err.Error()
				} else {
					code, err := logic.CreateAuthCode(context.Background(), token, redirectURI, challenge)
					if err != nil {
						log.Println("Failed to generate auth code")
						log.Println(err)
						data.RegisterError = "Failed to process request, please try again in a few minutes."
					} else {
						query := target.Query()
						query.Add("code", code)
						query.Add("state", state)
						target.RawQuery = query.Encode()
						w.Header().Add("Location", target.String())
						w.WriteHeader(302)
						return
					}
				}
			}

		default:
			w.Write([]byte("Invalid form"))
		}

		err = renderlogin(w, &data)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}
