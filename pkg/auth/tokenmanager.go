package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"log"
	"time"
)

var singingKey = viper.GetString("key.jwt")

type TokenManager struct {
}

func NewTokenManager() TokenManager {
	return TokenManager{}
}

type TokenClaims struct {
	*jwt.RegisteredClaims
	UserId   int    `json:"user_id"`
	UserRole string `json:"user_role"`
}

func (t *TokenManager) GenerateToken(id int, role string) (string, error) {
	if id == 0 || role == "" {
		return "", errors.New("unauthorized")
	}

	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAccess := jwt.NewNumericDate(time.Now().Add(60 * 24 * 365 * time.Minute))

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAccess,
		},
		id, role,
	})
	accessToken, err := accessClaims.SignedString([]byte(singingKey))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (t *TokenManager) ValidateToken(accessToken string) (int, string, error) {

	Token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid access-token")
		}
		return []byte(singingKey), nil
	})
	if err != nil {
		log.Print(err.Error())
	}

	claims, ok := Token.Claims.(*TokenClaims)
	if !ok {
		return 0, "none", errors.New("invalid token")
	}

	return claims.UserId, claims.UserRole, nil
}
