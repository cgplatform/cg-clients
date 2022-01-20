package interceptors

import (
	"fmt"
	"s2p-api/exceptions"

	"github.com/dgrijalva/jwt-go"
)

func IsAdmin(request interface{}, session jwt.MapClaims) (bool, error) {
	fmt.Printf("session: %v\n", session)
	if session["Kind"] != "administrator" {
		return false, exceptions.USER_NOT_AUTHORIZED
	}
	return true, nil
}
