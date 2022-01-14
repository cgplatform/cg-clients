package game

import (
	"context"
	"s2p-api/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
