package utils

import (
	"github.com/importcjj/sensitive"
	"log"
)

/*
*
使用方法 content = util.Filter.Replace(content, '#')
*/
var Filter *sensitive.Filter

const WordDictPath = "./public/sensitiveDict.txt"

func InitFilter() {
	Filter = sensitive.New()
	err := Filter.LoadWordDict(WordDictPath)
	if err != nil {
		log.Println("InitFilter Fail,Err=" + err.Error())
	}
}
