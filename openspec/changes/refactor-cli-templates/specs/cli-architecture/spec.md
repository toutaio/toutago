## ADDED Requirements

### Requirement: Template File Organization
The CLI SHALL organize project scaffolding templates in a dedicated directory structure that is embedded in the binary.

#### Scenario: Templates are embedded in binary
- **WHEN** the touta binary is built
- **THEN** all template files are embedded using Go embed directives and accessible at runtime

#### Scenario: Templates are organized by purpose
- **WHEN** a developer examines the template directory
- **THEN** templates are grouped logically (docker/, config/, code/) for easy discovery

### Requirement: Template Loader System
The CLI SHALL provide a template loading system that reads, parses, and renders templates with variable substitution.

#### Scenario: Load template from embedded filesystem
- **WHEN** the CLI needs to generate a file from a template
- **THEN** the template loader reads the template from the embedded filesystem

#### Scenario: Render template with variables
- **WHEN** a template contains variables like {{.ProjectName}}
- **THEN** the template loader substitutes variables with actual values

### Requirement: Separation of Command Logic and Templates
Command implementations SHALL contain only orchestration logic, with all content stored in template files.

#### Scenario: Commands are clean and focused
- **WHEN** a developer reads a command implementation
- **THEN** the code focuses on logic and flow, not multi-line string literals

#### Scenario: Template changes don't require code changes
- **WHEN** a template needs updating (e.g., new Dockerfile instruction)
- **THEN** only the template file is modified, not Go code

### Requirement: Template Discovery
Template files SHALL be stored in a well-known location that developers can easily find and modify.

#### Scenario: Developer finds template files
- **WHEN** a developer wants to customize project scaffolding
- **THEN** templates are in templates/project/ directory with clear naming

### Requirement: Backward Compatibility
The refactored CLI SHALL generate identical output to the previous inline implementation.

#### Scenario: New projects match previous structure
- **WHEN** running touta new project-name after refactoring
- **THEN** the generated project structure and file contents match the previous implementation

### Requirement: Template Variable Support
Templates SHALL support variable substitution for project-specific values.

#### Scenario: Project name substitution
- **WHEN** a template contains {{.ProjectName}} or similar placeholders
- **THEN** the actual project name is substituted when rendering

#### Scenario: Conditional template sections
- **WHEN** templates need optional content based on configuration
- **THEN** the template system supports conditional rendering

### Requirement: Error Handling
The template system SHALL provide clear error messages when templates are missing or fail to render.

#### Scenario: Missing template error
- **WHEN** a required template file is missing
- **THEN** a clear error message indicates which template is missing

#### Scenario: Template parse error
- **WHEN** a template has syntax errors
- **THEN** the error message shows the template name and line number
