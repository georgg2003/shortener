package jwthelper

import (
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("TOP_SECRET")

type tokenClaims struct {
	jwt.RegisteredClaims
	userID int64
}

func NewAccessToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		tokenClaims{
			userID: userID,
		},
	)
	return token.SignedString(secretKey)
}

func ReadAccessToken(encodedToken string) (int64, error) {
	parser := jwt.NewParser()

	claims := tokenClaims{}
	token, err := parser.ParseWithClaims(encodedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method == jwt.SigningMethodHS256 {
			return secretKey, jwt.NewValidationError(
				"signing method is not correct",
				jwt.ValidationErrorMalformed,
			)
		}

		return token, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, jwt.ErrTokenNotValidYet
	}

	return claims.userID, nil
}
