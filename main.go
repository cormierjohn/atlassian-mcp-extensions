package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cormierjohn/atlassian-mcp-extensions/config"
	"github.com/cormierjohn/atlassian-mcp-extensions/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

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
