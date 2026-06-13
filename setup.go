package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"golang.org/x/term"
)

func runSetup() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprintln(os.Stderr, "Atlassian MCP Extensions setup")
	fmt.Fprintln(os.Stderr, "Recommended: store credentials in your OS credential store.")
	fmt.Fprintln(os.Stderr, "Alternative: use environment variables for a shell, CI job, or MCP client config.")
	fmt.Fprintln(os.Stderr, "Existing MCP/client sessions must be restarted before they can see newly configured credentials.")
	fmt.Fprintln(os.Stderr)

	choice, err := prompt(reader, "Choose setup mode: [1] OS credential store (recommended), [2] show env var commands only", "1", false)
	if err != nil {
		return err
	}
	if strings.TrimSpace(choice) == "2" {
		printEnvInstructions()
		return nil
	}

	existing, _ := jira.LoadKeyringCredentials()
	siteURL, err := prompt(reader, "Atlassian site URL", existing.SiteURL, false)
	if err != nil {
		return err
	}
	email, err := prompt(reader, "Atlassian email", existing.Email, false)
	if err != nil {
		return err
	}
	apiToken, err := prompt(reader, "Atlassian API token", mask(existing.APIToken), true)
	if err != nil {
		return err
	}
	if apiToken == "" && existing.APIToken != "" {
		apiToken = existing.APIToken
	}

	if err := jira.SaveKeyringCredentials(jira.Credentials{
		SiteURL:  siteURL,
		Email:    email,
		APIToken: apiToken,
	}); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "\nCredentials saved to the OS credential store under service `atlassian-mcp-extensions`.")
	fmt.Fprintln(os.Stderr, "Restart any running MCP clients or agent sessions before using the updated credentials.")
	return nil
}

func prompt(reader *bufio.Reader, label, current string, secret bool) (string, error) {
	if current != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, current)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}
	var value string
	var err error
	if secret {
		value, err = readSecret()
	} else {
		value, err = reader.ReadString('\n')
	}
	if err != nil {
		return "", err
	}
	value = strings.TrimSpace(value)
	if value == "" && current != "" && !secret {
		return current, nil
	}
	return value, nil
}

func readSecret() (string, error) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		body, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}

func mask(value string) string {
	if value == "" {
		return ""
	}
	return "configured; press Enter to keep"
}

func printEnvInstructions() {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Set these variables in the environment that launches the MCP server:")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "PowerShell session only:")
	fmt.Fprintln(os.Stderr, `  $env:ATLASSIAN_SITE_URL = "https://your-site.atlassian.net"`)
	fmt.Fprintln(os.Stderr, `  $env:ATLASSIAN_EMAIL = "you@example.com"`)
	fmt.Fprintln(os.Stderr, `  $env:ATLASSIAN_API_TOKEN = "your-api-token"`)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "bash/zsh session only:")
	fmt.Fprintln(os.Stderr, `  export ATLASSIAN_SITE_URL="https://your-site.atlassian.net"`)
	fmt.Fprintln(os.Stderr, `  export ATLASSIAN_EMAIL="you@example.com"`)
	fmt.Fprintln(os.Stderr, `  export ATLASSIAN_API_TOKEN="your-api-token"`)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Environment variables are process-scoped unless you explicitly persist them.")
	fmt.Fprintln(os.Stderr, "If your MCP client is already running, restart it after setting variables.")
	fmt.Fprintln(os.Stderr, "Avoid globally persisting API tokens unless you understand the leakage risk.")
}
