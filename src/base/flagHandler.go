package base

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
)

// flagHandler handles command line flags for the ContainDB CLI
func FlagHandler() {
	if len(os.Args) > 1 && os.Args[1] == "--uninstall-docker" {
		if !Docker.IsDockerInstalled() {
			fmt.Println("Docker is not installed. Nothing to uninstall.")
		}
		fmt.Println("Uninstalling Docker...")
		err := Docker.UninstallDocker()
		if err != nil {
			fmt.Println("Failed to uninstall Docker:", err)
		}
		fmt.Println("Docker uninstalled successfully! Please restart the terminal or log out & log in again.")
		os.Exit(0) // Exit after handling flags
	} else if len(os.Args) > 1 && os.Args[1] == "--export" {
		fmt.Println("Exporting Docker Compose file with all running services...")
		fmt.Println("\n⚠️  IMPORTANT: The export functionality only exports container configurations, not the actual data.")
		fmt.Println("   Even if you used data persistence during installation, the exported compose file only")
		fmt.Println("   references local volume paths from your current machine which won't exist on other systems.")
		fmt.Println("   For data backup, please use each database's native backup tools.")

		filePath := Docker.MakeDockerComposeWithAllServices()
		if filePath == "" {
			fmt.Println("Failed to create Docker Compose file.")
		} else {
			fmt.Println("\n✅ Docker Compose file created successfully at:", filePath)
			fmt.Println("   This file contains only the configuration of your containers.")
		}
		os.Exit(0) // Exit after handling flags
	} else if len(os.Args) > 2 && os.Args[1] == "--import" {
		composeFile := os.Args[2]
		if _, err := os.Stat(composeFile); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist\n", composeFile)
			os.Exit(1)
		}
		fmt.Printf("Importing services from Docker Compose file: %s\n", composeFile)
		err := Docker.ImportDockerServices(composeFile)
		if err != nil {
			fmt.Printf("Failed to import services: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Services imported and started successfully!")
		os.Exit(0) // Exit after handling flags
	} else if len(os.Args) == 0 {
		return // No flags to handle, continue with normal execution
	}
}
