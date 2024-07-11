package test

import (
	"fmt"
	"testing"
	"worker-mesh/config"
)

func TestConfig(t *testing.T) {
	config.InitConfig("../")

	fmt.Println(config.ReadConfig().Nsq)
}
