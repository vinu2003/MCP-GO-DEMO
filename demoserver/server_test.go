// Test suite for the Demo2 Go MCP server.
//
// Every test connects an in-memory client to the same New() the real
// binary uses. No subprocess, no Claude Desktop, no LLM — so a tool ALWAYS runs
// when we call it, and we assert on the real return value. This is the Go
// counterpart of the FastMCP in-process pytest suite.
//
// Run:  go test -v ./...
package demoserver

import (
	"context"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// connect wires an in-memory client to a fresh server and returns the client
// session. It's the same helper pattern the SDK's own tests use.
func connect(t *testing.T) (context.Context, *mcp.ClientSession) {
	t.Helper()
	ctx := context.Background()

	serverT, clientT := mcp.NewInMemoryTransports()
	if _, err := New().Connect(ctx, serverT, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.0.1"}, nil)
	session, err := client.Connect(ctx, clientT, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })
	return ctx, session
}

func textOf(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		t.Fatal("result has no content")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("content[0] is %T, want *mcp.TextContent", res.Content[0])
	}
	return tc.Text
}

func TestDiscovery_ExpectedToolsRegistered(t *testing.T) {
	ctx, session := connect(t)
	res, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, tool := range res.Tools {
		got[tool.Name] = true
	}
	for _, want := range []string{"add", "current_time"} {
		if !got[want] {
			t.Errorf("tool %q not registered; got %v", want, got)
		}
	}
}

func TestAdd_ReturnsSum(t *testing.T) {
	ctx, session := connect(t)
	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "add",
		Arguments: map[string]any{"a": 2, "b": 3},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := textOf(t, res); got != `{"sum":5}` {
		t.Errorf("add(2,3) text = %s, want {\"sum\":5}", got)
	}
	// StructuredContent carries the typed output as a decoded JSON value.
	if sc, ok := res.StructuredContent.(map[string]any); !ok || sc["sum"] != 5.0 {
		t.Errorf("StructuredContent = %#v, want map with sum=5", res.StructuredContent)
	}
}

func TestAdd_Cases(t *testing.T) {
	cases := []struct {
		a, b, want int
	}{{0, 0, 0}, {100, 200, 300}, {-1, -1, -2}, {-4, 10, 6}}
	ctx, session := connect(t)
	for _, c := range cases {
		res, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name:      "add",
			Arguments: map[string]any{"a": c.a, "b": c.b},
		})
		if err != nil {
			t.Fatal(err)
		}
		sc := res.StructuredContent.(map[string]any)
		if sc["sum"] != float64(c.want) {
			t.Errorf("add(%d,%d) = %v, want %d", c.a, c.b, sc["sum"], c.want)
		}
	}
}

func TestCurrentTime_IsParseableAndRecent(t *testing.T) {
	ctx, session := connect(t)
	before := time.Now()
	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "current_time"})
	if err != nil {
		t.Fatal(err)
	}
	after := time.Now()

	sc := res.StructuredContent.(map[string]any)
	parsed, err := time.Parse(time.RFC3339Nano, sc["time"].(string))
	if err != nil {
		t.Fatalf("current_time returned unparseable time %q: %v", sc["time"], err)
	}
	if parsed.Before(before.Add(-time.Second)) || parsed.After(after.Add(time.Second)) {
		t.Errorf("current_time %v not within [%v, %v]", parsed, before, after)
	}
}

func TestGreetingResource(t *testing.T) {
	ctx, session := connect(t)
	rr, err := session.ReadResource(ctx, &mcp.ReadResourceParams{URI: "greeting://World"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rr.Contents) == 0 {
		t.Fatal("no contents returned")
	}
	if got := rr.Contents[0].Text; got != "Hello, World!" {
		t.Errorf("greeting = %q, want %q", got, "Hello, World!")
	}
}
