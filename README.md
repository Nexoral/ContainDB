# ContainDB

ContainDB is a Go-based CLI tool that automates the installation and management of popular database containers (MongoDB, Redis, MySQL, PostgreSQL, Cassandra) and PHPMyAdmin.  

ContainDB was born out of frustration: installing MongoDB and other databases on Linux can be cumbersome. This package streamlines the process by providing a pre-built `.deb` installer. Simply download our `.deb` file from the [Releases](https://github.com/yourusername/ContainDB/releases) section, install it, and start using `containDB` to manage your database containers effortlessly.

## Features

- Automatically installs Docker if missing.
- Creates a custom Docker network (`ContainDB-Network`).
- Interactive prompts for selecting databases or PHPMyAdmin.
- Pulls required Docker images and runs containers with optional:
  - Custom or default port mappings.
  - Auto-restart policies.
  - Environment variables for root credentials.
- Links PHPMyAdmin to an existing MySQL/PostgreSQL container.

## Prerequisites

- A Unix-like OS with Bash.
- Internet connection for downloading images.

## Installation

Download the latest `.deb` package from the [Releases](https://github.com/yourusername/ContainDB/releases) page:

```bash
wget https://github.com/AnkanSaha/ContainDB/releases/download/vX.Y.Z/containDB_X.Y.Z_amd64.deb
sudo dpkg -i containDB_X.Y.Z_amd64.deb
```

## Usage

```bash
# Run the tool
sudo containDB
```

- Select a service from the menu.
- Follow prompts to configure ports, credentials, and restart policies.
- For PHPMyAdmin, select an existing SQL container and expose it on a host port.

## How It Works

1. **Docker Check & Install**  
   Uses `installation` package to verify Docker; installs if absent.

2. **Network Setup**  
   Ensures a Docker network `ContainDB-Network` exists for inter-container communication.

3. **Interactive Selection**  
   Leverages `promptui` for a user-friendly menu to choose a service.

4. **Image Pull & Run**  
   - Pulls the selected database image (or PHPMyAdmin).  
   - Queries user for port mapping and restart policy.  
   - Builds and executes a `docker run` command with the chosen options.

5. **PHPMyAdmin Linking**  
   - Detects running MySQL/PostgreSQL containers.  
   - Prompts for which container to link.  
   - Runs PHPMyAdmin on the same network, pointing `PMA_HOST` to the selected container.

Enjoy seamless container setups with ContainDB!
