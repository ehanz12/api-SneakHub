package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// fungsi untuk generate token JWT
func GenerateJWT(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("not found JWT_SECRET in .env !")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)), // token berlaku selama 7 hari
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// fungsi untuk validasi token JWT
func ValidateJWT(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("not found JWT_SECRET in .env !")
	}
	// parsing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	// validasi token
	if err != nil {
		return "", err
	}
	// ambil user_id dari claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return userID, nil
	}
	return "", errors.New("invalid token")
}
