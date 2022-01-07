package database

import (
	"context"
	"fmt"
	"time"

	"s2p-api/config"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
)

func Connect() error {
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/?tlsInsecure=true",
		config.Mongo.User,
		config.Mongo.Password,
		config.Mongo.Host)

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if client, err := mongo.Connect(ctx, clientOptions); err != nil {
		log.Fatalln("[MongoDB]: Could not connect -", err)
		return err
	} else {
		log.Infoln("[MongoDB]: Connection to MongoDB was a success")
		mongoClient = client
	}

	return nil
}

func GetCollection(collection string) *mongo.Collection {
	return mongoClient.Database(config.Mongo.Database).Collection(collection)
}
