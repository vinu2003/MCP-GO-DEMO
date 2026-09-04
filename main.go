package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "minimal-mcp-server"
	serverVersion = "v0.1.0"
)

type PingOutput struct {
	Message string `json:"message" jsonschema:"response from the ping tool"`
}

func pingTool(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, PingOutput, error) {
	return nil, PingOutput{Message: "pong"}, nil
}

func newServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: serverName, Version: serverVersion}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "ping", Description: "Returns a pong response"}, pingTool)
	return server
}

func main() {
	if err := newServer().Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
