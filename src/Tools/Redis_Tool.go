package Tools;

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"ContainDB/src/Docker"
	"github.com/manifoldco/promptui"
)

func StartRedisInsight() {
	redisContainers := Docker.ListOfContainers([]string{"redis", "redis-stack", "redis-enterprise"})
	if len(redisContainers) == 0 {
		fmt.Println("No running Redis containers found.")
		return
	}

	prompt := promptui.Select{
		Label: "Select a Redis container to connect with RedisInsight",
		Items: redisContainers,
	}
	_, selectedContainer, _ := prompt.Run()

	port := askForInput("Enter host port to expose RedisInsight", "5540")

	fmt.Println("Pulling RedisInsight image...")
	cmd := exec.Command("docker", "pull", "redis/redisinsight:latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error pulling RedisInsight:", err)
		return
	}

	// Get IP address of selected redis container
	ipCmd := exec.Command("bash", "-c", fmt.Sprintf("docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' %s", selectedContainer))
	ipOut, err := ipCmd.Output()
	if err != nil {
		fmt.Println("Failed to get Redis container IP:", err)
		return
	}
	redisIP := strings.TrimSpace(string(ipOut))

	runCmd := fmt.Sprintf(
		`docker run -d --name redisinsight -p %s:5540 -e RI_REDIS_HOST=%s redis/redisinsight:latest`,
		port, redisIP,
	)

	fmt.Println("Running:", runCmd)
	cmd = exec.Command("bash", "-c", runCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting RedisInsight container:", err)
	} else {
		fmt.Printf("RedisInsight started. Access it at http://localhost:%s\n", port)
	}
}