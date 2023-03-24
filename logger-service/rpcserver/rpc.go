package rpcserver

import (
	"context"
	"log"
	"time"

	"github.com/slovojoe/logger-service/data"
	
)

//When working with an RPC  server first setup a type
 type RPCServer struct{}

//Define the kind of payload to be received from RPC
type RPCPayload struct {
    Name string
    Data string
}


//Will be using RPC to write logs to mongoDb
func (r *RPCServer) LogInfo (payload RPCPayload, rpcresponse *string)error{
collection:= data.Client.Database("logs").Collection("rpcLogs")
_,err :=collection.InsertOne(context.TODO(),data.LogEntry{
    Name:payload.Name,
    Data:payload.Data,
    CreatedAt: time.Now(),
})
if err !=nil {
    log.Printf("RPC error writing to mongo %s", err)
    return err

}
*rpcresponse= "Processed" + payload.Name + "payload via RPC:" 
return nil
}