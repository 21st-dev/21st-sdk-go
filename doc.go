// Package agents provides a Go client for the 21st relay API.
//
// It manages sandboxes, threads, tokens, file I/O, command execution,
// and git operations. This is the Go equivalent of @21st-sdk/node.
//
// Usage:
//
//	client := agents.NewAgentClient("your-api-key")
//	sandbox, err := client.Sandboxes.Create(ctx, &agents.CreateSandboxRequest{
//	    Agent: "my-agent",
//	})
package agents
