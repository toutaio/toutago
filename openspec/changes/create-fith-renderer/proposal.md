# Change: Create Fíth Template Renderer - Independent Project

## Why
Toutā needs a powerful, flexible template engine that aligns with its philosophy of clean architecture and Celtic-inspired terminology. Rather than relying on external template engines or building it into the core framework, creating Fíth as an independent, reusable Go library provides maximum flexibility and follows the nemeton (package) pattern. This allows the template engine to be used in any Go project, not just Toutā, while maintaining the framework's unique identity.

## What Changes
- Create a new independent Go project: `toutago-fith-renderer`
- Repository: https://github.com/toutaio/toutago-fith-renderer
- Location: `~/Proyects/toutago-fith-renderer`
- Implement a full-featured template engine with:
  - Template loading from directories using slug-based resolution
  - Support for maps and structs as data sources
  - Control flow: loops, conditionals (if/else)
  - Template composition: includes with parameters, parent/child layouts
  - Custom functions and filters
  - Multi-format output (HTML, CSV, text, markdown, etc.)
- Design as a standalone library that can be imported into Toutā or any Go project

## Impact
- Affected specs: fith-renderer (new capability - separate project)
- Affected code:
  - New repository: toutago-fith-renderer
  - Toutā can optionally integrate via Go modules
- Benefits:
  - Independent, reusable template engine
  - Clean separation from core framework
  - Can be used in any Go project
  - Follows Toutā's nemeton philosophy
  - Provides flexible, powerful templating
  - Supports multiple output formats
  - Extensible with custom functions/filters
  - Celtic-themed identity consistent with Toutā
