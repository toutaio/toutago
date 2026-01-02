# Toutā Terminology Guide

## Druidic-Inspired Terminology

Toutā uses terminology inspired by Celtic/Druidic concepts to create a unique identity:

### Core Terms

| Concept | Term | Origin | Meaning |
|---------|------|--------|---------|
| Project Name | **Toutā** | Proto-Celtic | "People" or "tribe" - representing a community of components working together |
| Commands | **Ogam** | Old Irish | Ogham script - the ancient Celtic writing system used by druids |
| Packages | **Nemeton** | Celtic | Sacred groves where druids gathered - representing collections of related components |
| Recipes | **Ritual** | Druidic practice | Complete ceremonial procedures - representing full application solutions |
| Recipe Repository | **Grove** | Celtic/English | Collection of sacred spaces - repository of ready-made rituals |
| Message Bus | **Scéla** | Old Irish | News, tidings - flow of messages between components |
| Auth & Access Control | **Breitheamh** | Old Irish | Judge, arbiter - authority and permission decisions |

### Usage Examples

#### Ogam (Commands)
```bash
# CLI commands for project management
touta new my-app       # Create a new project
touta init             # Initialize in current directory
touta serve            # Start development server
touta version          # Show version
```

#### Nemeton (Packages)
```yaml
# nemeton.yaml - Package metadata
name: auth-nemeton
version: 1.0.0
description: Authentication components
dependencies:
  - crypto-nemeton
  - session-nemeton
```

Nemetons are independent, reusable collections of components that can be:
- Developed locally in your project
- Imported from external sources
- Shared across multiple projects

#### Scéla (Message Bus)
```go
// Publishing and subscribing to messages
import "github.com/toutaio/toutago-scela-bus"

bus := scela.NewBus()

// Subscribe to messages
bus.Subscribe("user.registered", &UserHandler{})

// Publish messages
bus.Publish(ctx, &UserRegisteredMessage{
    Email: "user@example.com",
})
```

Messages flow through the Scéla bus like news traveling between Celtic tribes, enabling loose coupling and event-driven architecture.

#### Breitheamh (Authentication & Authorization)
```go
// Authentication and access control
import "github.com/toutaio/toutago-breitheamh-auth"

auth := breitheamh.New()

// Authenticate user
user, err := auth.Authenticate(credentials)

// Check permissions
if auth.Can(user, "posts:create") {
    // User has permission
}

// Role-based access
auth.AssignRole(user, "editor")
```

Breitheamh acts as the judge, deciding who can access what, just like Celtic Brehon judges made legal decisions.

#### Ritual Grove (Recipe Repository)
```yaml
# Complete application recipes
# Located at: github.com/toutaio/toutago-ritual-grove

rituals/
├── blog/              # Blog application ritual
│   └── ritual.yaml
├── wiki/              # Wiki application ritual
│   └── ritual.yaml
├── crm/               # CRM application ritual
│   └── ritual.yaml
└── erp/               # ERP application ritual
    └── ritual.yaml
```

The Ritual Grove is the sacred repository where complete, proven application solutions are stored and shared with the community.
```yaml
# ritual.yaml - Complete application solution
name: blog-ritual
description: A complete blog application
nemetons:
  - user-nemeton
  - post-nemeton
  - comment-nemeton
  - theme-nemeton
```

Rituals compose multiple nemetons into complete, deployable applications like:
- Blog systems
- E-commerce platforms
- Wiki applications
- Custom business solutions

## File and Directory Naming

### In Text and UI
- Use **Toutā** (with diacritical mark) for display and documentation
- Use **Ogam**, **Nemeton**, **Ritual** in documentation

### In Filenames and Code
- Use `touta` (lowercase, no diacritics) for executables, directories, and imports
- Use `touta.yaml` for configuration files
- Use `nemeton.yaml` for package manifests
- Use `ritual.yaml` for recipe definitions

### Examples
```
my-project/
├── touta.yaml          # Main configuration
├── nemetons/           # Local nemetons directory
│   └── auth/
│       └── nemeton.yaml
├── rituals/            # Ritual definitions
│   └── blog/
│       └── ritual.yaml
└── main.go
```

## Why These Terms?

The Celtic/Druidic theme reflects the framework's philosophy:

1. **Toutā (People/Tribe)**: Software development is about community and collaboration
2. **Ogam (Sacred Script)**: Commands are the sacred instructions that shape your application
3. **Nemeton (Sacred Grove)**: Packages are gathering places for related functionality
4. **Ritual (Ceremonial Procedure)**: Recipes are complete, proven processes for achieving goals
5. **Grove (Sacred Gathering)**: Repository where rituals are stored and shared
6. **Scéla (News/Tidings)**: Messages carry information between components like news between tribes
7. **Breitheamh (Judge)**: Authority that decides access and permissions, like Celtic judges

This terminology creates a cohesive, memorable identity while maintaining clarity and purpose.

## Component Libraries (Ecosystem)

The Toutā framework is built on specialized, production-ready component libraries:

| Component | Name | Origin | Purpose |
|-----------|------|--------|---------|
| DI Container | **Nasc** | Old Irish: "Link, Bond" | Dependency injection - linking components together |
| HTTP Router | **Cosan** | Irish: "Path, Pathway" | Routing HTTP requests to handlers |
| Template Engine | **Fíth** | Old Irish: "Weaving" | Weaving data into templates |
| Data Mapping | **DataMapper** | English | Mapping domain objects to data sources |
| Migrations | **Sil** | Old Irish: "Seed, Lineage" | Database schema evolution and seeding |
| Message Bus | **Scéla** | Old Irish: "News, Messages" | Pub/sub message passing between components |
| Auth/Access | **Breitheamh** | Old Irish: "Judge, Arbiter" | Authentication and authorization |
| Recipe Repository | **Ritual Grove** | Celtic/English | Collection of complete application solutions |

### Repository Names

All component libraries follow the `toutago-<name>` naming pattern:

- `toutago-nasc-dependency-injector` - Dependency injection container
- `toutago-cosan-router` - HTTP router with middleware
- `toutago-fith-renderer` - Template engine
- `toutago-datamapper` - Database abstraction layer
- `toutago-datamapper-mysql` - MySQL adapter
- `toutago-datamapper-postgres` - PostgreSQL adapter
- `toutago-sil-migrator` - Database migration tool
- `toutago-scela-bus` - Message bus implementation
- `toutago-breitheamh-auth` - Authentication and authorization
- `toutago-ritual-grove` - Repository of ready-made application recipes

### Pronunciation Guide

For non-Irish speakers:

- **Toutā**: TOO-tah
- **Ogam**: OH-gum (like "ogram")
- **Nemeton**: NEM-eh-ton
- **Nasc**: NASK (rhymes with "task")
- **Cosan**: KUH-sawn
- **Fíth**: FEEH (with soft "th")
- **Sil**: SHILL
- **Scéla**: SHKAY-la
- **Breitheamh**: BREH-hyuv (or simplified: BRAY-huv)
