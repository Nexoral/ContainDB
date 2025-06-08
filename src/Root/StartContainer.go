package main

import (
	"ContainDB/src/Docker"
	"ContainDB/src/tools"
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
)

func StartContainer(database string) {
	imageMap := map[string]string{
		"mongodb":    "mongo",
		"redis":      "redis",
		"mysql":      "mysql",
		"postgresql": "postgres",
		"cassandra":  "cassandra",
		"mariadb":    "mariadb",
	}

	defaultPorts := map[string]string{
		"mongodb":    "27017",
		"redis":      "6379",
		"mysql":      "3306",
		"postgresql": "5432",
		"cassandra":  "9042",
		"mariadb":    "3306",
	}

	image := imageMap[database]
	port := defaultPorts[database]

	if Docker.IsContainerRunning(image) {
		fmt.Printf("Database %s is already running on port %s\n", database, port)
		return
	}

	// Pull image
	fmt.Printf("Pulling image %s...\n", image)
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	// Ask for port mapping
	portMapping := ""
	if Docker.AskYesNo("Do you want to map container port with host?") {
		customPort := Docker.AskYesNo("Do you want to use custom host port?")
		if customPort {
			hostPort := tools.AskForInput("Enter custom host port", port)
			portMapping = fmt.Sprintf("-p %s:%s", hostPort, port)
		} else {
			portMapping = fmt.Sprintf("-p %s:%s", port, port)
		}
	}

	restartFlag := ""
	if Docker.AskYesNo("Do you want the container to auto-restart on system startup?") {
		restartFlag = "--restart unless-stopped"
	}

	// Ask for data persistence
	volumeMapping := ""
	if Docker.AskYesNo("Do you want to persist data?") {
		// map of container paths
		containerDirs := map[string]string{
			"mongodb":    "/data/db",
			"redis":      "/data",
			"mysql":      "/var/lib/mysql",
			"postgresql": "/var/lib/postgresql/data",
			"cassandra":  "/var/lib/cassandra",
			"mariadb":    "/var/lib/mysql",
		}
		volName := fmt.Sprintf("%s-data", database)
		// if already exists, ask reuse or recreate
		if Docker.VolumeExists(volName) {
			items := []string{"Use existing", "Create fresh", "Exit"}
			prompt := promptui.Select{
				Label: fmt.Sprintf("Volume '%s' exists. Use or recreate?", volName),
				Items: items,
			}
			_, choice, _ := prompt.Run()
			if choice == "Create fresh" {
				fmt.Println("Removing and recreating volume:", volName)
				_ = Docker.RemoveVolume(volName)
				_ = Docker.CreateVolume(volName)
			}
			if choice == "Exit" {
				fmt.Println("Exiting setup.")
				return
			}
		} else {
			_ = Docker.CreateVolume(volName)
		}
		volumeMapping = fmt.Sprintf("-v %s:%s", volName, containerDirs[database])
	}

	env := ""
	if database == "mysql" || database == "postgresql" || database == "mariadb" {
		fmt.Println("You need to set environment variables for the database.")
		user := tools.AskForInput("Enter root username", "root")
		pass := tools.AskForInput("Enter root password", "password")

		if database == "mysql" {
			env = fmt.Sprintf("-e MYSQL_ROOT_PASSWORD=%s", pass)
		} else if database == "postgresql" {
			env = fmt.Sprintf("-e POSTGRES_PASSWORD=%s -e POSTGRES_USER=%s", pass, user)
		} else if database == "mariadb" {
			env = fmt.Sprintf("-e MARIADB_ROOT_PASSWORD=%s", pass)
		}

	}

	containerName := fmt.Sprintf("%s-container", database)
	runCmd := fmt.Sprintf(
		"docker run -d --network ContainDB-Network %s %s %s %s --name %s %s",
		portMapping, restartFlag, volumeMapping, env, containerName, image,
	)
	fmt.Println("Running:", runCmd)
	cmd = exec.Command("bash", "-c", runCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting container:", err)
	} else {
		fmt.Println("Container started successfully.")

		if database == "mysql" || database == "postgresql" || database == "mariadb" {
			consentPhpMyAdmin := Docker.AskYesNo("Do you want to install phpMyAdmin for this database?")
			if consentPhpMyAdmin {
				tools.StartPHPMyAdmin()
			} else {
				fmt.Println("You can install phpMyAdmin later using the 'phpmyadmin' option.")
			}
		}
		if database == "mongodb" {
			consentCompass := Docker.AskYesNo("Do you want to install MongoDB Compass?")
			if consentCompass {
				tools.DownloadMongoDBCompass()
			} else {
				fmt.Println("You can install MongoDB Compass later using the 'mongodb compass' option.")
			}
		}
	}
}
