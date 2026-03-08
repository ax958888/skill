# Implementation Plan: Skill Automation Pipeline

## Overview

This plan implements two independent Go CLI tools (skill-analyzer and skill-builder) that work in sequence to transform GitHub repositories into zero-dependency static binaries. The implementation follows the Binary Arsenal specification with static compilation, CLI interfaces, JSON output, and comprehensive testing.

## Tasks

- [ ] 1. Set up project structure and shared components
  - Create directory structure for both tools (skill-analyzer/ and skill-builder/)
  - Initialize Go modules for each tool with Cobra dependency
  - Create shared types package for common data models (WorkflowContext, AnalysisReport, BuildReport)
  - Set up logging infrastructure with structured JSON logging
  - Create configuration loading utilities supporting file and environment variable override
  - _Requirements: 15.2, 15.3, 16.1, 16.2_

- [ ]* 1.1 Write property test for configuration management
  - **Property 38: Environment Variable Override**
  - **Validates: Requirements 16.2**

- [ ]* 1.2 Write property test for configuration fallback
  - **Property 39: Invalid Configuration Fallback**
  - **Validates: Requirements 16.4, 16.5**

- [ ] 2. Implement skill-analyzer: GitHub client and repository cloning
  - [ ] 2.1 Create GitHub client with repository cloning functionality
    - Implement Clone() method using git command execution
    - Implement GetMetadata() for repository information extraction
    - Add URL validation for GitHub HTTPS format
    - Handle authentication with optional GitHub token
    - _Requirements: 1.1, 1.4, 1.5_

  - [ ]* 2.2 Write property test for GitHub URL validation
    - **Property 1: Valid GitHub URL Acceptance**
    - **Property 2: Invalid URL Rejection**
    - **Validates: Requirements 1.1, 1.2, 1.4_

  - [ ]* 2.3 Write property test for repository cloning
    - **Property 3: Repository Cloning**
    - **Validates: Requirements 1.5**

  - [ ]* 2.4 Write unit tests for GitHub client
    - Test error handling for inaccessible repositories
    - Test clone timeout behavior
    - Test metadata extraction
    - _Requirements: 1.3_

- [ ] 3. Implement skill-analyzer: Security analyzer
  - [ ] 3.1 Create security pattern detection engine
    - Implement file scanning with pattern matching (regex-based)
    - Define malicious patterns (eval, exec, arbitrary code execution, unauthorized network, file system destruction)
    - Implement dependency analysis for suspicious packages
    - Calculate security score based on issue severity and count
    - Generate SecurityReport with issues categorized by severity
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

  - [ ]* 3.2 Write property test for complete file scanning
    - **Property 4: Complete File Scanning**
    - **Validates: Requirements 2.1**

  - [ ]* 3.3 Write property test for malicious pattern detection
    - **Property 5: Malicious Pattern Detection**
    - **Validates: Requirements 2.2**

  - [ ]* 3.4 Write property test for security score range
    - **Property 8: Security Score Range**
    - **Validates: Requirements 2.6**

  - [ ]* 3.5 Write property test for high-risk termination
    - **Property 7: High-Risk Termination**
    - **Validates: Requirements 2.4**

  - [ ]* 3.6 Write unit tests for security analyzer
    - Test specific malicious patterns (eval, exec, os.system)
    - Test dependency vulnerability detection
    - Test score calculation edge cases
    - _Requirements: 2.2, 2.3, 2.6_

- [ ] 4. Implement skill-analyzer: Language detector
  - [ ] 4.1 Create language detection engine
    - Implement file extension analysis for language identification
    - Detect build configuration files (go.mod, Cargo.toml, package.json, requirements.txt, etc.)
    - Calculate language confidence based on file count and size
    - Extract dependencies from build files
    - Support multi-language detection with primary/secondary classification
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

  - [ ]* 4.2 Write property test for language detection accuracy
    - **Property 9: Language Detection Accuracy**
    - **Validates: Requirements 3.2**

  - [ ]* 4.3 Write property test for build file detection
    - **Property 10: Build File Detection**
    - **Validates: Requirements 3.3**

  - [ ]* 4.4 Write property test for multi-language detection
    - **Property 11: Multi-Language Detection**
    - **Validates: Requirements 3.4**

  - [ ]* 4.5 Write unit tests for language detector
    - Test edge cases (empty repo, no code files, mixed languages)
    - Test confidence calculation
    - Test dependency extraction for each language
    - _Requirements: 3.2, 3.3, 3.4_

- [ ] 5. Implement skill-analyzer: Type analyzer and decision matrix
  - [ ] 5.1 Create type analyzer and decision matrix engine
    - Implement category classification (network_request, cloud_api, k3s_docker, file_processing, encryption, parsing, automation, computation)
    - Analyze code patterns to identify use cases and target users
    - Implement decision matrix rules (Go for network/API/automation, Rust for file/encryption/computation)
    - Generate language recommendation with rationale
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 18.1, 18.2, 18.3, 18.4, 18.5_

  - [ ]* 5.2 Write property test for type categorization
    - **Property 12: Type Categorization**
    - **Validates: Requirements 4.2**

  - [ ]* 5.3 Write property test for decision matrix (Go)
    - **Property 15: Decision Matrix Application for Network Skills**
    - **Validates: Requirements 18.2**

  - [ ]* 5.4 Write property test for decision matrix (Rust)
    - **Property 16: Decision Matrix Application for Computation Skills**
    - **Validates: Requirements 18.3**

  - [ ]* 5.5 Write property test for language recommendation rationale
    - **Property 42: Language Recommendation Rationale**
    - **Validates: Requirements 18.4, 18.5**

  - [ ]* 5.6 Write unit tests for type analyzer
    - Test category classification for various code patterns
    - Test use case identification
    - Test decision matrix edge cases
    - _Requirements: 4.2, 4.3, 18.2, 18.3_

- [ ] 6. Implement skill-analyzer: SOP generator
  - [ ] 6.1 Create SOP generation engine
    - Generate build steps based on target language (Go: CGO_ENABLED=0 static compilation, Rust: musl target)
    - Include setup, build, and validation phases
    - Define interface specifications (CLI input, JSON output, stderr errors)
    - List required dependencies
    - Estimate build time based on project complexity
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7_

  - [ ]* 6.2 Write property test for SOP completeness
    - **Property 17: SOP Completeness**
    - **Validates: Requirements 5.4, 5.5, 5.6, 5.7**

  - [ ]* 6.3 Write unit tests for SOP generator
    - Test Go SOP generation (verify CGO_ENABLED=0, ldflags, tags)
    - Test Rust SOP generation (verify musl target, release profile)
    - Test interface specification generation
    - _Requirements: 5.2, 5.3, 5.4_

- [ ] 7. Implement skill-analyzer: CLI and report generation
  - [ ] 7.1 Create Cobra CLI interface and analysis orchestration
    - Implement analyze command with flags (output, work-dir, timeout, config, verbose, workflow-id)
    - Orchestrate analysis pipeline (clone → security → language → type → SOP)
    - Generate unique workflow IDs (UUID v4)
    - Generate complete AnalysisReport with all sections
    - Output JSON to stdout, errors to stderr
    - Implement workflow state persistence
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 8.1, 8.3, 8.4, 13.1, 13.2, 13.6_

  - [ ]* 7.2 Write property test for unique workflow IDs
    - **Property 18: Unique Workflow IDs**
    - **Validates: Requirements 8.3, 13.1**

  - [ ]* 7.3 Write property test for analysis report completeness
    - **Property 13: Analysis Report Completeness**
    - **Validates: Requirements 6.2, 4.5, 3.5**

  - [ ]* 7.4 Write property test for JSON output format
    - **Property 14: JSON Output Format**
    - **Validates: Requirements 6.3**

  - [ ]* 7.5 Write property test for error output to stderr
    - **Property 34: Error Output to Stderr**
    - **Validates: Requirements 14.2, 19.5**

  - [ ]* 7.6 Write property test for workflow state persistence
    - **Property 31: Workflow State Persistence**
    - **Validates: Requirements 13.6**

  - [ ]* 7.7 Write unit tests for CLI and orchestration
    - Test command flag parsing
    - Test timeout enforcement
    - Test workflow ID generation
    - Test state persistence and retrieval
    - _Requirements: 6.1, 8.3, 13.1, 13.6_

- [ ] 8. Checkpoint - Ensure skill-analyzer tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. Implement skill-builder: Build environment and builder engine
  - [ ] 9.1 Create build environment manager
    - Implement isolated workspace creation with unique directories
    - Implement cleanup handlers for temporary files
    - Support concurrent builds with directory isolation
    - _Requirements: 9.1, 14.6, 17.2, 17.5_

  - [ ]* 9.2 Write property test for build environment isolation
    - **Property 19: Build Environment Isolation**
    - **Validates: Requirements 9.1, 17.2**

  - [ ]* 9.3 Write property test for temporary file cleanup
    - **Property 35: Temporary File Cleanup**
    - **Validates: Requirements 14.6**

  - [ ] 9.4 Create builder engine for Go and Rust compilation
    - Implement SOP step execution in order
    - Implement Go static compilation (CGO_ENABLED=0, ldflags, tags netgo)
    - Implement Rust static compilation (musl target, release profile)
    - Capture build logs and timing
    - Handle build failures with detailed error reporting
    - _Requirements: 9.2, 9.3, 9.4, 9.6, 9.7_

  - [ ]* 9.5 Write property test for build step execution order
    - **Property 20: Build Step Execution Order**
    - **Validates: Requirements 9.2**

  - [ ]* 9.6 Write property test for build report generation
    - **Property 22: Build Report Generation**
    - **Validates: Requirements 9.7, 10.6**

  - [ ]* 9.7 Write unit tests for builder engine
    - Test Go compilation with correct flags
    - Test Rust compilation with musl target
    - Test build failure handling
    - Test build log capture
    - _Requirements: 9.3, 9.4, 9.6_

- [ ] 10. Implement skill-builder: Binary validator
  - [ ] 10.1 Create binary validation engine
    - Implement single binary check (file exists and is executable)
    - Implement zero dependency check (ldd command execution)
    - Implement CLI interface check (--help flag test)
    - Implement JSON output check (execute with test args, parse output)
    - Implement size check (<100MB limit)
    - Generate ValidationReport with all check results
    - _Requirements: 9.5, 10.1, 10.2, 10.3, 19.1, 19.2, 19.3, 19.4, 19.5, 19.6_

  - [ ]* 10.2 Write property test for zero dependency validation
    - **Property 21: Zero Dependency Validation**
    - **Validates: Requirements 9.5, 10.3**

  - [ ]* 10.3 Write property test for binary interface validation
    - **Property 23: Binary Interface Validation**
    - **Validates: Requirements 10.2, 19.4**

  - [ ]* 10.4 Write property test for binary size limit
    - **Property 24: Binary Size Limit**
    - **Validates: Requirements 19.2**

  - [ ]* 10.5 Write property test for validation failure blocks deployment
    - **Property 43: Validation Failure Blocks Deployment**
    - **Validates: Requirements 19.6**

  - [ ]* 10.6 Write unit tests for validator
    - Test each validation check independently
    - Test validation report generation
    - Test edge cases (binary at size limit, missing --help)
    - _Requirements: 10.2, 10.3, 19.2, 19.4_

- [ ] 11. Implement skill-builder: Deployer
  - [ ] 11.1 Create deployment engine
    - Implement file copy to target directory (/root/workspace/agents/)
    - Implement permission setting (chmod 755)
    - Implement post-deployment executability test
    - Implement k3s deployment (optional, based on config)
    - Generate DeploymentResult with path, permissions, and timestamp
    - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5, 11.6_

  - [ ]* 11.2 Write property test for deployment file permissions
    - **Property 25: Deployment File Permissions**
    - **Validates: Requirements 11.2**

  - [ ]* 11.3 Write property test for post-deployment executability
    - **Property 26: Post-Deployment Executability**
    - **Validates: Requirements 11.5**

  - [ ]* 11.4 Write property test for deployment info recording
    - **Property 27: Deployment Info Recording**
    - **Validates: Requirements 11.6**

  - [ ]* 11.5 Write unit tests for deployer
    - Test file copy operation
    - Test permission setting
    - Test k3s deployment (if enabled)
    - Test deployment failure handling
    - _Requirements: 11.1, 11.2, 11.4, 11.5_

- [ ] 12. Implement skill-builder: GitHub backup
  - [ ] 12.1 Create GitHub backup engine
    - Implement git repository initialization
    - Implement commit creation with source code and binary
    - Implement version tag creation (timestamp or semantic versioning)
    - Implement push to backup repository with retry logic
    - Generate BackupResult with commit hash, tag, and repository URL
    - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5, 12.6_

  - [ ]* 12.2 Write property test for git backup commit
    - **Property 28: Git Backup Commit**
    - **Validates: Requirements 12.2**

  - [ ]* 12.3 Write property test for version tag creation
    - **Property 29: Version Tag Creation**
    - **Validates: Requirements 12.3**

  - [ ]* 12.4 Write property test for backup report completeness
    - **Property 30: Backup Report Completeness**
    - **Validates: Requirements 12.6**

  - [ ]* 12.5 Write unit tests for GitHub backup
    - Test git initialization
    - Test commit creation
    - Test tag format validation
    - Test push retry logic
    - _Requirements: 12.1, 12.2, 12.3, 12.4_

- [ ] 13. Implement skill-builder: CLI and orchestration
  - [ ] 13.1 Create Cobra CLI interface and build orchestration
    - Implement build command with flags (output, work-dir, deploy-dir, skip-deploy, skip-backup, config, verbose)
    - Load and parse AnalysisReport from input (file or stdin)
    - Orchestrate build pipeline (setup → build → validate → deploy → backup)
    - Generate complete BuildReport with all sections
    - Output JSON to stdout, errors to stderr
    - Implement workflow state updates
    - _Requirements: 8.2, 9.1, 10.1, 11.1, 12.1, 13.2, 13.6_

  - [ ]* 13.2 Write property test for stage timing records
    - **Property 32: Stage Timing Records**
    - **Validates: Requirements 13.2**

  - [ ]* 13.3 Write property test for error detail recording
    - **Property 33: Error Detail Recording**
    - **Validates: Requirements 13.4, 14.1**

  - [ ]* 13.4 Write unit tests for CLI and orchestration
    - Test command flag parsing
    - Test analysis report loading from file and stdin
    - Test skip flags (skip-deploy, skip-backup)
    - Test workflow state updates
    - _Requirements: 8.2, 13.2, 13.6_

- [ ] 14. Checkpoint - Ensure skill-builder tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 15. Implement error handling and retry logic
  - [ ] 15.1 Add retry logic for external service errors
    - Implement exponential backoff for GitHub operations (2s, 4s, 8s)
    - Implement retry for network timeouts
    - Implement retry counter and max attempts (3)
    - _Requirements: 8.5, 12.5, 14.3, 14.4_

  - [ ]* 15.2 Write unit tests for retry logic
    - Test retry count and timing
    - Test exponential backoff intervals
    - Test max retry limit
    - _Requirements: 14.3, 14.4_

  - [ ] 15.3 Implement error categorization and exit codes
    - Define error categories (user_input=1, external_service=2, security=3, build_validation=4)
    - Implement consistent error JSON format for stderr
    - Implement appropriate exit codes for each category
    - _Requirements: 14.1, 14.2_

  - [ ]* 15.4 Write unit tests for error handling
    - Test error categorization
    - Test exit codes
    - Test error JSON format
    - _Requirements: 14.1, 14.2_

- [ ] 16. Implement logging infrastructure
  - [ ] 16.1 Create structured JSON logger
    - Implement log levels (DEBUG, INFO, WARN, ERROR)
    - Implement JSON format with timestamp, level, message fields
    - Implement log output to file and stdout
    - Implement log rotation (daily or size-based)
    - _Requirements: 15.1, 15.2, 15.3, 15.4, 15.5_

  - [ ]* 16.2 Write property test for structured log format
    - **Property 36: Structured Log Format**
    - **Validates: Requirements 15.2, 15.3**

  - [ ]* 16.3 Write unit tests for logging
    - Test log level filtering
    - Test JSON format validation
    - Test log rotation
    - _Requirements: 15.2, 15.3, 15.5_

- [ ] 17. Implement concurrent workflow processing
  - [ ] 17.1 Add concurrency support to skill-analyzer
    - Implement workflow queue with configurable limit (default: 3)
    - Implement resource isolation for concurrent workflows
    - Implement queue status tracking
    - _Requirements: 17.1, 17.2, 17.3, 17.4, 17.5, 17.6_

  - [ ]* 17.2 Write property test for concurrent workflow processing
    - **Property 40: Concurrent Workflow Processing**
    - **Validates: Requirements 17.1**

  - [ ]* 17.3 Write property test for queue overflow handling
    - **Property 41: Queue Overflow Handling**
    - **Validates: Requirements 17.4**

  - [ ]* 17.4 Write unit tests for concurrency
    - Test parallel execution of multiple workflows
    - Test queue behavior at limit
    - Test resource isolation
    - _Requirements: 17.1, 17.3, 17.4_

- [ ] 18. Implement configuration management
  - [ ] 18.1 Create configuration loader with validation
    - Implement JSON configuration file parsing
    - Implement environment variable override logic
    - Implement configuration validation
    - Implement default value fallback
    - _Requirements: 16.1, 16.2, 16.3, 16.4, 16.5_

  - [ ]* 18.2 Write property test for configuration file loading
    - **Property 37: Configuration File Loading**
    - **Validates: Requirements 16.1**

  - [ ]* 18.3 Write unit tests for configuration
    - Test file parsing
    - Test environment variable override
    - Test validation
    - Test default values
    - _Requirements: 16.1, 16.2, 16.3, 16.4_

- [ ] 19. Create integration tests
  - [ ]* 19.1 Write integration test for happy path workflow
    - Test complete flow: valid repo → analysis → build → deploy → backup
    - Verify all intermediate outputs
    - _Requirements: 20.1_

  - [ ]* 19.2 Write integration test for security failure
    - Test malicious repo → analysis → termination
    - Verify no build or deployment occurs
    - _Requirements: 20.2_

  - [ ]* 19.3 Write integration test for build failure
    - Test invalid code → analysis → build failure → error report
    - _Requirements: 20.3_

  - [ ]* 19.4 Write integration test for concurrent processing
    - Test multiple repos → parallel analysis → all complete
    - _Requirements: 20.4_

  - [ ]* 19.5 Write integration test for idempotency
    - **Property 44: Integration Test Idempotency**
    - **Validates: Requirements 20.6**

- [ ] 20. Build and validate final binaries
  - [ ] 20.1 Compile skill-analyzer with static linking
    - Execute: CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -o skill-analyzer
    - Verify with ldd (should show "not a dynamic executable")
    - Test --help flag and basic functionality
    - _Requirements: 9.3, 9.5_

  - [ ] 20.2 Compile skill-builder with static linking
    - Execute: CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -o skill-builder
    - Verify with ldd (should show "not a dynamic executable")
    - Test --help flag and basic functionality
    - _Requirements: 9.3, 9.5_

  - [ ] 20.3 Validate both binaries meet Binary Arsenal specification
    - Verify single binary (no external files)
    - Verify zero dependencies (ldd check)
    - Verify CLI interface (--help works)
    - Verify JSON output format
    - Verify stderr error output
    - Verify size <100MB
    - _Requirements: 19.1, 19.2, 19.3, 19.4, 19.5_

- [ ] 21. Final checkpoint - Complete system validation
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation at key milestones
- Property tests validate universal correctness properties across all inputs
- Unit tests validate specific examples, edge cases, and error conditions
- Both tools use Go and follow the Binary Arsenal specification
- Integration tests verify end-to-end workflows in test mode (no production deployment)
- The implementation creates two independent CLI tools that can be executed separately
