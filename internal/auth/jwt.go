package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

func GenerateToken(userID int, email string, secret string, expiry string) (string, error) {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		duration = 24 * time.Hour
	}
	claims := &models.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (int, error) {
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func ValidateToken(tokenString string, secret string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
