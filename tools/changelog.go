package tools

import (
	"context"
	"fmt"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetJiraIssueChangelogInput struct {
	IssueKey     string `json:"issueKey" jsonschema:"Jira issue key"`
	PageSize     int    `json:"pageSize,omitempty" jsonschema:"Jira page size, default 100"`
	MaxHistories int    `json:"maxHistories,omitempty" jsonschema:"Maximum histories to return, default 500"`
}

type GetJiraIssueChangelogOutput struct {
	IssueKey string           `json:"issueKey" jsonschema:"Jira issue key"`
	Count    int              `json:"count" jsonschema:"Returned changelog history count"`
	Values   []map[string]any `json:"values" jsonschema:"Changelog histories"`
}

func GetJiraIssueChangelogHandler(ctx context.Context, req *mcp.CallToolRequest, input GetJiraIssueChangelogInput) (*mcp.CallToolResult, GetJiraIssueChangelogOutput, error) {
	if input.IssueKey == "" {
		return nil, GetJiraIssueChangelogOutput{}, fmt.Errorf("issueKey is required")
	}
	client, err := jira.NewFromEnv()
	if err != nil {
		return nil, GetJiraIssueChangelogOutput{}, err
	}
	values, err := client.GetChangelog(ctx, input.IssueKey, input.PageSize, input.MaxHistories)
	if err != nil {
		return nil, GetJiraIssueChangelogOutput{}, err
	}
	return nil, GetJiraIssueChangelogOutput{
		IssueKey: input.IssueKey,
		Count:    len(values),
		Values:   values,
	}, nil
}

func RegisterGetJiraIssueChangelog(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "getJiraIssueChangelog",
		Description: "Fetch paginated Jira issue changelog histories.",
	}, GetJiraIssueChangelogHandler)
}
