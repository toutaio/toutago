// Package integration provides adapters that integrate standalone Toutā component
// libraries with the main framework interfaces.
//
// This package bridges the gap between the framework's core interfaces (defined in
// pkg/touta/interfaces.go) and the specialized, production-ready component implementations:
//
//   - nasc: Dependency injection container
//   - cosan: HTTP router with path parameters and middleware
//   - fith: Template engine with Jinja2-style syntax
//   - datamapper: Database abstraction layer with multiple adapters
//   - sil: Database migration tool
//
// # Design Philosophy
//
// The integration layer follows the Adapter pattern to maintain loose coupling between
// the framework and component implementations. This allows:
//
//   - Components to evolve independently
//   - Easy swapping of implementations
//   - Clear separation between framework API and implementation details
//   - Component libraries to be used standalone in other projects
//
// # Basic Usage
//
//	import "github.com/toutaio/toutago/pkg/touta/integration"
//
//	// Create DI container
//	container := integration.NewContainer()
//	container.Bind((*Logger)(nil), &ConsoleLogger{})
//
//	// Create router
//	router := integration.NewRouter(container)
//	router.GET("/", HomeHandler)
//
//	// Create template renderer
//	renderer, _ := integration.NewTemplateRenderer(&fith.Config{
//	    TemplateDir: "templates",
//	})
//
// # Advanced Usage
//
// For scenarios where you need direct access to the underlying implementations,
// each adapter provides a Native() method:
//
//	adapter := integration.NewContainer().(*integration.NascContainerAdapter)
//	nascContainer := adapter.Native()
//	// Use nasc-specific features
//
// # Database Integration
//
// The datamapper integration supports multiple database backends through a plugin
// architecture:
//
//	// MySQL
//	import _ "github.com/toutaio/toutago-datamapper-mysql"
//	mapper, _ := integration.NewDataMapper("mysql", map[string]interface{}{
//	    "dsn": "user:pass@tcp(localhost:3306)/db",
//	})
//
//	// PostgreSQL
//	import _ "github.com/toutaio/toutago-datamapper-postgres"
//	mapper, _ := integration.NewDataMapper("postgres", map[string]interface{}{
//	    "dsn": "postgres://user:pass@localhost:5432/db",
//	})
//
// # Migrations
//
// Database migrations are handled by the sil-migrator integration:
//
//	migrator, _ := integration.NewMigrator(&integration.MigratorConfig{
//	    Driver: "mysql",
//	    DSN: "user:pass@tcp(localhost:3306)/db",
//	    MigrationsPath: "migrations",
//	})
//	migrator.Up()
package integration
