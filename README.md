# 🌌 stellar-status

A [Claude Code](https://claude.ai/code) status line that shows various astro information, such as:

- Current moon phase with emoji and illumination percentage [default]
- Next upcoming rocket launch from your chosen launch site [default]
- Visible planets
- Sun altitude
- Sunrise/sunset times

```
🌔 Waxing Gibbous 95% | 🚀[VBG] Falcon 9 Block 5 in 2d 16h 26m
```

## Installation

### Using `go install`

Requires Go 1.26.1+.

```bash
go install github.com/BradEwing/stellar-status@v1.0.2
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
| `--no-cache` | `-n` | `false` | Disable file-based cache for launch API responses |
| `--site` | `-s` | `VBG` | Launch site abbreviation (see supported sites below) |
| `--moon-ascii` | `-m` | `false` | Show 5x3 ASCII moon art (multi-line output) |
| `--no-moon` | | `false` | Disable moon phase display |
| `--no-launch` | | `false` | Disable launch tracking |
| `--solar` | `-o` | `false` | Show sun altitude |
| `--twilight` | `-t` | `false` | Show sunrise/sunset times |
| `--planets` | `-p` | `false` | Show visible planets |
| `--lat` | | `34.7420` | Observer latitude (degrees, positive north) |
| `--lon` | | `-120.5724` | Observer longitude (degrees, positive east) |

### Examples

```bash
# Use defaults (cache enabled, Vandenberg SFB)
stellar-status
# => 🌕 Full Moon 99% | 🚀[VBG] Falcon 9 Block 5 in 2d 19h 59m

# Track launches from Kennedy Space Center, no cache
stellar-status --site KSC --no-cache
# => 🌕 Full Moon 99% | 🚀[KSC] Falcon Heavy in 26d 20h 56m

# Launches only (no moon)
stellar-status --no-moon
# => 🚀[VBG] Falcon 9 Block 5 in 2d 19h 59m

# Moon only (no launches)
stellar-status --no-launch
# => 🌕 Full Moon 99%

# Show all widgets
stellar-status --solar --twilight --planets
# => 🌕 Full Moon 99% | 🚀[VBG] Falcon 9 Block 5 in 2d 19h 59m | 🌙 -9° | 🌃 twilight til 8:47pm | 🔭 Venus Jupiter

# Only astro widgets, no defaults
stellar-status --no-moon --no-launch --solar --twilight --planets
# => 🌙 -9° | 🌃 twilight til 8:47pm | 🔭 Venus Jupiter

# ASCII moon art
stellar-status -m --no-launch
# =>  ***
#    ***** Full Moon 99%
#     ***
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
