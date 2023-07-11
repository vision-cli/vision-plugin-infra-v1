package main

import (
	"fmt"
	"github.com/vision-cli/vision-plugin-infra-v1/plugin"
)

func maian(){
	err := plugin.EngageAzure()
	if err != nil {
		fmt.Println("error running app")
	}
}