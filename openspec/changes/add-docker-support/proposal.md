# Change: Add Docker Support

## Why
Enable developers to work with Toutā in containerized environments, simplifying development setup and ensuring consistency across development, testing, and production environments. This eliminates "works on my machine" issues and provides a standardized development experience.

## What Changes
- Add base Dockerfile with Go 1.21+ for developers starting from scratch
- Add docker-compose.yml for the Toutā framework repository to simplify development
- Modify `touta new` command to include docker-compose.yml in generated projects
- Include .dockerignore file for optimization
- Add hot-reload support in Docker development environment

## Impact
- Affected specs: docker-support (new capability)
- Affected code: 
  - internal/cli/commands.go (modify initProject function)
  - Root directory (add Dockerfile, docker-compose.yml, .dockerignore)
  - Project templates (include docker-compose.yml in new projects)
- Benefits:
  - Faster onboarding for new developers
  - Consistent development environments
  - Easier CI/CD integration
  - Production-ready containerization from day one
