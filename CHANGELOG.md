# Changelog

All notable changes to ContainDB will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [4.12.18-stable] - 2025-06-29

### Added
- Auto-rollback feature when container setup fails
- Support for exporting running containers as Docker Compose configuration
- Added detailed connection information after database setup

### Fixed
- Fixed issue with MongoDB container error handling on Debian systems
- Improved error handling for port conflicts
- Fixed volume mount permissions for PostgreSQL containers

### Changed
- Improved user interface with clearer prompts
- Enhanced cleanup process for temporary files
- Updated container health check implementation

## [4.12.17-stable] - 2025-05-15

### Added
- Support for Redis Insight as a management tool
- Added ability to configure custom network settings
- Improved volume management with persistence options

### Fixed
- Fixed container networking issues between databases
- Resolved Docker detection on newer Linux distributions
- Fixed credential handling for PostgreSQL containers

### Changed
- Enhanced error messages for better troubleshooting
- Improved Docker installation helper
- Updated default container configurations

## [4.11.0-stable] - 2025-03-20

### Added
- Initial support for Docker Compose export
- Added MongoDB Compass installation
- Implemented container health checks

### Fixed
- Fixed permission issues with volume mounts
- Improved error recovery during setup
- Fixed port detection for running containers

### Changed
- Refactored codebase for better maintainability
- Updated default configuration for MySQL containers
- Enhanced user experience with better progress indicators

## [4.0.0-stable] - 2024-12-15

### Added
- Initial public release with support for MySQL, PostgreSQL, MongoDB, MariaDB, and Redis
- Management tools integration (phpMyAdmin, pgAdmin)
- Container lifecycle management
- Volume persistence options
- Network configuration

### Security
- Secure credential management
- Isolated container network
