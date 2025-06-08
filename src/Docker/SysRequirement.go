package Docker;

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func CheckSystemRequirements() {
	// Check if Docker is installed and accessible
	if err := checkDockerInstallation(); err != nil {
		fmt.Printf("WARNING: %v\n", err)
		fmt.Println("Exiting program due to Docker requirement failure")
		os.Exit(1)
	}

	// Check if system has enough RAM (minimum 2GB recommended)
	if err := checkRAM(2); err != nil {
		fmt.Printf("WARNING: %v\n", err)
		fmt.Println("Exiting program due to RAM requirement failure")
		os.Exit(1)
	}

	// Check if system has enough disk space (minimum 10GB recommended)
	if err := checkDiskSpace(10); err != nil {
		fmt.Printf("WARNING: %v\n", err)
		fmt.Println("Exiting program due to disk space requirement failure")
		os.Exit(1)
	}

	fmt.Println("All system requirements checks passed!")
}

func checkDockerInstallation() error {
	cmd := exec.Command("docker", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker is not installed or not accessible: %v", err)
	}
	
	fmt.Printf("Docker is available: %s\n", strings.TrimSpace(string(output)))
	return nil
}

func checkRAM(minGB float64) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	totalRAM := float64(m.TotalAlloc) / (1024 * 1024 * 1024) // Convert to GB
	availableRAM := float64(m.Sys-m.Alloc) / (1024 * 1024 * 1024) // Convert to GB

	if availableRAM < minGB {
		return fmt.Errorf("insufficient RAM. Available: %.2f GB, Required: %.2f GB", availableRAM, minGB)
	}
	
	fmt.Printf("RAM check passed. Total: %.2f GB, Available: %.2f GB\n", totalRAM, availableRAM)
	return nil
}

func checkDiskSpace(minGB float64) error {
	var path string
	
	if runtime.GOOS == "windows" {
		path = "C:\\"
	} else {
		path = "/"
	}
	
	cmd := exec.Command("df", "-h", path)
	if runtime.GOOS == "windows" {
		// For Windows, use a different command
		cmd = exec.Command("powershell", "-Command", "Get-PSDrive C | Select-Object Free")
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check disk space: %v", err)
	}
	
	// Parse the output to get available space
	// This is a simplified approach and may need adjustments based on actual output format
	outputStr := string(output)
	freeGB := 0.0

	if runtime.GOOS == "windows" {
		// Parse Windows output
		lines := strings.Split(outputStr, "\n")
		if len(lines) >= 3 {
			freeBytes, err := strconv.ParseFloat(strings.TrimSpace(lines[2]), 64)
			if err == nil {
				freeGB = freeBytes / (1024 * 1024 * 1024)
			}
		}
	} else {
		// Parse Unix/Linux output
		lines := strings.Split(outputStr, "\n")
		for i, line := range lines {
			if i > 0 && strings.Contains(line, path) {
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					// Extract available space and convert to GB
					available := fields[3]
					if strings.HasSuffix(available, "G") {
						freeGB, _ = strconv.ParseFloat(available[:len(available)-1], 64)
					}
				}
				break
			}
		}
	}
	
	if freeGB < minGB {
		return fmt.Errorf("insufficient disk space. Available: %.2f GB, Required: %.2f GB", freeGB, minGB)
	}
	
	fmt.Printf("Disk space check passed. Available: %.2f GB\n", freeGB)
	return nil
}