package tools

import (
	"context"
	"fmt"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RankJiraIssuesInput struct {
	IssueKeys       []string `json:"issueKeys" jsonschema:"required,description=Jira issues to rank as a group"`
	RankBeforeIssue string   `json:"rankBeforeIssue,omitempty" jsonschema:"description=Place issues immediately before this issue"`
	RankAfterIssue  string   `json:"rankAfterIssue,omitempty" jsonschema:"description=Place issues immediately after this issue"`
	DryRun          bool     `json:"dryRun,omitempty" jsonschema:"description=Validate inputs but do not modify Jira"`
}

type RankJiraIssuesOutput struct {
	Status string   `json:"status" jsonschema:"description=success or dry_run"`
	Issues []string `json:"issues" jsonschema:"description=Ranked Jira issue keys"`
}

func RankJiraIssuesHandler(ctx context.Context, req *mcp.CallToolRequest, input RankJiraIssuesInput) (*mcp.CallToolResult, RankJiraIssuesOutput, error) {
	if len(input.IssueKeys) == 0 {
		return nil, RankJiraIssuesOutput{}, fmt.Errorf("issueKeys is required")
	}
	if (input.RankBeforeIssue == "") == (input.RankAfterIssue == "") {
		return nil, RankJiraIssuesOutput{}, fmt.Errorf("provide exactly one of rankBeforeIssue or rankAfterIssue")
	}
	if input.DryRun {
		return nil, RankJiraIssuesOutput{Status: "dry_run", Issues: input.IssueKeys}, nil
	}

	client, err := jira.NewFromEnv()
	if err != nil {
		return nil, RankJiraIssuesOutput{}, err
	}
	if err := client.RankIssues(ctx, input.IssueKeys, input.RankBeforeIssue, input.RankAfterIssue); err != nil {
		return nil, RankJiraIssuesOutput{}, err
	}
	return nil, RankJiraIssuesOutput{Status: "success", Issues: input.IssueKeys}, nil
}

func RegisterRankJiraIssues(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "rankJiraIssues",
		Description: "Rank Jira issues before or after another issue using Jira Software agile ranking.",
	}, RankJiraIssuesHandler)
}
