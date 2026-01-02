## 1. Template File Structure
- [x] 1.1 Create internal/cli/templates package
- [x] 1.2 Create templates/project directory for embedded templates
- [x] 1.3 Organize templates by type (docker, config, code)
- [x] 1.4 Add embed directives for template files

## 2. Template Loader Implementation
- [x] 2.1 Create TemplateLoader interface
- [x] 2.2 Implement embedded template loader
- [x] 2.3 Add template parsing and rendering functions
- [x] 2.4 Support variable substitution in templates

## 3. Extract Inline Templates
- [x] 3.1 Extract Dockerfile template
- [x] 3.2 Extract docker-compose.yml template
- [x] 3.3 Extract .dockerignore template
- [x] 3.4 Extract .air.toml template
- [x] 3.5 Extract main.go template
- [x] 3.6 Extract handler template
- [x] 3.7 Extract touta.yaml template

## 4. Refactor Commands
- [x] 4.1 Update createDockerFiles to use template loader
- [x] 4.2 Update initProject to use template loader
- [x] 4.3 Remove all inline template strings
- [x] 4.4 Simplify command logic to focus on orchestration

## 5. Testing and Validation
- [x] 5.1 Test template loading works
- [x] 5.2 Test new project creation
- [x] 5.3 Verify all files are generated correctly
- [x] 5.4 Ensure backward compatibility
- [x] 5.5 Validate code compiles and builds
