package game

import (
	"context"
	"s2p-api/database"

	"github.com/naamancurtis/mongo-go-struct-to-bson/mapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx = context.Background()
)

func Create(game *Game) (*Game, error) {
	collection := database.GetCollection("games")
	if result, err := collection.InsertOne(ctx, game); err != nil {
		return nil, err
	} else {
		game.ID = result.InsertedID.(primitive.ObjectID).Hex()
	}

	return game, nil
}

func UpdateByGame(game *Game) (*Game, error) {
	objectID, _ := primitive.ObjectIDFromHex(game.ID)

	filter := bson.M{"_id": objectID}

	update := mapper.ConvertStructToBSONMap(game, &mapper.MappingOpts{GenerateFilterOrPatch: true, RemoveID: true})

	update = bson.M{
		"$set": update,
	}

	return Update(filter, update)
}

func Update(filter bson.M, update bson.M) (*Game, error) {
	collection := database.GetCollection("games")

	var updated *Game

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	result := collection.FindOneAndUpdate(ctx, filter, update, &opt)

	if err := result.Decode(&updated); err != nil {
		return nil, err
	}

	return updated, nil
}

func Read(game Game) ([]Game, error) {
	collection := database.GetCollection("games")

	filter := mapper.ConvertStructToBSONMap(game, &mapper.MappingOpts{GenerateFilterOrPatch: true})

	var games []Game

	if cursor, err := collection.Find(ctx, filter); err != nil {
		return nil, err
	} else {
		for cursor.Next(ctx) {
			var game Game

			if err = cursor.Decode(&game); err != nil {
				return nil, err
			}

			games = append(games, game)
		}
	}

	return games, nil
}
