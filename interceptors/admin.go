package interceptors

import (
	"s2p-api/exceptions"

	"github.com/dgrijalva/jwt-go"
)

func IsAdmin(request interface{}, session jwt.MapClaims, key string) (bool, error) {

	if key == "" {
		return false, exceptions.USER_NOT_AUTHORIZED
	}
	if key != "s2p-awl-key" {
		return false, exceptions.USER_NOT_AUTHORIZED
	}
	return true, nil
}
