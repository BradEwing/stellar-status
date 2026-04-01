# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**stellar-status** is an opinionated Claude Code status line that displays:
- Current moon phase (using [ascii-moon-phase-python](https://github.com/asweigart/ascii-moon-phase-python))
- Upcoming rocket launches from Vandenberg Space Force Base (VBG)

The status line integrates with Claude Code's UI to provide contextual information while working.

## Development Commands

- **Build**: `go build -o stellar-status .` (outputs binary to current dir)
- **Run**: `go run .` or `./stellar-status` after building
- **Run tests**: `go test ./...` or `go test -v ./...` for verbose output (use testify/assert)
- **Run a specific test**: `go test -run TestName ./package`
- **Lint code**: `golangci-lint run` (if configured) or `go vet ./...`
- **Format code**: `gofmt -w .` or `go fmt ./...`

## Architecture

### Core Components

1. **Moon Phase Package** (`internal/moon/`)
   - Pure Go calculation of lunar phase (no external dependencies)
   - Returns phase percentage and ASCII representation
   - May reference ascii-moon-phase-python for output format reference

2. **Launch Tracking Package** (`internal/launches/`)
   - Fetches upcoming VBG rocket launches from external API
   - Implements caching layer to minimize API calls
   - Parses launch data and determines next scheduled launch
   - May use `time.Time` and file-based or in-memory cache

3. **Formatter Package** (`internal/format/`)
   - Combines moon phase and launch info into status string
   - Handles emoji/ASCII representation and truncation for terminal
   - Configurable output formats

4. **Main Package** (root `main.go`)
   - Main entry point
   - Orchestrates packages

### Data Sources

- **Moon phases**: Pure Go calculation (no external dependencies)
- **Launch data**: External API (e.g., LaunchLibrary, RocketLauncher, or custom VBG feed)
  - Document the chosen API, authentication, and rate limits as they're determined

## Code Style

- Prefer no comments within a struct or function body; code should be self-explanatory

## Key Decisions

- Status line is designed for clarity and at-a-glance consumption in the Claude Code UI
- Moon phase is a visual indicator; launch info is the primary actionable data
- Cache strategies should balance freshness with rate limiting on external APIs
