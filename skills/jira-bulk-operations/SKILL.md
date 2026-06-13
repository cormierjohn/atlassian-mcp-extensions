---
name: jira-bulk-operations
description: "Use safe Jira bulk operations exposed by atlassian-mcp-extensions, including bulk transitions, bulk updates, ranking, changelog retrieval, and issue-link removal."
---

# Jira Bulk Operations

Use this skill when a user needs to update multiple Jira issues or use Jira operations that are not currently exposed by the hosted Atlassian Rovo MCP server.

## Safety rules

1. Confirm the intended issue keys and operation before making bulk changes.
2. Prefer `dryRun: true` first for transitions, updates, ranking, and link removal.
3. Report per-issue results, including failures, instead of hiding partial success.
4. Do not delete Jira issues. This extension only removes issue links by explicit link id.

## Tools

- `bulkTransitionJiraIssues`
- `bulkUpdateJiraIssues`
- `rankJiraIssues`
- `removeJiraIssueLink`
- `getJiraIssueChangelog`
- `transitionJiraIssueByStatus`

Use Atlassian Rovo MCP for standard create/edit/search operations and this extension for missing Jira operations.
