package interceptors

import (
	"s2p-api/exceptions"

	"github.com/dgrijalva/jwt-go"
)

func IsAdmin(request interface{}, session jwt.MapClaims) (bool, error) {

	if session["type"] != "admin" {
		return false, exceptions.USER_NOT_AUTHORIZED
	}
	return true, nil
}
