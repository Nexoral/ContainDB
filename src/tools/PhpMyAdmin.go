package tools

import (
	"bufio"
	"ContainDB/src/Docker"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

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

func StartPHPMyAdmin() {
	sqlContainers := Docker.ListOfContainers([]string{"mysql", "postgres", "mariadb"})
	if len(sqlContainers) == 0 {
		fmt.Println("No running MySQL/PostgreSQL/MariaDB containers found.")
		return
	}

	prompt := promptui.Select{
		Label: "Select a SQL container to link with phpMyAdmin",
		Items: sqlContainers,
	}
	_, selectedContainer, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		Cleanup()
	}

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
