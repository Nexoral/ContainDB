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
		names, err := Docker.ListRunningDatabases();

		// Remove PhpMyAdmin if it exists from the list
		for i, name := range names {
			if name == "phpmyadmin" {
				names = append(names[:i], names[i+1:]...)
				break
			}
		}
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

	case "Exit":
		fmt.Println("Goodbye!")
		return
	}
}
