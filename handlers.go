package main

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github/Simple-Restaurant-API/models"
	"net/http"
	"reflect"
)

const (
	ContentType = "application/json;charset=UTF-8"
)

// ======================================HANDLERS========================================================

// handles the home route
func handleIndex(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	// fmt.Fprint(writer, "this is the home page ", request.URL.Path)
	writer.Header().Set("content-type", ContentType)
	var s string
	s = "MytestString"
	m := make(map[string]interface{})
	m["test1"] = 53
	m["Name"] = "ashani"
	m[s] = "testing 1,3"
	json.NewEncoder(writer).Encode(m)
}

//handles the /restaurant/:resaurantId GET route
func handleRestaurant(writer http.ResponseWriter, request *http.Request, parms httprouter.Params) {

	searchName := parms.ByName("restaurantSearchName") //retrieve the searchName
	// restaurant := models.Restaurant{Id: 1, Name: ""}
	restaurant := models.NewRestaurant()

	conf, err := restaurant.FindByField("searchName", searchName) // find the restaurant by its searchName
	if err != nil {                                               // if it isn't found return an error
		errorResponse := formatErrorResponse("error retrieving resource", http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	//if the restaurant was found display it
	if conf {
		writer.Header().Set("Content-Type", ContentType)
		writer.WriteHeader(http.StatusOK)
		// x, _ := json.Marshal(restaurant)
		// fmt.Fprintf(writer, "the restaurant name is: %s", x)
		response, _ := formatResourceResponse(restaurant)
		json.NewEncoder(writer).Encode(response)
		// fmt.Println("the param ", restaurant)
		// json.NewEncoder(writer).Encode(restaurant)
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
	// var err error
	// var restaurants []interface{}

	restaurants, err := res.FindAll() //find all the restaurants
	if err != nil {
		//format error response here 500 internal server error
		writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
		writer.WriteHeader(http.StatusInternalServerError)
		// var errResponse models.ErrorResponse
		// errResponse.Status = http.StatusInternalServerError
		// errResponse.Title = "error retrieving records"
		// errResponse.Details = "There was an error retrieving the list of restaurants"
		errResponse := formatErrorResponse("error retrieving resources", http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errResponse)
		return
	}
	//set the status and content type
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	// return the list of restaurants
	// json.NewEncoder(writer).Encode(restaurants)

	resList := make([]interface{}, len(restaurants))
	for v, t := range restaurants {
		resList[v] = t
	}

	resourceListResponse, _ := formatResourceListResponse(resList)
	json.NewEncoder(writer).Encode(resourceListResponse)
}

func formatErrorResponse(errorText string, statusCode int) models.ErrorResponse {
	var errorResponse = models.ErrorResponse{}
	errorResponse.Status = statusCode
	errorResponse.Title = errorText

	return errorResponse
}

//accepts a resource, formats it and returns the JSONAPI formatted ResourceResponse
func formatResourceResponse(resource interface{}) (models.FormattedResource, error) {

	formattedResource := models.FormattedResource{}
	formattedResource.Attributes = make(map[string]interface{})
	resourceResponse := models.ResourceResponse{}
	var resourceType string
	var resourceId int

	//get reflection of resource
	resourceRef := reflect.ValueOf(resource)
	// resourceRefType := resourceRef.Elem().Type() // get the pointer for the resource and its t
	// check if the resource is a struct before proceeeding
	if kind := resourceRef.Kind().String(); kind != "struct" {
		fmt.Println(kind)
		fmt.Println("this is an error")
		return formattedResource, errors.New("The resource to be formatted must be a struct")
	}

	// resourceRefType := resourceRef.Elem().Type() // get the pointer for the resource and its type
	resourceRefType := resourceRef.Type()
	//loop through resource Fields and add them to the attributes map
	//should be in the form "[field name] = field value"
	for i := 0; i < resourceRef.NumField(); i++ {
		field := resourceRefType.Field(i)
		formattedResource.Attributes[field.Name] = resourceRef.Field(i).Interface()

		//check for resource type and id
		if field.Name == "Type" {
			resourceType = resourceRef.Field(i).Interface().(string)
		} else if field.Name == "Id" {
			resourceId = resourceRef.Field(i).Interface().(int)
		}
	}
	//ensure the type and id get set outsite the attributes
	formattedResource.Type = resourceType
	formattedResource.Id = resourceId

	//create the resourceResponse and populate its Data
	resourceResponse.Data[0] = formattedResource

	return formattedResource, nil
}

func formatResourceListResponse(resourceList []interface{}) (models.ResourceListResponse, error) {
	resourceListResponse := models.ResourceListResponse{}

	for i := 0; i < len(resourceList); i++ {
		resourceResponse, err := formatResourceResponse(resourceList[i])

		if err != nil {
			fmt.Println("There was an error formatting the response - ", err)
			return resourceListResponse, err
		}
		resourceListResponse.Data = append(resourceListResponse.Data, resourceResponse)
	}

	return resourceListResponse, nil
}
