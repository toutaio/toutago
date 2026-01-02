## ADDED Requirements

### Requirement: Base Dockerfile for Development
The framework repository SHALL provide a Dockerfile that includes the latest stable Go version (1.21+) and all necessary dependencies for Toutā development.

#### Scenario: Developer builds base image
- **WHEN** a developer runs `docker build -t touta-dev .` in the framework repository
- **THEN** the image builds successfully with Go 1.21+, air for hot-reload, and all framework dependencies installed

#### Scenario: Multi-stage build optimization
- **WHEN** building for production with `docker build --target production`
- **THEN** the resulting image contains only the compiled binary and runtime dependencies, minimizing image size

### Requirement: Framework Development Docker Compose
The framework repository SHALL include a docker-compose.yml file that provides a complete development environment with hot-reload capabilities.

#### Scenario: Start framework development environment
- **WHEN** a developer runs `docker-compose up` in the framework repository
- **THEN** the development server starts on localhost:8080 with hot-reload enabled and source code mounted as volumes

#### Scenario: Code changes trigger hot-reload
- **WHEN** a developer modifies Go source files while docker-compose is running
- **THEN** the application automatically rebuilds and restarts within the container

### Requirement: Project Template Docker Files
The `touta new` command SHALL generate Docker configuration files (Dockerfile, docker-compose.yml, .dockerignore) in new projects.

#### Scenario: Create new project with Docker support
- **WHEN** a developer runs `touta new myproject`
- **THEN** the generated project includes Dockerfile, docker-compose.yml, and .dockerignore files configured for the project

#### Scenario: New project starts with Docker
- **WHEN** a developer runs `docker-compose up` in a newly created project directory
- **THEN** the project server starts successfully on the configured port with hot-reload enabled

### Requirement: Docker Ignore Configuration
All Docker configurations SHALL include .dockerignore files to optimize build context and reduce image size.

#### Scenario: Exclude unnecessary files from build context
- **WHEN** Docker builds an image
- **THEN** .git, node_modules, test files, and documentation are excluded from the build context per .dockerignore rules

### Requirement: Environment Variable Configuration
Docker configurations SHALL support environment-based configuration for flexibility across development, staging, and production environments.

#### Scenario: Override configuration via environment
- **WHEN** TOUTA_PORT environment variable is set in docker-compose.yml
- **THEN** the server binds to the specified port instead of the default 8080

#### Scenario: Development mode configuration
- **WHEN** TOUTA_ENV=development is set
- **THEN** hot-reload is enabled and debug logging is active

### Requirement: Volume Mounts for Development
Development docker-compose configurations SHALL mount source code as volumes to enable live code editing without rebuilding containers.

#### Scenario: Edit code without container rebuild
- **WHEN** a developer modifies source files on the host machine
- **THEN** changes are immediately visible inside the container and trigger hot-reload without rebuilding the image

### Requirement: Documentation Integration
Docker setup and usage SHALL be documented in project README and QUICKSTART guides.

#### Scenario: New developer follows Docker setup
- **WHEN** a developer reads the README Docker section
- **THEN** they can successfully build and run the framework or a new project using Docker with clear, step-by-step instructions

#### Scenario: Troubleshooting guidance
- **WHEN** a developer encounters common Docker issues (port conflicts, permission errors)
- **THEN** the documentation provides troubleshooting steps and solutions
