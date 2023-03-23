package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	//"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/slovojoe/broker/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct{
	Rabbit amqp.Connection
}

func main(){
	 //Try connecting to rabbitmq
	 rabbitCon,err:=connectTomq()
	 if err!=nil{
		 log.Println("An error occured connecting to mq ",err)
		 os.Exit(1)
	 }
 
	 //Close connection once the main function finishes
	 defer rabbitCon.Close()
	 handlers.NewMQ(rabbitCon)
	 //Start listening for messages
	 log.Println("Listening for and consuming RabbitMQ messages...")
    
    
	

    muxRouter := mux.NewRouter().StrictSlash(true)

	//specify who's allowe to connect
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
	err = http.ListenAndServe(":10000", handler) //Uncomment this line when using docker
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}

func connectTomq()(*amqp.Connection, error){
    //Rabbitmq initially takes a while to start
    //For this, we'll try to connect a fixed number of times
    var counts int64
    var backOff= 1*time.Second
    var connection *amqp.Connection

    //Dont connect until rabbit is ready
    for {
        c,err:=amqp.Dial("amqp://guest:guest@rabbitmq")
        if err!=nil {
            fmt.Println("RabbitMQ not yet ready...")
            counts++
        }else{
            connection=c
			fmt.Println("RabbitMQ is now ready for you daddy")
            break
        }


        //If connection fails after 5 attempts, something is wrong
        if counts>5{
            fmt.Println(err)
            return nil,err
        }

        //if counts is less than 5 increase the delay duration each time count is >5
        backOff=time.Duration(math.Pow(float64(counts),2))*time.Second
        log.Println("backing off...")
        time.Sleep(backOff)
        continue
    }
    return connection,nil
}