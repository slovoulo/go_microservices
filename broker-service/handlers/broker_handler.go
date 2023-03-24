package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/slovojoe/broker/event"
)

//The broker-service handlers are responsible for handling all requests coming
//to & from the other microservices

//Define a default payload struct
type RequestPayload struct{
	//Action is a string that we'll use to identify the type of 
	//Microservice the client wants
	Action string `json:"action"`
	//Payload to use when calling the authenitication microservice
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
}

//authentication microservice payload
type AuthPayload struct {
	Email string `json:"email"`
	Password string	`json:"password"`
}

//Logger microservice payload
type LogPayload struct {
	Name string	`json:"name"`
	Data string	`json:"data"`
}


var rabbit *amqp.Connection

func NewMQ(mq *amqp.Connection) {
    rabbit=mq

    
}



func HandleSubmission(w http.ResponseWriter, r *http.Request){
	// get the body of the  post request
	// unmarshal this into a new Request payload  struct
	
	reqBody, _ := ioutil.ReadAll(r.Body)
	var requestPayload RequestPayload
	json.Unmarshal(reqBody, &requestPayload)

	//Check the action string of the payload to determine which
	//service is being requested
	switch requestPayload.Action {
	case "auth":
		Authenticate(w,requestPayload.Auth)
	// case "log":
	// 	logEventViaRabbit(w,requestPayload.Log)
	case "log":
		logItemViaRPC(w,requestPayload.Log)

	default:
		log.Println("Unknown action/request")
		
	}


}

func LogItem(w http.ResponseWriter,  l LogPayload){
	//////Create a json to send to the log microservice///////
	//Marshal the request struct into a json
	jsonData,_:=json.Marshal(l)
	log.Printf("got json data as %s",jsonData)
	logServiceURL:= "http://logger-service/writelog"

	request,err := http.NewRequest("POST", logServiceURL,bytes.NewBuffer(jsonData))
	if err!=nil{
		log.Printf("Broker service couldnt connect to the logger service, %s",err)}
	request.Header.Set("Content-Type", "application/json")
	client:=&http.Client{}
	response, err := client.Do(request)
	if err!=nil{
		log.Printf("Broker service client error: , %s",err)

	}
	defer response.Body.Close()

		//////Make sure to get back the correct status code///////
		if response.StatusCode != http.StatusAccepted{
			log.Println("Error authenticating user")
			return
	
		}
	
		//Create a variable we'll read response.body (from log service ) into
		var jsonFromLogService jsonResponse
		//Decode the json from the log service
		err=json.NewDecoder(response.Body).Decode(&jsonFromLogService)
		if err!=nil{
			log.Printf("Error decoding json from log service: , %s",err)
	
		}
		var payload =jsonResponse{Error: false,Message: "logged",Data:jsonFromLogService.Data}
		json.Marshal(payload)

}

func Authenticate(w http.ResponseWriter,  a AuthPayload){
	//////Create a json to send to the auth microservice///////
	//Marshal the request struct into a json
	jsonData,_:=json.Marshal(a)
	log.Printf("got jsonm data as %s",jsonData)


	//////Call the auth microservice///////
	//http://{name of the authentication service in docker-compose}/{name of url we want from the auth service}
	request,err := http.NewRequest("POST", "http://authentication-service:8081/authenticateuser",bytes.NewBuffer(jsonData))
	if err!=nil{
		log.Printf("Broker service couldnt connect to the auth service, %s",err)
		//http.Error(w, "Internal server error", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Broker service couldnt connect to the auth service"))
		return

	}

	client:=&http.Client{}
	response, err := client.Do(request)
	if err!=nil{
		log.Printf("Broker service client error: , %s",err)
	//	http.Error(w, "Internal server error", http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Invalid password \n"))
		return

	}
	defer response.Body.Close()

	//////Make sure to get back the correct status code///////
	if response.StatusCode == http.StatusUnauthorized{
		b, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Invalid credentials")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(string(b)))
		return
	}else if response.StatusCode != http.StatusOK{
		log.Println("Error authenticating user error code ",response.StatusCode)
		w.WriteHeader(response.StatusCode)
		w.Write([]byte("Error authenticating user error code"))
		w.Write([]byte(string(rune(response.StatusCode))))
		return

	}else if response.StatusCode == http.StatusOK{
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User authenticated successfully"))
	}

	//Create a variable we'll read response.body (from auth service ) into
	var jsonFromAuthService jsonResponse
	//Decode the json from the auth service
	err=json.NewDecoder(response.Body).Decode(&jsonFromAuthService)
	if err!=nil{
		log.Printf("Error decoding json from auth service: , %s",err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}
	var payload =jsonResponse{Error: false,Message: "User Authenticated",Data:jsonFromAuthService.Data}
	json.NewEncoder(w).Encode(payload)

}


func logEventViaRabbit(w http.ResponseWriter,l LogPayload){
	
	err := pushToQueue(l.Name, l.Data)
	if err != nil {
		log.Println("An error occured pushing to queue",err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
		w.Write([]byte("Log submitted successfully"))

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	log.Printf("Payload is %s",payload)

	json.Marshal(payload)
}

func pushToQueue(name,msg string )error{
	log.Println("starting push operation")
	emitter, err := event.NewEventEmitter(rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.Marshal(&payload)
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}


type RPCPayload struct{
	Name string
	Data string
}
func logItemViaRPC(w http.ResponseWriter, l LogPayload){
	client,err:=rpc.Dial("tcp", "logger-service:5001")
	if err!=nil{
		log.Println("Broker RPC dial error", err)
	}

	//Create a payload with a type that matches what the remote RPC server expects
	rpcPayload:= RPCPayload{
		Name:l.Name,
		Data:l.Data,
	}

	//Declare a variable to store the rpcPayload result
	//It will be populated by the remote RPC call
	var result string

	//Make a call to the RPC server we created in rpc.go
	//The name should match the name of the server struct in this case RPCServer
	//Log.Info is the method we are calling also defined in rpc.go
	err=client.Call("RPCServer.LogInfo",rpcPayload,&result)
	if err!=nil{
		log.Println("RPCServer call error", err)
	}
	//Return a json to user
	payload:=jsonResponse{
		Error: false,
		Message: result,
	}

	json.Marshal(payload)
}































type jsonResponse struct{
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}


func HomeHandler(w http.ResponseWriter, r *http.Request){
	payload:=jsonResponse{
		Error: false,
		Message: "Welcome to Go broker microservice",
	}

	out,_:=json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)

	

}
// func HomeHandler(w http.ResponseWriter, r *http.Request){
	

// 	w.Write([]byte("Welcome to Go microservices!"))

// }

