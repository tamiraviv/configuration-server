package main

import (
	"fmt"
	"time"

	"configuration-server/internal/pkg/configuration"
)

func main() {
	conf, err := configuration.NewConfigurationService()
	if err != nil {
		fmt.Println("Failed to create configuration-service:", err)
		return
	}

	// Register hook function that will be run when we set the new value
	conf.RegisterHook("log.level", hookExample)

	//Sleep here to make sure we watch the changes before we set new value.
	time.Sleep(1 * time.Second)
	if err := conf.Set("log.level", "debug"); err != nil {
		fmt.Println("Failed to set key 'log.level' with value 'info':", err)
		return
	}

	//Sleep here to make sure that our watch changes set the new value before reading it.
	time.Sleep(1 * time.Second)
	value := conf.GetString("log.level")
	if value == "" {
		fmt.Println("Failed to get key 'log.level' from configuration-service:", err)
		return
	}

	fmt.Println("Got the value:", value)
}

func hookExample(v string) error {
	fmt.Println("hook Example that print the new set value:", v)
	return nil
}
