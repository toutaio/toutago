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

#### Ritual (Recipes)
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

This terminology creates a cohesive, memorable identity while maintaining clarity and purpose.
