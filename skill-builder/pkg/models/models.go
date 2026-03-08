package models

import "time"

type BuildReport struct {
	WorkflowID string             `json:"workflow_id"`
	Timestamp  time.Time          `json:"timestamp"`
	Build      *BuildResult       `json:"build"`
	Validation *ValidationReport  `json:"validation"`
	Deployment *DeploymentResult  `json:"deployment"`
	Backup     *BackupResult      `json:"backup"`
	Status     string             `json:"status"`
	Message    string             `json:"message"`
}

type BuildResult struct {
	Status     string `json:"status"`
	Language   string `json:"language"`
	SourcePath string `json:"source_path"`
	BinaryPath string `json:"binary_path"`
	BinarySize int64  `json:"binary_size"`
	BuildTime  string `json:"build_time"`
	BuildLog   string `json:"build_log"`
}

type ValidationReport struct {
	Status string             `json:"status"`
	Checks []ValidationCheck  `json:"checks"`
}

type ValidationCheck struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

type DeploymentResult struct {
	Status         string    `json:"status"`
	TargetPath     string    `json:"target_path"`
	Permissions    string    `json:"permissions"`
	K3sDeployed    bool      `json:"k3s_deployed"`
	DeploymentTime time.Time `json:"deployment_time"`
}

type BackupResult struct {
	Status     string    `json:"status"`
	Repository string    `json:"repository"`
	CommitHash string    `json:"commit_hash"`
	Tag        string    `json:"tag"`
	BackupTime time.Time `json:"backup_time"`
}

type AnalysisReport struct {
	WorkflowID string `json:"workflow_id"`
	GitHubURL  string `json:"github_url"`
	Repository struct {
		Name      string `json:"name"`
		ClonePath string `json:"clone_path"`
	} `json:"repository"`
	Recommendation struct {
		TargetLanguage string `json:"target_language"`
	} `json:"recommendation"`
	SOP SOP `json:"sop"`
}

type SOP struct {
	Language      string        `json:"language"`
	Steps         []BuildStep   `json:"steps"`
	InterfaceSpec InterfaceSpec `json:"interface_spec"`
	Dependencies  []string      `json:"dependencies"`
}

type BuildStep struct {
	Phase       string `json:"phase"`
	Command     string `json:"command"`
	Description string `json:"description"`
	Expected    string `json:"expected,omitempty"`
}

type InterfaceSpec struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Errors string `json:"errors"`
}

type ErrorOutput struct {
	Error       string    `json:"error"`
	Details     string    `json:"details"`
	Category    string    `json:"category"`
	WorkflowID  string    `json:"workflow_id"`
	Timestamp   time.Time `json:"timestamp"`
	Retryable   bool      `json:"retryable"`
	RetryCount  int       `json:"retry_count"`
	Suggestions []string  `json:"suggestions,omitempty"`
}
