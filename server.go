package main

import (
	"fmt"
	// "html"
	"database/sql"
	// "encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

func main() {

	var router *httprouter.Router = httprouter.New()
	router.GET("/", handleIndex)
	router.GET("/data/restaurants/:restaurantSearchName", handleRestaurant)
	router.GET("/data/restaurants", handleRestaurants)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// ===========================RESTAURANT=================================================================
type Restaurant struct {
	id         int
	name       string
	searchName string
}

//returns a string with the values of the restaurant
func (res *Restaurant) String() string {
	var fullString string = "{ " + "id: " + strconv.Itoa(res.id) + " | name: " + res.name + " | searchName: " + res.searchName + " }"
	return fullString
}

//finds the restaurant corresponding to the id and populates the struct
func (res *Restaurant) findById(id int) bool {
	var conf bool = false

	//connect to the database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/mydb")
	if err != nil {
		fmt.Println("error opening connection in model - ", err)
		return false
	}

	//query the restaurants table for a restaurant matching the id
	result, err := db.Query("SELECT * FROM restaurants where id = ?", id)
	if err != nil {
		fmt.Println("error opening connection in model - ", err)
		return false
	}

	//declare variables
	var resName string
	var resId int
	result.Next() //read the first record & check for errors
	err = result.Scan(&resId, &resName)
	if err != nil {
		fmt.Println(err)
	}

	//set restaurant values
	res.id = resId
	res.name = resName

	conf = true
	return conf
}

//finds the restaurant by a specified fieldname and populates the *Restaurant struct
func (res *Restaurant) findByField(fieldName string, fieldValue string) (bool, error) {

	var query string = "SELECT * FROM restaurants where restaurants." + fieldName + " = ?"

	//connect to the database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/mydb")
	if err != nil {
		fmt.Println("error opening connection in model - ", err)
		return false, err
	}

	//query the restaurants table for a restaurant matching the field
	result, err := db.Query(query, fieldValue)
	if err != nil {
		fmt.Println("error opening connection in model - ", err)
		return false, err
	}
	fmt.Println(result.Columns())
	conf := res.populate(result) // populate the restaurant and return confirmation

	return conf, nil
}

//finds and returns a slice of all *Restaurants
func (res *Restaurant) findAll() []*Restaurant {

	resList := make([]*Restaurant, 0) // create the slice of restaurants

	//create the connection and check for errors
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/mydb")
	if err != nil {
		fmt.Println("an error occurred while connecting to the db - ", err)
		return resList
	}

	//execute the query and check for errors
	result, err := db.Query("SELECT * FROM restaurants")
	if err != nil {
		fmt.Println("an error occurred while executing the query - ", err)
		return resList
	}

	//populate the list of *Restaurants and check for errors
	resList, err = res.populateList(result)
	if err != nil {
		fmt.Println("error populating the list - ", err)
		return nil
	}

	return resList
}

//accepts *sql.Rows and creates a list of *Restaruants
func (res *Restaurant) populateList(rows *sql.Rows) ([]*Restaurant, error) {

	var resList []*Restaurant

	// foreach row, read a restaurant from rows and add it to the list
	for rows.Next() {
		restaurant := Restaurant{}
		// populate the restaurant and check for errors
		err := rows.Scan(&restaurant.id, &restaurant.name, &restaurant.searchName)
		if err != nil {
			return nil, err
		}
		resList = append(resList, &restaurant) // add the restaurant to the list
	}
	return resList, nil
}

//accepts *sql.Rows and populates the calling restaurant struct
func (res *Restaurant) populate(rows *sql.Rows) bool {
	var resName string
	var resId int
	var resSearchName string

	conf := rows.Next() //read the first row, conf=false if none

	//read the column values into the variables
	err := rows.Scan(&resId, &resName, &resSearchName)
	if err != nil {
		fmt.Println("There was an error reading the rows - ", err)
		return false
	}

	//populate res
	res.id = resId
	res.name = resName
	res.searchName = resSearchName

	return conf
}

// ======================================END RESTAURANT==================================================

// ======================================HANDLERS========================================================

// handles the home route
func handleIndex(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	fmt.Fprint(writer, "this is the home page ", request.URL.Path)
	fmt.Println(testConnection())
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

// =====================================END HANDLERS=====================================================

//just checks if a connection to the database can be established
func testConnection() bool {
	var result bool = true
	//attempt to open a connection to the database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/mydb")
	if err != nil {
		fmt.Println("Could not connect to db - ", err)
		return false
	}

	//attempt to query the database
	res, err := db.Exec(
		"INSERT INTO restaurants (name, id) VALUES (?, ?)",
		"cook shop",
		87)
	if err != nil {
		fmt.Println("error executing query - ", err)
		return false
	}
	fmt.Println(res)
	db.Close()

	return result
}
