package tools

import "context"

type ToolType string

const (
	Front ToolType = "front"
	Rear  ToolType = "rear"
)

// Tool is a tool for the llm agent to interact with different applications.
type Tool interface {
	Name() string
	Types() ToolType
	Description() string
	Call(ctx context.Context, input string) (string, error)
}
