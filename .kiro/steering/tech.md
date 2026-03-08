# Technology Stack

## Primary Languages

The system uses a language decision matrix to select optimal implementation languages:

- **Go**: For network requests, cloud APIs, k3s/Docker interactions, simple automation
- **Rust**: For file processing, encryption/decryption, high-performance parsing, heavy computation

## Build Requirements

### Go Static Compilation
```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo
```

Key libraries:
- Cobra (CLI framework)

### Rust Static Compilation
```bash
cargo build --release --target x86_64-unknown-linux-musl
```

Build optimizations in Cargo.toml:
- `lto = true`
- `panic = "abort"`

Key libraries:
- Clap (CLI framework)

## Output Standards

All binaries must conform to:
- **Single Binary**: No external dependencies
- **Static Linking**: Must run on Ubuntu 24.04 without additional libraries
- **Standard Interface**:
  - Input: CLI arguments
  - Output: JSON to stdout on success
  - Errors: stderr with non-zero exit code

## Deployment Environment

- Target OS: Ubuntu 24.04
- Architecture: x86_64 (Linux)
- Container Platform: k3s
- Namespace: openclaw
- Installation Directory: /root/workspace/agents/

## Common Commands

### Verification
```bash
# Check for zero dependencies
ldd <binary>

# Verify binary size
ls -lh <binary>

# Test execution
./<binary> --help
```

### Deployment
```bash
# Set executable permissions
chmod 755 <binary>

# Deploy to root directory
cp <binary> /root/workspace/agents/
```
