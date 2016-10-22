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

func (menuItem *MenuItem) populate(rows *sql.Rows) bool {
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

func (menuItem *MenuItem) populateList(rows *sql.Rows) ([]MenuItem, error) {
	menuItems := []MenuItem{}

	var id int
	var name string
	var description string
	var price string
	var resId int

	for rows.Next() {
		menuItemCurr := MenuItem{}

		err := rows.Scan(&id, &name, &description, &price, &resId)
		if err != nil {
			fmt.Println("Error reading menuItem row - ", err)
			return menuItems, err
		}

		menuItemCurr.Id = id
		menuItemCurr.Name = name
		menuItemCurr.Description = description
		menuItemCurr.Price = price
		menuItemCurr.Type = "MenuItem"
		menuItems = append(menuItems, menuItemCurr)
	}
	return menuItems, nil
}
