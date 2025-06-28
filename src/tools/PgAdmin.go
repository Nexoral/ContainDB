package tools

import (
    "ContainDB/src/Docker"
    "fmt"
    "os"
    "os/exec"

    "github.com/manifoldco/promptui"
)

func StartPgAdmin() {
    // 1️⃣ Check if pgAdmin is already running
    if Docker.IsContainerRunning("pgadmin", true) {
        fmt.Println("pgAdmin container is already running.")
        if Docker.AskYesNo("Remove existing pgAdmin container and recreate?") {
            fmt.Println("Removing existing pgAdmin container...")
            cmd := exec.Command("docker", "rm", "-f", "pgadmin")
            cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
            if err := cmd.Run(); err != nil {
                fmt.Println("Error removing pgAdmin:", err)
                return
            }
            fmt.Println("Removed existing pgAdmin container.")
        } else {
            fmt.Println("Keeping existing container. Aborting setup.")
            return
        }
    }

    // 2️⃣ List running SQL containers to link with
    networks := Docker.ListOfContainers([]string{"postgres"})
    if len(networks) == 0 {
        fmt.Println("No running SQL containers (Postgres/MySQL/MariaDB) found.")
        return
    }
    items := append(networks, "Exit")
    prompt := promptui.Select{
        Label: "Select a DB container to link with pgAdmin",
        Items: items,
    }
    _, selected, err := prompt.Run()
    if err != nil {
        fmt.Println("\n⚠️ Interrupted. Rolling back...")
        return
    }
    if selected == "Exit" {
        fmt.Println("Exiting pgAdmin setup.")
        return
    }

    // 3️⃣ Ask port and credentials
    port := AskForInput("Enter host port for pgAdmin (e.g. 5050)", "5050")
    email := AskForInput("Enter PGADMIN_DEFAULT_EMAIL", "admin@local.com")
    password := AskForInput("Enter PGADMIN_DEFAULT_PASSWORD", "")

    // 4️⃣ Pull image
    fmt.Println("Pulling pgAdmin Docker image...")
    cmd := exec.Command("docker", "pull", "dpage/pgadmin4:latest")
    cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
    _ = cmd.Run()

    // 5️⃣ Run container
    fmt.Println("Creating pgAdmin container...")
    cmd = exec.Command("docker", "run",
        "-d",
        "--restart", "unless-stopped",
        "--network", "ContainDB-Network",
        "--name", "pgadmin",
        "-e", fmt.Sprintf("PGADMIN_DEFAULT_EMAIL=%s", email),
        "-e", fmt.Sprintf("PGADMIN_DEFAULT_PASSWORD=%s", password),
        "-p", fmt.Sprintf("%s:80", port),
        "dpage/pgadmin4:latest",
    )
    cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
    if err := cmd.Run(); err != nil {
        fmt.Println("Error starting pgAdmin:", err)
    } else {
        fmt.Printf("✅ pgAdmin started! Access it at http://localhost:%s\n", port)
        fmt.Printf("Link it to your DB container '%s' inside pgAdmin.\n", selected)
    }
}
