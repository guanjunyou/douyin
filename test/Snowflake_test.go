package test

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/utils"
	"testing"
)

func TestSnowflake(t *testing.T) {
	sf := utils.NewSnowflake()
	fmt.Println(sf.NextID())
}
