package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/WillCoates/FYP/auth/business"
	"github.com/WillCoates/FYP/common/util"
	"github.com/julienschmidt/httprouter"
)

func GetToken(logic *business.Logic) httprouter.Handle {
	invalidRequest := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"error\":\"invalid_request\"}"))
	}

	authorizationCode := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		redirect := query.Get("redirect_uri")
		clientID := query.Get("client_id")
		clientSecret := query.Get("client_secret")
		solution := query.Get("code_verifier")

		token, err := logic.RetrieveAuthCode(context.Background(), code, redirect, solution)

		if err != nil {
			log.Println("Failed to retrieve auth code")
			log.Println(err)
			invalidRequest(w, r)
			return
		}

		parsedToken, err := logic.DecodeTokenStr(token)

		if err != nil {
			log.Println("Failed to decode token")
			log.Println(err)
			invalidRequest(w, r)
			return
		}

		if parsedToken.Payload.Audience != clientID {
			log.Println("Audience doesn't match clientID")
			invalidRequest(w, r)
			return
		}

		audience, err := logic.GetAudience(context.Background(), clientID)
		if err != nil {
			log.Println("Failed to retrieve audience for auth code: ", audience)
			log.Println(err)
			invalidRequest(w, r)
			return
		}

		if audience.Secret != "" && !util.SecureEqualsStr(audience.Secret, clientSecret) {
			log.Println("Wrong secret")
			log.Println("Expected:", audience.Secret)
			log.Println("Provided:", clientSecret)
			invalidRequest(w, r)
			return
		}

		now := time.Now().Unix()

		w.Write([]byte(fmt.Sprintf("{\"access_token\":%q,\"expires_in\":%d}", token, parsedToken.Payload.Expires-now)))
	}

	handlers := map[string]http.HandlerFunc{
		"authorization_code": authorizationCode,
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		query := r.URL.Query()
		handler, success := handlers[query.Get("grant_type")]
		if success {
			handler(w, r)
		} else {
			invalidRequest(w, r)
		}
	}
}
