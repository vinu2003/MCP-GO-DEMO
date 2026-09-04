// Package demoserver builds the Demo2 MCP server. It's a normal importable
// package (not package main), so both the server binary and the client demo can
// construct the same server.
package demoserver

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AddInput struct {
	A int `json:"a" jsonschema:"the first addend"`
	B int `json:"b" jsonschema:"the second addend"`
}

type AddOutput struct {
	Sum int `json:"sum" jsonschema:"the sum of a and b"`
}

func Add(_ context.Context, _ *mcp.CallToolRequest, in AddInput) (*mcp.CallToolResult, AddOutput, error) {
	return nil, AddOutput{Sum: in.A + in.B}, nil
}

type CurrentTimeInput struct{}

type CurrentTimeOutput struct {
	Time string `json:"time" jsonschema:"the current server time in ISO-8601"`
}

func CurrentTime(_ context.Context, _ *mcp.CallToolRequest, _ CurrentTimeInput) (*mcp.CallToolResult, CurrentTimeOutput, error) {
	return nil, CurrentTimeOutput{Time: time.Now().Format(time.RFC3339Nano)}, nil
}

// New builds the server and registers every capability.
func New() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "Demo2", Version: "v0.1.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "Add two integers and return the sum.",
	}, Add)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "current_time",
		Description: "Return the current server time as an ISO-8601 string.",
	}, CurrentTime)

	server.AddResourceTemplate(&mcp.ResourceTemplate{
		Name:        "greeting",
		URITemplate: "greeting://{name}",
		Description: "A greeting for the given name.",
		MIMEType:    "text/plain",
	}, handleGreeting)

	return server
}

func handleGreeting(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := req.Params.URI
	name := uri[len("greeting://"):]
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{URI: uri, MIMEType: "text/plain", Text: "Hello, " + name + "!"},
		},
	}, nil
}
