package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/WillCoates/FYP/common/auth"
	"github.com/WillCoates/FYP/web/business"
	"github.com/julienschmidt/httprouter"
)

func AuthMiddleware(next httprouter.Handle, audience string, logic *business.Logic) httprouter.Handle {
	authEndpoint := logic.Config["auth_endpoint"]
	keys, err := auth.LoadBundleHTTP(authEndpoint+"/keys", logic.MasterKey)
	if err != nil {
		log.Println("Auth endpoint is", authEndpoint)
		log.Fatalln("Failed to load key bundle", err)
	}
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess, ok := GetSession(r.Context())
		if !ok {
			w.WriteHeader(500)
			fmt.Fprint(w, "<pre>500 Internal Server Error\nNo session provided</pre>")
			log.Println("Auth Middleware: No session found, missing middleware?", r.RequestURI)
			return
		}

		var token *auth.Token
		var err error

		_, ok = r.URL.Query()["code"]
		if ok {
			code := r.URL.Query().Get("code")
			_ = r.URL.Query().Get("state")

			// Remove code and state from query to make auth server happy
			requestQuery := r.URL.Query()
			requestQuery.Del("code")
			requestQuery.Del("state")
			r.URL.RawQuery = requestQuery.Encode()

			tokenRequest, err := http.NewRequest(http.MethodGet, authEndpoint+"/token", nil)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, "<pre>500 Internal Server Error\nUnexpected failure authenticating</pre>")
				log.Println("Failed to create token request", err)
				return
			}

			query := make(url.Values)
			query.Add("grant_type", "authorization_code")
			query.Add("code", code)
			query.Add("redirect_uri", r.Host+r.URL.String())
			query.Add("client_id", audience)
			tokenRequest.URL.RawQuery = query.Encode()

			log.Println(tokenRequest.URL.String())

			tokenResponse, err := http.DefaultClient.Do(tokenRequest)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, "<pre>500 Internal Server Error\nUnexpected failure authenticating</pre>")
				log.Println("Failed to send token request", err)
				return
			}

			defer tokenResponse.Body.Close()
			tokenRaw, err := ioutil.ReadAll(tokenResponse.Body)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, "<pre>500 Internal Server Error\nUnexpected failure authenticating</pre>")
				log.Println("Failed to read token response", err)
				return
			}

			tokenJson := make(map[string]interface{})
			err = json.Unmarshal(tokenRaw, &tokenJson)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, "<pre>500 Internal Server Error\nUnexpected failure authenticating</pre>")
				log.Println("Failed to read json", err)
				return
			}

			token, ok := tokenJson["access_token"]
			if !ok {
				w.WriteHeader(500)
				fmt.Fprint(w, "<pre>500 Internal Server Error\nUnexpected failure authenticating</pre>")
				log.Println("Token fetch failure", tokenJson)
				return
			}

			sess.Values["token_"+audience] = token.(string)
		}

		tokenRaw, ok := sess.Values["token_"+audience]
		if ok {
			token, err = auth.ParseToken([]byte(tokenRaw), keys.Keys)
			if err != nil {
				ok = false
			}
		}

		if !ok {
			redirectURL, err := url.Parse(authEndpoint + "/auth")
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, "<pre>500 Internal Server Error\nFailed to decode auth endpoint</pre>")
				log.Println("Failed to parse auth endpoint", err)
				return
			}
			query := make(url.Values)

			query.Add("response_type", "code")
			query.Add("client_id", audience)
			query.Add("redirect_uri", r.Host+r.RequestURI)
			query.Add("state", "test") // TODO: Randomize this
			redirectURL.RawQuery = query.Encode()
			http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), metaToken, token))

		next(w, r, p)
	}
}

// GetToken retrieves a JSON Web Token from a HTTP context
func GetToken(ctx context.Context) (*Session, bool) {
	session, ok := ctx.Value(metaSession).(*Session)
	return session, ok
}
