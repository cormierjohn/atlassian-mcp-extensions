package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	BaseURL  string
	Email    string
	APIToken string
	HTTP     *http.Client
}

type HTTPError struct {
	Method string
	URL    string
	Status int
	Body   string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s %s failed with HTTP %d: %s", e.Method, e.URL, e.Status, e.Body)
}

func NewFromEnv() (*Client, error) {
	creds, err := LoadCredentials()
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL:  creds.SiteURL,
		Email:    creds.Email,
		APIToken: creds.APIToken,
		HTTP:     &http.Client{Timeout: 60 * time.Second},
	}, nil
}

type Credentials struct {
	SiteURL  string `json:"siteUrl"`
	Email    string `json:"email"`
	APIToken string `json:"apiToken"`
}

func LoadCredentials() (Credentials, error) {
	creds := Credentials{
		SiteURL:  strings.TrimRight(os.Getenv("ATLASSIAN_SITE_URL"), "/"),
		Email:    os.Getenv("ATLASSIAN_EMAIL"),
		APIToken: os.Getenv("ATLASSIAN_API_TOKEN"),
	}
	if creds.SiteURL == "" || creds.Email == "" || creds.APIToken == "" {
		keyringCreds, err := LoadKeyringCredentials()
		if err == nil {
			if creds.SiteURL == "" {
				creds.SiteURL = keyringCreds.SiteURL
			}
			if creds.Email == "" {
				creds.Email = keyringCreds.Email
			}
			if creds.APIToken == "" {
				creds.APIToken = keyringCreds.APIToken
			}
		}
	}
	if creds.SiteURL == "" {
		return Credentials{}, errors.New("ATLASSIAN_SITE_URL is required; run `atlassian-mcp-extensions setup`")
	}
	if creds.Email == "" {
		return Credentials{}, errors.New("ATLASSIAN_EMAIL is required; run `atlassian-mcp-extensions setup`")
	}
	if creds.APIToken == "" {
		return Credentials{}, errors.New("ATLASSIAN_API_TOKEN is required; run `atlassian-mcp-extensions setup`")
	}
	return creds, nil
}

func (c *Client) Do(ctx context.Context, method, path string, query url.Values, payload any, out any) error {
	body, err := encodePayload(payload)
	if err != nil {
		return err
	}

	reqURL := c.BaseURL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	resp, respBody, err := c.doOnce(ctx, method, reqURL, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		wait := retryAfter(resp.Header.Get("Retry-After"))
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
		resp, respBody, err = c.doOnce(ctx, method, reqURL, body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode >= 400 {
		return &HTTPError{
			Method: method,
			URL:    reqURL,
			Status: resp.StatusCode,
			Body:   truncate(string(respBody), 2000),
		}
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode Jira response: %w", err)
	}
	return nil
}

func (c *Client) doOnce(ctx context.Context, method, reqURL string, body []byte) (*http.Response, []byte, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reader)
	if err != nil {
		return nil, nil, err
	}
	req.SetBasicAuth(c.Email, c.APIToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, nil, err
	}
	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		resp.Body.Close()
		return nil, nil, readErr
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	return resp, respBody, nil
}

func encodePayload(payload any) ([]byte, error) {
	if payload == nil {
		return nil, nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode Jira payload: %w", err)
	}
	return body, nil
}

func retryAfter(value string) time.Duration {
	seconds, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || seconds < 0 {
		seconds = 10
	}
	return time.Duration(seconds) * time.Second
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max] + fmt.Sprintf("... [truncated; full body %d chars]", len(value))
}
