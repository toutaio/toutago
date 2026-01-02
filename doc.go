/*
Package toutago provides a lightweight, convention-based web framework for Go.

Toutago ("tribe" in ancient Celtic) is designed to unite various Go web development
concerns under a cohesive, easy-to-use API. It provides routing, templating, middleware,
and dependency injection capabilities in a single framework.

# Features

  - Convention-based routing with automatic parameter binding
  - Template rendering with layout support
  - Middleware pipeline for cross-cutting concerns
  - Dependency injection for clean architecture
  - Database migration support
  - Zero configuration to get started

# Quick Start

Create a basic web application:

	package main

	import (
		"github.com/toutaio/toutago"
	)

	func main() {
		app := toutago.New()

		app.Get("/", func(ctx *toutago.Context) error {
			return ctx.String(200, "Hello, Toutā!")
		})

		app.Listen(":8080")
	}

# Architecture

Toutago is composed of several focused subpackages:

  - cosan: HTTP router with radix tree matching
  - fith: Template renderer with layout support
  - nasc: Dependency injection container
  - datamapper: Database abstraction layer
  - sil: Database migration tool

# Etymology

The name "Toutā" comes from Proto-Celtic *toutā, meaning "tribe" or "people".
This reflects the framework's goal of bringing together the Go web development
community with shared conventions and tools.

# Part of the Toutā Ecosystem

This package is part of the larger Toutā ecosystem of Go web development tools.
Each component can be used independently or together as a cohesive framework.

For more information, see https://github.com/toutaio/toutago
*/
package toutago
