package agents

import "net/http"

// --- Sandbox types ---

type CreateSandboxRequest struct {
	Agent string            `json:"agent"`
	Files map[string]string `json:"files,omitempty"`
	Envs  map[string]string `json:"envs,omitempty"`
	Setup []string          `json:"setup,omitempty"`
}

type Sandbox struct {
	ID        string `json:"id"`
	SandboxID string `json:"sandboxId"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type SandboxDetail struct {
	ID        string          `json:"id"`
	SandboxID string          `json:"sandboxId"`
	Status    string          `json:"status"`
	Error     *string         `json:"error,omitempty"`
	Agent     SandboxAgent    `json:"agent"`
	Threads   []ThreadSummary `json:"threads"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

type SandboxAgent struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// --- Thread types ---

type ThreadSummary struct {
	ID        string  `json:"id"`
	Name      *string `json:"name,omitempty"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"createdAt"`
}

type Thread struct {
	ID        string      `json:"id"`
	Name      *string     `json:"name,omitempty"`
	Status    string      `json:"status"`
	Messages  interface{} `json:"messages,omitempty"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
}

type RunThreadMessage struct {
	ID    string                   `json:"id,omitempty"`
	Role  string                   `json:"role"`
	Parts []map[string]interface{} `json:"parts"`
}

type RunThreadRequest struct {
	Agent     string             `json:"agent"`
	Messages  []RunThreadMessage `json:"messages"`
	SandboxID string             `json:"sandboxId,omitempty"`
	ThreadID  string             `json:"threadId,omitempty"`
	Name      string             `json:"name,omitempty"`
}

type RunThreadResult struct {
	SandboxID string         `json:"sandboxId"`
	ThreadID  string         `json:"threadId"`
	Response  *http.Response `json:"-"`
	ResumeURL string         `json:"resumeUrl"`
}

// --- Token types ---

type CreateTokenRequest struct {
	Agent     string `json:"agent,omitempty"`
	UserID    string `json:"userId,omitempty"`
	ExpiresIn string `json:"expiresIn,omitempty"`
}

type Token struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

// --- File types ---

type FileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// --- Exec types ---

type ExecRequest struct {
	SandboxID string            `json:"sandboxId"`
	Command   string            `json:"command"`
	Cwd       string            `json:"cwd,omitempty"`
	Envs      map[string]string `json:"envs,omitempty"`
	TimeoutMs *int              `json:"timeoutMs,omitempty"`
}

type ExecResult struct {
	ExitCode int    `json:"exitCode"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

// --- Git types ---

type GitCloneRequest struct {
	URL   string `json:"url"`
	Path  string `json:"path,omitempty"`
	Token string `json:"token,omitempty"`
	Depth *int   `json:"depth,omitempty"`
}

type GitCloneResult struct {
	Path string `json:"path"`
}
