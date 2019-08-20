package main

import (
	"configuration-server/internal/pkg/configuration"
	"fmt"
)

func main() {
	conf, err := configuration.NewConfigurationService()
	if err != nil {
		fmt.Println("Failed to create configuration-service:", err)
		return
	}

	value, err := conf.Get("log.level")
	if err != nil {
		fmt.Println("Failed to get key log.level from configuration-service:", err)
		return
	}

	fmt.Println("Got the value:", value)
}
