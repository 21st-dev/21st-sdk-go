package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const DefaultBaseURL = "https://relay.an.dev"

// Option configures an AgentClient.
type Option func(*AgentClient)

// WithBaseURL overrides the default relay URL.
func WithBaseURL(url string) Option {
	return func(c *AgentClient) { c.baseURL = url }
}

// WithHTTPClient sets a custom *http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *AgentClient) { c.http = hc }
}

// AgentClient is the main entry point for the 21st relay API.
type AgentClient struct {
	apiKey  string
	baseURL string
	http    *http.Client

	Sandboxes *SandboxesResource
	Threads   *ThreadsResource
	Tokens    *TokensResource
}

// NewAgentClient creates a new client with the given API key.
func NewAgentClient(apiKey string, opts ...Option) *AgentClient {
	c := &AgentClient{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		http:    http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.baseURL = strings.TrimRight(c.baseURL, "/")

	c.Sandboxes = &SandboxesResource{client: c}
	c.Sandboxes.Files = &FilesResource{client: c}
	c.Sandboxes.Git = &GitResource{client: c}
	c.Threads = &ThreadsResource{client: c}
	c.Tokens = &TokensResource{client: c}

	return c
}

// APIError represents an error response from the relay API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("request failed: %d", e.StatusCode)
}

// IsAPIError checks if err is an *APIError, optionally matching a status code.
func IsAPIError(err error, statusCode ...int) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if len(statusCode) > 0 {
		return apiErr.StatusCode == statusCode[0]
	}
	return true
}

// doRequest executes an HTTP request with auth and error handling.
// Returns the raw response on success. Caller must close the body.
func (c *AgentClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("agents: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("agents: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agents: %s %s: %w", method, path, err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		var envelope struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&envelope)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Code:       envelope.Error.Code,
			Message:    envelope.Error.Message,
		}
	}

	return resp, nil
}

// do executes a request and JSON-decodes the response into dest.
// For 204 No Content, dest is left untouched and nil is returned.
func (c *AgentClient) do(ctx context.Context, method, path string, body, dest interface{}) error {
	resp, err := c.doRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if dest != nil {
		if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
			return fmt.Errorf("agents: decode response: %w", err)
		}
	}
	return nil
}
