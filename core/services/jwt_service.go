package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtService struct {
	secretKey string
	isr       string
}

func NewJWTService() *jwtService {
	return &jwtService{
		secretKey: "s2play-awl",
		isr:       "s2p-api",
	}
}

type Claim struct {
	Sum string `bson:"sum"`
	jwt.StandardClaims
}

func (s *jwtService) GenerateToken(id string, duration time.Duration) (string, error) {
	claim := &Claim{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    s.isr,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := token.SignedString([]byte(s.secretKey))

	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *jwtService) ValidateToken(token string) jwt.MapClaims {

	keyFunction := func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(s.secretKey), nil
	}

	tk, err := jwt.Parse(token, keyFunction)

	if err != nil {
		return nil
	}
	claims := tk.Claims.(jwt.MapClaims)

	return claims
}
