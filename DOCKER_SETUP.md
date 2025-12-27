# Docker Setup Guide

This guide explains how to use Docker with Toutā for development and production.

## Quick Start

### Framework Development

```bash
# Start the development environment
docker-compose up

# The server starts on http://localhost:8080 with hot-reload
```

### New Project with Docker

```bash
# Create a new project (includes Docker files automatically)
touta new my-app
cd my-app

# Start with Docker
docker-compose up
```

## Docker Files Included

### Framework Repository

- `Dockerfile` - Multi-stage build with development and production targets
- `docker-compose.yml` - Development environment with hot-reload
- `.dockerignore` - Optimizes build context

### New Projects

When you run `touta new project-name`, these files are automatically created:

- `Dockerfile` - Multi-stage build for your project
- `docker-compose.yml` - Development environment
- `.dockerignore` - Build optimization
- `.air.toml` - Hot-reload configuration

## Docker Commands

### Development

```bash
# Start services
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild after dependency changes
docker-compose up --build
```

### Production Build

```bash
# Build production image
docker build --target production -t my-app:latest .

# Run production container
docker run -p 8080:8080 my-app:latest
```

## Features

### Hot-Reload

Both the framework and generated projects include automatic hot-reload using [air](https://github.com/cosmtrek/air):

- Code changes are detected automatically
- Application rebuilds and restarts
- No manual intervention needed

### Multi-Stage Builds

The Dockerfile uses multi-stage builds:

1. **Builder stage** - Compiles the Go application
2. **Development stage** - Includes hot-reload tools
3. **Production stage** - Minimal Alpine image with only the binary

### Volume Mounts

Development environments mount source code as volumes:

- Edit code on your host machine
- Changes appear instantly in the container
- No need to rebuild for code changes

### Go Module Caching

Docker volumes cache Go modules:

- Faster rebuilds
- Reduced bandwidth usage
- Persistent across container restarts

## Environment Variables

Configure your application using environment variables in `docker-compose.yml`:

```yaml
environment:
  - TOUTA_ENV=development
  - TOUTA_PORT=8080
  - TOUTA_HOST=0.0.0.0
```

## Troubleshooting

### Port Already in Use

```bash
# Change the port in docker-compose.yml
ports:
  - "3000:8080"  # Use port 3000 instead
```

### Permission Errors (Linux)

```bash
# Add your user to the docker group
sudo usermod -aG docker $USER
# Log out and back in
```

### Hot-Reload Not Working

```bash
# Ensure .air.toml exists
# Rebuild the container
docker-compose down
docker-compose up --build
```

### Build Fails

```bash
# Clear Docker cache
docker system prune -f
docker-compose up --build
```

## Production Deployment

### Build Optimized Image

```bash
docker build --target production -t my-app:v1.0.0 .
```

The production image:
- Uses Alpine Linux (minimal size)
- Includes only the compiled binary
- No development tools
- Optimized for security and performance

### Size Comparison

- **Development image**: ~500MB (includes Go tools, air)
- **Production image**: ~20MB (just binary + Alpine)

## Best Practices

1. **Use docker-compose for development** - Simplifies setup and configuration
2. **Use production target for deployment** - Minimal, secure images
3. **Mount source as volumes in dev** - Enable hot-reload
4. **Cache Go modules** - Faster builds
5. **Use .dockerignore** - Exclude unnecessary files

## Requirements

- Docker 20.10+
- Docker Compose 2.0+

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Air (Hot Reload)](https://github.com/cosmtrek/air)
- [Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)
