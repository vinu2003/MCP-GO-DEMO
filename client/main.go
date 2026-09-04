// Command client is an in-process MCP client demo.
//
// It builds the same server object and connects to it over an in-memory
// transport — no subprocess, no stdio pipes, no Claude Desktop, no LLM. This is
// the Go equivalent of FastMCP's `Client(mcp)`: YOU are the client, and you
// call tools directly and deterministically.
//
// Run:  go run ./client
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/vinu/mcp-go-demo/demoserver"
)

func main() {
	ctx := context.Background()

	// NewInMemoryTransports returns a connected pair: one end for the server,
	// one for the client. This is the in-memory transport the SDK's own tests use.
	serverT, clientT := mcp.NewInMemoryTransports()

	server := demoserver.New()
	if _, err := server.Connect(ctx, serverT, nil); err != nil {
		log.Fatal(err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "demo-client", Version: "v0.1.0"}, nil)
	session, err := client.Connect(ctx, clientT, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// 1. Discover tools.
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Tools available:")
	for _, t := range tools.Tools {
		fmt.Printf(" %s", t.Name)
	}
	fmt.Println()

	// 2. Call add. The structured output arrives in StructuredContent; the text
	//    block in Content mirrors it.
	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "add",
		Arguments: map[string]any{"a": 2, "b": 3},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("add(2, 3) -> %s\n", textOf(res))

	// 3. Call current_time.
	now, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "current_time"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("current_time() -> %s\n", textOf(now))

	// 4. Read the parameterised resource.
	rr, err := session.ReadResource(ctx, &mcp.ReadResourceParams{URI: "greeting://Vinu"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("greeting://Vinu -> %s\n", rr.Contents[0].Text)
}

// textOf returns the first text content block of a tool result.
func textOf(res *mcp.CallToolResult) string {
	if len(res.Content) == 0 {
		return ""
	}
	if tc, ok := res.Content[0].(*mcp.TextContent); ok {
		return tc.Text
	}
	return ""
}
