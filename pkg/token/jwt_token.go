package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWTMaker struct {
	Key string
}

func NewJWTMaker(signingKey string) Maker {
	return &JWTMaker{
		Key: signingKey,
	}
}

func (m JWTMaker) Create(userToken UserToken, duration time.Duration) (string, error) {
	jwtToken := NewJWTToken(userToken, duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtToken)
	signedString, err := token.SignedString([]byte(m.Key))
	if err != nil {
		fmt.Println("can not sign token, err:", err)
		return "", err
	}
	return signedString, nil
}

func (m JWTMaker) Verify(tokenStr string) (*JWT, error) {
	jwtToken := &JWT{}
	parsedToken, err := jwt.ParseWithClaims(tokenStr, jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //check if the token was signed with the same signing Method we used to create tokens
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.Key), nil
	})
	if err == nil && parsedToken.Valid {
		claims, ok := parsedToken.Claims.(*JWT) //check if the token is of our type token
		if !ok {
			return nil, ErrInvalidToken
		}
		if err = claims.Valid(); err != nil {
			return nil, err
		}
		return jwtToken, nil
	}

	return nil, err
}
