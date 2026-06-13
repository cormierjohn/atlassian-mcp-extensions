package tools

import (
	"context"
	"fmt"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TransitionJiraIssueByStatusInput struct {
	IssueKey       string `json:"issueKey" jsonschema:"Jira issue key to transition"`
	TargetStatus   string `json:"targetStatus,omitempty" jsonschema:"Destination status name to match case-insensitively"`
	TransitionName string `json:"transitionName,omitempty" jsonschema:"Transition name to match case-insensitively"`
	Comment        string `json:"comment,omitempty" jsonschema:"Optional plain-text comment to add during the transition"`
	DryRun         bool   `json:"dryRun,omitempty" jsonschema:"Resolve transition but do not modify Jira"`
}

type TransitionJiraIssueByStatusOutput struct {
	IssueKey     string `json:"issueKey" jsonschema:"Jira issue key"`
	TransitionID string `json:"transitionId" jsonschema:"Resolved transition id"`
	Status       string `json:"status" jsonschema:"success or dry_run"`
}

func TransitionJiraIssueByStatusHandler(ctx context.Context, req *mcp.CallToolRequest, input TransitionJiraIssueByStatusInput) (*mcp.CallToolResult, TransitionJiraIssueByStatusOutput, error) {
	if input.IssueKey == "" {
		return nil, TransitionJiraIssueByStatusOutput{}, fmt.Errorf("issueKey is required")
	}
	if (input.TargetStatus == "") == (input.TransitionName == "") {
		return nil, TransitionJiraIssueByStatusOutput{}, fmt.Errorf("provide exactly one of targetStatus or transitionName")
	}

	client, err := jira.NewFromEnv()
	if err != nil {
		return nil, TransitionJiraIssueByStatusOutput{}, err
	}
	transitions, err := client.GetTransitions(ctx, input.IssueKey)
	if err != nil {
		return nil, TransitionJiraIssueByStatusOutput{}, err
	}
	transition, err := jira.SelectTransition(transitions, input.TargetStatus, input.TransitionName)
	if err != nil {
		return nil, TransitionJiraIssueByStatusOutput{}, err
	}
	if input.DryRun {
		return nil, TransitionJiraIssueByStatusOutput{
			IssueKey:     input.IssueKey,
			TransitionID: transition.ID,
			Status:       "dry_run",
		}, nil
	}
	if err := client.TransitionIssue(ctx, input.IssueKey, *transition, input.Comment); err != nil {
		return nil, TransitionJiraIssueByStatusOutput{}, err
	}
	return nil, TransitionJiraIssueByStatusOutput{
		IssueKey:     input.IssueKey,
		TransitionID: transition.ID,
		Status:       "success",
	}, nil
}

func RegisterTransitionJiraIssueByStatus(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "transitionJiraIssueByStatus",
		Description: "Transition a Jira issue by destination status or transition name.",
	}, TransitionJiraIssueByStatusHandler)
}
