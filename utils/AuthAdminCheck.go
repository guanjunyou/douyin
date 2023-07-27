package utils

import (
	"errors"
	"github.com/RaymondCode/simple-demo/models"
	"log"
)

/*
*
鉴权
*/
func AuthAdminCheck(token string) error {
	claims, err := models.AnalyseToken(token)
	if err != nil || claims == nil {
		log.Printf("Can not find this token !")
		return errors.New("Can not find this token !")
	}
	return nil
}
