package tools

import (
	"context"
	"fmt"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type BulkUpdateJiraIssuesInput struct {
	IssueKeys             []string       `json:"issueKeys" jsonschema:"Jira issue keys to update"`
	Fields                map[string]any `json:"fields,omitempty" jsonschema:"Jira fields object to set"`
	Update                map[string]any `json:"update,omitempty" jsonschema:"Jira update operations object"`
	SuppressNotifications bool           `json:"suppressNotifications,omitempty" jsonschema:"Pass notifyUsers=false to Jira"`
	StopOnError           bool           `json:"stopOnError,omitempty" jsonschema:"Stop after the first failed issue"`
	DryRun                bool           `json:"dryRun,omitempty" jsonschema:"Validate inputs but do not modify Jira"`
}

func BulkUpdateJiraIssuesHandler(ctx context.Context, req *mcp.CallToolRequest, input BulkUpdateJiraIssuesInput) (*mcp.CallToolResult, BulkResult, error) {
	if len(input.IssueKeys) == 0 {
		return nil, BulkResult{}, fmt.Errorf("issueKeys is required")
	}
	if len(input.Fields) == 0 && len(input.Update) == 0 {
		return nil, BulkResult{}, fmt.Errorf("at least one of fields or update is required")
	}

	client, err := jira.NewFromEnv()
	if err != nil {
		return nil, BulkResult{}, err
	}

	output := BulkResult{DryRun: input.DryRun}
	for _, issueKey := range input.IssueKeys {
		result := ItemResult{Key: issueKey}
		if input.DryRun {
			result.Status = "dry_run"
		} else if err := client.UpdateIssue(ctx, issueKey, input.Fields, input.Update, input.SuppressNotifications); err != nil {
			result.Status = "failed"
			result.Error = err.Error()
			output.FailureCount++
			output.Results = append(output.Results, result)
			if input.StopOnError {
				break
			}
			continue
		} else {
			result.Status = "success"
		}
		output.SuccessCount++
		output.Results = append(output.Results, result)
	}

	return nil, output, nil
}

func RegisterBulkUpdateJiraIssues(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "bulkUpdateJiraIssues",
		Description: "Update fields or run update operations against multiple Jira issues with per-issue results.",
	}, BulkUpdateJiraIssuesHandler)
}
