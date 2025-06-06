package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func AskYesNo(label string) bool {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		// Handle cleanup locally or call a function that doesn't create an import cycle
		os.Exit(1)
	}
	return index == 0
}

func IsContainerRunning(image string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("docker ps --filter ancestor=%s --format '{{.Names}}'", image))
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) != ""
}

func ListOfContainers(images []string) []string {
	if len(images) == 0 {
		return []string{}
	}

	// Build grep pattern from image names (e.g., "mysql|postgres|mongo")
	pattern := strings.Join(images, "|")

	// Construct command
	cmd := exec.Command("bash", "-c", fmt.Sprintf("docker ps --format '{{.Names}} {{.Image}}' | grep -E '%s'", pattern))
	output, err := cmd.Output()
	if err != nil {
		// If grep fails (e.g., no match), return empty list
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var containers []string
	for _, line := range lines {
		if parts := strings.Fields(line); len(parts) > 0 {
			containers = append(containers, parts[0])
		}
	}
	return containers
}

// VolumeExists returns true if Docker volume with given name exists
func VolumeExists(name string) bool {
	cmd := exec.Command("docker", "volume", "inspect", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// CreateVolume creates a Docker volume with given name
func CreateVolume(name string) error {
	cmd := exec.Command("docker", "volume", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveVolume force-removes a Docker volume with given name
func RemoveVolume(name string) error {
	cmd := exec.Command("docker", "volume", "rm", "-f", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
