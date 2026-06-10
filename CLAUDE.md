# Pokedex CLI

A command-line REPL Pokedex built in Go, using the PokéAPI (https://pokeapi.co/).

## Project Overview
- REPL CLI with a command registry pattern; line editing and history via `chzyer/readline`
- Caching layer via `internal/pokecache` to avoid repeat API calls
- Pokedex persists to `~/.pokedex` (JSON) between sessions
- Gen 4-style battle system: real moves, type chart, crits, speed order

## Project Structure
- `main.go` — readline REPL loop, command dispatch
- `repl.go` — `config`, `cliCommand`, `getCommands()`, all command functions, `ballThresholds`
- `pokemon.go` — `Pokemon`/`Move` structs, `fetchPokemon`, `fetchMove`
- `location.go` — location-area API structs and fetch functions
- `battle.go` — `runBattle`, `damageRoll`, `pickMoves`, `getStat`
- `typechart.go` — Gen 4 type chart, `typeEffect`
- `persist.go` — `savePokedex`, `loadPokedex`, `clearPokedex`, `pokedexPath`
- `internal/pokecache/pokecache.go` — thread-safe TTL cache (`NewCache` returns `*Cache`)
- `commands_test.go`, `repl_test.go` — tests (cache-seeded, no network in `-short`)

## Commands
- `help`, `exit` — basics; exit saves the Pokedex
- `map` / `mapb` — paginate location areas
- `explore <location>` — list Pokemon in a location area
- `wander <location>` — random encounter: optional battle, then optional catch
- `catch <pokemon> [--ball pokeball|greatball|ultraball]` — catch chance vs base experience
- `inspect <pokemon>` — stats, types, nickname for a caught Pokemon
- `rename <pokemon> <nickname>` / `release <pokemon>` — manage caught Pokemon (release works by nickname too)
- `pokedex` — list all caught Pokemon
- `clear` — wipe the Pokedex and delete the save file

## Key Patterns
- All commands: `func(cfg *config, args []string) error`
- `config` holds pagination URLs, `*pokecache.Cache`, and the `Pokedex` map
- Cache key is the full URL (pokemon fetches use the bare name); stores raw response bytes
- Tests seed the cache with JSON bytes instead of hitting the network

## Development Environment
- Go project, run via WSL 2 on Windows; module name `pokedex`
- CI: GitHub Actions runs `go vet` + `go test -race -short` on push/PR

## Claude Code Guidelines
- **Syntax and general help only** — do not build features autonomously
- **Max 25 lines per response** — the app is simple, keep suggestions small
- Do not over-engineer or add abstractions beyond what's needed
- Run and test via WSL bash, not Windows PowerShell directly
- **Use Haiku for sub-agents** — when spinning up sub-agents for simple or data-heavy tasks, set the model to Haiku to keep costs down without sacrificing quality on the main agent

## Key Commands
- `go run .` — run the app
- `go test -short ./...` — run tests (skip network-dependent tests)
