# Pokedex CLI

A command-line REPL Pokedex built in Go, using the PokéAPI (https://pokeapi.co/).

## Project Overview
- Simple CLI tool to look up Pokemon info via GET requests to PokéAPI
- REPL-style interface with a command registry pattern
- Caching layer via `internal/pokecache` to avoid repeat API calls

## Project Structure
- `main.go` — REPL loop, scanner, command dispatch
- `repl.go` — `config` struct, `cliCommand` struct, `getCommands()`, command functions, `cleanInput()`
- `location_area.go` — API structs and fetch functions
- `pokemon.go` — `Pokemon` struct and `fetchPokemon` function
- `internal/pokecache/pokecache.go` — thread-safe cache with TTL reaping
- `repl_test.go` — unit tests for `cleanInput`

## Commands
- `help` — displays all commands and descriptions
- `exit` — exits the program
- `map` — displays next 20 location areas (paginates forward)
- `mapb` — displays previous 20 location areas (paginates back)
- `explore <location>` — lists all Pokemon in a given location area
- `catch <pokemon>` — attempts to catch a Pokemon; uses base experience to determine catch chance
- `inspect <pokemon>` — displays info (name, height, weight, stats, types) for a caught Pokemon
- `pokedex` — lists all caught Pokemon

## Key Patterns
- All commands have signature `func(cfg *config, args []string) error`
- `config` holds `Next`, `Previous` pagination URLs and the `Cache`
- Cache key is the full URL; stores raw response bytes

## Development Environment
- Go project, run via WSL 2 on Windows
- Module name: `pokedex`

## Claude Code Guidelines
- **Syntax and general help only** — do not build features autonomously
- **Max 25 lines per response** — the app is simple, keep suggestions small
- Do not over-engineer or add abstractions beyond what's needed
- Run and test via WSL bash, not Windows PowerShell directly
- **Use Haiku for sub-agents** — when spinning up sub-agents for simple or data-heavy tasks, set the model to Haiku to keep costs down without sacrificing quality on the main agent

## Key Commands
- `go run .` — run the app
- `go test ./...` — run tests
