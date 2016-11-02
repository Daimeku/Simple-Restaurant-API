package models

import (
	"database/sql"
	"fmt"
)

type MenuItem struct {
	Type        string `json:type`
	Id          int    `json:id`
	Name        string `json:name`
	Description string `json:description`
	Price       string `json:price`
}

//queries for all menuItems and returns them
//@ToDo paginate the results
func (menuItem *MenuItem) FindAll() ([]MenuItem, error) {

	menuItems := []MenuItem{}
	//open the DB connection
	db, err := sql.Open(Driver, ConnectionString)
	if err != nil {
		fmt.Println("There was an error connecting to the db - ", err)
		return menuItems, err
	}
	//select the menu Items
	result, err := db.Query("SELECT * FROM menuItems")
	if err != nil {
		fmt.Println("There was an error connecting to the db - ", err)
		return menuItems, err
	}

	menuItems, err = menuItem.PopulateList(result)
	if err != nil {
		fmt.Println("error populating the list - ", err)
		return menuItems, err
	}

	return menuItems, nil
}

//accepts a sql.Row and returns a populated *menuItem
func (menuItem *MenuItem) Populate(rows *sql.Rows) bool {
	var id int
	var name string
	var description string
	var price string
	var resId int

	conf := rows.Next()
	if conf != true {
		return false
	}

	err := rows.Scan(&id, &name, &description, &price, &resId)
	if err != nil {
		fmt.Println("Error reading menuItem row - ", err)
		return false
	}

	menuItem.Id = id
	menuItem.Name = name
	menuItem.Description = description
	menuItem.Price = price
	menuItem.Type = "MenuItem"

	return true
}

//accepts a sql.Row of menuItems and returns a populated list of menuItems
func (menuItem *MenuItem) PopulateList(rows *sql.Rows) ([]MenuItem, error) {
	menuItems := []MenuItem{}

	var id int
	var name string
	var description string
	var price string
	var resId int

	//for each row, create a menuItem and add it to the list
	for rows.Next() {
		menuItemCurr := MenuItem{}

		err := rows.Scan(&id, &name, &description, &price, &resId)
		if err != nil {
			return menuItems, err
		}

		//populate the menuItem
		menuItemCurr.Id = id
		menuItemCurr.Name = name
		menuItemCurr.Description = description
		menuItemCurr.Price = price
		menuItemCurr.Type = "MenuItem"
		menuItems = append(menuItems, menuItemCurr)
	}

	return menuItems, nil
}
