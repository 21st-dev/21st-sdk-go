package agent

const (
	defaultModel          = "claude-sonnet-4-6"
	defaultMaxTurns       = 50
	defaultPermissionMode = PermissionModeBypassPermissions
)

func New(config Config) Config {
	if config.Runtime == "" {
		config.Runtime = RuntimeClaudeCode
	}
	if config.Model == "" {
		config.Model = defaultModel
	}
	if config.PermissionMode == "" {
		config.PermissionMode = defaultPermissionMode
	}
	if config.MaxTurns == 0 {
		config.MaxTurns = defaultMaxTurns
	}
	if config.Tools == nil {
		config.Tools = Tools{}
	}
	return config
}

func (config Config) Manifest() Manifest {
	manifest := Manifest{
		Runtime:        config.Runtime,
		Model:          config.Model,
		SystemPrompt:   config.SystemPrompt,
		PermissionMode: config.PermissionMode,
		MaxTurns:       config.MaxTurns,
		MaxBudgetUsd:   config.MaxBudgetUsd,
		Sandbox:        config.Sandbox,
	}

	if len(config.Tools) > 0 {
		manifest.Tools = make(map[string]ToolManifest, len(config.Tools))
		for name, tool := range config.Tools {
			if tool == nil {
				continue
			}
			manifest.Tools[name] = tool.Manifest(name)
		}
	}

	return manifest
}
