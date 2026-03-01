# day

`day` is a small CLI to log work activities as "pings" and generate daily reports.

## Features

- Log activities with timestamps (`day ping`)
- Track scope/context per activity (ticket, meeting, etc.)
- Show a quick overview of today (`day status`)
- Print a grouped day report (`day report`)
- Store data locally in SQLite

## Installation

```bash
go install ./cmd/day
```

Optional: install `fzf` for interactive scope selection. Without it, `day` falls back to a text prompt.

## Commands

```bash
# record a ping now
day ping "Implemented OAuth callback"

# record retroactively
day ping "Standup" --at 09:15
day ping "Code review" --ago 30

# set scope/source explicitly (good for scripts)
day ping --silent "gc: fix auth edge case" --scope "PROJ-42" --source "git"

# overview for today
day status

# day report
day report
day report --back 1
day report --date 21.02.2026
```

## Configuration

Config file: `~/.config/day/config.toml` (optional)

```toml
[scopes]
predefined = ["daily", "meeting", "NOTICKET"]

[day]
data_dir = ".local/share/day"
```

If no config is present, sensible defaults are used.

## Data

- DB file: `~/.local/share/day/day.db` (unless `day.data_dir` is overridden)
- Table: `pings(ts, activity, scope, source)`
