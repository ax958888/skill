# Design Document: Skill Automation Pipeline

## Overview

The Skill Automation Pipeline is a two-phase automation system that transforms GitHub-hosted scripts and tools into zero-dependency, high-performance single binary executables. The system consists of two independent Go-based CLI tools that work in sequence:

1. **skill-analyzer**: Analyzes GitHub repositories for security, language, and functionality, then generates a Standard Operating Procedure (SOP)
2. **skill-builder**: Executes the SOP to rebuild the skill as a static binary, validates it, and deploys to the target environment

Both tools conform to the Binary Arsenal specification: static compilation, CLI argument input, JSON stdout output, stderr error reporting, and zero external dependencies.

### Design Goals

- **Security First**: Automated malicious code detection before any build process
- **Language Optimization**: Intelligent language selection based on skill characteristics (Go for network/API/automation, Rust for file processing/encryption/computation)
- **Zero Dependency**: All outputs are statically compiled single binaries
- **Workflow Traceability**: Unique IDs track work from analysis through deployment
- **Resilience**: Automatic retry logic, error recovery, and state persistence
- **Parallel Processing**: Support for concurrent skill processing with resource isolation

### Key Architectural Decisions

1. **Two Separate Tools**: Decoupling analysis from build allows independent execution, testing, and deployment
2. **Go Language Choice**: Both tools use Go for GitHub API interaction, network requests, process execution, and k3s integration
3. **JSON Communication**: Structured data exchange between tools via JSON files and stdout
4. **Stateless Execution**: Each tool can be invoked independently with all required context passed as input
5. **Cobra CLI Framework**: Consistent command-line interface across both tools

## Architecture

### System Context

```
┌─────────────┐
│   User      │
│  (GitHub    │
│   URL)      │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────────────────┐
│  skill-analyzer (Kiro Bot Phase)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   GitHub     │→ │  Security    │→ │  Language    │ │
│  │   Cloner     │  │  Analyzer    │  │  Detector    │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│         │                  │                  │         │
│         └──────────────────┴──────────────────┘         │
│                            ▼                            │
│                  ┌──────────────────┐                   │
│                  │  Type Analyzer   │                   │
│                  │  & SOP Generator │                   │
│                  └────────┬─────────┘                   │
│                           │                             │
│                           ▼                             │
│                  ┌──────────────────┐                   │
│                  │ Analysis Report  │                   │
│                  │     (JSON)       │                   │
│                  └────────┬─────────┘                   │
└───────────────────────────┼─────────────────────────────┘
                            │
                   [User Confirmation]
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│  skill-builder (Forge Agent Phase)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │  Build Env   │→ │   Builder    │→ │  Validator   │ │
│  │  Setup       │  │  (Go/Rust)   │  │              │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                            │                  │         │
│                            └──────────────────┘         │
│                                     ▼                   │
│                           ┌──────────────────┐          │
│                           │    Deployer      │          │
│                           │  (Root Dir +     │          │
│                           │   K3s/GitHub)    │          │
│                           └────────┬─────────┘          │
│                                    │                    │
│                                    ▼                    │
│                           ┌──────────────────┐          │
│                           │  Build Report    │          │
│                           │     (JSON)       │          │
│                           └──────────────────┘          │
└─────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                          ┌──────────────────┐
                          │  Deployed Binary │
                          │  /root/workspace │
                          │     /agents/     │
                          └──────────────────┘
```

### Component Architecture

#### skill-analyzer Components

1. **CLI Handler**: Cobra-based command parser and validator
2. **GitHub Client**: Repository cloning and metadata extraction
3. **Security Analyzer**: Pattern-based malicious code detection
4. **Language Detector**: File structure and configuration analysis
5. **Type Analyzer**: Functionality classification and use case identification
6. **SOP Generator**: Build instruction generation based on Binary Arsenal rules
7. **Report Generator**: JSON output formatter

#### skill-builder Components

1. **CLI Handler**: Cobra-based command parser and SOP loader
2. **Build Environment Manager**: Isolated workspace creation and cleanup
3. **Builder Engine**: Language-specific compilation orchestration
4. **Validator**: Binary verification (dependencies, interface, functionality)
5. **Deployer**: File system and k3s deployment
6. **GitHub Backup**: Git operations and version tagging
7. **Report Generator**: JSON output formatter

### Data Flow

```
GitHub URL → Clone → Security Scan → Language Detection → Type Analysis
     ↓
SOP Generation → Analysis Report (JSON) → User Confirmation
     ↓
Build Environment → Compile (Go/Rust) → Validate Binary
     ↓
Deploy to /root/workspace/agents/ → Backup to GitHub → Build Report (JSON)
```

## Components and Interfaces

### skill-analyzer CLI Interface

```bash
# Primary command
skill-analyzer analyze <github-url> [flags]

# Flags
--output, -o        Output file for analysis report (default: stdout)
--work-dir, -w      Working directory for cloning (default: /tmp/skill-analyzer-*)
--timeout, -t       Analysis timeout in seconds (default: 300)
--config, -c        Configuration file path
--verbose, -v       Enable verbose logging
--workflow-id       Workflow tracking ID (auto-generated if not provided)

# Examples
skill-analyzer analyze https://github.com/user/repo
skill-analyzer analyze https://github.com/user/repo -o report.json -v
```

**Output Format (stdout)**:
```json
{
  "workflow_id": "uuid-v4",
  "timestamp": "2024-01-15T10:30:00Z",
  "github_url": "https://github.com/user/repo",
  "repository": {
    "name": "repo",
    "owner": "user",
    "clone_path": "/tmp/skill-analyzer-xyz/repo",
    "commit_hash": "abc123..."
  },
  "security": {
    "score": 85,
    "status": "safe|warning|unsafe",
    "issues": [
      {
        "severity": "high|medium|low",
        "type": "arbitrary_execution|network_access|file_system",
        "file": "path/to/file.py",
        "line": 42,
        "description": "Detected eval() usage",
        "recommendation": "Replace with safe alternative"
      }
    ]
  },
  "language": {
    "primary": "python",
    "secondary": ["bash"],
    "confidence": 0.95,
    "build_files": ["requirements.txt", "setup.py"],
    "dependencies": ["requests", "click"]
  },
  "type_analysis": {
    "category": "network_request|cloud_api|k3s_docker|file_processing|encryption|parsing|automation",
    "use_cases": ["API client", "Data fetching"],
    "input_interface": "CLI arguments",
    "output_interface": "stdout text",
    "target_users": "developers"
  },
  "recommendation": {
    "target_language": "go|rust",
    "rationale": "Network-heavy operations suit Go's concurrency model",
    "decision_matrix_rule": "network_request → Go"
  },
  "sop": {
    "language": "go",
    "steps": [
      {
        "phase": "setup",
        "command": "go mod init skill-name",
        "description": "Initialize Go module"
      },
      {
        "phase": "build",
        "command": "CGO_ENABLED=0 go build -ldflags=\"-s -w\" -tags netgo -o skill-name",
        "description": "Static compilation"
      },
      {
        "phase": "validate",
        "command": "ldd skill-name",
        "expected": "not a dynamic executable"
      }
    ],
    "interface_spec": {
      "input": "CLI arguments via Cobra",
      "output": "JSON to stdout",
      "errors": "stderr with non-zero exit code"
    },
    "dependencies": ["github.com/spf13/cobra"],
    "estimated_build_time": "2-5 minutes"
  },
  "status": "success|failed",
  "message": "Analysis completed successfully"
}
```

**Error Output (stderr)**:
```json
{
  "error": "failed to clone repository",
  "details": "authentication required",
  "workflow_id": "uuid-v4",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### skill-builder CLI Interface

```bash
# Primary command
skill-builder build <analysis-report> [flags]

# Flags
--output, -o        Output file for build report (default: stdout)
--work-dir, -w      Working directory for building (default: /tmp/skill-builder-*)
--deploy-dir, -d    Deployment directory (default: /root/workspace/agents/)
--skip-deploy       Skip deployment step (build and validate only)
--skip-backup       Skip GitHub backup step
--config, -c        Configuration file path
--verbose, -v       Enable verbose logging

# Examples
skill-builder build report.json
skill-builder build report.json --skip-backup -v
echo '{"workflow_id":"...","sop":{...}}' | skill-builder build -
```

**Output Format (stdout)**:
```json
{
  "workflow_id": "uuid-v4",
  "timestamp": "2024-01-15T10:35:00Z",
  "build": {
    "status": "success|failed",
    "language": "go",
    "source_path": "/tmp/skill-builder-xyz/src",
    "binary_path": "/tmp/skill-builder-xyz/skill-name",
    "binary_size": 8388608,
    "build_time": "45s",
    "build_log": "..."
  },
  "validation": {
    "status": "passed|failed",
    "checks": [
      {
        "name": "single_binary",
        "passed": true,
        "message": "Output is a single executable file"
      },
      {
        "name": "zero_dependency",
        "passed": true,
        "message": "ldd: not a dynamic executable"
      },
      {
        "name": "cli_interface",
        "passed": true,
        "message": "Accepts --help flag"
      },
      {
        "name": "json_output",
        "passed": true,
        "message": "Outputs valid JSON"
      },
      {
        "name": "size_check",
        "passed": true,
        "message": "Binary size: 8MB (< 100MB limit)"
      }
    ]
  },
  "deployment": {
    "status": "success|failed|skipped",
    "target_path": "/root/workspace/agents/skill-name",
    "permissions": "755",
    "k3s_deployed": false,
    "deployment_time": "2024-01-15T10:35:30Z"
  },
  "backup": {
    "status": "success|failed|skipped",
    "repository": "https://github.com/backup-org/skill-name",
    "commit_hash": "def456...",
    "tag": "v1.0.0-20240115",
    "backup_time": "2024-01-15T10:36:00Z"
  },
  "status": "success|failed",
  "message": "Build, validation, deployment, and backup completed successfully"
}
```

### Internal Component Interfaces

#### GitHub Client Interface

```go
type GitHubClient interface {
    Clone(url string, destPath string) (*Repository, error)
    GetMetadata(url string) (*RepoMetadata, error)
    CreateBackupRepo(name string) (*Repository, error)
    PushBackup(localPath string, remoteURL string, tag string) error
}

type Repository struct {
    Name       string
    Owner      string
    ClonePath  string
    CommitHash string
    CloneTime  time.Time
}
```

#### Security Analyzer Interface

```go
type SecurityAnalyzer interface {
    Scan(repoPath string) (*SecurityReport, error)
    DetectPatterns(filePath string) ([]SecurityIssue, error)
    CalculateScore(issues []SecurityIssue) int
}

type SecurityReport struct {
    Score  int
    Status SecurityStatus // safe, warning, unsafe
    Issues []SecurityIssue
}

type SecurityIssue struct {
    Severity       string // high, medium, low
    Type           string // arbitrary_execution, network_access, file_system
    File           string
    Line           int
    Description    string
    Recommendation string
}
```

#### Language Detector Interface

```go
type LanguageDetector interface {
    Detect(repoPath string) (*LanguageInfo, error)
    AnalyzeBuildFiles(repoPath string) ([]string, error)
    ExtractDependencies(repoPath string, language string) ([]string, error)
}

type LanguageInfo struct {
    Primary     string
    Secondary   []string
    Confidence  float64
    BuildFiles  []string
    Dependencies []string
}
```

#### Type Analyzer Interface

```go
type TypeAnalyzer interface {
    Analyze(repoPath string, langInfo *LanguageInfo) (*TypeInfo, error)
    ClassifyCategory(repoPath string) (string, error)
    IdentifyUseCases(repoPath string) ([]string, error)
}

type TypeInfo struct {
    Category        string   // network_request, cloud_api, file_processing, etc.
    UseCases        []string
    InputInterface  string
    OutputInterface string
    TargetUsers     string
}
```

#### SOP Generator Interface

```go
type SOPGenerator interface {
    Generate(langInfo *LanguageInfo, typeInfo *TypeInfo) (*SOP, error)
    ApplyDecisionMatrix(typeInfo *TypeInfo) (string, string, error) // language, rationale
    GenerateBuildSteps(language string) ([]BuildStep, error)
}

type SOP struct {
    Language          string
    Steps             []BuildStep
    InterfaceSpec     InterfaceSpec
    Dependencies      []string
    EstimatedBuildTime string
}

type BuildStep struct {
    Phase       string // setup, build, validate
    Command     string
    Description string
    Expected    string // for validation steps
}
```

#### Builder Engine Interface

```go
type BuilderEngine interface {
    Build(sop *SOP, sourcePath string, outputPath string) (*BuildResult, error)
    ExecuteStep(step BuildStep, workDir string) error
    CompileGo(sourcePath string, outputPath string) error
    CompileRust(sourcePath string, outputPath string) error
}

type BuildResult struct {
    Status     string
    BinaryPath string
    BinarySize int64
    BuildTime  time.Duration
    BuildLog   string
}
```

#### Validator Interface

```go
type Validator interface {
    Validate(binaryPath string, sop *SOP) (*ValidationReport, error)
    CheckSingleBinary(binaryPath string) error
    CheckZeroDependency(binaryPath string) error
    CheckCLIInterface(binaryPath string) error
    CheckJSONOutput(binaryPath string) error
    CheckSize(binaryPath string, maxSize int64) error
}

type ValidationReport struct {
    Status string
    Checks []ValidationCheck
}

type ValidationCheck struct {
    Name    string
    Passed  bool
    Message string
}
```

#### Deployer Interface

```go
type Deployer interface {
    Deploy(binaryPath string, targetDir string, k3sConfig *K3sConfig) (*DeploymentResult, error)
    CopyToDirectory(src string, dest string) error
    SetPermissions(path string, mode os.FileMode) error
    DeployToK3s(binaryPath string, config *K3sConfig) error
}

type DeploymentResult struct {
    Status       string
    TargetPath   string
    Permissions  string
    K3sDeployed  bool
    DeploymentTime time.Time
}
```

## Data Models

### Core Domain Models

#### WorkflowContext

```go
type WorkflowContext struct {
    ID          string    `json:"workflow_id"`
    GitHubURL   string    `json:"github_url"`
    StartTime   time.Time `json:"start_time"`
    Status      string    `json:"status"` // pending, analyzing, building, deployed, failed
    CurrentPhase string   `json:"current_phase"`
    Error       *ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
    Message    string    `json:"message"`
    Details    string    `json:"details"`
    Timestamp  time.Time `json:"timestamp"`
    Retryable  bool      `json:"retryable"`
    RetryCount int       `json:"retry_count"`
}
```

#### AnalysisReport

```go
type AnalysisReport struct {
    WorkflowID   string              `json:"workflow_id"`
    Timestamp    time.Time           `json:"timestamp"`
    GitHubURL    string              `json:"github_url"`
    Repository   *Repository         `json:"repository"`
    Security     *SecurityReport     `json:"security"`
    Language     *LanguageInfo       `json:"language"`
    TypeAnalysis *TypeInfo           `json:"type_analysis"`
    Recommendation *Recommendation   `json:"recommendation"`
    SOP          *SOP                `json:"sop"`
    Status       string              `json:"status"`
    Message      string              `json:"message"`
}

type Recommendation struct {
    TargetLanguage      string `json:"target_language"`
    Rationale           string `json:"rationale"`
    DecisionMatrixRule  string `json:"decision_matrix_rule"`
}
```

#### BuildReport

```go
type BuildReport struct {
    WorkflowID   string              `json:"workflow_id"`
    Timestamp    time.Time           `json:"timestamp"`
    Build        *BuildResult        `json:"build"`
    Validation   *ValidationReport   `json:"validation"`
    Deployment   *DeploymentResult   `json:"deployment"`
    Backup       *BackupResult       `json:"backup"`
    Status       string              `json:"status"`
    Message      string              `json:"message"`
}

type BackupResult struct {
    Status     string    `json:"status"`
    Repository string    `json:"repository"`
    CommitHash string    `json:"commit_hash"`
    Tag        string    `json:"tag"`
    BackupTime time.Time `json:"backup_time"`
}
```

### Configuration Models

#### AnalyzerConfig

```go
type AnalyzerConfig struct {
    WorkDir          string        `json:"work_dir"`
    Timeout          time.Duration `json:"timeout"`
    SecurityPatterns []Pattern     `json:"security_patterns"`
    MaxRepoSize      int64         `json:"max_repo_size"`
    GitHubToken      string        `json:"github_token,omitempty"`
    LogLevel         string        `json:"log_level"`
}

type Pattern struct {
    Name        string   `json:"name"`
    Regex       string   `json:"regex"`
    Severity    string   `json:"severity"`
    Type        string   `json:"type"`
    Description string   `json:"description"`
}
```

#### BuilderConfig

```go
type BuilderConfig struct {
    WorkDir       string        `json:"work_dir"`
    DeployDir     string        `json:"deploy_dir"`
    Timeout       time.Duration `json:"timeout"`
    MaxBinarySize int64         `json:"max_binary_size"`
    K3sConfig     *K3sConfig    `json:"k3s_config,omitempty"`
    BackupRepo    string        `json:"backup_repo"`
    GitHubToken   string        `json:"github_token,omitempty"`
    LogLevel      string        `json:"log_level"`
}

type K3sConfig struct {
    Enabled     bool   `json:"enabled"`
    Namespace   string `json:"namespace"`
    KubeConfig  string `json:"kubeconfig"`
}
```

### State Persistence Models

#### WorkflowState

```go
type WorkflowState struct {
    Context       *WorkflowContext `json:"context"`
    AnalysisReport *AnalysisReport `json:"analysis_report,omitempty"`
    BuildReport   *BuildReport     `json:"build_report,omitempty"`
    UpdatedAt     time.Time        `json:"updated_at"`
}
```

State is persisted to JSON files in a state directory:
- Location: `/var/lib/skill-automation-pipeline/state/`
- Filename: `{workflow_id}.json`
- Retention: 30 days

### Language Decision Matrix Model

```go
type DecisionMatrix struct {
    Rules []DecisionRule `json:"rules"`
}

type DecisionRule struct {
    Categories []string `json:"categories"`
    Language   string   `json:"language"`
    Rationale  string   `json:"rationale"`
    Priority   int      `json:"priority"`
}

// Default matrix
var DefaultDecisionMatrix = DecisionMatrix{
    Rules: []DecisionRule{
        {
            Categories: []string{"network_request", "cloud_api", "k3s_docker", "automation"},
            Language:   "go",
            Rationale:  "Go excels at network operations, API clients, and system automation",
            Priority:   1,
        },
        {
            Categories: []string{"file_processing", "encryption", "parsing", "computation"},
            Language:   "rust",
            Rationale:  "Rust provides memory safety and performance for intensive operations",
            Priority:   1,
        },
    },
}
```


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property Reflection

After analyzing all acceptance criteria, I identified several areas of redundancy:

1. **Report Content Properties**: Multiple criteria check that reports contain specific sections (security, language, type, SOP). These can be combined into comprehensive report completeness properties.

2. **Validation Properties**: Several criteria check different aspects of binary validation (single file, no dependencies, correct interface). These overlap and can be consolidated.

3. **Logging Properties**: Multiple criteria about logging (format, levels, content) can be combined into fewer comprehensive properties.

4. **Build Command Properties**: Separate criteria for Go and Rust build commands are examples, not properties that need separate testing.

5. **State Tracking Properties**: Multiple criteria about recording timestamps, status, and errors can be consolidated into comprehensive state tracking properties.

The following properties represent the unique, non-redundant validation requirements:

### Property 1: Valid GitHub URL Acceptance

*For any* valid GitHub HTTPS URL in the format `https://github.com/{owner}/{repo}`, the skill-analyzer should accept it and return a workflow ID.

**Validates: Requirements 1.1, 1.4**

### Property 2: Invalid URL Rejection

*For any* string that is not a valid GitHub HTTPS URL format, the skill-analyzer should reject it with an error message to stderr and non-zero exit code.

**Validates: Requirements 1.2**

### Property 3: Repository Cloning

*For any* valid and accessible GitHub repository, after successful cloning, all repository files should exist in the designated work directory.

**Validates: Requirements 1.5**

### Property 4: Complete File Scanning

*For any* cloned repository, the security analyzer should scan all code files (not just a subset).

**Validates: Requirements 2.1**

### Property 5: Malicious Pattern Detection

*For any* code file containing known malicious patterns (arbitrary code execution, unauthorized network access, file system destruction), the security analyzer should detect and flag them.

**Validates: Requirements 2.2**

### Property 6: Dependency Risk Detection

*For any* project with suspicious or known-vulnerable dependencies, the security analyzer should flag them in the report.

**Validates: Requirements 2.3**

### Property 7: High-Risk Termination

*For any* repository with high-risk security issues detected, the workflow should terminate with status "unsafe" and not proceed to build.

**Validates: Requirements 2.4**

### Property 8: Security Score Range

*For any* repository analyzed, the security score should be an integer in the range [0, 100].

**Validates: Requirements 2.6**

### Property 9: Language Detection Accuracy

*For any* repository with a clear primary language (>70% of code in one language), the language detector should correctly identify it.

**Validates: Requirements 3.2**

### Property 10: Build File Detection

*For any* repository containing standard build configuration files (go.mod, Cargo.toml, package.json, requirements.txt), they should be identified and listed in the analysis report.

**Validates: Requirements 3.3**

### Property 11: Multi-Language Detection

*For any* repository with multiple languages, both primary and secondary languages should be identified in the report.

**Validates: Requirements 3.4**

### Property 12: Type Categorization

*For any* analyzed repository, it should be assigned to at least one valid category from the set: {network_request, cloud_api, k3s_docker, file_processing, encryption, parsing, automation, computation}.

**Validates: Requirements 4.2**

### Property 13: Analysis Report Completeness

*For any* completed analysis, the report should contain all required sections: security results, language info, type analysis, recommendation, and SOP.

**Validates: Requirements 6.2, 4.5, 3.5**

### Property 14: JSON Output Format

*For any* analysis report output to stdout, parsing it as JSON should succeed without errors.

**Validates: Requirements 6.3**

### Property 15: Decision Matrix Application for Network Skills

*For any* skill categorized as network_request, cloud_api, k3s_docker, or automation, the recommendation should suggest Go as the target language.

**Validates: Requirements 18.2**

### Property 16: Decision Matrix Application for Computation Skills

*For any* skill categorized as file_processing, encryption, parsing, or computation, the recommendation should suggest Rust as the target language.

**Validates: Requirements 18.3**

### Property 17: SOP Completeness

*For any* generated SOP, it should contain all required sections: build steps, interface specification, dependencies, and validation steps.

**Validates: Requirements 5.4, 5.5, 5.6, 5.7**

### Property 18: Unique Workflow IDs

*For any* two workflow instances created at different times, their workflow IDs should be unique (not equal).

**Validates: Requirements 8.3, 13.1**

### Property 19: Build Environment Isolation

*For any* two concurrent build workflows, their work directories should be different and non-overlapping.

**Validates: Requirements 9.1, 17.2**

### Property 20: Build Step Execution Order

*For any* SOP with multiple build steps, they should be executed in the order specified in the SOP steps array.

**Validates: Requirements 9.2**

### Property 21: Zero Dependency Validation

*For any* successfully built binary, running `ldd` on it should indicate no dynamic dependencies (output contains "not a dynamic executable" or "statically linked").

**Validates: Requirements 9.5, 10.3**

### Property 22: Build Report Generation

*For any* build attempt (success or failure), a build report should be generated with status, build results, and validation results.

**Validates: Requirements 9.7, 10.6**

### Property 23: Binary Interface Validation

*For any* built binary, it should accept a `--help` flag and output valid JSON when executed with valid arguments.

**Validates: Requirements 10.2, 19.4**

### Property 24: Binary Size Limit

*For any* successfully built binary, its file size should be less than 100MB.

**Validates: Requirements 19.2**

### Property 25: Deployment File Permissions

*For any* binary deployed to the target directory, its file permissions should be set to 755 (owner: rwx, group: rx, others: rx).

**Validates: Requirements 11.2**

### Property 26: Post-Deployment Executability

*For any* deployed binary, executing it with the `--help` flag should succeed (exit code 0).

**Validates: Requirements 11.5**

### Property 27: Deployment Info Recording

*For any* successful deployment, the build report should contain the deployment path and timestamp.

**Validates: Requirements 11.6**

### Property 28: Git Backup Commit

*For any* successful deployment with backup enabled, a git commit should be created containing the source code and binary.

**Validates: Requirements 12.2**

### Property 29: Version Tag Creation

*For any* git backup, a version tag should be created in the format matching either semantic versioning or timestamp-based versioning.

**Validates: Requirements 12.3**

### Property 30: Backup Report Completeness

*For any* completed backup, the backup report should contain commit hash, repository URL, and tag.

**Validates: Requirements 12.6**

### Property 31: Workflow State Persistence

*For any* workflow, its state should be persisted to disk and retrievable by workflow ID.

**Validates: Requirements 13.6**

### Property 32: Stage Timing Records

*For any* workflow stage (analysis, build, validation, deployment), start time, end time, and status should be recorded in the workflow state.

**Validates: Requirements 13.2**

### Property 33: Error Detail Recording

*For any* failed workflow stage, the error message, details, and timestamp should be recorded in the workflow state.

**Validates: Requirements 13.4, 14.1**

### Property 34: Error Output to Stderr

*For any* error condition in either tool, the error message should be output to stderr (not stdout).

**Validates: Requirements 14.2, 19.5**

### Property 35: Temporary File Cleanup

*For any* failed workflow, temporary files in the work directory should be cleaned up (directory removed or emptied).

**Validates: Requirements 14.6**

### Property 36: Structured Log Format

*For any* log entry generated by the system, it should be valid JSON with at minimum: timestamp, level, and message fields.

**Validates: Requirements 15.2, 15.3**

### Property 37: Configuration File Loading

*For any* valid configuration file in JSON format, the system should successfully load all parameters from it.

**Validates: Requirements 16.1**

### Property 38: Environment Variable Override

*For any* configuration parameter that has both a file value and an environment variable value, the environment variable value should take precedence.

**Validates: Requirements 16.2**

### Property 39: Invalid Configuration Fallback

*For any* invalid or missing configuration parameter, the system should use the default value and log a warning.

**Validates: Requirements 16.4, 16.5**

### Property 40: Concurrent Workflow Processing

*For any* set of N workflow requests where N is less than or equal to the concurrency limit, all workflows should be processed in parallel (not sequentially).

**Validates: Requirements 17.1**

### Property 41: Queue Overflow Handling

*For any* workflow request that exceeds the concurrency limit, it should be added to a queue and processed when a slot becomes available.

**Validates: Requirements 17.4**

### Property 42: Language Recommendation Rationale

*For any* analysis report where the recommended language differs from the detected language, the report should include a rationale explaining the recommendation.

**Validates: Requirements 18.4, 18.5**

### Property 43: Validation Failure Blocks Deployment

*For any* binary that fails any validation check, deployment should not proceed and an error report should be generated.

**Validates: Requirements 19.6**

### Property 44: Integration Test Idempotency

*For any* valid test input, running the integration test twice should produce equivalent results (same analysis outcomes, same build success/failure).

**Validates: Requirements 20.6**

## Error Handling

### Error Categories

The system defines four categories of errors:

1. **User Input Errors**: Invalid URLs, malformed configuration
   - Response: Immediate failure with descriptive error message
   - Retry: Not applicable
   - Exit code: 1

2. **External Service Errors**: GitHub API failures, network timeouts
   - Response: Automatic retry up to 3 times with exponential backoff
   - Retry: Yes (3 attempts)
   - Exit code: 2

3. **Security Errors**: Malicious code detected, high-risk patterns
   - Response: Immediate termination, detailed security report
   - Retry: Not applicable
   - Exit code: 3

4. **Build/Validation Errors**: Compilation failures, validation failures
   - Response: Detailed error report with logs
   - Retry: Not applicable (requires code changes)
   - Exit code: 4

### Error Response Format

All errors follow a consistent JSON structure on stderr:

```json
{
  "error": "short error message",
  "details": "detailed explanation",
  "category": "user_input|external_service|security|build_validation",
  "workflow_id": "uuid",
  "timestamp": "2024-01-15T10:30:00Z",
  "retryable": true|false,
  "retry_count": 0,
  "suggestions": ["suggestion 1", "suggestion 2"]
}
```

### Retry Logic

For retryable errors (external services):
- Initial attempt: immediate
- Retry 1: after 2 seconds
- Retry 2: after 4 seconds
- Retry 3: after 8 seconds
- After 3 failures: report error and exit

### Error Recovery

The system supports recovery from certain error states:

1. **Partial Analysis Failure**: If language detection fails but security passes, continue with manual language specification
2. **Build Failure**: Preserve analysis report and allow rebuild with modified SOP
3. **Deployment Failure**: Binary remains in build directory for manual deployment
4. **Backup Failure**: Deployment succeeds but backup is marked as failed for manual retry

### Cleanup on Error

Both tools implement cleanup handlers:
- Temporary directories are removed on exit (success or failure)
- Partial clones are deleted on clone failure
- Incomplete builds are removed on build failure
- File locks are released on any exit path

## Testing Strategy

### Dual Testing Approach

The system requires both unit testing and property-based testing for comprehensive coverage:

**Unit Tests**: Focus on specific examples, edge cases, and error conditions
- Example: Test that `https://github.com/user/repo` is accepted
- Example: Test that Go SOP contains `CGO_ENABLED=0`
- Example: Test that timeout occurs after exactly 300 seconds
- Example: Test that retry happens exactly 3 times
- Edge case: Empty repository
- Edge case: Repository with no code files
- Edge case: Binary exactly at 100MB size limit

**Property Tests**: Verify universal properties across all inputs
- Property: All valid GitHub URLs are accepted
- Property: All security scores are in range [0, 100]
- Property: All analysis reports are valid JSON
- Property: All built binaries have no dynamic dependencies

### Property-Based Testing Configuration

**Library Selection**:
- Go: Use `gopter` (https://github.com/leanovate/gopter) for property-based testing
- Alternative: `rapid` (https://github.com/flyingmutant/rapid)

**Test Configuration**:
- Minimum iterations per property test: 100
- Seed: Use fixed seed for reproducibility in CI
- Shrinking: Enable automatic shrinking to find minimal failing cases

**Test Tagging**:
Each property test must include a comment tag referencing the design property:

```go
// Feature: skill-automation-pipeline, Property 1: Valid GitHub URL Acceptance
func TestProperty_ValidGitHubURLAcceptance(t *testing.T) {
    properties := gopter.NewProperties(nil)
    properties.Property("valid GitHub URLs are accepted", prop.ForAll(
        func(owner, repo string) bool {
            url := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
            result := analyzer.Analyze(url)
            return result.WorkflowID != ""
        },
        gen.AlphaString(),
        gen.AlphaString(),
    ))
    properties.TestingRun(t, gopter.ConsoleReporter(os.Stdout))
}
```

### Test Organization

```
skill-analyzer/
├── cmd/
│   └── analyze.go
├── internal/
│   ├── github/
│   │   ├── client.go
│   │   ├── client_test.go          # Unit tests
│   │   └── client_property_test.go # Property tests
│   ├── security/
│   │   ├── analyzer.go
│   │   ├── analyzer_test.go
│   │   └── analyzer_property_test.go
│   └── ...
├── test/
│   ├── integration/
│   │   └── end_to_end_test.go
│   └── fixtures/
│       └── test_repos/
└── go.mod

skill-builder/
├── cmd/
│   └── build.go
├── internal/
│   ├── builder/
│   │   ├── engine.go
│   │   ├── engine_test.go
│   │   └── engine_property_test.go
│   └── ...
└── go.mod
```

### Integration Testing

Integration tests verify the complete workflow:

1. **Happy Path Test**: Valid repo → analysis → build → deploy → backup
2. **Security Failure Test**: Malicious repo → analysis → termination
3. **Build Failure Test**: Invalid code → analysis → build failure → error report
4. **Concurrent Processing Test**: Multiple repos → parallel analysis → all complete
5. **Recovery Test**: Failure → state persistence → recovery → completion

Integration tests run in a test mode that:
- Uses test GitHub repositories
- Deploys to test directories (not production)
- Skips actual GitHub backup (or uses test repo)
- Validates all intermediate outputs

### Test Coverage Goals

- Unit test coverage: >80% of lines
- Property test coverage: All 44 correctness properties
- Integration test coverage: All major workflows
- Error path coverage: All error categories and retry logic

### Continuous Integration

Tests run on every commit:
1. Unit tests (fast, <1 minute)
2. Property tests (medium, 2-5 minutes with 100 iterations)
3. Integration tests (slow, 5-10 minutes)

Property tests use a fixed seed in CI for reproducibility but can use random seeds locally for exploration.

