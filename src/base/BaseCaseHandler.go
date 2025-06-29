package base

import (
	"ContainDB/src/Docker"
	"ContainDB/src/tools"
	"fmt"

	"github.com/manifoldco/promptui"
)

func BaseCaseHandler() {
	// Top-level action menu
	actionPrompt := promptui.Select{
		Label: "What do you want to do?",
		Items: []string{"Install Database", "List Databases", "Remove Database", "Remove Image", "Remove Volume", "Export Services", "Exit"},
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
		switch database {
		case "phpmyadmin":
			tools.StartPHPMyAdmin()
		case "MongoDB Compass":
			tools.DownloadMongoDBCompass()
		case "PgAdmin":
			tools.StartPgAdmin()
		case "Redis Insight":
			tools.StartRedisInsight()
		default:
			StartContainer(database)
		}

	case "List Databases":
		names, err := Docker.ListRunningDatabases()

		// Remove PgAdmin, phpmyadmin if it exists from the list
		for i, name := range names {
			if name == "phpmyadmin" {
				names = append(names[:i], names[i+1:]...)
				break
			} else if name == "pgadmin" {
				names = append(names[:i], names[i+1:]...)
				break
			} else if name == "redisinsight" {
				names = append(names[:i], names[i+1:]...)
				break
			}
		}
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
				fmt.Println("✅ Database", name, "removed successfully")
			}
		}

	case "Remove Image":
		images, err := Docker.ListDatabaseImages()
		if err != nil {
			fmt.Println("Error listing database images:", err)
			return
		}

		if len(images) == 0 {
			fmt.Println("No database images found.")
			return
		}

		items := append(images, "Exit")
		sel := promptui.Select{
			Label: "Select image to remove",
			Items: items,
		}
		_, selected, cerr := sel.Run()
		if cerr != nil || selected == "Exit" {
			fmt.Println("\n⚠️ Cancelled")
			return
		}

		// Check if image is in use
		inUse, containerName, err := Docker.IsImageInUse(selected)
		if err != nil {
			fmt.Println("Error checking if image is in use:", err)
			return
		}

		if inUse {
			fmt.Printf("⚠️ Cannot remove image: it's currently used by container '%s'\n", containerName)
			fmt.Println("Please remove the container first using 'Remove Database' option.")
			return
		}

		// Confirm removal
		if !Docker.AskYesNo(fmt.Sprintf("Are you sure you want to remove image '%s'?", selected)) {
			fmt.Println("\n⚠️ Image removal cancelled")
			return
		}

		if err := Docker.RemoveImage(selected); err != nil {
			fmt.Println("Error removing image:", err)
		} else {
			fmt.Printf("✅ Image '%s' removed successfully\n", selected)
		}

	case "Remove Volume":
		volumes, err := Docker.ListContainDBVolumes()
		if err != nil {
			fmt.Println("Error listing volumes:", err)
			return
		}

		if len(volumes) == 0 {
			fmt.Println("No database volumes found.")
			return
		}

		items := append(volumes, "Exit")
		sel := promptui.Select{
			Label: "Select volume to remove",
			Items: items,
		}
		_, selected, cerr := sel.Run()
		if cerr != nil || selected == "Exit" {
			fmt.Println("\n⚠️ Cancelled")
			return
		}

		// Check if volume is in use
		inUse, containerName, err := Docker.IsVolumeInUse(selected)
		if err != nil {
			fmt.Println("Error checking if volume is in use:", err)
			return
		}

		if inUse {
			fmt.Printf("⚠️ Cannot remove volume: it's currently used by container '%s'\n", containerName)
			fmt.Println("Please remove the container first using 'Remove Database' option.")
			return
		}

		// Confirm removal
		if !Docker.AskYesNo(fmt.Sprintf("Are you sure you want to remove volume '%s'? This will delete ALL DATA in this volume!", selected)) {
			fmt.Println("\n⚠️ Volume removal cancelled")
			return
		}

		if err := Docker.RemoveVolume(selected); err != nil {
			fmt.Println("Error removing volume:", err)
		} else {
			fmt.Printf("✅ Volume '%s' removed successfully\n", selected)
		}
	case "Export Services":
		fmt.Println("Exporting Docker Compose file with all running services...")
		fmt.Println("\n⚠️  IMPORTANT: The export functionality only exports container configurations, not the actual data.")
		fmt.Println("   Even if you used data persistence during installation, the exported compose file only")
		fmt.Println("   references local volume paths from your current machine which won't exist on other systems.")
		fmt.Println("   For data backup, please use each database's native backup tools.\n")

		filePath := Docker.MakeDockerComposeWithAllServices()
		if filePath == "" {
			fmt.Println("Failed to create Docker Compose file.")
		} else {
			fmt.Println("\n✅ Docker Compose file created successfully at:", filePath)
			fmt.Println("   This file contains only the configuration of your containers.")
		}

	// Handle exit case
	case "Exit":
		fmt.Println("Goodbye!")
		return
	}
}
