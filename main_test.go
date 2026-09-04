package main

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNewServer(t *testing.T) {
	if newServer() == nil {
		t.Fatal("newServer returned nil")
	}
}

func TestPingTool(t *testing.T) {
	_, output, err := pingTool(context.Background(), &mcp.CallToolRequest{}, struct{}{})
	if err != nil {
		t.Fatalf("pingTool returned error: %v", err)
	}
	if output.Message != "pong" {
		t.Fatalf("unexpected ping response: %q", output.Message)
	}
}
