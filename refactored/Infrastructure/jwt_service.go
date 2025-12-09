package Infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	Generate(userID int, username, role string) (string, error)
	Validate(tokenStr string) (*jwt.Token, error)
	ParseClaims(tokenStr string) (*CustomClaims, error)
}

type jwtService struct {
	secret []byte
}

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string) JWTService {
	return &jwtService{secret: []byte(secret)}
}

func (j *jwtService) Generate(userID int, username, role string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

func (j *jwtService) Validate(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
}

func (j *jwtService) ParseClaims(tokenStr string) (*CustomClaims, error) {
	t, err := j.Validate(tokenStr)
	if err != nil || !t.Valid {
		return nil, errors.New("invalid token")
	}
	if claims, ok := t.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	return nil, errors.New("cannot parse claims")
}
