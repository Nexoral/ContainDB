package tools

import (
	"fmt"
	"os"
	"os/exec"
)

// Cleanup stops and removes any created containers and temporary artifacts.
func Cleanup() {
	fmt.Println("🧹 Cleaning up resources...")

	// remove exited/dead containers
	fmt.Println("- Removing failed containers...")
	exec.Command("bash", "-c", "docker rm -f $(docker ps -a --filter \"status=exited\" -q)").Run()
	exec.Command("bash", "-c", "docker rm -f $(docker ps -a --filter \"status=dead\" -q)").Run()
	exec.Command("bash", "-c", "docker rm -f $(docker ps -a --filter \"status=created\" -q)").Run()

	// remove dangling images
	fmt.Println("- Removing dangling images...")
	exec.Command("bash", "-c", "docker image prune -f").Run()

	// clean up MongoDB Compass download
	os.Remove("/tmp/mongodb-compass.deb")

	fmt.Println("✅ Cleanup completed.")
	// exit immediately to prevent any further interactive prompts
	os.Exit(1)
}
