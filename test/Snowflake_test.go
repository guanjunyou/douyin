package test

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/utils"
	"testing"
)

func TestSnowflake(t *testing.T) {
	sf := utils.NewSnowflake(1)
	fmt.Println(sf.NextID())
}
