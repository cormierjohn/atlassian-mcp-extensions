package tools

import (
	"context"
	"fmt"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type BulkTransitionJiraIssuesInput struct {
	IssueKeys      []string `json:"issueKeys" jsonschema:"required,description=Jira issue keys to transition"`
	TargetStatus   string   `json:"targetStatus,omitempty" jsonschema:"description=Destination status name to match case-insensitively"`
	TransitionName string   `json:"transitionName,omitempty" jsonschema:"description=Transition name to match case-insensitively"`
	Comment        string   `json:"comment,omitempty" jsonschema:"description=Optional plain-text comment to add during the transition"`
	StopOnError    bool     `json:"stopOnError,omitempty" jsonschema:"description=Stop after the first failed issue"`
	DryRun         bool     `json:"dryRun,omitempty" jsonschema:"description=Resolve transitions but do not modify Jira"`
}

func BulkTransitionJiraIssuesHandler(ctx context.Context, req *mcp.CallToolRequest, input BulkTransitionJiraIssuesInput) (*mcp.CallToolResult, BulkResult, error) {
	if len(input.IssueKeys) == 0 {
		return nil, BulkResult{}, fmt.Errorf("issueKeys is required")
	}
	if (input.TargetStatus == "") == (input.TransitionName == "") {
		return nil, BulkResult{}, fmt.Errorf("provide exactly one of targetStatus or transitionName")
	}

	client, err := jira.NewFromEnv()
	if err != nil {
		return nil, BulkResult{}, err
	}

	output := BulkResult{DryRun: input.DryRun}
	for _, issueKey := range input.IssueKeys {
		result := ItemResult{Key: issueKey}
		transitions, err := client.GetTransitions(ctx, issueKey)
		if err == nil {
			var transition *jira.Transition
			transition, err = jira.SelectTransition(transitions, input.TargetStatus, input.TransitionName)
			if err == nil {
				result.TransitionID = transition.ID
				if input.DryRun {
					result.Status = "dry_run"
				} else {
					err = client.TransitionIssue(ctx, issueKey, *transition, input.Comment)
					if err == nil {
						result.Status = "success"
					}
				}
			}
		}
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
			output.FailureCount++
			output.Results = append(output.Results, result)
			if input.StopOnError {
				break
			}
			continue
		}
		output.SuccessCount++
		output.Results = append(output.Results, result)
	}

	return nil, output, nil
}

func RegisterBulkTransitionJiraIssues(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "bulkTransitionJiraIssues",
		Description: "Transition multiple Jira issues by target status or transition name with per-issue results.",
	}, BulkTransitionJiraIssuesHandler)
}
