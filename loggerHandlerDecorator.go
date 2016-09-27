package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

type HandleFunc func(http.ResponseWriter, *http.Request, httprouter.Params)

// func (handlerFunction httprouter.Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
// 	handlerFuncton(writer, request, params)
// }

//wraps the handler function in a logger
func setLogger(innerHandler httprouter.Handle, name string) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//I need to return a function that fits the type httprouter.Handle
	//my type HandleFunc also matches that type so I declare a new HandleFunc and have it accept all the httprouter.Handle params
	//then I call the original handler (innerHandler) and log the request details
	return HandleFunc(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		startTime := time.Now()
		innerHandler(writer, request, params)
		log.Println(request.Method, request.RequestURI, name, time.Since(startTime), request.RemoteAddr) // log the request
	})

}
