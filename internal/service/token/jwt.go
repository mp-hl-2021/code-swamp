package token

import (
	"github.com/dgrijalva/jwt-go"

	"crypto/rsa"
	"errors"
	"fmt"
	"time"
)

type JwtHandler struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	expire time.Duration
}

type Claims struct {
	Id uint
	jwt.StandardClaims
}

func NewJwt(privateBytes, publicBytes []byte, keyExpiration time.Duration) (*JwtHandler, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return nil, err
	}
	return &JwtHandler{
		publicKey:  publicKey,
		privateKey: privateKey,
		expire:     keyExpiration,
	}, nil
}

func (j JwtHandler) IssueToken(userId uint) (string, error) {
	claims := Claims{
		Id: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.expire).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j JwtHandler) UserIdByToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected token signing method")
		}
		return j.publicKey, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	return claims.Id, nil
}