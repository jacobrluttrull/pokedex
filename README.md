# Pokedex CLI

A command-line Pokedex built in Go. Uses the [PokéAPI](https://pokeapi.co/) to explore locations, catch Pokemon, and battle wild encounters.

## Requirements

- Go 1.21+

## Installation

```bash
git clone https://github.com/mrdetails/pokedexcli
cd pokedexcli
go build .
```

## Usage

```bash
./pokedexcli
```

Launches an interactive REPL. Type `help` to list available commands.

## Commands

| Command | Description |
|---|---|
| `help` | List all commands |
| `exit` | Save and exit |
| `map` | Show next 20 location areas |
| `mapb` | Show previous 20 location areas |
| `explore <location>` | List all Pokemon in a location |
| `wander <location>` | Encounter a random wild Pokemon — battle or catch it |
| `catch <pokemon> [--ball <type>]` | Attempt to catch a Pokemon |
| `inspect <pokemon>` | Show stats, types, height, and weight |
| `rename <pokemon> <nickname>` | Give a caught Pokemon a nickname |
| `release <pokemon>` | Release a caught Pokemon |
| `pokedex` | List all caught Pokemon |
| `clear` | Reset your Pokedex and delete save data |

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
go test ./...
```

To skip network calls:

```bash
go test -short ./...
```

## Project Structure

```
main.go           — REPL loop and entry point
repl.go           — Commands and config
pokemon.go        — Pokemon struct and API fetch
location_area.go  — Location area API structs and fetch
battle.go         — Battle logic
balls.go          — Pokeball catch thresholds
save.go           — Save and load Pokedex
clear.go          — Clear Pokedex and delete save file
internal/
  pokecache/      — Thread-safe cache with TTL
```
