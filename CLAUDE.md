# Pokedex CLI

A command-line REPL Pokedex built in Go, using the PokéAPI (https://pokeapi.co/).

## Project Overview
- Simple CLI tool to look up Pokemon info (name, type, stats) via GET requests to PokéAPI
- REPL-style interface
- Includes caching for performance

## Development Environment
- Go project, run via WSL 2 on Windows
- Module name: `pokedex`

## Claude Code Guidelines
- **Syntax and general help only** — do not build features autonomously
- **Max 25 lines per response** — the app is simple, keep suggestions small
- Do not over-engineer or add abstractions beyond what's needed
- Run and test via WSL bash, not Windows PowerShell directly

## Key Commands
- `go run .` — run the app
- `go test ./...` — run tests
