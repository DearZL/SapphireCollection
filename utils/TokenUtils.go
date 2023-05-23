package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
)

type TokenUtils interface {
	GToken() (string, error)
	PToken(string string) *jwt.Token
	MapClaims(token *jwt.Token)
}

type UserClaims struct {
	UserId    string `json:"userId"`
	UserGroup string `json:"userGroup"`
	UserEmail string `json:"userEmail"`
	jwt.RegisteredClaims
}

// GToken 生成token
func (t *UserClaims) GToken() (string, error) {
	key := []byte("Zhoul0722")
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	signedString, err := claims.SignedString(key)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return signedString, nil
}

// PToken 解析token
func (t *UserClaims) PToken(s string) (*jwt.Token, error) {
	claims, err := jwt.ParseWithClaims(s, t, func(token *jwt.Token) (interface{}, error) {
		return []byte("Zhoul0722"), nil
	})
	if err != nil || claims == nil {
		return nil, err
	}
	return claims, nil
}

// MapClaims 将PToken解析后的结果映射为UserClaims结构体
func (t *UserClaims) MapClaims(token *jwt.Token) {
	m := token.Claims.(*UserClaims)
	*t = UserClaims{
		UserId:    m.UserId,
		UserGroup: m.UserGroup,
		UserEmail: m.UserEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        m.ID,
			Issuer:    m.Issuer,
			Subject:   m.Subject,
			IssuedAt:  m.IssuedAt,
			Audience:  m.Audience,
			NotBefore: m.NotBefore,
			ExpiresAt: m.ExpiresAt,
		},
	}
	return
}
