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