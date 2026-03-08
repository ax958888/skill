# Project Structure

## Directory Organization

```
.
├── .kiro/
│   ├── specs/
│   │   └── skill-automation-pipeline/    # Current spec for automation pipeline
│   │       ├── .config.kiro              # Spec configuration
│   │       └── requirements.md           # Requirements document
│   └── steering/                         # Project steering rules
│       ├── product.md                    # Product overview
│       ├── tech.md                       # Tech stack and build info
│       └── structure.md                  # This file
└── skill重塑.md                          # Binary Arsenal Generator skill definition
```

## Key Files

- **skill重塑.md**: Defines the Binary Arsenal Generator skill logic, language decision matrix, and execution standards
- **.kiro/specs/**: Contains specification documents for features and workflows
- **.kiro/steering/**: Contains project conventions and guidelines for AI assistance

## Workflow Components

The automation pipeline consists of two main agents:

1. **Kiro Bot** (@Pojun_kirobot): Analysis phase
   - Receives GitHub URLs
   - Performs security analysis
   - Detects language and skill type
   - Generates SOP

2. **Forge Agent** (@Forge_coderxbot): Build and deployment phase
   - Rebuilds skills according to SOP
   - Reviews build results
   - Deploys to k3s cluster
   - Backs up to GitHub

## Naming Conventions

- Specs: Use kebab-case for feature names (e.g., skill-automation-pipeline)
- Binaries: Single executable with descriptive name
- Documentation: Markdown format with clear section headers
