package main

import (
	"github.com/julienschmidt/httprouter"
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
		// create a handler function wrapped in a logger
		handlerFunction := setLogger(route.handlerFunc, route.name)
		//assign the route to the function
		router.Handle(route.method, route.path, handlerFunction)
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
