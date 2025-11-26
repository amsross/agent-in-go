package tools

import (
	"time"

	"google.golang.org/genai"
)

type TimeTool struct{}

func (t *TimeTool) FunctionDeclaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "get_current_time",
		Description: "Get the current time in the current location",
		Parameters:  &genai.Schema{Type: "object"},
	}
}

func (t *TimeTool) Execute(args map[string]any) (map[string]any, error) {
	return map[string]any{
		"time": time.Now().Format(time.Kitchen),
	}, nil
}
