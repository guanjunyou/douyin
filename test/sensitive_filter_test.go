package test

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/utils"
	"testing"
)

func TestSensitiveFilter(t *testing.T) {
	utils.InitFilter()
	content := "氟abcdfuck"
	contentfiltered := utils.Filter.Replace(content, '*')
	fmt.Println(contentfiltered)
}
