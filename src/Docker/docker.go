package Docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ListRunningDatabases returns names of containers on the ContainDB-Network
func ListRunningDatabases() ([]string, error) {
	cmd := exec.Command("docker", "ps",
		"--filter", "network=ContainDB-Network",
		"--format", "{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}

// RemoveDatabase forcibly removes the given container,
// optionally deleting its associated data volumes.
func RemoveDatabase(name string) error {
	// ask user whether to delete attached volumes
	deleteVolumes := AskYesNo("Do you want to delete associated data volumes?")

	// Get container type from name (assuming container name follows pattern like "mongodb-container")
	containerType := ""
	if parts := strings.Split(name, "-"); len(parts) > 0 {
		containerType = parts[0] // Extract database type from container name
	}

	// First remove the container itself
	args := []string{"rm", "-f"}
	if deleteVolumes {
		args = append(args, "-v") // This only removes anonymous volumes
	}
	args = append(args, name)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error removing container: %v", err)
	}

	// If user chose to delete volumes and we could identify the database type, remove named volume
	if deleteVolumes && containerType != "" {
		volumeName := fmt.Sprintf("%s-data", containerType)
		if VolumeExists(volumeName) {
			fmt.Printf("Removing associated volume: %s\n", volumeName)
			if err := RemoveVolume(volumeName); err != nil {
				return fmt.Errorf("error removing volume %s: %v", volumeName, err)
			}
		}
	}

	return nil
}
