package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client


func New(mongo *mongo.Client) Models{
    Client=mongo

    return Models {
        LogEntry:LogEntry{},
    }
}

type Models struct{
    LogEntry LogEntry
}

type LogEntry struct{
    ID string `bson:"_id,omitempty" json:"id,omitempty"`
    Name string `bson:"name" json:"name"`
    Data string `bson:"data" json:"data"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

//A function to insert a single record to mongodb
func  Insert(entry LogEntry)error{
	log.Printf("Attempting to insert %s", entry)
	//db:=client.Database("logs")
	//log.Printf("database  is  %s", db)
	
   collection:=Client.Database("logs").Collection("logsCollection")
	log.Printf("collection is  %s", collection)

    _,err :=collection.InsertOne(context.TODO(), LogEntry{
        Name: entry.Name,
        Data: entry.Data,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    })
    if err!=nil{
        log.Println("Error inserting new record into logs collection",err)
        return err
    }
	log.Println("insert successful")
    return nil
}

func  getAll () ([]*LogEntry,error) {
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := Client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}


func  GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := Client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func  DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := Client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func  Update(entry LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := Client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(entry.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: entry.Name},
				{Key: "data", Value: entry.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}