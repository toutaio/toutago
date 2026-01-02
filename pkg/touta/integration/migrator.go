package integration

import (
	sil "github.com/toutaio/toutago-sil-migrator/pkg/sil"
)

// MigratorConfig holds configuration for creating a migrator instance.
type MigratorConfig struct {
	Driver         string // mysql, postgres, sqlite, etc.
	DSN            string // Data source name
	MigrationsPath string // Path to migration files
	TableName      string // Migration table name (optional)
	Verbose        bool   // Enable verbose logging
}

// NewMigrator creates a new database migrator using sil-migrator.
//
// Example:
//
//	migrator, err := NewMigrator(&MigratorConfig{
//	    Driver: "mysql",
//	    DSN: "user:pass@tcp(localhost:3306)/dbname",
//	    MigrationsPath: "migrations",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Run migrations
//	if err := migrator.Up(ctx); err != nil {
//	    log.Fatal(err)
//	}
//
// Note: You need to provide the appropriate DatabaseAdapter for your driver.
// This is a simplified wrapper - for full control, use sil.NewMigrator directly.
func NewMigrator(config *MigratorConfig, adapter sil.DatabaseAdapter) (sil.Migrator, error) {
	silConfig := &sil.Config{
		MigrationsDir: config.MigrationsPath,
		TableName:     config.TableName,
	}
	
	if silConfig.TableName == "" {
		silConfig.TableName = "schema_migrations"
	}
	
	return sil.NewMigrator(silConfig, adapter)
}
