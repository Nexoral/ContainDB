package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"ContainDB/src/Docker"
	"ContainDB/src/Tools"
	"github.com/manifoldco/promptui"
)

func selectDatabase() string {
	prompt := promptui.Select{
		Label: "Select the service to start",
		Items: []string{"mongodb", "redis", "mysql", "postgresql", "cassandra", "MariaDB", "phpmyadmin", "MongoDB Compass", "RedisInsight"},
	}
	_, result, _ := prompt.Run()
	return result
}



func askForInput(label, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [%s]: ", label, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func startContainer(database string) {
	imageMap := map[string]string{
		"mongodb":       "mongo",
		"redis":         "redis",
		"mysql":         "mysql",
		"postgresql":    "postgres",
		"cassandra":     "cassandra",
		"mariaDB":      "mariadb",
	}

	defaultPorts := map[string]string{
		"mongodb":       "27017",
		"redis":         "6379",
		"mysql":         "3306",
		"postgresql":    "5432",
		"cassandra":     "9042",
		"mariaDB":      "3306",
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
			hostPort := askForInput("Enter custom host port", port)
			portMapping = fmt.Sprintf("-p %s:%s", hostPort, port)
		} else {
			portMapping = fmt.Sprintf("-p %s:%s", port, port)
		}
	}

	restartFlag := ""
	if Docker.AskYesNo("Do you want the container to auto-restart on system startup?") {
		restartFlag = "--restart unless-stopped"
	}

	env := ""
	if database == "mysql" || database == "postgresql" || database == "mariaDB" {
		fmt.Println("You need to set environment variables for the database.")
		user := askForInput("Enter root username", "root")
		pass := askForInput("Enter root password", "password")

		if database == "mysql" {
			env = fmt.Sprintf("-e MYSQL_ROOT_PASSWORD=%s", pass)
		} else if database == "postgresql" {
			env = fmt.Sprintf("-e POSTGRES_PASSWORD=%s -e POSTGRES_USER=%s", pass, user)
		} else if database == "mariaDB" {
			env = fmt.Sprintf("-e MARIADB_ROOT_PASSWORD=%s", pass)
		}

	}

	containerName := fmt.Sprintf("%s-container", database)
	runCmd := fmt.Sprintf("docker run -d --network ContainDB-Network %s %s %s --name %s %s", portMapping, restartFlag, env, containerName, image)
	fmt.Println("Running:", runCmd)
	cmd = exec.Command("bash", "-c", runCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting container:", err)
	} else {
		fmt.Println("Container started successfully.")

		if database == "mysql" || database == "postgresql" || database == "mariaDB" {
			consentPhpMyAdmin := Docker.AskYesNo("Do you want to install phpMyAdmin for this database?")
			if consentPhpMyAdmin {
				Tools.StartPHPMyAdmin()
			} else {
				fmt.Println("You can install phpMyAdmin later using the 'phpmyadmin' option.")
			}
		}
		if database == "mongodb" {
			consentCompass := Docker.AskYesNo("Do you want to install MongoDB Compass?")
			if consentCompass {
				Tools.DownloadMongoDBCompass()
			} else {
				fmt.Println("You can install MongoDB Compass later using the 'mongodb compass' option.")
			}
		}

		if database == "redis" {
			consentRedisInsight := Docker.AskYesNo("Do you want to install Redis Insight?")
			if consentRedisInsight {
				Tools.StartRedisInsight()
			} else {
				fmt.Println("You can install RedisInsight later using the 'redis insight' option.")
			}
		}
	}
}

func main() {
	if !Docker.IsDockerInstalled() {
		err := Docker.InstallDocker()
		if err != nil {
			fmt.Println("Failed to install Docker:", err)
			return
		}
		fmt.Println("Docker installed successfully! Please restart the terminal or log out & log in again.")
	}

	err := Docker.CreateDockerNetworkIfNotExists()
	if err != nil {
		fmt.Println("Failed to create Docker network:", err)
		return
	}

	database := selectDatabase()
	if database == "phpmyadmin" {
		Tools.StartPHPMyAdmin()
	}
	if database == "MongoDB Compass" {
		Tools.DownloadMongoDBCompass();
	} 
	if database == "RedisInsight" {
		Tools.StartRedisInsight();
	} else {
		startContainer(database)
	}
}
