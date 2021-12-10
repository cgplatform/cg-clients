package user

type User struct {
	ID        string `bson:"_id,omitempty" gql:"id"`
	Name      string `bson:"name"`
	Mail      string `bson:"mail"`
	Phone     string `bson:"phone"`
	Password  string `bson:"password"`
	Birthdate string `bson:"birthdate"`
}

var UserInstance = User{}
