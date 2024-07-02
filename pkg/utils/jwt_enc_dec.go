package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
)

func CreateJWT(claims interface{}, secret string) (string, error) {
	jwtClaims := claims.(jwt.Claims)

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// Sign the token with a secret key
	secretKey := []byte(secret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, secret string, out interface{}) error {
	// Define the secret key
	secretKey := []byte(secret)

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// Use reflection to set the values in the output struct
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("claims are not of type jwt.MapClaims")
	}

	err = assignMapToStructFields(out, mapClaims)
	if err != nil {
		return err
	}

	return nil
}
