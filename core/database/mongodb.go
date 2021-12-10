package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"es-api/config"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
)

func Connect() error {

	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s",
		url.QueryEscape(config.Mongo.User),
		url.QueryEscape(config.Mongo.Password),
		url.QueryEscape(config.Mongo.Host))

	log.Infoln("Connecting to MongoDB", "[", config.Mongo.Host, "]")

	if client, err := mongo.NewClient(options.Client().ApplyURI(uri)); err != nil {
		log.Fatalln("[MongoDB]: Could not connect -", err)
		return err
	} else {
		mongoClient = client
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err := mongoClient.Connect(ctx); err != nil {
		log.Fatalln("[MongoDB]: Could not connect -", err)
		return err
	} else {
		log.Infoln("[MongoDB]: Connection to MongoDB was a success")
	}

	return nil
}

func GetCollection(collection string) *mongo.Collection {
	return mongoClient.Database(config.Mongo.Database).Collection(collection)
}
