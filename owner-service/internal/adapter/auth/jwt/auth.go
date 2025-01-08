package jwt

import (
	"errors"
	"fmt"
	"time"

	"owner-service/internal/adapter/config"
	"owner-service/internal/core/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JwtToken struct {
	Duration  time.Duration
	SecretKey string
}

func New(config *config.TokenConfiguration) *JwtToken {
	// Parse the token duration
	tokenDuration, err := time.ParseDuration(config.Duration)
	if err != nil {
		tokenDuration = 1 * time.Hour
	}

	return &JwtToken{
		Duration:  tokenDuration,
		SecretKey: config.Secret,
	}
}

func (jt *JwtToken) CreateToken(id, email string, role string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jt.Duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    config.GetConfig().App.Name,
		Subject:   email + "," + role,
		ID:        id,
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(jt.SecretKey))
}

func (jt *JwtToken) VerifyToken(tokenString string) (domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Validate issuer
		if iss, err := token.Claims.GetIssuer(); iss != config.GetConfig().App.Name || err != nil {
			return nil, fmt.Errorf("unknown issuer: %v", token.Header["iss"])
		}

		return []byte(config.GetConfig().Token.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return domain.Claims{}, errors.New(domain.ErrExpiredToken.Error())
		} else {
			return domain.Claims{}, errors.New("invalid token")
		}
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !(ok && token.Valid) {
		return domain.Claims{}, errors.New("invalid token")
	}

	return domain.Claims{
		Email:  claims.Subject,
		Issuer: claims.Issuer,
		ID:     claims.ID,
	}, nil
}
