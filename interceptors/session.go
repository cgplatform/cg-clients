package interceptors

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

func IsLoggedIn(request interface{}, session jwt.MapClaims) (bool, error) {
	if session == nil {
		return false, errors.New("not authorized")
	} else {
		return true, nil
	}
}
