// Command mcp-go-demo runs the Demo2 MCP server over stdio.
//
// Run over stdio:  go run .
// This is what Claude Desktop launches.
package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/vinu/mcp-go-demo/demoserver"
)

func main() {
	if err := demoserver.New().Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
