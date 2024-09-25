package helpers

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectMongoDb() {

	connectionUrl := os.Getenv("DB_CONNECTION_URL")
	dbName := os.Getenv("DB_NAME")

	// serverApi := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(connectionUrl)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("monogdb successfully connected !!")
	DB = client.Database(dbName)
}
