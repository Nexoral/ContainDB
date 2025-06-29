package main

import (
	"ContainDB/src/Docker"
	"ContainDB/src/base"
	"ContainDB/src/tools"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func main() {
	VERSION := "4.12.19-stable"

	// handle version flag without requiring sudo
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("ContainDB CLI Version:", VERSION)
		return
	}

	// Replace Ctrl+C handler to avoid triggering on normal exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\n⚠️ Interrupt received, rolling back...")
		tools.Cleanup()
		os.Exit(1)
	}()

	// require sudo
	if os.Geteuid() != 0 {
		fmt.Println("❌ Please run this program with sudo")
		os.Exit(1)
	}
	
	// Check if running on Ubuntu or Debian-based system
	osReleaseBytes, err := os.ReadFile("/etc/os-release")
	if err != nil {
		fmt.Println("❌ Unable to determine OS distribution")
		os.Exit(1)
	}
	osRelease := string(osReleaseBytes)
	if !strings.Contains(strings.ToLower(osRelease), "ubuntu") && 
	   !strings.Contains(strings.ToLower(osRelease), "debian") {
		fmt.Println("❌ This program requires Ubuntu or Debian-based system")
		os.Exit(1)
	}
	
	// Check if bash shell is available
	_, err = exec.LookPath("bash")
	if err != nil {
		fmt.Println("❌ bash shell not found. This program requires bash to be installed")
		os.Exit(1)
	}



	// Check if Docker is installed and if not, prompt to install it
	base.DockerStarter()

	// Handle command line flags with FlagHandler function
	base.FlagHandler()

	errs := Docker.CreateDockerNetworkIfNotExists()
	if errs != nil {
		fmt.Println("Failed to create Docker network:", err)
		return
	}

	// Show welcome banner
	base.ShowBanner()

	// Start the base case handler which contains the main menu
	base.BaseCaseHandler()
}
