package agents

import (
	"context"
	"fmt"
	"net/url"
)

// ThreadsResource manages threads within sandboxes.
type ThreadsResource struct {
	client *AgentClient
}

func (t *ThreadsResource) List(ctx context.Context, sandboxID string) ([]ThreadSummary, error) {
	var out []ThreadSummary
	err := t.client.do(ctx, "GET", fmt.Sprintf("/v1/sandboxes/%s/threads", sandboxID), nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ThreadsResource) Create(ctx context.Context, sandboxID string, name string) (*ThreadSummary, error) {
	body := struct {
		Name string `json:"name,omitempty"`
	}{Name: name}
	var out ThreadSummary
	err := t.client.do(ctx, "POST", fmt.Sprintf("/v1/sandboxes/%s/threads", sandboxID), body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (t *ThreadsResource) Get(ctx context.Context, sandboxID, threadID string) (*Thread, error) {
	var out Thread
	err := t.client.do(ctx, "GET", fmt.Sprintf("/v1/sandboxes/%s/threads/%s", sandboxID, threadID), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (t *ThreadsResource) Delete(ctx context.Context, sandboxID, threadID string) error {
	return t.client.do(ctx, "DELETE", fmt.Sprintf("/v1/sandboxes/%s/threads/%s", sandboxID, threadID), nil, nil)
}

// Run sends messages to an agent. It auto-creates a sandbox and thread if not provided.
// The returned RunThreadResult contains a streaming *http.Response — caller must close it.
func (t *ThreadsResource) Run(ctx context.Context, req *RunThreadRequest) (*RunThreadResult, error) {
	if req.ThreadID != "" && req.SandboxID == "" {
		return nil, fmt.Errorf("agents: threadId requires sandboxId")
	}

	sandboxID := req.SandboxID
	if sandboxID == "" {
		sandbox, err := t.client.Sandboxes.Create(ctx, &CreateSandboxRequest{Agent: req.Agent})
		if err != nil {
			return nil, fmt.Errorf("agents: auto-create sandbox: %w", err)
		}
		sandboxID = sandbox.ID
	}

	threadID := req.ThreadID
	if threadID == "" {
		thread, err := t.Create(ctx, sandboxID, req.Name)
		if err != nil {
			return nil, fmt.Errorf("agents: auto-create thread: %w", err)
		}
		threadID = thread.ID
	}

	encodedAgent := url.PathEscape(req.Agent)
	body := struct {
		Messages  []RunThreadMessage `json:"messages"`
		SandboxID string             `json:"sandboxId"`
		ThreadID  string             `json:"threadId"`
	}{
		Messages:  req.Messages,
		SandboxID: sandboxID,
		ThreadID:  threadID,
	}

	resp, err := t.client.doRequest(ctx, "POST", "/v1/chat/"+encodedAgent, body)
	if err != nil {
		return nil, err
	}

	return &RunThreadResult{
		SandboxID: sandboxID,
		ThreadID:  threadID,
		Response:  resp,
		ResumeURL: fmt.Sprintf("%s/v1/chat/%s/%s/stream", t.client.baseURL, encodedAgent, sandboxID),
	}, nil
}
