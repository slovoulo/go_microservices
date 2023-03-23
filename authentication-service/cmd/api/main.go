package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/slovojoe/authentication-service/cmd/database"
)

func main(){
    //Connect to postgres db
    database.ConnectDB()
    log.Println("Starting auth service")
    muxRouter := mux.NewRouter().StrictSlash(true)

	//specify who's allowed to connect
	c:=cors.New(cors.Options{ 
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
})
	router := AddRoutes(muxRouter)
	handler := c.Handler(router)
	
	//err := http.ListenAndServe(constants.CONN_HOST+":"+constants.CONN_PORT, router)//Uncomment this line when using localhost
	err := http.ListenAndServe(":8081", handler) //Uncomment this line when using docker
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}

	log.Println("Auth service started")
}