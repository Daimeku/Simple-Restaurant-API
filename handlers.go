package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github/Simple-Restaurant-API/models"
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
	// restaurant := models.Restaurant{Id: 1, Name: ""}
	restaurant := models.NewRestaurant()
	// var restaurant *models.Restaurant

	conf, err := restaurant.FindByField("searchName", searchName) // find the restaurant by its searchName
	if err != nil {                                               // if it isn't found return an error
		fmt.Fprint(writer, "There was an error trying to find the requested restaurant")
		return
	}
	//if the restaurant was found display it
	if conf {
		writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
		writer.WriteHeader(http.StatusOK)
		// x, _ := json.Marshal(restaurant)
		// fmt.Fprintf(writer, "the restaurant name is: %s", x)
		fmt.Println("the param ", restaurant)
		json.NewEncoder(writer).Encode(restaurant)
	} else {
		writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, "the restaurant ", searchName, " was not found")
		fmt.Println("the restaurant %s wasn't found", searchName)

	}

}

// handles the /restaurants GET route
func handleRestaurants(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var res *models.Restaurant
	restaurants := res.FindAll() //find all the restaurants
	// display an error if there are no restaurants
	if len(restaurants) < 1 {
		fmt.Fprint(writer, "No restaurants available")
		return
	}
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	// jsList, _ := json.Marshal(restaurants)
	// display the list of restaurants
	// fmt.Fprint(writer, "restaurants", jsList)
	json.NewEncoder(writer).Encode(restaurants)
}
