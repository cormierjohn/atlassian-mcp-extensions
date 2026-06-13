package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cormierjohn/atlassian-mcp-extensions/config"
	"github.com/cormierjohn/atlassian-mcp-extensions/internal/jira"
	"github.com/cormierjohn/atlassian-mcp-extensions/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "setup":
			if err := runSetup(); err != nil {
				slog.Error("setup failed", "error", err)
				os.Exit(1)
			}
			return
		case "check-auth":
			if _, err := jira.LoadCredentials(); err != nil {
				slog.Error("authentication is not configured", "error", err)
				os.Exit(2)
			}
			slog.Info("authentication is configured")
			return
		}
	}

	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	if _, err := jira.LoadCredentials(); err != nil {
		slog.Error(
			"authentication is not configured; run `atlassian-mcp-extensions setup` and restart your MCP client",
			"error", err,
		)
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		slog.Info("shutting down")
		cancel()
	}()

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    cfg.ServerName,
			Version: cfg.Version,
		},
		&mcp.ServerOptions{HasTools: true},
	)

	tools.RegisterTools(server)

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
