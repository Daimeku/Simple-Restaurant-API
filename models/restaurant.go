package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

const (
	Driver           = "mysql"
	ConnectionString = "root:@tcp(localhost:3306)/mydb"
)

// ===========================RESTAURANT=================================================================
type Restaurant struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	SearchName string     `json:searchName`
	Type       string     `json:type`
	Menu       []MenuItem `json:menu`
}

type Restaurantl *struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	SearchName string `json:searchName`
}

type Restaurants []Restaurant

// type RestaurantList []*Restaurant

//returns a string with the values of the restaurant
func (res *Restaurant) String() string {
	var fullString string = "{ " + "id: " + strconv.Itoa(res.Id) + " | name: " + res.Name + " | searchName: " + res.SearchName + " }"
	return fullString
}

//finds the restaurant corresponding to the id and populates the struct
func (res *Restaurant) FindById(id int) bool {
	var conf bool = false

	//connect to the database
	db, err := sql.Open(Driver, ConnectionString)
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
	res.Id = resId
	res.Name = resName

	conf = true
	return conf
}

//finds the restaurant by a specified fieldname and populates the *Restaurant struct
func (res *Restaurant) FindByField(fieldName string, fieldValue string) (bool, error) {

	var query string = "SELECT * FROM restaurants where restaurants." + fieldName + " = ?"

	//connect to the database
	db, err := sql.Open(Driver, ConnectionString)
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
	// fmt.Println(result.Columns())
	conf := res.Populate(result) // populate the restaurant and return confirmation

	return conf, nil
}

//finds and returns a slice of all *Restaurants
func (res *Restaurant) FindAllPtr() ([]*Restaurant, error) {

	resList := make([]*Restaurant, 0) // create the slice of restaurants

	//create the connection and check for errors
	db, err := sql.Open(Driver, ConnectionString)
	if err != nil {
		fmt.Println("an error occurred while connecting to the db - ", err)
		return nil, err
	}

	//execute the query and check for errors
	result, err := db.Query("SELECT * FROM restaurants")
	if err != nil {
		fmt.Println("an error occurred while executing the query - ", err)
		return nil, err
	}

	//populate the list of *Restaurants and check for errors
	resList, err = res.PopulateListPtr(result)
	if err != nil {
		fmt.Println("error populating the list - ", err)
		return nil, err
	}

	return resList, nil
}

func (res *Restaurant) FindAll() ([]Restaurant, error) {
	resList := make([]Restaurant, 0) // create the slice of restaurants

	//create the connection and check for errors
	db, err := sql.Open(Driver, ConnectionString)
	if err != nil {
		fmt.Println("an error occurred while connecting to the db - ", err)
		return nil, err
	}

	//execute the query and check for errors
	result, err := db.Query("SELECT * FROM restaurants")
	if err != nil {
		fmt.Println("an error occurred while executing the query - ", err)
		return nil, err
	}

	//populate the list of *Restaurants and check for errors
	resList, err = res.PopulateList(result)
	if err != nil {
		fmt.Println("error populating the list - ", err)
		return nil, err
	}

	return resList, nil
}

//accepts *sql.Rows and creates a list of *Restaruants
func (res *Restaurant) PopulateListPtr(rows *sql.Rows) ([]*Restaurant, error) {

	var resList []*Restaurant

	// foreach row, read a restaurant from rows and add it to the list
	for rows.Next() {
		restaurant := Restaurant{}
		restaurant.Type = "Restaurant"
		// restaurant := NewRestaurant()
		// populate the restaurant and check for errors
		err := rows.Scan(&restaurant.Id, &restaurant.Name, &restaurant.SearchName)
		if err != nil {
			return nil, err
		}
		resList = append(resList, &restaurant) // add the restaurant to the list
	}
	return resList, nil
}

//populates and returns a list of Restaurants
func (res *Restaurant) PopulateList(rows *sql.Rows) ([]Restaurant, error) {
	var resList []Restaurant

	// foreach row, read a restaurant from rows and add it to the list
	for rows.Next() {
		restaurant := Restaurant{}
		restaurant.Type = "Restaurant"

		// populate the restaurant and check for errors
		err := rows.Scan(&restaurant.Id, &restaurant.Name, &restaurant.SearchName)
		if err != nil {
			return nil, err
		}
		resList = append(resList, restaurant) // add the restaurant to the list
	}
	return resList, nil
}

//accepts *sql.Rows and populates the calling restaurant struct
func (res *Restaurant) Populate(rows *sql.Rows) bool {
	var resName string
	var resId int
	var resSearchName string

	conf := rows.Next() //read the first row, conf=false if none
	if conf != true {
		return false
	}
	//read the column values into the variables
	err := rows.Scan(&resId, &resName, &resSearchName)
	if err != nil {
		fmt.Println("There was an error reading the rows - ", err)
		return false
	}

	//populate res
	res.Id = resId
	res.Name = resName
	res.SearchName = resSearchName
	conf = res.LoadMenuItems()
	if conf != true {
		fmt.Println("failed to load items for restaurant")
	}

	return true
}

//populates res.Menu with a list of menuItems
//must be called after res.Id has been set
func (res *Restaurant) LoadMenuItems() bool {
	//open the connection
	db, err := sql.Open(Driver, ConnectionString)
	if err != nil {
		fmt.Println("error opening connection for loading menuItems")
		return false
	}
	//select the list of menuItems
	result, err := db.Query("select * FROM menuItems WHERE restaurantId = ?", res.Id)
	if err != nil {
		fmt.Println("Error selecting menuItems - ", err)
		return false
	}

	menuItem := MenuItem{}
	//populate the restaurant's menu
	menu, err := menuItem.PopulateList(result)
	if err != nil {
		fmt.Println("failed to populate menu - ", err)
		return false
	}
	res.Menu = menu

	return true
}

//use as a constructor
func NewRestaurantPtr() *Restaurant {
	return &Restaurant{Type: "restaurant"}
}

func NewRestaurant() Restaurant {
	return Restaurant{Type: "restaurant"}
}

// ======================================END RESTAURANT==================================================
