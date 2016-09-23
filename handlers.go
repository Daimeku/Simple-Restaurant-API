package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// ======================================HANDLERS========================================================

// handles the home route
func handleIndex(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	fmt.Fprint(writer, "this is the home page ", request.URL.Path)
}

//handles the /restaurant/:resaurantId GET route
func handleRestaurant(writer http.ResponseWriter, request *http.Request, parms httprouter.Params) {

	searchName := parms.ByName("restaurantSearchName") //retrieve the searchName
	restaurant := Restaurant{id: 1, name: ""}

	conf, err := restaurant.findByField("searchName", searchName) // find the restaurant by its searchName
	if err != nil {                                               // if it isn't found return an error
		fmt.Fprint(writer, "There was an error trying to find the requested restaurant")
		return
	}
	//if the restaurant was found display it
	if conf {
		fmt.Fprintf(writer, "the restaurant name is: %s", restaurant.name)
		fmt.Println("the param %s", restaurant.String())
	} else {
		fmt.Fprint(writer, "the restaurant %s was not found", searchName)
		fmt.Println("the restaurant %s wasn't found", searchName)
	}

}

// handles the /restaurants GET route
func handleRestaurants(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var res *Restaurant
	restaurants := res.findAll() //find all the restaurants
	// display an error if there are no restaurants
	if len(restaurants) < 1 {
		fmt.Fprint(writer, "No restaurants available")
		return
	}
	// display the list of restaurants
	fmt.Fprint(writer, "restaurants", restaurants)
}
