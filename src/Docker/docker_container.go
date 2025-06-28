package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func AskYesNo(label string) bool {
	items := []string{"Yes", "No", "Exit"}
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		// Handle cleanup locally or call a function that doesn't create an import cycle
		os.Exit(1)
	}
	if index == len(items)-1 {
		fmt.Println("Exiting...")
		os.Exit(0)
	}
	return index == 0
}

func IsContainerRunning(nameOrImage string, checkByName bool) bool {
	var cmd *exec.Cmd
	if checkByName {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("docker ps --filter name=%s --format '{{.Names}}'", nameOrImage))
	} else {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("docker ps --filter ancestor=%s --format '{{.Names}}'", nameOrImage))
	}
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
	err := cmd.Run()
	return err == nil
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
	// First check if volume exists
	if !VolumeExists(name) {
		return fmt.Errorf("volume %s does not exist", name)
	}

	fmt.Printf("Removing volume %s...\n", name)
	cmd := exec.Command("docker", "volume", "rm", "-f", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove volume: %v", err)
	}
	return nil
}
