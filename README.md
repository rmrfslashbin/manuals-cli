# manuals-cli

Command-line interface for the Manuals documentation platform.

## Author

Robert Sigler (code@sigler.io)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Installation

### From Source

```bash
go install github.com/rmrfslashbin/manuals-cli/cmd/manuals@latest
```

### Build Locally

```bash
git clone https://github.com/rmrfslashbin/manuals-cli.git
cd manuals-cli
go build -o manuals ./cmd/manuals
```

## Configuration

Configure the CLI via environment variables or a config file.

### Environment Variables

```bash
export MANUALS_API_URL="http://manuals.local:8080"
export MANUALS_API_KEY="your-api-key"
```

### Config File

Create `~/.manuals.yaml`:

```yaml
api_url: http://manuals.local:8080
api_key: your-api-key
output_format: table  # table, json, or text
```

## Usage

### Search

```bash
# Search for devices
manuals search "raspberry pi gpio"
manuals search "uart protocol" --limit 5
manuals search esp32 -o json
```

### Devices

```bash
# List devices
manuals devices list
manuals devices list --domain hardware
manuals devices list --type dev-boards --limit 10

# Get device details
manuals devices get <device-id>
```

### Documents

```bash
# List documents
manuals documents list
manuals docs list --device <device-id>

# Get document details
manuals docs get <document-id>

# Download a document
manuals docs download <document-id>
manuals docs download <document-id> -o ~/Downloads/
```

### Output Formats

Use `-o` or `--output` to change the output format:

- `table` - Formatted table (default)
- `json` - JSON output for scripting
- `text` - Plain text

```bash
manuals devices list -o json | jq '.data[].name'
```

## Commands

| Command | Description |
|---------|-------------|
| `search <query>` | Search for devices and documentation |
| `devices list` | List all devices |
| `devices get <id>` | Get device details |
| `docs list` | List all documents |
| `docs get <id>` | Get document details |
| `docs download <id>` | Download a document |
| `version` | Show version information |

## Examples

```bash
# Find ESP32 documentation
manuals search "esp32 pinout"

# List all hardware devices
manuals devices list --domain hardware

# Download a datasheet
manuals docs download abc12345 -o ~/Documents/

# Get JSON output for scripting
manuals devices list -o json | jq '.data[] | select(.type == "sensors")'
```
