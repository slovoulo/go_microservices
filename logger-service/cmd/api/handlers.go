package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slovojoe/logger-service/data"
)

type jsonResponse struct{
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

type logJSONPayload struct{
    Name string `json:"name"`
    Data string `json:"data"`
}
func  WriteLog(w http.ResponseWriter, r *http.Request){
    // get the body of the  POST request
	// unmarshal this into a new logJSONPayload struct
	log.Println("Starting log processsssssss")
	reqBody, _ := ioutil.ReadAll(r.Body)
	log.Printf("Starting writelog process with %s",reqBody)
	var loginput logJSONPayload
	json.Unmarshal(reqBody, &loginput)
	log.Printf("Unmarshalled body is %s",loginput)


	//Insert data
	log.Println("Starting log proces")
	event:= data.LogEntry{
		Name: loginput.Name,
		Data:loginput.Data,
	}

	err:=data.Insert(event)
	if err!=nil {
		log.Println("An error occured inserting log ",err)
		return
	}
	log.Println("Successfully submitted a log")

	
} 