package user

import (
	"context"
	"s2p-api/database"

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

func Read(filter bson.M) ([]User, error) {
	collection := database.GetCollection("users")

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

func Update(user *User) (*User, error) {
	collection := database.GetCollection("users")

	objectID, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{"_id": objectID}

	document := bson.M{}

	if user.Name != "" {
		document["name"] = user.Name
	}

	if user.Email != "" {
		document["email"] = user.Email
	}

	if user.Phone != "" {
		document["phone"] = user.Phone
	}

	if user.Birthdate != "" {
		document["birthdate"] = user.Birthdate
	}

	update := bson.M{
		"$set": document,
	}

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

func FilterFrom(args map[string]interface{}) bson.M {
	result := make(bson.M, len(args))

	for k, v := range args {

		if k == "id" {
			k = "_id"
			v, _ = primitive.ObjectIDFromHex(v.(string))
		}

		result[k] = v
	}

	return result
}

func TryLogin(email string, password string) (bool, string) {

	collection := database.GetCollection("users")

	filter := bson.M{"email": email}

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return false, ""
	}

	id := result["_id"].(primitive.ObjectID).Hex()

	if result["password"].(string) == password {
		return true, id
	}

	return false, ""

}
