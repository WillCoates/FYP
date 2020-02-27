package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/WillCoates/FYP/web/framework"
	"github.com/julienschmidt/httprouter"
)

func Counter() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		sess, ok := framework.GetSession(r.Context())
		if !ok {
			w.WriteHeader(500)
			fmt.Fprint(w, "<pre>500 Internal Server Error\nNo session provided</pre>")
			log.Println("Counter Session Test Route: No session found, missing middleware?")
			return
		}
		count, ok := sess.Values["count"]
		if ok {
			countInt, _ := strconv.ParseInt(count, 10, 64)
			countInt++
			count = strconv.FormatInt(countInt, 10)
		} else {
			count = "1"
		}
		sess.Values["count"] = count
		fmt.Fprint(w, "You have loaded this page", count, "times!")
	}
}
