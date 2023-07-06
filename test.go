package main

import (
	"fmt"
	"plugin"
)

func main(){
	err := plugin.EngageAzure()
	if err != nil {
		fmt.Println("error running app")
	}
}