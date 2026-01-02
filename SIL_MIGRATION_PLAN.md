# Síl Migration and Seeding System - Implementation Plan

## Overview

**Síl** (Old Irish: "seed" or "lineage") is an independent database migration and seeding tool for Go projects, designed to work standalone or integrated with Toutā.

## Repository Information

- **GitHub Repository**: https://github.com/toutaio/toutago-sil-migrator
- **Development Path**: /home/nestor/Proyects/toutago-sil-migrator
- **Package Import**: github.com/toutaio/toutago-sil-migrator
- **CLI Binary**: `sil`

## Core Principles

1. **SOLID Compliance**: Every component follows Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, and Dependency Inversion principles
2. **Database Agnostic**: Interface-driven design supports PostgreSQL, MySQL, SQLite, and custom adapters
3. **Zero Toutā Dependencies**: Works as standalone library in any Go project
4. **Production Ready**: Migration locking, transaction support, error recovery
5. **Developer Friendly**: Intuitive CLI, helpful errors, comprehensive docs

## Inspired By Best Practices From

- **Rails ActiveRecord Migrations**: Sequential versioning, up/down migrations
- **Laravel Migrations**: Fluent schema builder, rollback capabilities  
- **Alembic (Python)**: Dependency management, migration paths
- **Flyway**: Version-based migrations with checksums
- **Knex.js**: Transaction-wrapped migrations, seed management

## Phase-Based Implementation

### Phase 1: Foundation (Weeks 1-3)
**Deliverables**: Core migration engine + PostgreSQL adapter

- Version-based migration system
- PostgreSQL database adapter
- Transaction support with auto-rollback
- Migration locking (prevent concurrent runs)
- Basic CLI: `migrate`, `rollback`, `status`, `create`
- 80%+ test coverage
- Basic documentation

**Key Commands**:
```bash
sil init                    # Initialize migrations directory
sil create add_users_table  # Create new migration
sil migrate                 # Run all pending migrations
sil rollback                # Rollback last batch
sil status                  # Show migration status
```

### Phase 2: Multi-Database (Weeks 4-5)
**Deliverables**: MySQL + SQLite support

- MySQL adapter with GET_LOCK() locking
- SQLite adapter with file-based locking
- Adapter factory and auto-detection
- Database-specific documentation
- Cross-database testing
- 85%+ test coverage

**Support Matrix**:
| Database   | Transaction DDL | Locking Method         |
|------------|-----------------|------------------------|
| PostgreSQL | ✅ Yes          | pg_advisory_lock()     |
| MySQL      | ❌ No (implicit)| GET_LOCK()             |
| SQLite     | ⚠️ Limited      | File-based locking     |

### Phase 3: Seeding System (Weeks 6-7)
**Deliverables**: Data seeding with dependencies

- Seeder interface and engine
- Dependency graph resolution
- Topological sort for execution order
- Idempotent seeder support
- Environment-specific seeds (dev/test/staging/prod)
- Seeder CLI commands
- 90%+ test coverage

**Key Commands**:
```bash
sil seed:create UserSeeder      # Create new seeder
sil seed:run --all              # Run all seeders
sil seed:run UserSeeder         # Run specific seeder
sil seed:run --env=development  # Run dev-only seeders
sil seed:status                 # Show seeder status
```

### Phase 4: Advanced Features (Weeks 8-9)
**Deliverables**: Production optimization

- Dry-run mode (preview without executing)
- Migration squashing (combine multiple migrations)
- Schema builder helpers (create_table, add_column, etc.)
- Progress reporting for long operations
- Performance optimization
- Programmatic API with callbacks
- Comprehensive examples
- Complete documentation

**Advanced Commands**:
```bash
sil migrate --dry-run           # Preview migrations
sil migrate --steps=3           # Run specific number
sil squash 001..005             # Combine migrations
sil lock:status                 # Check lock status
sil lock:release                # Emergency unlock
```

## Architecture

### Core Interfaces

```go
// Migration represents a database migration
type Migration interface {
    Version() string
    Description() string
    Up(adapter DatabaseAdapter) error
    Down(adapter DatabaseAdapter) error
}

// DatabaseAdapter handles database-specific operations
type DatabaseAdapter interface {
    Connect(config Config) error
    Close() error
    Exec(query string, args ...interface{}) error
    Query(query string, args ...interface{}) (Rows, error)
    BeginTx() (Transaction, error)
    CreateMigrationsTable() error
    GetAppliedMigrations() ([]MigrationRecord, error)
    RecordMigration(version, description string) error
    RemoveMigration(version string) error
    Lock() (Lock, error)
}

// Seeder represents a data seeder
type Seeder interface {
    Name() string
    Dependencies() []string
    Seed(adapter DatabaseAdapter) error
    ShouldRun(adapter DatabaseAdapter) (bool, error)
}
```

### Directory Structure

```
toutago-sil-migrator/
├── cmd/sil/                    # CLI binary
├── pkg/sil/
│   ├── migrator.go             # Core engine
│   ├── seeder.go               # Seeding engine
│   ├── interfaces.go           # All interfaces
│   ├── config.go               # Configuration
│   ├── lock.go                 # Locking mechanism
│   └── adapters/
│       ├── postgres.go         # PostgreSQL
│       ├── mysql.go            # MySQL
│       ├── sqlite.go           # SQLite
│       └── mock.go             # Testing
├── examples/
│   ├── basic/                  # Basic usage
│   ├── multi-db/               # Multi-database
│   └── with-touta/             # Toutā integration
└── tests/
    ├── unit/
    ├── integration/
    └── e2e/
```

## Migration Flow

```
1. Initialize Migrator
   ↓
2. Acquire migration lock (database-level advisory lock)
   ↓
3. Load migration files from disk
   ↓
4. Query applied migrations from database
   ↓
5. Calculate pending migrations (not yet applied)
   ↓
6. For each pending migration:
   ├─ Begin database transaction
   ├─ Execute Up() function
   ├─ Record migration in migrations table
   ├─ Commit transaction
   └─ (on error) Rollback transaction & stop
   ↓
7. Release migration lock
```

## Seeding Flow

```
1. Initialize SeedManager
   ↓
2. Load all seeder files
   ↓
3. Build dependency graph from Dependencies()
   ↓
4. Topological sort (detect circular deps)
   ↓
5. For each seeder (in dependency order):
   ├─ Check ShouldRun() for idempotency
   ├─ (if true) Execute Seed() function
   └─ Record execution timestamp
   ↓
6. Report results
```

## Example Usage

### Creating a Migration

```go
// 20251229120000_create_users_table.go
package migrations

import "github.com/toutaio/toutago-sil-migrator/pkg/sil"

type CreateUsersTable struct {
    sil.BaseMigration
}

func (m *CreateUsersTable) Up(db sil.DatabaseAdapter) error {
    return db.Exec(`
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
}

func (m *CreateUsersTable) Down(db sil.DatabaseAdapter) error {
    return db.Exec(`DROP TABLE users`)
}
```

### Creating a Seeder

```go
// user_seeder.go
package seeders

import "github.com/toutaio/toutago-sil-migrator/pkg/sil"

type UserSeeder struct {
    sil.BaseSeeder
}

func (s *UserSeeder) Name() string {
    return "UserSeeder"
}

func (s *UserSeeder) Dependencies() []string {
    return []string{} // No dependencies
}

func (s *UserSeeder) Seed(db sil.DatabaseAdapter) error {
    users := []string{
        "admin@example.com",
        "user@example.com",
    }
    
    for _, email := range users {
        err := db.Exec(
            "INSERT INTO users (email, name) VALUES (?, ?)",
            email, email,
        )
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *UserSeeder) ShouldRun(db sil.DatabaseAdapter) (bool, error) {
    // Check if already seeded
    rows, err := db.Query("SELECT COUNT(*) FROM users")
    if err != nil {
        return false, err
    }
    defer rows.Close()
    
    var count int
    rows.Next()
    rows.Scan(&count)
    
    return count == 0, nil // Only run if table is empty
}
```

### Programmatic API Usage

```go
package main

import (
    "github.com/toutaio/toutago-sil-migrator/pkg/sil"
    "github.com/toutaio/toutago-sil-migrator/pkg/sil/adapters"
)

func main() {
    // Configure
    config := sil.Config{
        Database: sil.DatabaseConfig{
            Driver: "postgres",
            DSN:    "postgres://user:pass@localhost/mydb",
        },
        MigrationsDir: "./migrations",
        LockTimeout:   300, // 5 minutes
    }
    
    // Create adapter
    adapter := adapters.NewPostgresAdapter()
    
    // Create migrator
    migrator := sil.NewMigrator(adapter, config)
    
    // Run migrations
    if err := migrator.Migrate(); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Migrations completed successfully")
}
```

## Integration with Toutā (Optional)

After Síl is stable, Toutā can optionally integrate it:

```bash
# Toutā CLI wrappers (optional future enhancement)
touta migrate              # Wrapper for sil migrate
touta migrate:rollback     # Wrapper for sil rollback
touta migrate:status       # Wrapper for sil status
touta seed                 # Wrapper for sil seed:run --all
```

## Success Criteria

- ✅ **Independent**: Works without any Toutā dependencies
- ✅ **Multi-Database**: PostgreSQL, MySQL, SQLite support
- ✅ **Production Ready**: Migration locking, transaction support
- ✅ **Well Tested**: 90%+ code coverage
- ✅ **Well Documented**: Comprehensive docs and examples
- ✅ **Developer Friendly**: Intuitive CLI and clear errors

## Timeline Summary

| Phase | Duration | Key Deliverable |
|-------|----------|-----------------|
| 1     | 2-3 weeks | Core + PostgreSQL |
| 2     | 2 weeks   | MySQL + SQLite |
| 3     | 2 weeks   | Seeding system |
| 4     | 1-2 weeks | Advanced features |
| **Total** | **7-9 weeks** | **Production-ready v0.1.0** |

## Next Steps

1. **Review this proposal** for approval
2. **Create GitHub repository** at https://github.com/toutaio/toutago-sil-migrator
3. **Initialize project structure** at /home/nestor/Proyects/toutago-sil-migrator
4. **Begin Phase 1 implementation** following tasks.md checklist
5. **Track progress** through OpenSpec change management

## Documentation

Full specification and design documents are available in:
- `/home/nestor/Proyects/toutago/openspec/changes/create-sil-migrator/proposal.md`
- `/home/nestor/Proyects/toutago/openspec/changes/create-sil-migrator/design.md`
- `/home/nestor/Proyects/toutago/openspec/changes/create-sil-migrator/tasks.md`
- `/home/nestor/Proyects/toutago/openspec/changes/create-sil-migrator/specs/sil-migrator/spec.md`

Validate with: `openspec validate create-sil-migrator --strict`
