package tools

type ItemResult struct {
	Key          string `json:"key,omitempty" jsonschema:"description=Jira issue key or link id"`
	Status       string `json:"status" jsonschema:"description=success, dry_run, or failed"`
	TransitionID string `json:"transitionId,omitempty" jsonschema:"description=Resolved Jira transition id"`
	Error        string `json:"error,omitempty" jsonschema:"description=Error message for failed items"`
}

type BulkResult struct {
	SuccessCount int          `json:"successCount" jsonschema:"description=Number of successful items"`
	FailureCount int          `json:"failureCount" jsonschema:"description=Number of failed items"`
	DryRun       bool         `json:"dryRun" jsonschema:"description=Whether the request only planned changes"`
	Results      []ItemResult `json:"results" jsonschema:"description=Per-item operation results"`
}
