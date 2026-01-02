// Package main demonstrates using the Fíth template renderer with the Toutā framework.
package main

import (
	"log"

	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago/pkg/touta"
	"github.com/toutaio/toutago/pkg/touta/integration"
)

func main() {
	// Create DI container
	container := integration.NewContainer()
	
	// Create template renderer using fith
	renderer, err := integration.NewTemplateRenderer(&fith.Config{
		TemplateDir: "templates",
		Extensions:  []string{".html", ".fith"},
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Bind renderer to container
	container.Singleton((*touta.TemplateRenderer)(nil), renderer)
	
	// Create router
	router := integration.NewRouter(container)
	
	// Home page with template
	router.GET("/", func(ctx touta.Context) error {
		data := map[string]interface{}{
			"Title": "Welcome to Toutā",
			"Name":  "Developer",
			"Features": []string{
				"Integrated component architecture",
				"Fíth template engine",
				"Cosan HTTP router",
				"Nasc dependency injection",
			},
		}
		
		output, err := renderer.Render("home", data)
		if err != nil {
			return err
		}
		
		return ctx.HTML(200, string(output))
	})
	
	log.Println("Server starting on :8080")
	log.Println("Visit http://localhost:8080/")
	log.Fatal(router.Listen(":8080"))
}
