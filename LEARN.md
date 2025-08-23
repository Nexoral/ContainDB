# ContainDB - Learning Guide

## Introduction

ContainDB is a CLI tool that simplifies the management of containerized database systems using Docker. It provides an intuitive interface for deploying, configuring, and managing various database containers, eliminating the need to remember complex Docker commands and configuration options.

## Core Concepts

### Database Containerization
ContainDB uses Docker to containerize databases, providing isolation, portability, and consistent environments regardless of the host system. This solves common issues like conflicting dependencies and simplifies database setup.

### Container Networking
All database containers created by ContainDB run on a dedicated Docker network (`ContainDB-Network`), allowing seamless communication between databases and management tools using container names as hostnames.

### Data Persistence
ContainDB automatically sets up Docker volumes for data persistence, ensuring your database data survives container restarts and removals.

## Supported Databases and Tools

### Databases
- MongoDB
- MySQL
- PostgreSQL
- MariaDB
- Redis

### Management Tools
- MongoDB Compass
- phpMyAdmin (for MySQL/MariaDB)
- pgAdmin (for PostgreSQL)
- RedisInsight

## Key Features

1. **Simple CLI Interface**: Interactive menus guide you through database setup and management
2. **Automated Container Setup**: Handles Docker image pulling, container creation, and networking
3. **Management Tool Integration**: One-click setup for database administration interfaces
4. **Data Persistence**: Automatic volume management for data durability
5. **Docker Compose Export**: Export your database configurations as docker-compose.yml
6. **Docker Compose Import**: Import and deploy services from existing compose files
7. **Intelligent Conflict Resolution**: Detects and manages port conflicts automatically
8. **Auto-Rollback**: Cleans up resources if errors occur during installation

## Getting Started

### Installation

ContainDB can be installed using the Debian package or built from source.

#### Using the Debian Package
```bash
# Download latest .deb release
wget https://github.com/AnkanSaha/ContainDB/releases/download/v5.12.30-stable/containDB_5.12.30-stable_amd64.deb

# Install the package
sudo dpkg -i containDB_5.12.30-stable_amd64.deb
```

#### Build from Source
```bash
# Clone the repository
git clone https://github.com/AnkanSaha/ContainDB.git
cd ContainDB

# Build the CLI
./Scripts/BinBuilder.sh

# Install binary to /usr/local/bin
sudo mv ./bin/containDB /usr/local/bin/
```

### Basic Usage

ContainDB requires root privileges to manage Docker resources:

```bash
sudo containDB
```

This will display the main menu where you can:
- Install databases
- Install management tools
- List running databases
- Remove containers, images, or volumes
- Export or import Docker Compose configurations

## Command-line Options

ContainDB supports several command-line options:

```bash
sudo containdb --version      # Display version information
sudo containdb --help         # Show help message
sudo containdb --install-docker  # Install Docker if not present
sudo containdb --export       # Export Docker Compose configuration
sudo containdb --import ./docker-compose.yml  # Import services from Docker Compose
```

## Workflow Examples

### Installing MongoDB with Compass

1. Start ContainDB:
   ```bash
   sudo containdb
   ```

2. Select "Install Database" from the main menu
3. Choose "MongoDB" from the database options
4. Follow the prompts to set port, credentials, and persistence options
5. After MongoDB is installed, return to the main menu
6. Select "Install Database" again
7. Choose "MongoDB Compass" from the options
8. Connect Compass to your MongoDB container using the provided connection details

### Exporting Your Database Environment

After setting up your databases, you can export the configuration:

1. Run ContainDB with the export flag:
   ```bash
   sudo containdb --export
   ```
2. A `docker-compose.yml` file will be created in your current directory
3. This file contains all the configuration needed to recreate your database environment

## Architecture

ContainDB follows a layered architecture:

1. **Core CLI Interface**: The main entry point that handles user interaction
2. **Base Operations**: Contains the primary functionality and menu system
3. **Docker Interface**: Abstracts Docker operations (container management, networking, etc.)
4. **System Utilities**: Provides helper functions for system interactions

The application creates a Docker network and manages containers, volumes, and images through the Docker Engine API.

## Docker Integration

ContainDB uses Docker's capabilities to:
- Create and manage a dedicated network
- Pull database images
- Create containers with appropriate configurations
- Mount volumes for data persistence
- Manage container lifecycle (start, stop, remove)
- Connect management tools to databases

## Common Use Cases

### Development Environment Setup
Quickly set up databases for new projects without complex configuration.

### Multi-Database Applications
Deploy and manage multiple database systems that need to work together.

### Database Management
Use integrated management tools to visualize and manipulate your data.

### Environment Portability
Export configurations to move your database setup between machines.

### Database Testing
Easily create isolated databases for testing different scenarios.

## Best Practices

1. **Data Backups**: Even with persistent volumes, regular backups are recommended
2. **Resource Management**: Monitor container resource usage (RAM, CPU)
3. **Security**: Change default passwords for production-like environments
4. **Network Access**: Restrict port exposure for sensitive databases

## Troubleshooting

### Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| Permission Denied | Run ContainDB with `sudo` |
| Docker Not Found | Use `--install-docker` flag or install manually |
| Port Conflicts | Choose alternative ports when prompted |
| Container Fails to Start | Check logs with `docker logs [container-name]` |
| Volume Already Exists | Choose to reuse or recreate when prompted |

### Logs and Debugging

ContainDB provides feedback on operations. For deeper issues, Docker logs can be accessed:

```bash
docker logs [container-name]
```

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## License

ContainDB is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
