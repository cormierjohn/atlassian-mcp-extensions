package jira

import (
	"fmt"
	"strings"

	"github.com/zalando/go-keyring"
)

const keyringService = "atlassian-mcp-extensions"

func LoadKeyringCredentials() (Credentials, error) {
	siteURL, siteErr := keyring.Get(keyringService, "site-url")
	email, emailErr := keyring.Get(keyringService, "email")
	apiToken, tokenErr := keyring.Get(keyringService, "api-token")
	if siteErr != nil || emailErr != nil || tokenErr != nil {
		return Credentials{}, fmt.Errorf(
			"credentials not found in OS credential store; run `atlassian-mcp-extensions setup` or set ATLASSIAN_SITE_URL, ATLASSIAN_EMAIL, and ATLASSIAN_API_TOKEN",
		)
	}
	return Credentials{
		SiteURL:  strings.TrimRight(siteURL, "/"),
		Email:    email,
		APIToken: apiToken,
	}, nil
}

func SaveKeyringCredentials(creds Credentials) error {
	creds.SiteURL = strings.TrimRight(strings.TrimSpace(creds.SiteURL), "/")
	creds.Email = strings.TrimSpace(creds.Email)
	creds.APIToken = strings.TrimSpace(creds.APIToken)

	if creds.SiteURL == "" {
		return fmt.Errorf("site URL is required")
	}
	if creds.Email == "" {
		return fmt.Errorf("email is required")
	}
	if creds.APIToken == "" {
		return fmt.Errorf("API token is required")
	}

	if err := keyring.Set(keyringService, "site-url", creds.SiteURL); err != nil {
		return fmt.Errorf("store site URL in OS credential store: %w", err)
	}
	if err := keyring.Set(keyringService, "email", creds.Email); err != nil {
		return fmt.Errorf("store email in OS credential store: %w", err)
	}
	if err := keyring.Set(keyringService, "api-token", creds.APIToken); err != nil {
		return fmt.Errorf("store API token in OS credential store: %w", err)
	}
	return nil
}
