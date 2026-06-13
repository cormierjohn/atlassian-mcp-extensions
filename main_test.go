package main

import (
	"context"
	"testing"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/cormierjohn/atlassian-mcp-extensions/tools"
)

func TestSelectTransitionByTargetStatus(t *testing.T) {
	transitions := []jira.Transition{
		{ID: "11", Name: "Start Progress"},
		{ID: "21", Name: "Done"},
	}
	transitions[0].To.Name = "In Progress"
	transitions[1].To.Name = "Done"

	got, err := jira.SelectTransition(transitions, "done", "")
	if err != nil {
		t.Fatalf("SelectTransition failed: %v", err)
	}
	if got.ID != "21" {
		t.Fatalf("expected transition 21, got %s", got.ID)
	}
}

func TestRankJiraIssuesValidation(t *testing.T) {
	_, _, err := tools.RankJiraIssuesHandler(context.Background(), nil, tools.RankJiraIssuesInput{
		IssueKeys:       []string{"PROJ-2"},
		RankBeforeIssue: "PROJ-1",
		RankAfterIssue:  "PROJ-3",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestTextToADF(t *testing.T) {
	adf := jira.TextToADF("hello")
	if adf["type"] != "doc" {
		t.Fatalf("expected doc node, got %#v", adf["type"])
	}
}
