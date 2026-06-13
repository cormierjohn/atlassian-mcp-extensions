# Atlassian MCP Extensions

Community extensions for Atlassian MCP. This project is designed to be installed beside the hosted Atlassian Rovo MCP server and expose focused Jira tools that are useful for automation but not currently available from the hosted tool set.

## What this provides

- `bulkTransitionJiraIssues`
- `bulkUpdateJiraIssues`
- `rankJiraIssues`
- `removeJiraIssueLink`
- `getJiraIssueChangelog`
- `transitionJiraIssueByStatus`

The plugin config also points to the hosted Atlassian Rovo MCP server so clients can use both tool sets together.

## Authentication

The first version uses Jira API token authentication through environment variables:

```powershell
$env:ATLASSIAN_SITE_URL = "https://your-site.atlassian.net"
$env:ATLASSIAN_EMAIL = "you@example.com"
$env:ATLASSIAN_API_TOKEN = "your-api-token"
```

Do not commit tokens to this repository or place them directly in `.mcp.json`. OAuth PKCE with OS keyring storage is a planned auth mode.

## Build

```powershell
go mod download
go test ./...
go build -o atlassian-mcp-extensions.exe
```

## MCP configuration

`.mcp.json` exposes two servers:

```json
{
  "mcpServers": {
    "atlassian-rovo": {
      "type": "http",
      "url": "https://mcp.atlassian.com/v1/mcp"
    },
    "atlassian-extensions": {
      "command": "atlassian-mcp-extensions"
    }
  }
}
```

## Development

Run tests:

```powershell
go test ./...
```

Run locally over stdio:

```powershell
go run .
```

## License

MIT
