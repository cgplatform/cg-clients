package user

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

func Create(user *User) (*User, error) {
	collection := database.GetCollection("users")
	if result, err := collection.InsertOne(ctx, user); err != nil {
		return nil, err
	} else {
		user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	}

	return user, nil
}

func Read(user User) ([]User, error) {
	collection := database.GetCollection("users")

	filter := mapper.ConvertStructToBSONMap(user, &mapper.MappingOpts{GenerateFilterOrPatch: true})

	var users []User

	if cursor, err := collection.Find(ctx, filter); err != nil {
		return nil, err
	} else {
		for cursor.Next(ctx) {
			var user User

			if err = cursor.Decode(&user); err != nil {
				return nil, err
			}

			users = append(users, user)
		}
	}

	return users, nil
}

func FindById(id string) (*User, error) {
	collection := database.GetCollection("users")

	objectID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objectID}

	var user *User

	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByTokenAndAlias(alias string, user *User, token string) (*User, error) {
	collection := database.GetCollection("users")

	objectID, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{"_id": objectID}

	filter["tokens."+alias] = token

	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateByUser(user *User) (*User, error) {
	objectID, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{"_id": objectID}

	update := mapper.ConvertStructToBSONMap(user, &mapper.MappingOpts{GenerateFilterOrPatch: true, RemoveID: true})

	update = bson.M{
		"$set": update,
	}

	return Update(filter, update)
}

func UpdateTokenByAlias(alias string, user *User, token string) (*User, error) {
	documentSet := bson.M{}

	objectID, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{"_id": objectID}

	documentSet["tokens."+alias] = token

	update := bson.M{
		"$set": documentSet,
	}

	return Update(filter, update)
}

func Update(filter bson.M, update bson.M) (*User, error) {
	collection := database.GetCollection("users")

	var updated *User

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

func Delete(user *User) (*DeleteResponse, error) {
	collection := database.GetCollection("users")
	objectID, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{"_id": objectID}

	if result, err := collection.DeleteOne(ctx, filter); err != nil {
		return nil, err
	} else {
		response := &DeleteResponse{ID: user.ID, DeletedCount: result.DeletedCount}
		return response, nil
	}

}

func DeleteTokenByAlias(alias string, user *User) (*User, error) {
	documentUnset := bson.M{}

	objectID, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{"_id": objectID}

	documentUnset["tokens."+alias] = 1

	update := bson.M{
		"$unset": documentUnset,
	}

	return Update(filter, update)
}

func TryLogin(login LoginRequest) (bool, *User) {

	collection := database.GetCollection("users")

	filter := bson.M{"email": login.Email}

	var user User
	err := collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return false, nil
	}

	if user.Password == login.Password {
		return true, &user
	}

	return false, nil

}
