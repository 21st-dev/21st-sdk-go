# 21st SDK for Go

Go SDK for [21st Agents](https://21st.dev/agents). Manage sandboxes, threads, and tokens programmatically.

## Install

```bash
go get github.com/21st-dev/21st-sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	agents "github.com/21st-dev/21st-sdk-go"
)

func main() {
	client := agents.NewAgentClient("an_sk_...") // your API key

	ctx := context.Background()

	// Create a sandbox for your agent
	sandbox, err := client.Sandboxes.Create(ctx, &agents.CreateSandboxRequest{
		Agent: "my-agent",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a thread
	thread, err := client.Threads.Create(ctx, sandbox.ID, "Review PR #42")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Thread:", thread.ID)

	// Generate a short-lived token for browser clients
	token, err := client.Tokens.Create(ctx, &agents.CreateTokenRequest{
		Agent:     "my-agent",
		ExpiresIn: "1h",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token expires at:", token.ExpiresAt)
}
```

## API

### `agents.NewAgentClient(apiKey, ...opts)`

```go
client := agents.NewAgentClient(
	"an_sk_...",                                       // API key (required)
	agents.WithBaseURL("https://custom-relay.example"), // optional
	agents.WithHTTPClient(customHTTPClient),            // optional
)
```

### `client.Sandboxes`

| Method | Description |
|--------|-------------|
| `Create(ctx, *CreateSandboxRequest)` | Create a new sandbox for an agent |
| `Get(ctx, sandboxID)` | Get sandbox details (status, threads, agent info) |
| `Delete(ctx, sandboxID)` | Delete a sandbox |
| `Exec(ctx, *ExecRequest)` | Run a command in a sandbox |

### `client.Sandboxes.Files`

| Method | Description |
|--------|-------------|
| `Write(ctx, sandboxID, files)` | Write files to a sandbox |
| `Read(ctx, sandboxID, path)` | Read a file from a sandbox |

### `client.Sandboxes.Git`

| Method | Description |
|--------|-------------|
| `Clone(ctx, sandboxID, *GitCloneRequest)` | Clone a repository into a sandbox |

### `client.Threads`

| Method | Description |
|--------|-------------|
| `List(ctx, sandboxID)` | List all threads in a sandbox |
| `Create(ctx, sandboxID, name)` | Create a new thread |
| `Get(ctx, sandboxID, threadID)` | Get thread with messages |
| `Delete(ctx, sandboxID, threadID)` | Delete a thread |
| `Run(ctx, *RunThreadRequest)` | Send messages to an agent (streaming) |

### `client.Tokens`

| Method | Description |
|--------|-------------|
| `Create(ctx, *CreateTokenRequest)` | Create a short-lived JWT (default: 1h) |

### Error Handling

```go
detail, err := client.Sandboxes.Get(ctx, "nonexistent")
if agents.IsAPIError(err, 404) {
	fmt.Println("sandbox not found")
}
```

## License

MIT
