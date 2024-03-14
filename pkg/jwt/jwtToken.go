package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenFileds struct {
	jwt.StandardClaims
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

func VerifyToken(tokenString, key string) (*TokenFileds, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenFileds)
	if !ok {
		return nil, fmt.Errorf("token calims are not of type *TokenClaims")
	}

	return claims, nil
}

func GenerateRefreshToken(fields TokenFileds, key string, expire time.Duration) (string, error) {
	payload := jwt.MapClaims{
		"sub":  fields.ID,
		"exp":  time.Now().Add(expire).Unix(),
		"role": fields.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(key))
}
