# Proclens

**Filter and optionally kill processes using a lot of resources with ease.**

**Built with Go**

## Features
- Lists all running processes with CPU and memory usage.
- Filters processes by CPU percentage and RAM (MB) thresholds.
- Optionally kills matching processes with a single flag.


## Installation
Using Homebrew:
```bash
brew tap mj307/proclens
brew install proclens
```
## Run the CLI
```bash
proclens --cpu 1.5 --mem 20 
```

## Flags
| Flag        | Description                               | Default |
| ----------- | ----------------------------------------- | ------- |
| `--cpu`     | Minimum CPU percentage to filter          | `0.0`   |
| `--mem`     | Minimum RAM usage in MB to filter         | `0.0`   |
| `--dry-run` | Only show matching processes (no killing) | `true`  |
| `--kill`    | Actually kill matching processes          | `false` |

## Examples

### List all processes using more than 1.5% CPU and 20MB RAM:

```bash
proclens --cpu 1.5 --mem 20
```
### Kill all processes exceeding 5% CPU and 100MB RAM (use with caution!):

```bash
proclens --cpu 5 --mem 100 --dry-run=false --kill
```
