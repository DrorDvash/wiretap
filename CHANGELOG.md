# Changelog

## [Unreleased]

### Security Enhancements

#### Private Key Protection in `serve` Command

**Problem**: The `serve` command printed full configuration including private keys to stdout by default, potentially exposing sensitive data to other users or in logs.

**Solution**: 
- Configuration is now hidden by default
- Added `--show-config` / `-s` flag to explicitly display configuration
- When displayed, private keys are masked showing only first/last 10 characters

**Examples**:
```bash
# Default: No config output (production)
./wiretap serve -f config.conf

# Show masked config (debug)
./wiretap serve -f config.conf -s
# Output: private_key=e0f1eb06f2..REDACTED..cf911451
```

**Files Changed**:
- `src/cmd/serve.go` - Added showConfig flag and conditional output
- `src/peer/config.go` - Added maskKey() and AsIPCMasked()
- `src/peer/peer_config.go` - Added AsIPCMasked()

### New Features

#### Config Data Loading

**Problem**: Need flexible ways to load configuration in restricted environments without writing files to disk.

**Solution**: Added support for encoded configuration loading via multiple methods:

**Usage**:
```bash
# Method 1: Inline flag
./wiretap serve --config-data "W1JlbGF5..."

# Method 2: Environment variable
$env:WIRETAP_CONFIG_DATA="W1JlbGF5..."
./wiretap serve

# Method 3: .enc file
./wiretap serve -f wiretap_server.conf.enc
```

**Generate encoded configs**:
```bash
./wiretap configure --endpoint X:Y --routes Z
# Creates: wiretap_server.conf.enc + prints WIRETAP_CONFIG_DATA
```

**Files Changed**:
- `src/config/format.go` - NEW: Encoding/decoding logic
- `src/config/loader.go` - NEW: Config loading dispatch
- `src/cmd/serve.go` - Added --config-data flag
- `src/cmd/configure.go` - Generate .enc files and print encoded data

#### HTTP(S) Config Fetching

**Problem**: Need to fetch configuration from remote servers (C2 infrastructure) without manual file transfers.

**Solution**: Added HTTP(S) fetching with authentication and TLS options:

**Usage**:
```bash
# Basic HTTP(S) fetch
./wiretap serve --config-url https://c2.example.com/config

# With authentication header
./wiretap serve --config-url https://... --auth-header "X-API-Key: secret"
./wiretap serve --config-url https://... --auth-header "Authorization: Bearer TOKEN"

# Skip TLS verification (self-signed certs)
./wiretap serve --config-url https://... --insecure
```

**Features**:
- 3 automatic retries with 2s backoff
- 30s timeout per attempt
- Supports any HTTP header for auth
- Works with encoded (.enc) configs from remote servers

**Files Changed**:
- `src/config/remote.go` - NEW: HTTP fetcher with retry logic
- `src/config/loader.go` - Added URL loading support
- `src/cmd/serve.go` - Added --config-url, --auth-header, --insecure flags

#### UNC Path Support

**Note**: Windows UNC paths (SMB/WebDAV) work automatically with `-f` flag:
```bash
# SMB share
.\wiretap.exe serve -f \\server\share\config.conf

# WebDAV
.\wiretap.exe serve -f \\server@SSL\DavWWWRoot\config.conf
```
No special flags needed - Go's file handling supports UNC paths natively on Windows.
