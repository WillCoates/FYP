package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Redirect(target string, permanent bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Add("Location", target)
		if permanent {
			w.WriteHeader(301)
		} else {
			w.WriteHeader(302)
		}
	}
}
