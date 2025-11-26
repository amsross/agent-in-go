package tools

import (
	"time"

	"google.golang.org/genai"
)

type DateTool struct{}

func (t *DateTool) FunctionDeclaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "get_current_date",
		Description: "Get the current date in the current location",
		Parameters:  &genai.Schema{Type: "object"},
	}
}

func (t *DateTool) Execute(args map[string]any) (map[string]any, error) {
	return map[string]any{
		"date": time.Now().Format(time.DateOnly),
	}, nil
}
