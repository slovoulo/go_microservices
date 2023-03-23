package event

import (
	

	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"

	
)

//The type used for receiving amqp events from the queue
type Consumer struct {
    conn *amqp.Connection
    queueName string
}

//Create an instance of a consumer
func NewConsumer(conn *amqp.Connection)(Consumer,error){
	log.Println("Starting to newconsumer")
    consumer :=Consumer{
        conn: conn,
    }

    err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

//A function that listens to the queue for specific topics
func (consumer *Consumer) Listen(topics []string) error {
	log.Println("Starting to listen")
    //Get the channel from consumer struct
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
    //close the consumer when done
	defer ch.Close()

    //Declare and store a random queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}


    //Loop through the topics slice
    //Bind the channel to each topic
	for _, s := range topics {
		ch.QueueBind(
			q.Name, 
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

    //Look for messages coming from rabbitmq
    
    messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

    //Keep consuming everything coming from rabbitmq until the program is closed
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

//A function that takes an action based on the nAME of an event received from the queue
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
        //If the payload name is "log" or "event"
		//Call the logger service to log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// authenticate

	// you can have as many cases as you want, as long as you write the logic

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	//////Create a json to send to the log microservice///////
	//Marshal the request struct into a json
	jsonData,_:=json.Marshal(entry)
	log.Printf("log event got json data as %s",jsonData)
	logServiceURL:= "http://logger-service:9090/writelog"

	request,err := http.NewRequest("POST", logServiceURL,bytes.NewBuffer(jsonData))
	if err!=nil{
		log.Printf("Broker service couldnt connect to the logger service, %s",err)
    return err}

	request.Header.Set("Content-Type", "application/json")

	client:=&http.Client{}

	response, err := client.Do(request)
	if err!=nil{
		log.Printf("Broker service client error: , %s",err)
        return err
	}
	defer response.Body.Close()

		//////Make sure to get back the correct status code///////
		if response.StatusCode != http.StatusAccepted{
			log.Println("Error authenticating user")
			return err
	
		}
	
	

        return nil
}

