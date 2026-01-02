## ADDED Requirements

### Requirement: TemplateRenderer Interface

The framework SHALL provide a TemplateRenderer interface for rendering templates.

#### Scenario: Render template
- **WHEN** developer calls Render() with template name and data
- **THEN** the renderer SHALL execute the template
- **AND** return rendered output as bytes
- **AND** return an error for missing templates

#### Scenario: Register template function
- **WHEN** developer registers a custom function via RegisterFunction()
- **THEN** the function SHALL be available in all templates
- **AND** templates can call the function
- **AND** function SHALL receive template data as context

### Requirement: HTML Template Wrapper

The framework SHALL provide an html/template-based renderer implementation.

#### Scenario: Wrap Go html/template
- **WHEN** HTMLRenderer is instantiated
- **THEN** it SHALL wrap Go's html/template
- **AND** provide XSS protection automatically
- **AND** support all html/template features

#### Scenario: Template caching
- **WHEN** template is parsed
- **THEN** compiled template SHALL be cached
- **AND** subsequent renders SHALL use cached version
- **AND** cache SHALL be invalidated on file changes (in dev mode)

### Requirement: Template Hot Reload

The framework SHALL support template hot reload in development mode.

#### Scenario: Template file changes
- **WHEN** template file is modified in development mode
- **THEN** the renderer SHALL detect the change
- **AND** reparse the template
- **AND** next render SHALL use updated template

### Requirement: Built-in Template Functions

The framework SHALL provide common template helper functions.

#### Scenario: URL generation
- **WHEN** template uses url() function
- **THEN** it SHALL generate proper URLs for routes
- **AND** include configured base URL
- **AND** handle parameters correctly

#### Scenario: Asset URL generation
- **WHEN** template uses asset() function
- **THEN** it SHALL generate URLs for static assets
- **AND** include asset versioning for cache busting
- **AND** handle relative and absolute paths

#### Scenario: Safe HTML rendering
- **WHEN** template uses safe() function
- **THEN** HTML content SHALL not be escaped
- **AND** developer SHALL explicitly mark content as safe
- **AND** XSS protection SHALL still apply by default
