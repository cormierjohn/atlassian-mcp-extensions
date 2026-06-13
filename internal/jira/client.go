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
	baseURL := strings.TrimRight(os.Getenv("ATLASSIAN_SITE_URL"), "/")
	email := os.Getenv("ATLASSIAN_EMAIL")
	token := os.Getenv("ATLASSIAN_API_TOKEN")
	if baseURL == "" {
		return nil, errors.New("ATLASSIAN_SITE_URL is required")
	}
	if email == "" {
		return nil, errors.New("ATLASSIAN_EMAIL is required")
	}
	if token == "" {
		return nil, errors.New("ATLASSIAN_API_TOKEN is required")
	}
	return &Client{
		BaseURL:  baseURL,
		Email:    email,
		APIToken: token,
		HTTP:     &http.Client{Timeout: 60 * time.Second},
	}, nil
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
