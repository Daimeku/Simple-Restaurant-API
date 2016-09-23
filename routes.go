package main

import (
	"github.com/julienschmidt/httprouter"
	// "net/http"
)

type Route struct {
	name        string
	method      string
	path        string
	handlerFunc httprouter.Handle
}

type Routes []Route

func makeRouter() *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		router.Handle(route.method, route.path, route.handlerFunc)
	}

	return router
}

var routes = Routes{

	Route{
		"home",
		"GET",
		"/",
		handleIndex,
	},
	Route{
		"restaurants",
		"GET",
		"/data/restaurants",
		handleRestaurants,
	},
	Route{
		"restaurant",
		"GET",
		"/data/restaurants/:restaurantSearchName",
		handleRestaurant,
	},
}
