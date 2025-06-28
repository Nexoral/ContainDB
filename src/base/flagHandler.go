package base

import (
	"ContainDB/src/Docker"
	"fmt"
	"os"
)

// flagHandler handles command line flags for the ContainDB CLI
func FlagHandler() {
	VERSION := "4.11.17-stable"
	// handle version flag without requiring sudo
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("ContainDB CLI Version:", VERSION)
		os.Exit(0) // Exit after handling flags
	} else if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("ContainDB CLI - A tool for managing Docker databases")
		fmt.Println("Usage: sudo containdb")
		fmt.Println("Options:")
		fmt.Println("  --version   Show version information")
		fmt.Println("  --help             Show this help message")
		fmt.Println("  --install-docker   Install Docker if not installed")
		fmt.Println("  --uninstall-docker Uninstall Docker if installed")
		os.Exit(0) // Exit after handling flags
	} else if len(os.Args) > 1 && os.Args[1] == "--install-docker" {
		if !Docker.IsDockerInstalled() {
			fmt.Println("Docker is not installed. Installing Docker...")
			err := Docker.InstallDocker()
			if err != nil {
				fmt.Println("Failed to install Docker:", err)
			}
			fmt.Println("Docker installed successfully! Please restart the terminal or log out & log in again.")
		} else {
			fmt.Println("Docker is already installed.")
		}
		os.Exit(0) // Exit after handling flags
	} else if len(os.Args) > 1 && os.Args[1] == "--uninstall-docker" {
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
	} else if len(os.Args) == 0 {
		return // No flags to handle, continue with normal execution
	}
}
