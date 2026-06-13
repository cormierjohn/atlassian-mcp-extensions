package tools

type ItemResult struct {
	Key          string `json:"key,omitempty" jsonschema:"Jira issue key or link id"`
	Status       string `json:"status" jsonschema:"success, dry_run, or failed"`
	TransitionID string `json:"transitionId,omitempty" jsonschema:"Resolved Jira transition id"`
	Error        string `json:"error,omitempty" jsonschema:"Error message for failed items"`
}

type BulkResult struct {
	SuccessCount int          `json:"successCount" jsonschema:"Number of successful items"`
	FailureCount int          `json:"failureCount" jsonschema:"Number of failed items"`
	DryRun       bool         `json:"dryRun" jsonschema:"Whether the request only planned changes"`
	Results      []ItemResult `json:"results" jsonschema:"Per-item operation results"`
}
