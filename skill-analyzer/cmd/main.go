package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/binary-arsenal/skill-analyzer/pkg/models"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	workDir    string
	timeout    int
	verbose    bool
	workflowID string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "skill-analyzer",
		Short: "Analyze GitHub repositories for skill transformation",
	}

	analyzeCmd := &cobra.Command{
		Use:   "analyze <github-url>",
		Short: "Analyze a GitHub repository",
		Args:  cobra.ExactArgs(1),
		RunE:  runAnalyze,
	}

	analyzeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	analyzeCmd.Flags().StringVarP(&workDir, "work-dir", "w", "", "Working directory")
	analyzeCmd.Flags().IntVarP(&timeout, "timeout", "t", 300, "Timeout in seconds")
	analyzeCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	analyzeCmd.Flags().StringVar(&workflowID, "workflow-id", "", "Workflow ID")

	rootCmd.AddCommand(analyzeCmd)

	if err := rootCmd.Execute(); err != nil {
		outputError("command_execution", err.Error(), "", true, 0)
		os.Exit(1)
	}
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	githubURL := args[0]

	if workflowID == "" {
		workflowID = uuid.New().String()
	}

	if workDir == "" {
		workDir = filepath.Join(os.TempDir(), fmt.Sprintf("skill-analyzer-%s", workflowID[:8]))
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Starting analysis: %s\n", githubURL)
		fmt.Fprintf(os.Stderr, "Workflow ID: %s\n", workflowID)
	}

	// Validate GitHub URL
	if !isValidGitHubURL(githubURL) {
		return outputError("user_input", "Invalid GitHub URL", githubURL, false, 0)
	}

	// Clone repository
	repo, err := cloneRepository(githubURL, workDir)
	if err != nil {
		return outputError("external_service", "Failed to clone repository", err.Error(), true, 0)
	}

	// Security analysis
	security := analyzeSecurity(repo.ClonePath)

	if security.Status == "unsafe" {
		return outputError("security", "High-risk security issues detected", "Repository contains malicious patterns", false, 0)
	}

	// Language detection
	language := detectLanguage(repo.ClonePath)

	// Type analysis
	typeInfo := analyzeType(repo.ClonePath, language)

	// Generate recommendation
	recommendation := generateRecommendation(typeInfo)

	// Generate SOP
	sop := generateSOP(recommendation.TargetLanguage, repo.Name)

	// Build report
	report := models.AnalysisReport{
		WorkflowID:     workflowID,
		Timestamp:      time.Now(),
		GitHubURL:      githubURL,
		Repository:     repo,
		Security:       security,
		Language:       language,
		TypeAnalysis:   typeInfo,
		Recommendation: recommendation,
		SOP:            sop,
		Status:         "success",
		Message:        "Analysis completed successfully",
	}

	return outputReport(report)
}

func isValidGitHubURL(url string) bool {
	pattern := `^https://github\.com/[\w-]+/[\w.-]+/?$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

func cloneRepository(url, destDir string) (*models.Repository, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	owner := parts[len(parts)-2]
	clonePath := filepath.Join(destDir, repoName)

	cmd := exec.Command("git", "clone", "--depth", "1", url, clonePath)
	if verbose {
		fmt.Fprintf(os.Stderr, "Cloning: %s\n", url)
	}

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git clone failed: %w", err)
	}

	// Get commit hash
	hashCmd := exec.Command("git", "-C", clonePath, "rev-parse", "HEAD")
	hashOut, _ := hashCmd.Output()
	commitHash := strings.TrimSpace(string(hashOut))

	return &models.Repository{
		Name:       repoName,
		Owner:      owner,
		ClonePath:  clonePath,
		CommitHash: commitHash,
	}, nil
}

func analyzeSecurity(repoPath string) *models.SecurityReport {
	issues := []models.SecurityIssue{}
	
	// Simple pattern matching for common security issues
	patterns := map[string]string{
		"eval\\(":           "arbitrary_execution",
		"exec\\(":           "arbitrary_execution",
		"__import__":        "arbitrary_execution",
		"os\\.system":       "arbitrary_execution",
		"subprocess\\.call": "arbitrary_execution",
	}

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".py" && ext != ".js" && ext != ".ts" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		for pattern, issueType := range patterns {
			re := regexp.MustCompile(pattern)
			if re.Match(content) {
				issues = append(issues, models.SecurityIssue{
					Severity:       "high",
					Type:           issueType,
					File:           path,
					Line:           0,
					Description:    fmt.Sprintf("Detected pattern: %s", pattern),
					Recommendation: "Review and replace with safe alternative",
				})
			}
		}

		return nil
	})

	score := 100 - (len(issues) * 20)
	if score < 0 {
		score = 0
	}

	status := "safe"
	if score < 50 {
		status = "unsafe"
	} else if score < 80 {
		status = "warning"
	}

	return &models.SecurityReport{
		Score:  score,
		Status: status,
		Issues: issues,
	}
}

func detectLanguage(repoPath string) *models.LanguageInfo {
	langCounts := make(map[string]int)
	buildFiles := []string{}
	dependencies := []string{}

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		name := info.Name()
		ext := filepath.Ext(name)

		// Count language files
		switch ext {
		case ".py":
			langCounts["python"]++
		case ".js":
			langCounts["javascript"]++
		case ".ts":
			langCounts["typescript"]++
		case ".go":
			langCounts["go"]++
		case ".rs":
			langCounts["rust"]++
		}

		// Detect build files
		switch name {
		case "requirements.txt", "setup.py", "pyproject.toml":
			buildFiles = append(buildFiles, name)
		case "package.json", "package-lock.json":
			buildFiles = append(buildFiles, name)
		case "go.mod", "go.sum":
			buildFiles = append(buildFiles, name)
		case "Cargo.toml", "Cargo.lock":
			buildFiles = append(buildFiles, name)
		}

		return nil
	})

	// Find primary language
	primary := "unknown"
	maxCount := 0
	for lang, count := range langCounts {
		if count > maxCount {
			maxCount = count
			primary = lang
		}
	}

	secondary := []string{}
	for lang, count := range langCounts {
		if lang != primary && count > 0 {
			secondary = append(secondary, lang)
		}
	}

	confidence := 0.95
	if maxCount < 5 {
		confidence = 0.7
	}

	return &models.LanguageInfo{
		Primary:      primary,
		Secondary:    secondary,
		Confidence:   confidence,
		BuildFiles:   buildFiles,
		Dependencies: dependencies,
	}
}

func analyzeType(repoPath string, lang *models.LanguageInfo) *models.TypeInfo {
	category := "automation"
	useCases := []string{"Script automation"}

	// Simple heuristics based on file content
	hasNetwork := false
	hasFileOps := false

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		contentStr := string(content)
		if strings.Contains(contentStr, "http") || strings.Contains(contentStr, "requests") || strings.Contains(contentStr, "fetch") {
			hasNetwork = true
		}
		if strings.Contains(contentStr, "open(") || strings.Contains(contentStr, "readFile") || strings.Contains(contentStr, "writeFile") {
			hasFileOps = true
		}

		return nil
	})

	if hasNetwork {
		category = "network_request"
		useCases = []string{"API client", "Data fetching"}
	} else if hasFileOps {
		category = "file_processing"
		useCases = []string{"File manipulation", "Data processing"}
	}

	return &models.TypeInfo{
		Category:        category,
		UseCases:        useCases,
		InputInterface:  "CLI arguments",
		OutputInterface: "stdout text",
		TargetUsers:     "developers",
	}
}

func generateRecommendation(typeInfo *models.TypeInfo) *models.Recommendation {
	targetLang := "go"
	rationale := "Go excels at automation and system operations"
	rule := "automation → Go"

	switch typeInfo.Category {
	case "network_request", "cloud_api", "k3s_docker", "automation":
		targetLang = "go"
		rationale = "Go excels at network operations and system automation"
		rule = fmt.Sprintf("%s → Go", typeInfo.Category)
	case "file_processing", "encryption", "parsing", "computation":
		targetLang = "rust"
		rationale = "Rust provides memory safety and performance for intensive operations"
		rule = fmt.Sprintf("%s → Rust", typeInfo.Category)
	}

	return &models.Recommendation{
		TargetLanguage:     targetLang,
		Rationale:          rationale,
		DecisionMatrixRule: rule,
	}
}

func generateSOP(language, skillName string) *models.SOP {
	if language == "rust" {
		return &models.SOP{
			Language: "rust",
			Steps: []models.BuildStep{
				{Phase: "setup", Command: "cargo init --name " + skillName, Description: "Initialize Rust project"},
				{Phase: "build", Command: "cargo build --release --target x86_64-unknown-linux-musl", Description: "Static compilation"},
				{Phase: "validate", Command: "ldd target/x86_64-unknown-linux-musl/release/" + skillName, Expected: "not a dynamic executable"},
			},
			InterfaceSpec: models.InterfaceSpec{
				Input:  "CLI arguments via Clap",
				Output: "JSON to stdout",
				Errors: "stderr with non-zero exit code",
			},
			Dependencies:       []string{"clap"},
			EstimatedBuildTime: "3-7 minutes",
		}
	}

	return &models.SOP{
		Language: "go",
		Steps: []models.BuildStep{
			{Phase: "setup", Command: "go mod init " + skillName, Description: "Initialize Go module"},
			{Phase: "build", Command: "CGO_ENABLED=0 go build -ldflags=\"-s -w\" -tags netgo -o " + skillName, Description: "Static compilation"},
			{Phase: "validate", Command: "ldd " + skillName, Expected: "not a dynamic executable"},
		},
		InterfaceSpec: models.InterfaceSpec{
			Input:  "CLI arguments via Cobra",
			Output: "JSON to stdout",
			Errors: "stderr with non-zero exit code",
		},
		Dependencies:       []string{"github.com/spf13/cobra"},
		EstimatedBuildTime: "2-5 minutes",
	}
}

func outputReport(report models.AnalysisReport) error {
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
		WorkflowID: workflowID,
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
	case "security":
		return fmt.Errorf("exit code 3")
	default:
		return fmt.Errorf("exit code 4")
	}
}
