# skill-analyzer

Analyzes GitHub repositories for skill transformation into zero-dependency binaries.

## Build

```bash
chmod +x build.sh
./build.sh
```

## Usage

```bash
# Analyze a repository
./skill-analyzer analyze https://github.com/user/repo

# Save to file
./skill-analyzer analyze https://github.com/user/repo -o analysis.json

# Verbose mode
./skill-analyzer analyze https://github.com/user/repo -v
```

## Output

JSON format to stdout containing:
- Security analysis
- Language detection
- Type classification
- Build recommendations
- Standard Operating Procedure (SOP)

## Requirements

- Go 1.21+
- git command available
