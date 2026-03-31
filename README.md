# 🌌 stellar-status

An opinionated [Claude Code](https://claude.ai/code) status line that shows:

- Current moon phase with emoji and illumination percentage
- Next upcoming rocket launch from Vandenberg Space Force Base (VBG)

```
🌔 Waxing Gibbous 95% | 🚀[VBG] Falcon 9 Block 5 in 2d 16h 26m
```

## Installation

### Using `go install`

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

## Data Sources

- **Moon phase**: Pure Go calculation (no external API)
- **Launch data**: [Launch Library 2 API](https://ll.thespacedevs.com) filtered to Vandenberg SFB
  - Cached locally for 10 minutes at `~/.cache/stellar-status/launches.json`
  - 15-second timeout with graceful fallback
