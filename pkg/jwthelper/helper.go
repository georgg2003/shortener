package jwthelper

import (
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("TOP_SECRET")

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID int64
}

func NewAccessToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		TokenClaims{
			UserID: userID,
		},
	)
	return token.SignedString(secretKey)
}

func ReadAccessToken(encodedToken string) (int64, error) {
	parser := jwt.NewParser()

	claims := TokenClaims{}
	token, err := parser.ParseWithClaims(encodedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return secretKey, jwt.NewValidationError(
				"signing method is not correct",
				jwt.ValidationErrorUnverifiable,
			)
		}

		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, jwt.ErrTokenNotValidYet
	}

	return claims.UserID, nil
}
