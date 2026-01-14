# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run

```bash
go build -o mlwcli cmd/mlwcli/main.go   # Build the binary
go run cmd/mlwcli/main.go               # Run without building
go install github.com/goofansu/mlwcli/cmd/mlwcli@latest  # Install globally
```

## Architecture

This is a Go CLI tool for managing links (Linkding), RSS feeds (Miniflux), and pages (Wallabag). It uses `go-flags` for command parsing with a nested subcommand structure.

### Package Structure

- **cmd/mlwcli/main.go** - Entry point and command definitions. All CLI commands are defined as structs with `Execute()` methods. Commands embed `BaseCommand` for shared app access and `JSONOutputOptions` for output filtering.
- **internal/app/** - Business logic layer. `App` struct holds config and provides methods like `AddBookmark()`, `AddFeed()`, `ListEntries()`, `ListBookmarks()`.
- **internal/auth/** - Authentication flow. `Login()` and `Logout()` use `huh` TUI library for interactive service selection and credential input. `PromptServiceTUI()` presents radio button menu, service-specific TUI functions collect credentials with validation.
- **internal/config/** - TOML config management at `~/.config/mlwcli/auth.toml`. Stores endpoint/API key per service.
- **internal/linkding/** - Linkding API client wrapper using `github.com/piero-vic/go-linkding`.
- **internal/miniflux/** - Miniflux API client wrapper using `miniflux.app/v2/client`.
- **internal/format/** - Output formatting with `--json` field filtering and `--jq` expression support via `gojq`.

### Command Flow

1. `main.go` loads config, creates `App`, wires it to command structs
2. `go-flags` parser routes to appropriate `Execute()` method
3. Command calls app layer method with options struct
4. App layer calls service client, formats output

### Authentication Commands

- **`mlwcli auth login`** - Interactive TUI login. Displays radio button menu for service selection (Linkding, Miniflux, Wallabag) with ✓ check marks next to already signed-in services, then shows form-based credential input with password masking and validation. Endpoint URLs are normalized (trailing slashes removed) before saving.
- **`mlwcli auth logout`** - Interactive TUI logout. Displays radio button menu showing only currently signed-in services, then removes selected credentials from config.

Uses `github.com/charmbracelet/huh` for beautiful terminal forms with arrow key navigation, password fields, and input validation. Config state awareness ensures users can see which services are configured. Input normalization (trimming whitespace, removing trailing slashes) ensures clean configuration.

### Output Format

List commands return `{total, items}` structure. Use `--json "field1,field2"` to filter fields or `--jq` for complex queries.
