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

## Install

You need both pieces:

1. The **server binary** (`atlassian-mcp-extensions`) on your `PATH`.
2. The **plugin** installed in your MCP client so it registers the Atlassian Rovo MCP server and this extension MCP server.

### Option A: Install with Go

If you have Go installed:

```powershell
go install github.com/cormierjohn/atlassian-mcp-extensions@latest
```

Make sure Go's bin directory is on your `PATH`.

PowerShell:

```powershell
$goBin = "$(go env GOPATH)\bin"
$env:Path = "$goBin;$env:Path"
atlassian-mcp-extensions check-auth
```

bash/zsh:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
atlassian-mcp-extensions check-auth
```

### Option B: Install from a release binary

Download the matching archive from the GitHub Releases page for your operating system:

- Windows: `atlassian-mcp-extensions-windows-amd64.zip`
- macOS Apple Silicon: `atlassian-mcp-extensions-darwin-arm64.tar.gz`
- macOS Intel: `atlassian-mcp-extensions-darwin-amd64.tar.gz`
- Linux: `atlassian-mcp-extensions-linux-amd64.tar.gz`

Extract it and place the binary somewhere on your `PATH`.

Example Windows PowerShell:

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\bin" | Out-Null
Expand-Archive .\atlassian-mcp-extensions-windows-amd64.zip -DestinationPath "$env:USERPROFILE\bin" -Force
$env:Path = "$env:USERPROFILE\bin;$env:Path"
atlassian-mcp-extensions check-auth
```

Example macOS/Linux:

```bash
mkdir -p ~/.local/bin
tar -xzf atlassian-mcp-extensions-linux-amd64.tar.gz -C ~/.local/bin
chmod +x ~/.local/bin/atlassian-mcp-extensions
export PATH="$HOME/.local/bin:$PATH"
atlassian-mcp-extensions check-auth
```

### Install the Copilot CLI plugin

After the server binary is on `PATH`:

```powershell
copilot plugin install cormierjohn/atlassian-mcp-extensions
```

Restart any running Copilot CLI session after installing the plugin or changing `PATH`.

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

## Build from source

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
