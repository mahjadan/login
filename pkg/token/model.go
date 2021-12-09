package token

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

const ISSUER = "back-end"

var (
	ErrTokenExpired = errors.New("expired token")
	ErrInvalidToken = errors.New("invalid token")
)

type JWT struct {
	Username string
	Roles    []string
	jwt.RegisteredClaims
}

func (j *JWT) Valid() error {
	if err := j.RegisteredClaims.Valid(); err != nil {
		return ErrTokenExpired
	}
	if j.RegisteredClaims.VerifyIssuer(ISSUER, true) {
		return ErrInvalidToken
	}
	return nil
}

func NewJWTToken(userToken UserToken, duration time.Duration) JWT {
	return JWT{
		Username: userToken.Username,
		Roles:    userToken.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ISSUER,
			Subject:   userToken.Username,
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        uuid.NewV4().String(),
		},
	}
}

type UserToken struct {
	Username string
	Roles    []string
}
