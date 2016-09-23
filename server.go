package main

import (
	"fmt"
	// "html"
	"database/sql"
	// "encoding/json"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {

	router := makeRouter()
	// router.GET("/", handleIndex)
	// router.GET("/data/restaurants/:restaurantSearchName", handleRestaurant)
	// router.GET("/data/restaurants", handleRestaurants)
	log.Fatal(http.ListenAndServe(":8080", router))
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
