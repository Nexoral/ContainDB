# ContainDB

ContainDB is a CLI tool written in Go that automates the installation and management of popular database containers (MongoDB, Redis, MySQL, PostgreSQL, Cassandra) as well as PHPMyAdmin and MongoDB Compass. It handles Docker installation (if missing), network setup, container orchestration, optional port mapping, restart policies, and data persistence.

---

## Features

- ✅ Automatic Docker installation and setup  
- 🌐 Dedicated Docker network (`ContainDB-Network`) for seamless inter-container communication  
- 📦 Pulls and runs database images with sensible defaults  
- 🔌 Interactive prompts for:
  - Port mapping (default or custom)
  - Auto-restart policies
  - Environment variables for credentials
  - Volume persistence and management
- 🔗 One-click PHPMyAdmin linking to existing SQL containers  
- 🧭 Optional MongoDB Compass download and install  

---

## Prerequisites

- Linux or macOS (with Bash)  
- Internet connection to download Docker images  

---

## Installation

Option 1: Download and install via Debian package  
```bash
# download latest .deb from Releases
wget https://github.com/AnkanSaha/ContainDB/releases/download/v3.11.18-stable/containDB_3.11.18-stable_amd64.deb

# install the package
sudo dpkg -i containDB_3.11.18-stable_amd64.deb
```

Option 2: Build from source  
```bash
# clone the repository
git clone https://github.com/AnkanSaha/ContainDB.git
cd ContainDB

# build the CLI
./Scripts/BinBuilder.sh

# install binary to /usr/local/bin
sudo mv ./bin/containDB /usr/local/bin/
```

---

## Quick Start

Run the CLI with root privileges:
```bash
sudo containDB
```
1. **Welcome Banner** – Displays version and basic info.  
2. **Main Menu** – Choose one of:
   - Install Database  
   - List Databases  
   - Remove Database  
   - Exit  
3. **Follow Prompts** – Configure and launch containers in a few keystrokes.

---

## Usage Examples

### 1. Install a Database
```bash
sudo containDB
# Select "Install Database"
# Pick "mysql" (or any supported service)
# Answer port mapping, restart policy, persistence, credentials
```

### 2. Launch PHPMyAdmin
```bash
sudo containDB
# Select "Install Database"
# Choose "phpmyadmin"
# Pick an existing SQL container and assign host port
```

### 3. Install MongoDB Compass
```bash
sudo containDB
# Select "Install Database"
# Choose "MongoDB Compass"
```

### 4. List Running Containers
```bash
sudo containDB --version
sudo containDB
# Select "List Databases"
```

### 5. Remove a Database
```bash
sudo containDB
# Select "Remove Database"
# Pick the container to delete (with optional volume cleanup)
```

---

## Flags

- `--version` : Print ContainDB CLI version and exit.

---

## Troubleshooting

- **Permission Denied**  
  Ensure you run with `sudo`.  
- **Docker Not Found**  
  The CLI auto-installs Docker on first run; please restart your terminal after installation.  
- **Port Already in Use**  
  Choose a different host port when prompted.  
- **Volume Conflicts**  
  On existing volumes, you can opt to reuse or recreate fresh ones.  

---

## Contributing

1. Fork the repo  
2. Create your feature branch (`git checkout -b feature/…`)  
3. Commit your changes (`git commit -m 'Add …'`)  
4. Push to the branch (`git push origin feature/…`)  
5. Open a Pull Request  

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
