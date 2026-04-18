# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a ShortLink (URL shortening) service built using the SDD (Spec-Driven Development) methodology. The project is written in Go and implements a URL shortening service.

## Project Structure

The project is currently minimal with:
- Go source files (location TBD)
- Standard Go project structure with .gitignore configured for Go development

## Commands

### Building
- Use `go build` to build the project
- Use `go build ./cmd/shortlink` if there's a cmd/shortlink directory for the main executable

### Running
- After building: `./shortlink` to run the service
- Or `go run ./cmd/shortlink` to run directly from source

### Testing
- `go test ./...` to run all tests
- `go test -v ./...` for verbose output
- `go test ./path/to/package` to test specific packages

### Development
- `go mod tidy` to manage dependencies
- `go fmt ./...` to format code
- `go vet ./...` to examine code for common errors
- `golint ./...` for linting (if golint is installed)

## Architecture

Since this is a ShortLink service built with SDD, the architecture likely includes:

- An API layer for accepting URL shortening requests
- A service layer for business logic
- A storage layer for persisting short URL mappings
- A handler for redirecting short URLs to original URLs

The SDD (Spec-Driven Development) approach suggests the project likely has specification files that drive the implementation, focusing on defining behavior before implementing it.

## Development Guidelines

- Follow Go idiomatic code patterns and naming conventions
- Write tests for new functionality
- Ensure proper error handling
- Use structured logging where appropriate