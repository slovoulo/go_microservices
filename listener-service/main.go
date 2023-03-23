package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/slovojoe/listener-service/event"
    
)

func main(){

    //Try connecting to rabbitmq
    rabbitCon,err:=connectTomq()
    if err!=nil{
        log.Println("An error occured connecting to mq ",err)
        os.Exit(1)
    }

    //Close connection once the main function finishes
    defer rabbitCon.Close()
    

    //Start listening for messages
    log.Println("Listening for and consuming RabbitMQ messages...")
    
    
    //Create consumer
    consumer, err:=event.NewConsumer(rabbitCon)
    if err!=nil{
        panic(err)
    }

    //watch the queue and consume events
    err=consumer.Listen([]string{"log.INFO","log.WARNING", "log.ERROR"})
    if err !=nil{
        log.Println(err)
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