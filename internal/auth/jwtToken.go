package auth

import (
	"math/rand"
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
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"role": fields.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(key))
}

func GenerateRefreshToken() (string, error) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	data := make([]rune, 90)
	for i := range data {
		data[i] = letters[r.Intn(len(letters))]
	}

	return string(data), nil
}
