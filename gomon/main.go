package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Trainer struct {
	Name string
	Age int
	City string
}

func main(){
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatalf("error loading the .env file: %s", err)
	}

	mongoUri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("trainers")

	ash := Trainer{"Ash", 10, "pallet town"}
	misty := Trainer{"Misty", 10, "cerulean city"}
	brock := Trainer{"Brock", 15, "Pewter city"}

	//insert one
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil{
		log.Fatal(err)

	}
	fmt.Println("inserted a single document", insertResult.InsertedID)

	//insert many 
	trainers := []interface{}{misty, brock}
	insertManyResult, err := collection.InsertMany(context.TODO(),trainers)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)


	//update documents
	filter := bson.D{{"name", "Ash"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	//find a single document
	var result Trainer
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("found a single document: %+v\n", result)

	//find multiple documents
	findOptions := options.Find()
	findOptions.SetLimit(2)

	var results []*Trainer

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil{
		log.Fatal(err)
	}

	for cur.Next(context.TODO()){
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil{
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil{
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	fmt.Printf("found multiple documents (array of pointers): %+v\n", results )

	//Delete

	deleteResult, err  := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	err = client.Disconnect(context.TODO())
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("connection to mongodb closed")
}