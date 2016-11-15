package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

func main() {
	router := makeRouter()
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(router)))
}

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
