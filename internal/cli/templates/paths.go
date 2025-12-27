package templates

// Template file paths
const (
	// Docker templates
	TemplateDockerfile       = "project/docker/Dockerfile.tmpl"
	TemplateDockerCompose    = "project/docker/docker-compose.yml.tmpl"
	TemplateDockerIgnore     = "project/docker/dockerignore.tmpl"
	
	// Config templates
	TemplateAirConfig        = "project/config/air.toml.tmpl"
	TemplateToutaConfig      = "project/config/touta.yaml.tmpl"
	
	// Code templates
	TemplateMainGo           = "project/code/main.go.tmpl"
	TemplateHelloHandler     = "project/code/hello.go.tmpl"
)

// Template mapping for quick reference
var TemplateMap = map[string]string{
	"Dockerfile":         TemplateDockerfile,
	"docker-compose.yml": TemplateDockerCompose,
	".dockerignore":      TemplateDockerIgnore,
	".air.toml":          TemplateAirConfig,
	"touta.yaml":         TemplateToutaConfig,
	"main.go":            TemplateMainGo,
	"hello.go":           TemplateHelloHandler,
}
