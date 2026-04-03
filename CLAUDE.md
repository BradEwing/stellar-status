# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**stellar-status** is an opinionated Claude Code status line that displays moon phase, upcoming rocket launches, sun position, twilight times, planet visibility, meteor showers, deep sky objects, aurora/Kp index, and NASA APOD. It integrates with Claude Code's UI to provide at-a-glance astronomical and launch information.

## Development Commands

- **Build**: `go build -o stellar-status .`
- **Run**: `go run .` or `./stellar-status` after building
- **Run tests**: `go test ./...` or `go test -v ./...` for verbose output
- **Run a specific test**: `go test -run TestName ./package`
- **Lint**: `go vet ./...`
- **Format**: `gofmt -w .` or `go fmt ./...`

## Architecture

### Entry Point

- `main.go` — calls `cmd.Execute()`
- `cmd/root.go` — Cobra CLI with Viper flag binding, orchestrates all widgets

### CLI Flags

`--site` (launch site slug), `--lat`/`--lon` (observer coordinates), `--solar`, `--twilight`, `--planets`, `--meteors`, `--deepsky`, `--aurora`, `--apod`, `--nasa-key`, `--moon-ascii`, `--no-moon`, `--no-launch`, `--no-cache`

### Packages

1. **`internal/astro/`** — shared astronomical utilities: Julian Date, altitude/azimuth, local sidereal time, `Location` type. Based on Meeus, *Astronomical Algorithms*.

2. **`internal/moon/`** — lunar phase calculation and ASCII art representation. Pure Go.

3. **`internal/launches/`** — Launch Library 2 client (`https://ll.thespacedevs.com/2.2.0`). File-based cache at `~/.cache/stellar-status/launches-{site}.json` with 10-minute TTL. Supports 15 launch sites.

4. **`internal/solar/`** — sun position and formatted status output.

5. **`internal/twilight/`** — sunrise/sunset, civil/astronomical twilight, golden hour.

6. **`internal/planets/`** — planet orbital elements and visibility calculations.

7. **`internal/meteors/`** — meteor shower activity from embedded IMO data. Static dataset of 12 major showers with peak dates, active windows, and ZHR. No external API.

8. **`internal/deepsky/`** — deep sky object visibility. Curated catalog of 15 Messier/NGC showpiece objects. Uses `internal/astro` to compute alt/az; shows highest object above 15° when sky is dark. No external API.

9. **`internal/aurora/`** — aurora/geomagnetic activity via NOAA SWPC Kp index. File-based cache at `~/.cache/stellar-status/aurora.json` with 30-minute TTL. No auth required.

10. **`internal/apod/`** — NASA Astronomy Picture of the Day title. File-based cache at `~/.cache/stellar-status/apod.json` with 6-hour TTL. Uses `DEMO_KEY` by default or `NASA_API_KEY` env var.

### Widget Patterns

**Pure computation widgets** (moon, solar, twilight, planets, meteors, deepsky) follow:
```
Current(loc) → Type        // uses time.Now()
ForTime(t, loc) → Type     // parameterized for testing
(Type) FormatStatus() string
```

**Network-dependent widgets** (launches, aurora, apod) use a `Fetch(ctx, ...) → (*Type, error)` pattern with file-based caching.

### Output Format

Widgets produce emoji-based status indicators joined with ` | `. Times use `"3:04pm"` in local timezone. Countdowns use smart truncation (`"3d 4h 15m"`).

### Dependencies

- `github.com/spf13/cobra` + `github.com/spf13/viper` — CLI framework
- `github.com/stretchr/testify` — test assertions
- All astronomy algorithms are pure Go

## Code Style

- No comments inside function or struct bodies
- Package names: lowercase singular nouns
- Error wrapping with `fmt.Errorf` + `%w`

## Testing

- Use `testify/assert` and `testify/require`
- Use `httptest` for HTTP mocking
- Write individual test functions, not table-driven tests

## CI/CD

- **`tag.yml`**: manual dispatch workflow, computes next semver tag, creates and pushes tag (requires `RELEASE_TOKEN`)
- **`release.yml`**: triggered by `v*.*.*` tags, builds 5-platform matrix (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64), creates GitHub release, auto-updates README install version
