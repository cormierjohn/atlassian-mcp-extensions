package tools

import (
	"context"
	"fmt"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RemoveJiraIssueLinkInput struct {
	LinkID string `json:"linkId" jsonschema:"required,description=Jira issue link id to remove"`
	DryRun bool   `json:"dryRun,omitempty" jsonschema:"description=Validate inputs but do not modify Jira"`
}

type RemoveJiraIssueLinkOutput struct {
	LinkID string `json:"linkId" jsonschema:"description=Jira issue link id"`
	Status string `json:"status" jsonschema:"description=success or dry_run"`
}

func RemoveJiraIssueLinkHandler(ctx context.Context, req *mcp.CallToolRequest, input RemoveJiraIssueLinkInput) (*mcp.CallToolResult, RemoveJiraIssueLinkOutput, error) {
	if input.LinkID == "" {
		return nil, RemoveJiraIssueLinkOutput{}, fmt.Errorf("linkId is required")
	}
	if input.DryRun {
		return nil, RemoveJiraIssueLinkOutput{LinkID: input.LinkID, Status: "dry_run"}, nil
	}

	client, err := jira.NewFromEnv()
	if err != nil {
		return nil, RemoveJiraIssueLinkOutput{}, err
	}
	if err := client.RemoveIssueLink(ctx, input.LinkID); err != nil {
		return nil, RemoveJiraIssueLinkOutput{}, err
	}
	return nil, RemoveJiraIssueLinkOutput{LinkID: input.LinkID, Status: "success"}, nil
}

func RegisterRemoveJiraIssueLink(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "removeJiraIssueLink",
		Description: "Remove a Jira issue link by link id.",
	}, RemoveJiraIssueLinkHandler)
}
