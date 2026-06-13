package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

func RegisterTools(server *mcp.Server) {
	RegisterBulkTransitionJiraIssues(server)
	RegisterBulkUpdateJiraIssues(server)
	RegisterRankJiraIssues(server)
	RegisterRemoveJiraIssueLink(server)
	RegisterGetJiraIssueChangelog(server)
	RegisterTransitionJiraIssueByStatus(server)
}
