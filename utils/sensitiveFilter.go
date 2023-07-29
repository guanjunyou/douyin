package utils

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/importcjj/sensitive"
	"log"
)

/*
*
使用方法 content = util.Filter.Replace(content, '#')
*/
var Filter *sensitive.Filter

func InitFilter() {
	Filter = sensitive.New()
	err := Filter.LoadWordDict(config.WordDictPath)
	if err != nil {
		log.Println("InitFilter Fail,Err=" + err.Error())
	}
}
