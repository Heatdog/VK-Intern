package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenFileds struct {
	ID   string
	Role string
}

func GenerateToken(fields TokenFileds, key string) (string, error) {
	payload := jwt.MapClaims{
		"sub":  fields.ID,
		"exp":  time.Now().Add(time.Hour * 2).Unix(),
		"role": fields.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString(key)
}
