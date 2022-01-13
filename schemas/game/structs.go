package game

type Game struct {
	ID          string `bson:"_id,omitempty" gql:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
	Developer   string `bson:"developer"`
	PlatForm    string `bson:"platform"`
	Category    string `bson:"category"`
}

var GameInstance = Game{}
