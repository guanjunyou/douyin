package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	CommonEntity
	Name string `json:"name"`
	jwt.StandardClaims
}

var myKey = []byte("douyin")

// GenerateToken
// 生成 token
func GenerateToken(name string, commonEntity CommonEntity) (string, error) {
	UserClaim := &UserClaims{
		CommonEntity: commonEntity,
		Name:         name,
		//IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}
