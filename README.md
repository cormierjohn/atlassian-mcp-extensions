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
atlassian-mcp-extensions setup
atlassian-mcp-extensions check-auth
```

bash/zsh:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
atlassian-mcp-extensions setup
atlassian-mcp-extensions check-auth
```

### Option B: Install from a release binary

Download the matching archive from the GitHub Releases page for your operating system:

```text
https://github.com/cormierjohn/atlassian-mcp-extensions/releases
```

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
atlassian-mcp-extensions setup
atlassian-mcp-extensions check-auth
```

Example macOS/Linux:

```bash
mkdir -p ~/.local/bin
tar -xzf atlassian-mcp-extensions-linux-amd64.tar.gz -C ~/.local/bin
chmod +x ~/.local/bin/atlassian-mcp-extensions
export PATH="$HOME/.local/bin:$PATH"
atlassian-mcp-extensions setup
atlassian-mcp-extensions check-auth
```

### Install the Copilot CLI plugin

After the server binary is on `PATH` and `atlassian-mcp-extensions check-auth` succeeds:

Marketplace-style install:

```powershell
copilot plugin marketplace add cormierjohn/atlassian-mcp-extensions
copilot plugin install atlassian-mcp-extensions@atlassian-tools
```

Direct GitHub install, useful for development while marketplace workflows evolve:

```powershell
copilot plugin install cormierjohn/atlassian-mcp-extensions
```

Start a new Copilot CLI session after installing the plugin or changing `PATH`.

## Authentication

If you skipped setup during installation, run:

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

## Plugin manifests

This repository includes manifests for the major local-agent clients:

- Copilot CLI: `.github/plugin/marketplace.json` and `.github/plugin/plugin.json`
- Claude-style plugin layout: `.claude-plugin/marketplace.json` and `.claude-plugin/plugin.json`
- Cursor-style plugin layout: `.cursor-plugin/marketplace.json` and `.cursor-plugin/plugin.json`

All manifests point back to the same root `.mcp.json`, so the server binary still needs to be installed locally and available on `PATH`.

## Development

Run tests:

```powershell
go test ./...
```

Run locally over stdio:

```powershell
go run .
```

## Publishing releases

This repository has a GitHub Actions release workflow that can do two different things:

- **Manual workflow run:** builds binaries and attaches them as artifacts to that workflow run. This is useful for testing the build, but it does not create a GitHub Release page.
- **Version tag push:** builds binaries and publishes them to the GitHub Releases page for users to download.

To publish a real release:

```powershell
git switch main
git pull
git tag v0.1.0
git push origin v0.1.0
```

Use the next semantic version tag for future releases, for example `v0.1.1` or `v0.2.0`.

After the workflow finishes, release downloads are available at:

```text
https://github.com/cormierjohn/atlassian-mcp-extensions/releases
```

## License

MIT
