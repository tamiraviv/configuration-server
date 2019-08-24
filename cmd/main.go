package main

import (
	"fmt"
	"time"

	"configuration-server/internal/pkg/configuration"
)

const (
	logLevelKey   = "log.level"
	logLevelValue = "debug"
)

func main() {
	conf, err := configuration.NewConfigurationService()
	if err != nil {
		fmt.Println("Failed to create configuration-service:", err)
		return
	}

	// Register hook function that will be run when we set the new value
	fmt.Println("Registering hook...")
	conf.RegisterHook(logLevelKey, hookExample)

	//Sleep here to make sure we watch the changes before we set new value.
	time.Sleep(1 * time.Second)
	fmt.Printf("Setting key (%s) value (%s) to configuration...\n", logLevelKey, logLevelValue)
	if err := conf.Set(logLevelKey, logLevelValue); err != nil {
		fmt.Printf("Failed to set key '%s' with value '%s': %s\n", logLevelKey, logLevelValue, err)
		return
	}
	fmt.Println("Set succecful!")

	//Sleep here to make sure that our watch changes set the new value before reading it.
	time.Sleep(1 * time.Second)
	fmt.Printf("Getting key (%s) from configuration...\n", logLevelKey)
	value := conf.GetString(logLevelKey)
	if value == "" {
		fmt.Println("Failed to get key '%s' from configuration-service: %s\n", logLevelKey, err)
		return
	}

	fmt.Printf("Expected value: %s, Got the value: %s\n", logLevelValue, value)
	time.Sleep(1 * time.Second)
}

func hookExample(v string) error {
	fmt.Println("hook Example that print the new set value:", v)
	return nil
}
