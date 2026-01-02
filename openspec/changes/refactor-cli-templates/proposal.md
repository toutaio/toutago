# Change: Refactor CLI to Use Template System

## Why
The `internal/cli/commands.go` file has grown to 420 lines with large inline strings for Dockerfiles, docker-compose files, and other project templates. This makes the code hard to maintain, test, and extend. Extracting templates to separate files will improve code clarity, enable easier template updates, and follow the clean architecture principles of the framework.

## What Changes
- Extract all inline template strings to separate template files
- Create a `templates/` directory structure for project scaffolding templates
- Implement a template loading system for CLI commands
- Separate concerns: commands handle logic, templates handle content
- Make templates easily discoverable and editable
- Support future template customization and extension

## Impact
- Affected specs: cli-architecture (new capability)
- Affected code:
  - internal/cli/commands.go (refactor to use template loader)
  - internal/cli/templates/ (new package for template management)
  - templates/project/ (new directory for project templates)
- Benefits:
  - Cleaner, more maintainable code
  - Easier to customize project templates
  - Better separation of concerns
  - Simpler testing of template logic
  - Future support for custom template sets
