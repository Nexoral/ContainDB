package Docker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// DockerComposeConfig represents the structure of a docker-compose.yml file
type DockerComposeConfig struct {
	Version  string                          `yaml:"version"`
	Services map[string]DockerComposeService `yaml:"services"`
	Volumes  map[string]interface{}          `yaml:"volumes"`
}

// DockerComposeService represents a service in a docker-compose.yml file
type DockerComposeService struct {
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
	Environment map[string]string `yaml:"environment"`
	// Add other fields as needed
}

// ImportDockerServices imports services from a docker-compose.yml file
func ImportDockerServices(composeFilePath string) error {
	// Check if Docker is installed
	if !IsDockerInstalled() {
		return errors.New("docker is not installed. Please install Docker first")
	}

	// Read and parse docker-compose.yml
	composeData, err := ioutil.ReadFile(composeFilePath)
	if err != nil {
		return fmt.Errorf("failed to read compose file: %v", err)
	}

	var composeConfig DockerComposeConfig
	err = yaml.Unmarshal(composeData, &composeConfig)
	if err != nil {
		return fmt.Errorf("failed to parse compose file: %v", err)
	}

	// Check if services are already running
	fmt.Println("Checking if services are already running...")
	runningContainers, err := getRunningContainers()
	if err != nil {
		return fmt.Errorf("failed to get running containers: %v", err)
	}

	for serviceName := range composeConfig.Services {
		for _, containerInfo := range runningContainers {
			if strings.Contains(containerInfo, serviceName) {
				fmt.Printf("Warning: Service '%s' appears to be already running\n", serviceName)
			}
		}
	}

	// Check if ports are available
	fmt.Println("Checking if ports are available...")
	for serviceName, serviceConfig := range composeConfig.Services {
		for _, portMapping := range serviceConfig.Ports {
			hostPort := strings.Split(strings.Split(portMapping, ":")[0], "-")[0]
			if !isPortAvailable(hostPort) {
				fmt.Printf("Warning: Port %s required by service '%s' is already in use\n", hostPort, serviceName)
			}
		}
	}

	// Create volumes if they don't exist
	fmt.Println("Creating volumes...")
	if len(composeConfig.Volumes) > 0 {
		for volumeName := range composeConfig.Volumes {
			if !volumeExists(volumeName) {
				fmt.Printf("Creating volume '%s'...\n", volumeName)
				cmd := exec.Command("docker", "volume", "create", volumeName)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("failed to create volume '%s': %v", volumeName, err)
				}
			} else {
				fmt.Printf("Volume '%s' already exists\n", volumeName)
			}
		}
	}

	// Start services using docker-compose up -d
	fmt.Println("Starting services...")
	cmd := exec.Command("docker", "compose", "-f", composeFilePath, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start services: %v", err)
	}

	return nil
}

// Helper function to check if a volume exists
func volumeExists(name string) bool {
	cmd := exec.Command("docker", "volume", "inspect", name)
	err := cmd.Run()
	return err == nil
}

// getRunningContainers returns a list of running Docker containers
func getRunningContainers() ([]string, error) {
	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(output), "\n"), nil
}

// isPortAvailable checks if a port is available
func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}
