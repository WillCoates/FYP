package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func CreateRouter() http.Handler {
	router := httprouter.New()
	return router
}
