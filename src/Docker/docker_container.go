package Docker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func AskYesNo(label string) bool {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	index, _, _ := prompt.Run()
	return index == 0
}

func IsContainerRunning(image string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("docker ps --filter ancestor=%s --format '{{.Names}}'", image))
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) != ""
}
