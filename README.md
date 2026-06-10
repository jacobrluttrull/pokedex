# Pokedex CLI

![CI](https://github.com/jacobrluttrull/pokedex/actions/workflows/ci.yml/badge.svg)

A command-line Pokedex built in Go. Explore locations, catch Pokemon, and fight Gen 4-style battles against wild encounters using the [PokéAPI](https://pokeapi.co/).

## Requirements

- Go 1.22+

## Installation

```bash
git clone https://github.com/jacobrluttrull/pokedex
cd pokedex
go build .
```

## Usage

```bash
./pokedex
```

Launches an interactive REPL with line editing and up-arrow command history. Type `help` to list available commands.

## Commands

| Command | Description |
|---|---|
| `help` | List all commands |
| `exit` | Save and exit |
| `map` | Show next 20 location areas |
| `mapb` | Show previous 20 location areas |
| `explore <location>` | List all Pokemon in a location |
| `wander <location>` | Encounter a random wild Pokemon — battle and/or catch it |
| `catch <pokemon> [--ball <type>]` | Attempt to catch a Pokemon |
| `inspect <pokemon>` | Show stats, types, height, and weight |
| `rename <pokemon> <nickname>` | Give a caught Pokemon a nickname |
| `release <pokemon>` | Release a caught Pokemon (by name or nickname) |
| `pokedex` | List all caught Pokemon |
| `clear` | Reset your Pokedex and delete save data |

## Battle System

`wander` encounters can turn into turn-based battles modeled on the Gen 4 games:

- Each Pokemon fights with up to 4 real damaging moves pulled from the PokéAPI
- Pick a move each turn (or `run` to flee); the wild Pokemon picks its own
- Faster Pokemon attacks first; speed ties are a coin flip
- Gen 4 damage formula at level 50 with STAB, the full type effectiveness chart, 1/16 critical hits, accuracy checks, and 85–100% damage variance
- Win the battle and you get a chance to catch the wild Pokemon

### Ball Types

| Ball | Catch Threshold |
|---|---|
| `pokeball` (default) | 50 |
| `greatball` | 75 |
| `ultraball` | 90 |

## Data Persistence

Your Pokedex is saved to `~/.pokedex` on exit and loaded automatically on startup.

## Running Tests

```bash
go test -short ./...   # skips network-dependent tests
go test -race ./...    # full suite with the race detector
```

CI runs `go vet` and `go test -race -short` on every push and pull request.

## Project Structure

```
main.go           — readline REPL loop and entry point
repl.go           — Commands, config, ball thresholds
persist.go        — Save, load, and clear Pokedex
internal/
  pokeapi/        — Pokemon/Move/location structs and API fetches
  game/           — Battle loop, damage formula, type chart
  pokecache/      — Thread-safe cache with TTL
```
