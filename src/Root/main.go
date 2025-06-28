package main

import (
	"fmt"
	"os"
	"os/signal"

	"ContainDB/src/Docker"
	"ContainDB/src/tools"

	"github.com/manifoldco/promptui"
)

func main() {
	VERSION := "3.11.16-stable"

	// handle version flag without requiring sudo
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("ContainDB CLI Version:", VERSION)
		return
	}

	// Replace Ctrl+C handler to avoid triggering on normal exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		tools.Cleanup()
		os.Exit(1)
	}()

	// require sudo
	if os.Geteuid() != 0 {
		fmt.Println("❌ Please run this program with sudo")
		os.Exit(1)
	}

	if !Docker.IsDockerInstalled() {
		fmt.Println("❌ Docker is not installed. Without Docker the tool cannot run.")
		installPrompt := promptui.Select{
			Label: "Would you like to install Docker now?",
			Items: []string{"Yes", "No", "Exit"},
		}
		_, choice, err := installPrompt.Run()
		if err != nil || choice != "Yes" {
			fmt.Println("Exiting. Please install Docker manually and rerun.")
			os.Exit(1)
		}
		fmt.Println(("Checking system requirements..."))
		Docker.CheckSystemRequirements() // Check system requirements before installing Docker
		err = Docker.InstallDocker()
		if err != nil {
			fmt.Println("Failed to install Docker:", err)
			return
		}
		fmt.Println("Docker installed successfully! Please restart the terminal or log out & log in again.")
		return
	}

	err := Docker.CreateDockerNetworkIfNotExists()
	if err != nil {
		fmt.Println("Failed to create Docker network:", err)
		return
	}

	// Show welcome banner
	ShowBanner()

	// Top-level action menu
	actionPrompt := promptui.Select{
		Label: "What do you want to do?",
		Items: []string{"Install Database", "List Databases", "Remove Database", "Exit"},
	}
	_, action, err := actionPrompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		tools.Cleanup()
		return
	}

	switch action {
	case "Install Database":
		database := SelectDatabase()
		if database == "phpmyadmin" {
			tools.StartPHPMyAdmin()
		} else if database == "MongoDB Compass" {
			tools.DownloadMongoDBCompass()
		} else {
			StartContainer(database)
		}

	case "List Databases":
		names, err := Docker.ListRunningDatabases()
		if err != nil {
			fmt.Println("Error listing databases:", err)
			return
		}
		if len(names) == 0 {
			fmt.Println("No running databases found.")
		} else {
			fmt.Println("Running databases:")
			for _, n := range names {
				fmt.Println(" -", n)
			}
		}

	case "Remove Database":
		names, err := Docker.ListRunningDatabases()
		if err != nil {
			fmt.Println("Error listing databases:", err)
			return
		}
		if len(names) == 0 {
			fmt.Println("No running databases to remove.")
		} else {
			items := append(names, "Exit")
			sel := promptui.Select{
				Label: "Select database to remove",
				Items: items,
			}
			_, name, cerr := sel.Run()
			if cerr != nil || name == "Exit" {
				fmt.Println("\n⚠️ Cancelled")
				return
			}
			if err := Docker.RemoveDatabase(name); err != nil {
				fmt.Println("Error removing database:", err)
			} else {
				fmt.Println("Database removed:", name)
			}
		}

	case "Exit":
		fmt.Println("Goodbye!")
		return
	}
}
