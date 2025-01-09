package jwt

import (
	"errors"
	"fmt"
	"time"

	"product-service/internal/adapter/config"
	"product-service/internal/core/domain"

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

func (jt *JwtToken) VerifyToken(tokenString string) (domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Validate issuer
		if iss, err := token.Claims.GetIssuer(); iss != config.GetConfig().Token.Issuer || err != nil {
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
