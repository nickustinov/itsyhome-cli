# itsyhome

[![tests](https://github.com/nickustinov/itsyhome-cli/actions/workflows/test.yml/badge.svg)](https://github.com/nickustinov/itsyhome-cli/actions/workflows/test.yml)

A command-line tool to control your HomeKit devices through the [Itsyhome](https://github.com/nickustinov/itsyhome-macos) macOS app.

## Install

```bash
brew install nickustinov/tap/itsyhome
```

Or build from source:

```bash
go install github.com/nickustinov/itsyhome-cli@latest
```

## Prerequisites

- [Itsyhome](https://itsyhome.app) macOS app with Pro subscription
- Webhooks/CLI server enabled in Settings → Webhooks/CLI

## Usage

### Control commands

```bash
itsyhome toggle Office/Spotlights
itsyhome on Kitchen/Light
itsyhome off Bedroom/Lamp
itsyhome brightness 50 Office/Lamp
itsyhome position 75 "Living Room/Blinds"
itsyhome temp 22 Hallway/Thermostat
itsyhome color FF6600 Bedroom/Light
itsyhome scene Goodnight
itsyhome lock "Front Door"
itsyhome unlock "Front Door"
itsyhome open Garage/Door
itsyhome close Bedroom/Blinds
itsyhome toggle "group.All Lights"            # Control a global group
itsyhome toggle "Office/group.All Lights"    # Control a room-scoped group
itsyhome off "group.Office Lights"
```

### Query commands

```bash
itsyhome status                  # Home summary
itsyhome status Office           # Device states for a room
itsyhome status "Living Room"    # Use quotes for spaces
itsyhome list rooms              # List all rooms
itsyhome list devices            # List all devices
itsyhome list devices Office     # List devices in a room
itsyhome list scenes             # List all scenes
itsyhome list groups             # List all groups
itsyhome info Office/Lamp        # Device info with state
itsyhome info "Living Room"      # All devices in a room
itsyhome info "group.All Lights"         # Global group info
itsyhome info "Office/group.All Lights" # Room-scoped group info
```

### Example output

```
$ itsyhome status Office
Device         | State | Value
---------------|-------|------
Office AC      | on    | 22.5°
Spotlights     | on    | 80%
Blinds         | off   | —

$ itsyhome list devices
Device     | Type       | Room       | Status
-----------|------------|------------|------
Lamp       | light      | Office     | ok
AC Unit    | thermostat | Bedroom    | ok
Blinds     | blind      | Living Room| unreachable

$ itsyhome info Office/Lamp
Property   | Value
-----------|------
Name       | Lamp
Type       | light
Room       | Office
Status     | reachable
brightness | 80
on         | true
```

### JSON output

Add `--json` to any command for machine-readable output:

```bash
itsyhome status --json
itsyhome list devices --json
itsyhome info Office/Lamp --json
```

### Configuration

```bash
itsyhome config                        # Show current config
itsyhome config set --host 192.168.1.5 # Connect to remote Mac
itsyhome config set --port 9000        # Use custom port
```

Config file: `~/.config/itsyhome/config.json`

Default: `localhost:8423`

### Shell completions

```bash
itsyhome completion bash > /etc/bash_completion.d/itsyhome
itsyhome completion zsh > "${fpath[1]}/_itsyhome"
itsyhome completion fish > ~/.config/fish/completions/itsyhome.fish
```

## Target formats

| Format | Example | Description |
|--------|---------|-------------|
| `Room/Device` | `Office/Spotlights` | Device in a specific room |
| `Device` | `Lamp` | Device by name (if unique) |
| `Room` | `Office` | All devices in room (for info) |
| `Room/group.Name` | `Office/group.All Lights` | Group scoped to a room |
| `group.Name` | `group.Office Lights` | Global group |
| `scene.Name` | `scene.Goodnight` | Scene by name |

## API reference

The CLI communicates with the Itsyhome webhook server over HTTP.

### Control endpoints

```
GET /<action>/<target>
GET /<action>/<value>/<target>
```

| Action | Format | Example |
|--------|--------|---------|
| `toggle` | `/toggle/<target>` | `/toggle/Office/Lamp` |
| `on` | `/on/<target>` | `/on/Kitchen/Light` |
| `off` | `/off/<target>` | `/off/Bedroom/Lamp` |
| `brightness` | `/brightness/<0-100>/<target>` | `/brightness/50/Office/Lamp` |
| `position` | `/position/<0-100>/<target>` | `/position/75/Living%20Room/Blinds` |
| `temp` | `/temp/<mireds>/<target>` | `/temp/300/Office/Lamp` |
| `color` | `/color/<hex>/<target>` | `/color/FF6600/Bedroom/Light` |
| `scene` | `/scene/<name>` | `/scene/Goodnight` |
| `lock` | `/lock/<target>` | `/lock/Front%20Door` |
| `unlock` | `/unlock/<target>` | `/unlock/Front%20Door` |
| `open` | `/open/<target>` | `/open/Garage/Door` |
| `close` | `/close/<target>` | `/close/Bedroom/Blinds` |

### Query endpoints

| Endpoint | Response |
|----------|----------|
| `/status` | `{"rooms":3,"devices":10,"accessories":5,"reachable":8,"unreachable":2,"scenes":4,"groups":2}` |
| `/list/rooms` | `[{"name":"Office"},{"name":"Bedroom"}]` |
| `/list/devices` | `[{"name":"Lamp","type":"light","room":"Office","reachable":true}]` |
| `/list/devices/<room>` | Devices filtered by room |
| `/list/scenes` | `[{"name":"Goodnight"}]` |
| `/list/groups` | `[{"name":"All Lights","icon":"lightbulb","devices":5,"room":"Office"}]` |
| `/list/groups/<room>` | Groups available in room (room-scoped + global) |
| `/info/<target>` | `{"name":"Lamp","type":"light","room":"Office","reachable":true,"state":{"on":true,"brightness":80}}` |

### Response format

Success:
```json
{"status": "success"}
```

Error:
```json
{"status": "error", "message": "device not found"}
```

HTTP 403 when Pro is not active.

## License

MIT License © 2026 Nick Ustinov
