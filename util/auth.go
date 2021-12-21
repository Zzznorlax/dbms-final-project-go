package util

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// custom claims
type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func NewToken(config *Config, userID int) (string, error) {

	now := time.Now()

	account := strconv.Itoa(userID)

	claims := Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			Audience:  account,
			ExpiresAt: now.Add(600 * time.Minute).Unix(),
			Id:        account + now.String(),
			IssuedAt:  now.Unix(),
			Issuer:    "ginJWT",
			NotBefore: now.Add(1 * time.Second).Unix(),
			Subject:   account,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(config.JWTSecret))

	if err != nil {
		return token, fmt.Errorf("generating token %w", err)
	}

	return token, err
}

func ValidateToken(config *Config, token string) (*Claims, error) {
	// parse and validate token for six things:
	// validationErrorMalformed => token is malformed
	// validationErrorUnverifiable => token could not be verified because of signing problems
	// validationErrorSignatureInvalid => signature validation failed
	// validationErrorExpired => exp validation failed
	// validationErrorNotValidYet => nbf validation failed
	// validationErrorIssuedAt => iat validation failed
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}

		err = &ErrUnauthorized{
			Reason: message,
		}

		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, err
	}

	return nil, &ErrUnauthorized{
		Reason: "invalid token",
	}
}
