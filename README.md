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

Run setup:

```powershell
atlassian-mcp-extensions setup
```

Setup offers two modes:

1. Store credentials in your OS credential store (recommended).
2. Print environment variable instructions only.

Credential-store mode prompts for:

- Atlassian site URL
- Atlassian email
- Atlassian API token

Credentials are stored using the operating system credential store: Windows Credential Manager, macOS Keychain, or Linux Secret Service. They are not written to this repository, `.mcp.json`, or a project config file.

Existing MCP clients and agent sessions must be restarted after setup because already-running processes cannot see newly configured credentials.

For headless environments, you can use environment variables instead. Environment variables are checked before the OS credential store so they can act as explicit process-level overrides:

```powershell
$env:ATLASSIAN_SITE_URL = "https://your-site.atlassian.net"
$env:ATLASSIAN_EMAIL = "you@example.com"
$env:ATLASSIAN_API_TOKEN = "your-api-token"
```

Environment variables must be present in the process that launches the MCP server. Setting them in a terminal does not affect already-running MCP clients. Do not commit tokens to this repository or place them directly in `.mcp.json`. OAuth PKCE is a planned auth mode.

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
