package handlers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type customJWTClaim struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

func createToken(userID string, secret string, lifespanHours int) (string, error) {
	claims := customJWTClaim{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(lifespanHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func parseToken(tokenString string, secret string) (*customJWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&customJWTClaim{},
		func(t *jwt.Token) (any, error) { return []byte(secret), nil },
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*customJWTClaim)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
