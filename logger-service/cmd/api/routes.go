package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Create a Route struct defining all the parameters a route should have
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//type Routes []Route
//var h = handlers.New(database.Db)
// Define a slice of routes to handle all of the apps routing
var Routes = []Route{
   

	//Write a log 
	{
		Name:        "writeLog",
		Method:      "POST",
		Pattern:     "/writelog",
        HandlerFunc: WriteLog,
		
	},

	
}

//Loop through the specified routes
func AddRoutes(router *mux.Router) *mux.Router {
	for _, route := range Routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

			log.Println("added route ", route.Name)
	}
	return router
}