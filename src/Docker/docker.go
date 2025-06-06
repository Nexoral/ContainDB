package Docker

import (
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
	// build docker rm args
	args := []string{"rm", "-f"}
	if deleteVolumes {
		args = append(args, "-v")
	}
	args = append(args, name)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
