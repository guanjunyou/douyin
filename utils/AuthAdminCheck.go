package utils

import (
	"errors"
	"log"
)

/*
*
鉴权
*/
func AuthAdminCheck(token string) error {
	claims, err := AnalyseToken(token)
	if err != nil || claims == nil {
		log.Printf("Can not find this token !")
		return errors.New("Can not find this token !")
	}
	return nil
}
