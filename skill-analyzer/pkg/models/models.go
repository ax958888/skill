package models

import "time"

type AnalysisReport struct {
	WorkflowID     string          `json:"workflow_id"`
	Timestamp      time.Time       `json:"timestamp"`
	GitHubURL      string          `json:"github_url"`
	Repository     *Repository     `json:"repository"`
	Security       *SecurityReport `json:"security"`
	Language       *LanguageInfo   `json:"language"`
	TypeAnalysis   *TypeInfo       `json:"type_analysis"`
	Recommendation *Recommendation `json:"recommendation"`
	SOP            *SOP            `json:"sop"`
	Status         string          `json:"status"`
	Message        string          `json:"message"`
}

type Repository struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	ClonePath  string `json:"clone_path"`
	CommitHash string `json:"commit_hash"`
}

type SecurityReport struct {
	Score  int             `json:"score"`
	Status string          `json:"status"`
	Issues []SecurityIssue `json:"issues"`
}

type SecurityIssue struct {
	Severity       string `json:"severity"`
	Type           string `json:"type"`
	File           string `json:"file"`
	Line           int    `json:"line"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
}

type LanguageInfo struct {
	Primary      string   `json:"primary"`
	Secondary    []string `json:"secondary"`
	Confidence   float64  `json:"confidence"`
	BuildFiles   []string `json:"build_files"`
	Dependencies []string `json:"dependencies"`
}

type TypeInfo struct {
	Category        string   `json:"category"`
	UseCases        []string `json:"use_cases"`
	InputInterface  string   `json:"input_interface"`
	OutputInterface string   `json:"output_interface"`
	TargetUsers     string   `json:"target_users"`
}

type Recommendation struct {
	TargetLanguage     string `json:"target_language"`
	Rationale          string `json:"rationale"`
	DecisionMatrixRule string `json:"decision_matrix_rule"`
}

type SOP struct {
	Language           string         `json:"language"`
	Steps              []BuildStep    `json:"steps"`
	InterfaceSpec      InterfaceSpec  `json:"interface_spec"`
	Dependencies       []string       `json:"dependencies"`
	EstimatedBuildTime string         `json:"estimated_build_time"`
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
