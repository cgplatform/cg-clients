package interceptors

import (
	"s2p-api/exceptions"

	"github.com/dgrijalva/jwt-go"
)

func IsLoggedIn(request interface{}, session jwt.MapClaims, key string) (bool, error) {
	if session == nil {
		return false, exceptions.USER_NOT_AUTHORIZED
	} else {
		return true, nil
	}
}
