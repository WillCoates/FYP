package framework

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type key int

var metaSession key = 0

const SESSION_COOKIE string = "session_id"

// SessionMiddleware adds a Session to a HTTP Request, saving after the request has been processed
func SessionMiddleware(next httprouter.Handle, manager *SessionManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cookie, err := r.Cookie(SESSION_COOKIE)
		var session *Session
		if err != nil {
			cookie = nil
		}

		if cookie == nil {
			session, err = manager.Create()
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, "<pre>500 Internal Server Error\nFailed to create session</pre>")
				log.Println("Failed to create session", err)
				return
			}
		} else {
			session, err = manager.Load(cookie.Value)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, "<pre>500 Internal Server Error\nFailed to load session</pre>")
				log.Println("Failed to load session", err)
				return
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), metaSession, session))

		cookie = &http.Cookie{
			Name:    SESSION_COOKIE,
			Value:   session.Id(),
			Expires: time.Now().AddDate(0, 0, 1),
		}
		http.SetCookie(w, cookie)

		next(w, r, p)

		err = session.Save()
		if err != nil {
			log.Println("Failed to save session", session.Id(), err)
		}
	}
}

// GetSession retrieves a session from a HTTP context
func GetSession(ctx context.Context) (*Session, bool) {
	session, ok := ctx.Value(metaSession).(*Session)
	return session, ok
}
