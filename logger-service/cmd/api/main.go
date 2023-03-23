package main

import (
	"context"
	"os"

	"log"

	"net/http"

	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/slovojoe/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


const (
    webPort="9090"
    rpcPort= "5001"
    mongoURL="mongodb://mongodb:27017"
    grpcPort="50001"
   
)


var client *mongo.Client


type Config struct{
    Models data.Models

}


func main(){
   
      //Connect to mongo
      mongoClient,err:=connectToMongo()
      if err!=nil{
        log.Printf("The error is %s", err)
          log.Panic(err)
      }
      client=mongoClient
      
      //Create a context that mongo needs in order to disconnect
      ctx,cancel:=context.WithTimeout(context.Background(), 15*time.Second)
     // ctx,cancel:=context.WithTimeout(context.Background(), 15*time.Second)
      defer cancel()
      err = client.Ping(ctx, readpref.Primary())
    
  
      //close connection
      defer func ()  {
          if err =client.Disconnect(ctx); err!=nil{
              panic(err)
          }
      }() //last two brackets immediately call this function after creating it
  
         // app:=Config{Models:data.New(client)}

         data.New(client)

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
    log.Println("Service stratring at o  port ",webPort)

    sterr := http.ListenAndServe(":9090", handler) //Uncomment this line when using docker
	if sterr != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}

    log.Println("Service started at port ",webPort)


    
  

}

// func (app *Config)serve(){
//     srv:=&http.Server{
//         Addr: fmt.Sprintf(":%s",webPort),
//     }
// }

func connectToMongo()(*mongo.Client,error){
    mongoUsername := os.Getenv("MONGOUSERNAME")
    mongoPassword := os.Getenv("MONGOPASSWORD")

    log.Printf("username is %s and password is %s",mongoUsername,mongoPassword)
    //create connection options
    clientOptions:=options.Client().ApplyURI(mongoURL)
    clientOptions.SetAuth(options.Credential{
        Username: mongoUsername,
        Password: mongoPassword,
    })

    //connect
    c,err:=mongo.Connect(context.TODO(),clientOptions)
    if err!=nil{
        log.Println("Error connecting to mongo",err)
        return nil,err
    }
    log.Println("Connected to mongo ")
    return c,nil
}