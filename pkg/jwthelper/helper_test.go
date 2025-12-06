package jwthelper

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

var testUserID int64 = 1

func TestNewAccessToken(t *testing.T) {
	token, err := NewAccessToken(testUserID)
	assert.NoError(t, err)

	gotUserID, err := ReadAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, testUserID, gotUserID)
}

func TestInvalidSigningMethod(t *testing.T) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		TokenClaims{
			UserID: testUserID,
		},
	)
	encodedToken, err := token.SignedString(secretKey)
	assert.NoError(t, err)

	_, err = ReadAccessToken(encodedToken)

	var tgt *jwt.ValidationError
	assert.ErrorAs(t, err, &tgt)
	assert.True(t, tgt.Is(jwt.ErrTokenUnverifiable))
}
