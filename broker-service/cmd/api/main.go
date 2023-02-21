package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main(){
    muxRouter := mux.NewRouter().StrictSlash(true)
	router := AddRoutes(muxRouter)
	//err := http.ListenAndServe(constants.CONN_HOST+":"+constants.CONN_PORT, router)//Uncomment this line when using localhost
	err := http.ListenAndServe(":10000", router) //Uncomment this line when using docker
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}