# 🌌 stellar-status

A [Claude Code](https://claude.ai/code) status line that shows various astro information, such as:

- Current moon phase with emoji and illumination percentage [default]
- Next upcoming rocket launch from your chosen launch site [default]
- Visible planets
- Sun altitude
- Sunrise/sunset times
- Meteor shower activity
- Deep sky object visibility
- Aurora/geomagnetic activity (Kp index)
- NASA Astronomy Picture of the Day

```
🌔 Waxing Gibbous 95% | 🚀[VBG] Falcon 9 Block 5 in 2d 16h 26m
```

## Installation

### Pre-built binary

Download the latest binary for your platform from the [Releases](https://github.com/BradEwing/stellar-status/releases) page.

### Using `go install`

Requires Go 1.26.1+.

```bash
go install github.com/BradEwing/stellar-status@v1.0.3
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

To enable additional widgets, pass flags in the command:

```json
{
  "statusLine": {
    "type": "command",
    "command": "stellar-status --solar --twilight --planets --meteors --deepsky --aurora --apod"
  }
}
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--no-cache` | `-n` | `false` | Disable file-based cache for API responses |
| `--site` | `-s` | `VBG` | Launch site abbreviation (see supported sites below) |
| `--moon-ascii` | `-m` | `false` | Show 5x3 ASCII moon art (multi-line output) |
| `--no-moon` | | `false` | Disable moon phase display |
| `--no-launch` | | `false` | Disable launch tracking |
| `--solar` | `-o` | `false` | Show sun altitude |
| `--twilight` | `-t` | `false` | Show sunrise/sunset times |
| `--planets` | `-p` | `false` | Show visible planets |
| `--meteors` | `-e` | `false` | Show meteor shower activity |
| `--deepsky` | `-d` | `false` | Show best visible deep sky object |
| `--aurora` | `-a` | `false` | Show aurora/Kp geomagnetic index |
| `--apod` | | `false` | Show NASA Astronomy Picture of the Day title |
| `--nasa-key` | | | NASA API key (or set `NASA_API_KEY` env var) |
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
stellar-status --solar --twilight --planets --meteors --deepsky --aurora --apod
# => 🌕 Full Moon 99% | 🚀[VBG] Falcon 9 in 2d 19h | ☀️ 45° | 🌅 sunset 7:42pm | 🔭 Venus Jupiter | ☄️ Perseids peak in 3d (ZHR ~100) | 🌌 M13 (Hercules Cluster) alt 68° | 🌌 Aurora possible (Kp=5) | 🔭 APOD: "Pillars of Creation"

# Meteor showers only
stellar-status --no-moon --no-launch --meteors
# => ☄️ Geminids peak tonight (ZHR ~150)

# Night sky combo: deep sky objects + planets + aurora
stellar-status --no-moon --no-launch --deepsky --planets --aurora
# => 🔭 Venus Jupiter Saturn | 🌌 M42 (Orion Nebula) alt 62° | ☀️ Solar quiet (Kp=2)

# Aurora monitoring
stellar-status --no-moon --no-launch --aurora
# => 🌌 Aurora likely! (Kp=8)

# APOD with custom NASA API key
stellar-status --no-moon --no-launch --apod --nasa-key YOUR_KEY
# => 🔭 APOD: "The Horsehead Nebula in Infrared"

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
  - Cached locally at `~/.cache/stellar-status/launches-{site}.json` (10-minute TTL)
  - 10-second HTTP timeout with graceful fallback
- **Meteor showers**: Static dataset from IMO working list (no external API)
  - 12 major showers with peak dates, active windows, and ZHR
- **Deep sky objects**: Pure Go calculation using `internal/astro` (no external API)
  - Curated catalog of 15 Messier/NGC showpiece objects
  - Computes altitude/azimuth from observer location; shows highest object above 15° when sky is dark
- **Aurora/Kp index**: [NOAA Space Weather Prediction Center](https://www.swpc.noaa.gov/)
  - Cached locally at `~/.cache/stellar-status/aurora.json` (30-minute TTL)
  - No authentication required
- **APOD**: [NASA Astronomy Picture of the Day API](https://api.nasa.gov/)
  - Cached locally at `~/.cache/stellar-status/apod.json` (6-hour TTL)
  - Uses `DEMO_KEY` by default (30 req/hr); set `NASA_API_KEY` for higher limits
