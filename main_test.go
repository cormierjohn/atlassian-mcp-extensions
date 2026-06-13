package main

import (
	"context"
	"testing"

	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/cormierjohn/atlassian-mcp-extensions/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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

func TestMCPToolsList(t *testing.T) {
	ctx := context.Background()
	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "test"}, &mcp.ServerOptions{HasTools: true})
	tools.RegisterTools(server)
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect failed: %v", err)
	}
	defer serverSession.Close()

	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}
	defer clientSession.Close()

	result, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	want := map[string]bool{
		"bulkTransitionJiraIssues":    false,
		"bulkUpdateJiraIssues":        false,
		"rankJiraIssues":              false,
		"removeJiraIssueLink":         false,
		"getJiraIssueChangelog":       false,
		"transitionJiraIssueByStatus": false,
	}
	for _, tool := range result.Tools {
		if _, ok := want[tool.Name]; ok {
			want[tool.Name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Fatalf("tool %s not listed", name)
		}
	}
}
