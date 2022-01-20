package services

import (
	"fmt"
	"s2p-api/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	Sum  string `bson:"sum"`
	Kind string `bson:"kind"`
	jwt.StandardClaims
}

func GenerateToken(id string, kind string, duration time.Duration) (string, error) {
	claim := &Claim{
		id,
		kind,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    config.Jwt.Isr,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := token.SignedString([]byte(config.Jwt.SecretKey))

	if err != nil {
		return "", err
	}

	return t, nil
}

func GenerateTokenDefault(id string, kind string) (string, error) {
	return GenerateToken(id, kind, time.Hour*12)
}

func ValidateToken(token string) jwt.MapClaims {
	keyFunction := func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(config.Jwt.SecretKey), nil
	}

	tk, err := jwt.Parse(token, keyFunction)

	if err != nil {
		return nil
	}
	claims := tk.Claims.(jwt.MapClaims)

	return claims
}
