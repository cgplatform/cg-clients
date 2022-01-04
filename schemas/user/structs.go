package user

type User struct {
	ID        string `bson:"_id,omitempty" gql:"id"`
	Name      string `bson:"name"`
	Email     string `bson:"email"`
	Phone     string `bson:"phone"`
	Password  string `bson:"password"`
	Birthdate string `bson:"birthdate"`
}

var UserInstance = User{}

type LoginRequest struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

var LoginRequestInstance = LoginRequest{}

type LoginResponse struct {
	Token string `bson:"token"`
}

var LoginResponseInstance = LoginResponse{}

type MessageResponse struct {
	Message string `bson:"message"`
}

var MessageResponseInstance = MessageResponse{}
