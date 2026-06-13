package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	To   struct {
		Name string `json:"name"`
	} `json:"to"`
}

type ChangelogPage struct {
	StartAt    int              `json:"startAt"`
	MaxResults int              `json:"maxResults"`
	Total      int              `json:"total"`
	IsLast     bool             `json:"isLast"`
	Values     []map[string]any `json:"values"`
}

func (c *Client) GetTransitions(ctx context.Context, issueKey string) ([]Transition, error) {
	var response struct {
		Transitions []Transition `json:"transitions"`
	}
	path := fmt.Sprintf("/rest/api/3/issue/%s/transitions", url.PathEscape(issueKey))
	if err := c.Do(ctx, http.MethodGet, path, nil, nil, &response); err != nil {
		return nil, err
	}
	return response.Transitions, nil
}

func SelectTransition(transitions []Transition, targetStatus, transitionName string) (*Transition, error) {
	if (targetStatus == "") == (transitionName == "") {
		return nil, fmt.Errorf("provide exactly one of targetStatus or transitionName")
	}
	if targetStatus != "" {
		for i := range transitions {
			if strings.EqualFold(transitions[i].To.Name, targetStatus) {
				return &transitions[i], nil
			}
		}
		return nil, fmt.Errorf("no transition to target status %q; available: %s", targetStatus, transitionSummary(transitions))
	}
	for i := range transitions {
		if strings.EqualFold(transitions[i].Name, transitionName) {
			return &transitions[i], nil
		}
	}
	return nil, fmt.Errorf("no transition named %q; available: %s", transitionName, transitionSummary(transitions))
}

func (c *Client) TransitionIssue(ctx context.Context, issueKey string, transition Transition, comment string) error {
	body := map[string]any{
		"transition": map[string]string{"id": transition.ID},
	}
	if strings.TrimSpace(comment) != "" {
		body["update"] = map[string]any{
			"comment": []map[string]any{{
				"add": map[string]any{"body": TextToADF(comment)},
			}},
		}
	}
	path := fmt.Sprintf("/rest/api/3/issue/%s/transitions", url.PathEscape(issueKey))
	return c.Do(ctx, http.MethodPost, path, nil, body, nil)
}

func (c *Client) UpdateIssue(ctx context.Context, issueKey string, fields, update map[string]any, suppressNotifications bool) error {
	body := map[string]any{}
	if len(fields) > 0 {
		body["fields"] = fields
	}
	if len(update) > 0 {
		body["update"] = update
	}
	query := url.Values{}
	if suppressNotifications {
		query.Set("notifyUsers", "false")
	}
	path := fmt.Sprintf("/rest/api/3/issue/%s", url.PathEscape(issueKey))
	return c.Do(ctx, http.MethodPut, path, query, body, nil)
}

func (c *Client) RankIssues(ctx context.Context, issueKeys []string, before, after string) error {
	body := map[string]any{"issues": issueKeys}
	if before != "" {
		body["rankBeforeIssue"] = before
	} else {
		body["rankAfterIssue"] = after
	}
	return c.Do(ctx, http.MethodPut, "/rest/agile/1.0/issue/rank", nil, body, nil)
}

func (c *Client) RemoveIssueLink(ctx context.Context, linkID string) error {
	path := fmt.Sprintf("/rest/api/3/issueLink/%s", url.PathEscape(linkID))
	return c.Do(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (c *Client) GetChangelog(ctx context.Context, issueKey string, pageSize, maxHistories int) ([]map[string]any, error) {
	if pageSize <= 0 {
		pageSize = 100
	}
	if maxHistories <= 0 {
		maxHistories = 500
	}

	histories := []map[string]any{}
	startAt := 0
	for len(histories) < maxHistories {
		query := url.Values{}
		query.Set("startAt", fmt.Sprintf("%d", startAt))
		query.Set("maxResults", fmt.Sprintf("%d", pageSize))

		var page ChangelogPage
		path := fmt.Sprintf("/rest/api/3/issue/%s/changelog", url.PathEscape(issueKey))
		if err := c.Do(ctx, http.MethodGet, path, query, nil, &page); err != nil {
			return nil, err
		}

		if len(page.Values) == 0 {
			break
		}
		for _, history := range page.Values {
			if len(histories) >= maxHistories {
				break
			}
			histories = append(histories, history)
		}
		if page.IsLast {
			break
		}
		startAt += len(page.Values)
	}
	return histories, nil
}

func TextToADF(text string) map[string]any {
	return map[string]any{
		"type":    "doc",
		"version": 1,
		"content": []map[string]any{{
			"type": "paragraph",
			"content": []map[string]string{{
				"type": "text",
				"text": text,
			}},
		}},
	}
}

func transitionSummary(transitions []Transition) string {
	parts := make([]string, 0, len(transitions))
	for _, transition := range transitions {
		parts = append(parts, fmt.Sprintf("%s -> %s", transition.Name, transition.To.Name))
	}
	return strings.Join(parts, ", ")
}
