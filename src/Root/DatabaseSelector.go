package main

import (
	"ContainDB/src/tools"
	"fmt"

	"github.com/manifoldco/promptui"
)

func SelectDatabase() string {
	prompt := promptui.Select{
		Label: "Select the service to start",
		Items: []string{"mongodb", "redis", "mysql", "postgresql", "cassandra", "mariadb", "phpmyadmin", "MongoDB Compass"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		tools.Cleanup() // perform cleanup and exit
	}
	return result
}