package test

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/utils"
	"testing"
)

func TestSnowflake(t *testing.T) {
	for i := 0; i < 1000; i++ {
		sf := utils.NewSnowflake()
		fmt.Println(sf.NextID())
	}
}
