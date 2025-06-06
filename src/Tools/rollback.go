package Tools

import (
	"fmt"
	"os"
	"os/exec"
)

// Cleanup stops and removes any created containers and temporary artifacts.
func Cleanup() {
	fmt.Println("🧹 Cleaning up resources...")

	// stop & remove all containers named "*-container"
	exec.Command("bash", "-c", "docker rm -f $(docker ps -a --filter \"name=-container\" -q)").Run()

	// remove phpMyAdmin and RedisInsight containers if present
	exec.Command("docker", "rm", "-f", "phpmyadmin").Run()
	exec.Command("docker", "rm", "-f", "redisinsight").Run()

	// clean up MongoDB Compass download
	os.Remove("/tmp/mongodb-compass.deb")

	fmt.Println("✅ Cleanup completed.")
	// exit immediately to prevent any further interactive prompts
	os.Exit(1)
}
