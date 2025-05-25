package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"ContainDB/installation"

	"github.com/manifoldco/promptui"
)

func selectDatabase() string {
	prompt := promptui.Select{
		Label: "Select the service to start",
		Items: []string{"mongodb", "redis", "mysql", "postgresql", "cassandra", "phpmyadmin"},
	}
	_, result, _ := prompt.Run()
	return result
}

func listSQLContainers() []string {
	cmd := exec.Command("bash", "-c", "docker ps --format '{{.Names}} {{.Image}}' | grep -E 'mysql|postgres'")
	output, _ := cmd.Output()

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	containers := []string{}
	for _, line := range lines {
		if parts := strings.Fields(line); len(parts) > 0 {
			containers = append(containers, parts[0])
		}
	}
	return containers
}

func startPHPMyAdmin() {
	sqlContainers := listSQLContainers()
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
		"docker run -d --network ContainDB-Network --name phpmyadmin -e PMA_HOST=%s -p %s:80 phpmyadmin/phpmyadmin",
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

func isContainerRunning(image string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("docker ps --filter ancestor=%s --format '{{.Names}}'", image))
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) != ""
}

func askYesNo(label string) bool {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	index, _, _ := prompt.Run()
	return index == 0
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
		"mongodb":    "mongo",
		"redis":      "redis",
		"mysql":      "mysql",
		"postgresql": "postgres",
		"cassandra":  "cassandra",
	}

	defaultPorts := map[string]string{
		"mongodb":    "27017",
		"redis":      "6379",
		"mysql":      "3306",
		"postgresql": "5432",
		"cassandra":  "9042",
	}

	image := imageMap[database]
	port := defaultPorts[database]

	if isContainerRunning(image) {
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
	if askYesNo("Do you want to map container port with host?") {
		customPort := askYesNo("Do you want to use custom host port?")
		if customPort {
			hostPort := askForInput("Enter custom host port", port)
			portMapping = fmt.Sprintf("-p %s:%s", hostPort, port)
		} else {
			portMapping = fmt.Sprintf("-p %s:%s", port, port)
		}
	}

	restartFlag := ""
	if askYesNo("Do you want the container to auto-restart on system startup?") {
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
	}
}

func main() {
	if !Installation.IsDockerInstalled() {
		err := Installation.InstallDocker()
		if err != nil {
			fmt.Println("Failed to install Docker:", err)
			return
		}
		fmt.Println("Docker installed successfully! Please restart the terminal or log out & log in again.")
	}

	err := Installation.CreateDockerNetworkIfNotExists()
	if err != nil {
		fmt.Println("Failed to create Docker network:", err)
		return
	}

	database := selectDatabase()
	if database == "phpmyadmin" {
		startPHPMyAdmin()
	} else {
		startContainer(database)
	}
}
