package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"ContainDB/src/Docker"
	"github.com/manifoldco/promptui"
)

func selectDatabase() string {
	prompt := promptui.Select{
		Label: "Select the service to start",
		Items: []string{"mongodb", "redis", "mysql", "postgresql", "cassandra", "phpmyadmin", "MongoDB Compass", "RedisInsight"},
	}
	_, result, _ := prompt.Run()
	return result
}

func startPHPMyAdmin() {
	sqlContainers := Docker.ListSQLContainers()
	if len(sqlContainers) == 0 {
		fmt.Println("No running MySQL/PostgreSQL containers found.")
		return
	}

	prompt := promptui.Select{
		Label: "Select a SQL container to link with phpMyAdmin",
		Items: sqlContainers,
	}
	_, selectedContainer, _ := prompt.Run()

	port := askForInput("Enter host port to expose phpMyAdmin", "8080")

	fmt.Printf("Pulling phpMyAdmin image...\n")
	cmd := exec.Command("docker", "pull", "phpmyadmin/phpmyadmin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	runCmd := fmt.Sprintf(
		"docker run -d --restart unless-stopped --network ContainDB-Network --name phpmyadmin -e PMA_HOST=%s -p %s:80 phpmyadmin/phpmyadmin",
		selectedContainer, port,
	)

	fmt.Println("Running:", runCmd)
	cmd = exec.Command("bash", "-c", runCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting phpMyAdmin:", err)
	} else {
		fmt.Printf("phpMyAdmin started. Access it at http://localhost:%s\n", port)
	}
}

func DownloadMongoDBCompass() {
	fmt.Println("Downloading MongoDB Compass...")
	// Download the deb file first
	cmd := exec.Command("bash", "-c", "wget -q https://downloads.mongodb.com/compass/mongodb-compass_1.46.2_amd64.deb -O /tmp/mongodb-compass.deb")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error downloading MongoDB Compass:", err)
		return
	}
	
	// Install the downloaded deb file
	installCmd := exec.Command("sudo", "dpkg", "-i", "/tmp/mongodb-compass.deb")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Println("Error installing MongoDB Compass:", err)
	} else {
		fmt.Println("MongoDB Compass downloaded and installed successfully.")

		// Clean up the downloaded file
		cleanupCmd := exec.Command("rm", "/tmp/mongodb-compass.deb")
		cleanupCmd.Stdout = os.Stdout
		cleanupCmd.Stderr = os.Stderr
		if err := cleanupCmd.Run(); err != nil {
			fmt.Println("Error cleaning up downloaded file:", err)
		} else {
			fmt.Println("Temporary files cleaned up successfully.")
		}
		fmt.Println("You can now launch MongoDB Compass from your applications menu or by running 'mongodb-compass' in the terminal.")
		fmt.Println("Note: If you encounter any issues, please ensure you have the necessary dependencies installed.")
		fmt.Println("For more information, visit: https://www.mongodb.com/docs/compass/current/install/")
		fmt.Println("Enjoy using MongoDB Compass!")
		
	}
}

func DownloadRedisInsight() {
	fmt.Println("Downloading RedisInsight...")
	cmd := exec.Command("docker", "pull", "redis/redisinsight:latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error downloading RedisInsight:", err)
	} else {
		command := exec.Command("docker", "run", "-d", "--name", "redisinsight", "-p", "8001:8001", "redis/redisinsight:latest")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			fmt.Println("Error starting RedisInsight container:", err)
		} else {
			fmt.Println("RedisInsight started successfully. Access it at http://localhost:8001")

		}
	}
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
		"redis insight": "redis/redisinsight:latest",
	}

	defaultPorts := map[string]string{
		"mongodb":       "27017",
		"redis":         "6379",
		"mysql":         "3306",
		"postgresql":    "5432",
		"cassandra":     "9042",
		"redis insight": "8001",
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
	if database == "mysql" || database == "postgresql" {
		user := askForInput("Enter root username", "root")
		pass := askForInput("Enter root password", "password")

		if database == "mysql" {
			env = fmt.Sprintf("-e MYSQL_ROOT_PASSWORD=%s", pass)
		} else if database == "postgresql" {
			env = fmt.Sprintf("-e POSTGRES_PASSWORD=%s -e POSTGRES_USER=%s", pass, user)
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

		if database == "mysql" || database == "postgresql" {
			consentPhpMyAdmin := Docker.AskYesNo("Do you want to install phpMyAdmin for this database?")
			if consentPhpMyAdmin {
				startPHPMyAdmin()
			} else {
				fmt.Println("You can install phpMyAdmin later using the 'phpmyadmin' option.")
			}
		}
		if database == "mongodb" {
			consentCompass := Docker.AskYesNo("Do you want to install MongoDB Compass?")
			if consentCompass {
				DownloadMongoDBCompass()
			} else {
				fmt.Println("You can install MongoDB Compass later using the 'mongodb compass' option.")
			}
		}

		if database == "redis" {
			consentRedisInsight := Docker.AskYesNo("Do you want to install RedisInsight?")
			if consentRedisInsight {
				DownloadRedisInsight()
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
		startPHPMyAdmin()
	}
	if database == "MongoDB Compass" {
		DownloadMongoDBCompass()
	} else {
		startContainer(database)
	}
}
