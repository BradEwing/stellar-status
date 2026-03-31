# stellar-status

An opinionated [Claude Code](https://claude.ai/code) status line that shows:

- Current moon phase with emoji and illumination percentage
- Next upcoming rocket launch from Vandenberg Space Force Base (VBG)

```
🌔 Waxing Gibbous 95% | 🚀[VBG] Falcon 9 Block 5 in 2d 16h 26m [Go]
```

## Installation

### Using `go install`

```bash
go install github.com/bradewing/stellar-status/cmd/stellar-status@latest
```

### From GitHub Releases

Download the latest binary for your platform from the [Releases](https://github.com/bradewing/stellar-status/releases) page and place it somewhere on your `$PATH`.

## Claude Code Statusline Setup

### Option 1: Standalone (stellar-status only)

Add to `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "stellar-status"
  }
}
```

### Option 2: Combined with context info (recommended)

Copy the bundled wrapper script and configure it:

```bash
cp scripts/statusline.sh ~/.claude/stellar-statusline.sh
chmod +x ~/.claude/stellar-statusline.sh
```

Add to `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "bash ~/.claude/stellar-statusline.sh"
  }
}
```

This shows stellar-status output alongside a context usage bar and git branch info.

## Data Sources

- **Moon phase**: Pure Go calculation (no external API)
- **Launch data**: [Launch Library 2 API](https://ll.thespacedevs.com) filtered to Vandenberg SFB
  - Cached locally for 10 minutes at `~/.cache/stellar-status/launches.json`
  - 15-second timeout with graceful fallback

## Launch Status Codes

| Code  | Meaning         |
|-------|-----------------|
| `Go`  | Go for Launch   |
| `TBC` | To Be Confirmed |
| `TBD` | To Be Determined |
