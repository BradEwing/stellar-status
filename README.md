# 🌌 stellar-status

An opinionated [Claude Code](https://claude.ai/code) status line that shows:

- Current moon phase with emoji and illumination percentage
- Next upcoming rocket launch from your chosen launch site

```
🌔 Waxing Gibbous 95% | 🚀[VBG] Falcon 9 Block 5 in 2d 16h 26m
```

## Installation

### Using `go install`

Requires Go 1.26.1+.

```bash
go install github.com/BradEwing/stellar-status@latest
```

## Claude Code Statusline Setup

Add to `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "stellar-status"
  }
}
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--cache` | `-c` | `true` | Enable file-based cache for launch API responses |
| `--site` | `-s` | `VBG` | Launch site abbreviation (see supported sites below) |
| `--moon-ascii` | `-m` | `false` | Show 5x3 ASCII moon art (multi-line output) |

### Examples

```bash
# Use defaults (cache enabled, Vandenberg SFB)
stellar-status

# Track launches from Kennedy Space Center, no cache
stellar-status --site KSC --cache=false

# Short flags with ASCII moon art
stellar-status -s STARBASE -m
```

### Supported Sites

| Abbreviation | Name |
|-------------|------|
| `BAIKONUR` | Baikonur Cosmodrome |
| `CCSFS` | Cape Canaveral SFS |
| `CORNRANCH` | Blue Origin Corn Ranch |
| `CSG` | Guiana Space Centre |
| `JSLC` | Jiuquan Satellite Launch Center |
| `KSC` | Kennedy Space Center |
| `LC1` | Rocket Lab LC-1 |
| `PLESETSK` | Plesetsk Cosmodrome |
| `PSCA` | Pacific Spaceport Complex Alaska |
| `SDSC` | Satish Dhawan Space Centre |
| `STARBASE` | SpaceX Starbase |
| `TANEGASHIMA` | Tanegashima Space Center |
| `VBG` | Vandenberg SFB |
| `WFF` | Wallops Flight Facility |
| `XSLC` | Xichang Satellite Launch Center |

## Data Sources

- **Moon phase**: Pure Go calculation (no external API)
- **Launch data**: [Launch Library 2 API](https://ll.thespacedevs.com) filtered by launch site
  - Cached locally at `~/.cache/stellar-status/launches-{site}.json` (enabled by default, 10-minute TTL)
  - 10-second HTTP timeout with graceful fallback
