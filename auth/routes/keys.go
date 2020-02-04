package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WillCoates/FYP/auth/business"
	"github.com/julienschmidt/httprouter"
)

func GetKeys(logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		keys, err := logic.GetKeyBundle()

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, err)
			log.Println(err)
		} else {
			w.Write(keys)
		}
	}
}

func GetKeysSignature(logic *business.Logic) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sig, err := logic.GetKeyBundleSignature()

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, err)
			log.Println(err)
		} else {
			w.Write(sig)
		}
	}
}
