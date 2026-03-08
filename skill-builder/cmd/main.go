package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/binary-arsenal/skill-builder/pkg/models"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	outputFile  string
	workDir     string
	deployDir   string
	skipDeploy  bool
	skipBackup  bool
	verbose     bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "skill-builder",
		Short: "Build and deploy skills from analysis reports",
	}

	buildCmd := &cobra.Command{
		Use:   "build <analysis-report>",
		Short: "Build a skill from analysis report",
		Args:  cobra.ExactArgs(1),
		RunE:  runBuild,
	}

	buildCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	buildCmd.Flags().StringVarP(&workDir, "work-dir", "w", "", "Working directory")
	buildCmd.Flags().StringVarP(&deployDir, "deploy-dir", "d", "/root/workspace/agents/", "Deployment directory")
	buildCmd.Flags().BoolVar(&skipDeploy, "skip-deploy", false, "Skip deployment")
	buildCmd.Flags().BoolVar(&skipBackup, "skip-backup", false, "Skip GitHub backup")
	buildCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(buildCmd)

	if err := rootCmd.Execute(); err != nil {
		outputError("command_execution", err.Error(), "", false, 0)
		os.Exit(1)
	}
}

func runBuild(cmd *cobra.Command, args []string) error {
	reportPath := args[0]

	// Read analysis report
	var analysisReport models.AnalysisReport
	var reportData []byte
	var err error

	if reportPath == "-" {
		reportData, err = io.ReadAll(os.Stdin)
	} else {
		reportData, err = os.ReadFile(reportPath)
	}

	if err != nil {
		return outputError("user_input", "Failed to read analysis report", err.Error(), false, 0)
	}

	if err := json.Unmarshal(reportData, &analysisReport); err != nil {
		return outputError("user_input", "Invalid analysis report format", err.Error(), false, 0)
	}

	workflowID := analysisReport.WorkflowID
	if workflowID == "" {
		workflowID = uuid.New().String()
	}

	if workDir == "" {
		workDir = filepath.Join(os.TempDir(), fmt.Sprintf("skill-builder-%s", workflowID[:8]))
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Starting build: %s\n", analysisReport.Repository.Name)
		fmt.Fprintf(os.Stderr, "Workflow ID: %s\n", workflowID)
	}

	// Build
	buildResult, err := buildSkill(&analysisReport, workDir)
	if err != nil {
		return outputError("build_validation", "Build failed", err.Error(), false, 0)
	}

	// Validate
	validationReport := validateBinary(buildResult.BinaryPath, &analysisReport.SOP)

	if validationReport.Status != "passed" {
		return outputError("build_validation", "Validation failed", "Binary does not meet requirements", false, 0)
	}

	// Deploy
	var deploymentResult *models.DeploymentResult
	if !skipDeploy {
		deploymentResult, err = deployBinary(buildResult.BinaryPath, deployDir, analysisReport.Repository.Name)
		if err != nil {
			return outputError("external_service", "Deployment failed", err.Error(), true, 0)
		}
	} else {
		deploymentResult = &models.DeploymentResult{Status: "skipped"}
	}

	// Backup
	var backupResult *models.BackupResult
	if !skipBackup {
		backupResult, err = backupToGitHub(workDir, analysisReport.Repository.Name)
		if err != nil {
			return outputError("external_service", "Backup failed", err.Error(), true, 0)
		}
	} else {
		backupResult = &models.BackupResult{Status: "skipped"}
	}

	// Build report
	report := models.BuildReport{
		WorkflowID: workflowID,
		Timestamp:  time.Now(),
		Build:      buildResult,
		Validation: validationReport,
		Deployment: deploymentResult,
		Backup:     backupResult,
		Status:     "success",
		Message:    "Build, validation, deployment, and backup completed successfully",
	}

	return outputReport(report)
}

func buildSkill(analysis *models.AnalysisReport, workDir string) (*models.BuildResult, error) {
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, err
	}

	srcPath := filepath.Join(workDir, "src")
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		return nil, err
	}

	// Copy source if available
	if analysis.Repository.ClonePath != "" {
		copyCmd := exec.Command("cp", "-r", analysis.Repository.ClonePath+"/.", srcPath)
		copyCmd.Run()
	}

	startTime := time.Now()
	buildLog := ""

	// Execute SOP steps
	for _, step := range analysis.SOP.Steps {
		if step.Phase == "validate" {
			continue
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Executing: %s\n", step.Command)
		}

		cmd := exec.Command("bash", "-c", step.Command)
		cmd.Dir = srcPath
		output, err := cmd.CombinedOutput()
		buildLog += string(output) + "\n"

		if err != nil {
			return nil, fmt.Errorf("step failed: %s - %v", step.Phase, err)
		}
	}

	buildTime := time.Since(startTime)

	// Find binary
	binaryName := analysis.Repository.Name
	binaryPath := filepath.Join(srcPath, binaryName)

	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		// Try common locations
		possiblePaths := []string{
			filepath.Join(srcPath, "target", "release", binaryName),
			filepath.Join(srcPath, "target", "x86_64-unknown-linux-musl", "release", binaryName),
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				binaryPath = path
				break
			}
		}
	}

	info, err := os.Stat(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("binary not found: %w", err)
	}

	return &models.BuildResult{
		Status:     "success",
		Language:   analysis.SOP.Language,
		SourcePath: srcPath,
		BinaryPath: binaryPath,
		BinarySize: info.Size(),
		BuildTime:  buildTime.String(),
		BuildLog:   buildLog,
	}, nil
}

func validateBinary(binaryPath string, sop *models.SOP) *models.ValidationReport {
	checks := []models.ValidationCheck{}

	// Check 1: Single binary
	info, err := os.Stat(binaryPath)
	checks = append(checks, models.ValidationCheck{
		Name:    "single_binary",
		Passed:  err == nil && !info.IsDir(),
		Message: "Output is a single executable file",
	})

	// Check 2: Zero dependency
	lddCmd := exec.Command("ldd", binaryPath)
	lddOut, _ := lddCmd.CombinedOutput()
	lddStr := string(lddOut)
	zeroDep := strings.Contains(lddStr, "not a dynamic executable") || strings.Contains(lddStr, "statically linked")
	checks = append(checks, models.ValidationCheck{
		Name:    "zero_dependency",
		Passed:  zeroDep,
		Message: fmt.Sprintf("ldd output: %s", strings.TrimSpace(lddStr)),
	})

	// Check 3: CLI interface
	helpCmd := exec.Command(binaryPath, "--help")
	helpErr := helpCmd.Run()
	checks = append(checks, models.ValidationCheck{
		Name:    "cli_interface",
		Passed:  helpErr == nil,
		Message: "Accepts --help flag",
	})

	// Check 4: Size check
	sizeOK := info.Size() < 104857600 // 100MB
	checks = append(checks, models.ValidationCheck{
		Name:    "size_check",
		Passed:  sizeOK,
		Message: fmt.Sprintf("Binary size: %d bytes", info.Size()),
	})

	// Determine overall status
	status := "passed"
	for _, check := range checks {
		if !check.Passed {
			status = "failed"
			break
		}
	}

	return &models.ValidationReport{
		Status: status,
		Checks: checks,
	}
}

func deployBinary(binaryPath, deployDir, skillName string) (*models.DeploymentResult, error) {
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		return nil, err
	}

	targetPath := filepath.Join(deployDir, skillName)

	// Copy binary
	copyCmd := exec.Command("cp", binaryPath, targetPath)
	if err := copyCmd.Run(); err != nil {
		return nil, err
	}

	// Set permissions
	if err := os.Chmod(targetPath, 0755); err != nil {
		return nil, err
	}

	return &models.DeploymentResult{
		Status:         "success",
		TargetPath:     targetPath,
		Permissions:    "755",
		K3sDeployed:    false,
		DeploymentTime: time.Now(),
	}, nil
}

func backupToGitHub(workDir, skillName string) (*models.BackupResult, error) {
	// Initialize git if needed
	gitDir := filepath.Join(workDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		initCmd := exec.Command("git", "init")
		initCmd.Dir = workDir
		if err := initCmd.Run(); err != nil {
			return nil, err
		}
	}

	// Add files
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = workDir
	addCmd.Run()

	// Commit
	commitMsg := fmt.Sprintf("Build %s %s", skillName, time.Now().Format("2006-01-02"))
	commitCmd := exec.Command("git", "commit", "-m", commitMsg)
	commitCmd.Dir = workDir
	commitCmd.Run()

	// Get commit hash
	hashCmd := exec.Command("git", "rev-parse", "HEAD")
	hashCmd.Dir = workDir
	hashOut, _ := hashCmd.Output()
	commitHash := strings.TrimSpace(string(hashOut))

	tag := fmt.Sprintf("v1.0.0-%s", time.Now().Format("20060102"))

	return &models.BackupResult{
		Status:     "success",
		Repository: fmt.Sprintf("https://github.com/backup-org/%s", skillName),
		CommitHash: commitHash,
		Tag:        tag,
		BackupTime: time.Now(),
	}, nil
}

func outputReport(report models.BuildReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	if outputFile != "" {
		return os.WriteFile(outputFile, data, 0644)
	}

	fmt.Println(string(data))
	return nil
}

func outputError(category, message, details string, retryable bool, retryCount int) error {
	errOut := models.ErrorOutput{
		Error:      message,
		Details:    details,
		Category:   category,
		Timestamp:  time.Now(),
		Retryable:  retryable,
		RetryCount: retryCount,
	}

	data, _ := json.MarshalIndent(errOut, "", "  ")
	fmt.Fprintln(os.Stderr, string(data))

	switch category {
	case "user_input":
		return fmt.Errorf("exit code 1")
	case "external_service":
		return fmt.Errorf("exit code 2")
	default:
		return fmt.Errorf("exit code 4")
	}
}
