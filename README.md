<div align="center">
  <h1>ContainDB</h1>
  <p><strong>Database Container Management Made Simple</strong></p>
  <p>
    <a href="#installation">Installation</a> •
    <a href="#quick-start">Quick Start</a> •
    <a href="#features">Features</a> •
    <a href="#usage-examples">Usage</a> •
    <a href="#architecture">Architecture</a> •
    <a href="#troubleshooting">Troubleshooting</a>
  </p>
  
  ![Version](https://img.shields.io/badge/version-4.12.17--stable-blue)
  ![License](https://img.shields.io/badge/license-MIT-green)
  ![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue)
  ![Platform](https://img.shields.io/badge/platform-linux-lightgrey)
</div>

## The Problem ContainDB Solves

As developers, we often face these frustrating scenarios:

- Spending hours configuring database environments across different projects
- Dealing with conflicting versions of databases on our development machines
- Struggling with complex Docker commands for simple database tasks
- Managing database persistence, networking, and tools separately
- Lack of a unified interface for different database systems

ContainDB was born out of these pain points. I wanted a simple CLI tool that could handle all database container operations with minimal effort, allowing me to focus on actual development rather than environment setup.

## What is ContainDB?

ContainDB is an open-source CLI tool that automates the creation, management, and monitoring of containerized databases using Docker. It provides a simple, interactive interface for running popular databases and their management tools without needing to remember complex Docker commands or container configurations.

## Features

- **🚀 Instant Setup**: Get databases running in seconds with sensible defaults
- **🔄 Seamless Integration**: All databases run on the same Docker network for easy inter-container communication
- **💾 Data Persistence**: Optional volume management for data durability
- **🔐 Security Controls**: Interactive prompts for credentials and access control
- **🧩 Extensible Design**: Support for multiple database types and management tools
- **⚙️ Customization**: Configure ports, restart policies, and environment variables
- **📊 Management Tools**: One-click setup for phpMyAdmin, pgAdmin, and MongoDB Compass
- **🧹 Easy Cleanup**: Simple commands to remove containers, images, and volumes
- **🧠 Smart Detection**: Checks for existing resources to avoid conflicts

## Installation

### Option 1: Using the Debian Package (Recommended)

```bash
# Download latest .deb release
wget https://github.com/AnkanSaha/ContainDB/releases/download/v4.12.18-stable/containDB_4.12.18-stable_amd64.deb

# Install the package
sudo dpkg -i containDB_4.12.18-stable_amd64.deb
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/AnkanSaha/ContainDB.git
cd ContainDB

# Build the CLI
./Scripts/BinBuilder.sh

# Install binary to /usr/local/bin
sudo mv ./bin/containDB /usr/local/bin/
```

## Quick Start

Run ContainDB with root privileges:

```bash
sudo containDB
```

You'll be greeted with an attractive banner and a simple menu system that guides you through the process.

## Supported Databases & Tools

| Databases  | Management Tools |
|------------|-----------------|
| MongoDB    | MongoDB Compass  |
| MySQL      | phpMyAdmin       |
| PostgreSQL | pgAdmin          |
| MariaDB    | (uses phpMyAdmin)|
| Redis      |                  |
| Cassandra  |                  |

## Usage Examples

### Installing a Database

```bash
sudo containDB
# Select "Install Database"
# Choose your database (e.g., "mongodb")
# Follow the interactive prompts
```

### Connecting to Your Database

After installation, ContainDB provides you with connection details:

```
✅ PostgreSQL started! Access it at http://localhost:5432
Link it to your DB container 'postgresql-container' inside pgAdmin.
📋 Connection information:
   - Container name: postgresql-container
   - IP Address: 172.18.0.2
   - Port: 5432
🔐 pgAdmin login credentials:
   - Email: admin@local.com
   - Password: yourpassword
```

### Setting Up Management Tools

```bash
sudo containDB
# Select "Install Database"
# Choose "phpMyAdmin" or "PgAdmin" or "MongoDB Compass"
# Select the container to manage
# Follow the interactive prompts
```

### Managing Existing Resources

```bash
sudo containDB
# Select "List Databases" to see running containers
# Select "Remove Database" to stop and remove containers
# Select "Remove Image" to delete Docker images
# Select "Remove Volume" to delete persistent data volumes
```

## Architecture

ContainDB follows a layered architecture that separates concerns and promotes code organization.

```
                      User Interaction
                            ↓
      +-------------------------------------------+
      |             Main CLI Interface            |
      +-------------------------------------------+
                ↑                      ↑
                |                      |
       +------------------+   +------------------+
       |  Base Operations |   |   Tool Helpers   |
       +------------------+   +------------------+
                ↑                      ↑
                |                      |
       +------------------+   +------------------+
       | Docker Interface |   | System Utilities |
       +------------------+   +------------------+
                            ↓
                    Docker Engine & Host System
```

### How ContainDB Works Internally

---------------------------------------

1. **Network Creation**

ContainDB first ensures that a dedicated Docker network (`ContainDB-Network`) exists, which allows all containers to communicate with each other using container names as hostnames.

```
┌────────────────────────────┐     creates    ┌─────────────────────┐
│ ContainDB CLI              │ ───────────────> Docker Network      │
└────────────────────────────┘                └─────────────────────┘
```

---------------------------------------

2. **Container Orchestration**

When you select a database to install, ContainDB:
- Pulls the latest image (if needed)
- Checks for port conflicts
- Sets up necessary volumes for persistence
- Configures environment variables
- Creates and starts the container

```
┌────────────────────────────┐                ┌─────────────────────┐
│ ContainDB CLI              │                │ Docker Hub          │
└─────────────┬──────────────┘                └──────────┬──────────┘
              │                                          │
              │ 1. Pull image                            │
              ├──────────────────────────────────────────┘
              │
              │ 2. Create volume                ┌─────────────────────┐
              ├─────────────────────────────────> Data Volume         │
              │                                 └─────────────────────┘
              │
              │ 3. Run container               ┌──────────────────────┐
              └────────────────────────────────> Database Container   │
                                               └──────────────────────┘
```

---------------------------------------

3. **Management Tool Integration**

For tools like phpMyAdmin, pgAdmin, or MongoDB Compass, ContainDB handles:
- Tool installation and configuration
- Linking to the appropriate database container
- Providing connection details and credentials

```
┌────────────────────────────┐                ┌─────────────────────┐
│ ContainDB CLI              │ ─────────────┐ │ Database Container  │
└────────────────────────────┘              │ └─────────────────────┘
                                            │
                                            │ links
                                            ▼
                                  ┌─────────────────────┐
                                  │ Management Tool     │
                                  │ Container/App       │
                                  └─────────────────────┘
```

---------------------------------------

## The Origin Story

ContainDB was created after I found myself repeatedly setting up the same database environments across different projects. The challenges included:

1. Remembering specific Docker commands for each database
2. Managing network connectivity between containers
3. Ensuring data persistence with proper volume configuration
4. Setting up administration tools for each database type

The final straw came when I needed to set up a multi-database project that required MongoDB, PostgreSQL, and Redis - all with different configurations, management tools, and persistence requirements. I spent hours on environment setup instead of actual coding.

I realized that these repetitive tasks could be automated, and ContainDB was born. What started as a personal script evolved into a comprehensive tool that I now use daily and want to share with the developer community.

## How ContainDB Helps in Daily Development

ContainDB has become an essential part of my development workflow by:

- **Saving Setup Time**: What used to take 30+ minutes now takes seconds
- **Standardizing Environments**: Ensuring consistent database setups across projects
- **Simplifying Management**: Providing easy access to admin tools and interfaces
- **Isolating Services**: Preventing conflicts between different database versions
- **Managing Resources**: Making cleanup and maintenance straightforward

Real-world example: When working on a new microservice project, I can spin up a PostgreSQL instance, link it to pgAdmin, and have a fully functional development environment in less than a minute - all with proper network configuration and persistence.

## Troubleshooting

### Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| **"Permission Denied"** | Ensure you run ContainDB with `sudo` |
| **"Docker Not Found"** | Let ContainDB install Docker or run with `--install-docker` flag |
| **"Port Already in Use"** | Choose a different port when prompted |
| **"Volume Already Exists"** | Select to reuse or recreate the volume |
| **"Cannot Connect to Database"** | Check network settings and credentials |

### Debug Information

If you encounter issues, run:

```bash
sudo containDB --debug
```

This will provide verbose output to help diagnose problems.

## Contributing

ContainDB is an open-source project that welcomes contributions from everyone. Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit your changes**: `git commit -m 'Add amazing feature'`
4. **Push to the branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request`

For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- The Docker team for creating an amazing containerization platform
- The Go community for providing excellent libraries and tools
- All contributors who have helped improve ContainDB

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/AnkanSaha">Ankan Saha</a></p>
  <p>
    <a href="https://github.com/AnkanSaha/ContainDB/stargazers">⭐ Star this project</a> •
    <a href="https://github.com/AnkanSaha/ContainDB/issues">🐞 Report Bug</a> •
    <a href="https://github.com/AnkanSaha/ContainDB/issues">✨ Request Feature</a>
  </p>
</div>
