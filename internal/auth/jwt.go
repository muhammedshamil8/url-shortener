package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, email string, secret string, expiry string) (string, error) {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		duration = 24 * time.Hour
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(duration).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["user_id"].(float64)), nil
	}
	return 0, errors.New("invalid token")
}
