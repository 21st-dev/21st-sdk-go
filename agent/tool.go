package agent

import (
	"context"
	"fmt"
	"reflect"
)

type ToolConfig[T any] struct {
	Description string
	Execute     func(context.Context, T, ToolContext) (CallToolResult, error)
}

type toolDefinition[T any] struct {
	description string
	execute     func(context.Context, T, ToolContext) (CallToolResult, error)
}

func Tool[T any](config ToolConfig[T]) ToolDefinition {
	return &toolDefinition[T]{
		description: config.Description,
		execute:     config.Execute,
	}
}

func (tool *toolDefinition[T]) InputType() reflect.Type {
	return reflect.TypeFor[T]()
}

func (tool *toolDefinition[T]) ExecuteValue(
	ctx context.Context,
	input any,
	toolContext ToolContext,
) (CallToolResult, error) {
	typedInput, ok := input.(T)
	if !ok {
		return CallToolResult{}, fmt.Errorf("invalid tool input type")
	}
	return tool.execute(ctx, typedInput, toolContext)
}

func (tool *toolDefinition[T]) Manifest(_ string) ToolManifest {
	return ToolManifest{Description: tool.description}
}
