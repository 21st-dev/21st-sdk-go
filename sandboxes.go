package agents

import (
	"context"
	"fmt"
	"net/url"
)

// SandboxesResource manages sandbox lifecycle and sub-resources.
type SandboxesResource struct {
	client *AgentClient
	Files  *FilesResource
	Git    *GitResource
}

func (s *SandboxesResource) Create(ctx context.Context, req *CreateSandboxRequest) (*Sandbox, error) {
	var out Sandbox
	err := s.client.do(ctx, "POST", "/v1/sandboxes", req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *SandboxesResource) Get(ctx context.Context, sandboxID string) (*SandboxDetail, error) {
	var out SandboxDetail
	err := s.client.do(ctx, "GET", "/v1/sandboxes/"+sandboxID, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *SandboxesResource) Delete(ctx context.Context, sandboxID string) error {
	return s.client.do(ctx, "DELETE", "/v1/sandboxes/"+sandboxID, nil, nil)
}

func (s *SandboxesResource) Exec(ctx context.Context, req *ExecRequest) (*ExecResult, error) {
	var out ExecResult
	err := s.client.do(ctx, "POST", fmt.Sprintf("/v1/sandboxes/%s/exec", req.SandboxID), req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// FilesResource manages files inside a sandbox.
type FilesResource struct {
	client *AgentClient
}

func (f *FilesResource) Write(ctx context.Context, sandboxID string, files map[string]string) error {
	body := struct {
		Files map[string]string `json:"files"`
	}{Files: files}
	return f.client.do(ctx, "POST", fmt.Sprintf("/v1/sandboxes/%s/files", sandboxID), body, nil)
}

func (f *FilesResource) Read(ctx context.Context, sandboxID string, path string) (*FileContent, error) {
	var out FileContent
	encodedPath := url.QueryEscape(path)
	err := f.client.do(ctx, "GET", fmt.Sprintf("/v1/sandboxes/%s/files?path=%s", sandboxID, encodedPath), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GitResource manages git operations inside a sandbox.
type GitResource struct {
	client *AgentClient
}

func (g *GitResource) Clone(ctx context.Context, sandboxID string, req *GitCloneRequest) (*GitCloneResult, error) {
	var out GitCloneResult
	err := g.client.do(ctx, "POST", fmt.Sprintf("/v1/sandboxes/%s/git/clone", sandboxID), req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
