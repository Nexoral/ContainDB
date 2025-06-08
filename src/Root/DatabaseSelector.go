package main

import (
	"ContainDB/src/tools"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func SelectDatabase() string {
	options := []string{"mongodb", "redis", "mysql", "postgresql", "cassandra", "mariadb", "phpmyadmin", "MongoDB Compass", "Exit"}
	prompt := promptui.Select{
		Label: "Select the service to start",
		Items: options,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		tools.Cleanup() // perform cleanup and exit
	}
	if result == "Exit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}
	return result
}
