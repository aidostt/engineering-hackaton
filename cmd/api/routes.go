package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) Router() *httprouter.Router {
	router := httprouter.New()
	return router
}
