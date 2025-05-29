package Tools;

import (
	"fmt"
	"os"
	"os/exec"
)

func DownloadMongoDBCompass() {
	fmt.Println("Downloading MongoDB Compass...")
	
	// Download with progress display using wget with progress bar
	cmd := exec.Command("bash", "-c", "wget https://downloads.mongodb.com/compass/mongodb-compass_1.46.2_amd64.deb -O /tmp/mongodb-compass.deb --progress=bar:force 2>&1")
	
	// Create a pipe to capture and process the output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}
	
	cmd.Stderr = cmd.Stdout
	
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting download:", err)
		return
	}
	
	// Read and display progress output
	buf := make([]byte, 100)
	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			fmt.Print(string(buf[:n]))
		}
		if err != nil {
			break
		}
	}
	
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error downloading MongoDB Compass:", err)
		return
	}
	
	// Install the downloaded deb file
	installCmd := exec.Command("sudo", "dpkg", "-i", "/tmp/mongodb-compass.deb")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Println("Error installing MongoDB Compass:", err)
	} else {
		fmt.Println("MongoDB Compass downloaded and installed successfully.")

		// Clean up the downloaded file
		cleanupCmd := exec.Command("rm", "/tmp/mongodb-compass.deb")
		cleanupCmd.Stdout = os.Stdout
		cleanupCmd.Stderr = os.Stderr
		if err := cleanupCmd.Run(); err != nil {
			fmt.Println("Error cleaning up downloaded file:", err)
		} else {
			fmt.Println("Temporary files cleaned up successfully.")
		}
		fmt.Println("You can now launch MongoDB Compass from your applications menu or by running 'mongodb-compass' in the terminal.")
		fmt.Println("Note: If you encounter any issues, please ensure you have the necessary dependencies installed.")
		fmt.Println("For more information, visit: https://www.mongodb.com/docs/compass/current/install/")
		fmt.Println("Enjoy using MongoDB Compass!")
	}
}