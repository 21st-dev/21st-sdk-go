package agent

import (
	"context"
	"reflect"
)

type PermissionMode string

const (
	PermissionModeDefault           PermissionMode = "default"
	PermissionModeAcceptEdits       PermissionMode = "acceptEdits"
	PermissionModeBypassPermissions PermissionMode = "bypassPermissions"
	PermissionModePlan              PermissionMode = "plan"
	PermissionModeDontAsk           PermissionMode = "dontAsk"
)

type Runtime string

const (
	RuntimeClaudeCode Runtime = "claude-code"
	RuntimeCodex      Runtime = "codex"
)

type ToolResultContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MIMEType string `json:"mimeType,omitempty"`
}

type CallToolResult struct {
	Content []ToolResultContent `json:"content"`
	IsError bool                `json:"isError,omitempty"`
}

type ToolContext struct {
	Env map[string]string `json:"-"`
}

type ToolManifest struct {
	Description string `json:"description"`
}

type ToolDefinition interface {
	Manifest(name string) ToolManifest
	InputType() reflect.Type
	ExecuteValue(context.Context, any, ToolContext) (CallToolResult, error)
}

type Tools map[string]ToolDefinition

type Sandbox struct {
	Apt             []string          `json:"apt,omitempty"`
	Build           []string          `json:"build,omitempty"`
	Setup           []string          `json:"setup,omitempty"`
	Files           map[string]string `json:"files,omitempty"`
	Cwd             string            `json:"cwd,omitempty"`
	TimeoutMs       *int              `json:"timeoutMs,omitempty"`
	NetworkAllowOut []string          `json:"networkAllowOut,omitempty"`
	NetworkDenyOut  []string          `json:"networkDenyOut,omitempty"`
}

type Config struct {
	Runtime        Runtime `json:"runtime"`
	Model          string  `json:"model"`
	SystemPrompt   any     `json:"systemPrompt,omitempty"`
	PermissionMode PermissionMode
	MaxTurns       int     `json:"maxTurns"`
	MaxBudgetUsd   float64 `json:"maxBudgetUsd,omitempty"`
	Sandbox        *Sandbox
	Tools          Tools
}

type Manifest struct {
	Runtime        Runtime                 `json:"runtime"`
	Model          string                  `json:"model"`
	SystemPrompt   any                     `json:"systemPrompt,omitempty"`
	PermissionMode PermissionMode          `json:"permissionMode"`
	MaxTurns       int                     `json:"maxTurns"`
	MaxBudgetUsd   float64                 `json:"maxBudgetUsd,omitempty"`
	Sandbox        *Sandbox                `json:"sandbox,omitempty"`
	Tools          map[string]ToolManifest `json:"tools,omitempty"`
}
