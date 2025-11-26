package tools

import (
	"fmt"

	"google.golang.org/genai"
)


// Tool is the interface for a tool that can be called by the model.
type Tool interface {
	FunctionDeclaration() *genai.FunctionDeclaration
	Execute(args map[string]any) (map[string]any, error)
}

// ToolBox is a collection of tools.
type ToolBox struct {
	tools map[string]Tool
}

// NewToolBox creates a new ToolBox.
func NewToolBox() *ToolBox {
	return &ToolBox{
		tools: make(map[string]Tool),
	}
}

// AddTool adds a tool to the ToolBox.
func (tb *ToolBox) AddTool(tool Tool) {
	tb.tools[tool.FunctionDeclaration().Name] = tool
}

// FunctionDeclarations returns the FunctionDeclarations for all tools in the ToolBox.
func (tb *ToolBox) FunctionDeclarations() []*genai.FunctionDeclaration {
	var decls []*genai.FunctionDeclaration
	for _, tool := range tb.tools {
		decls = append(decls, tool.FunctionDeclaration())
	}
	return decls
}

// Execute executes a tool by name.
func (tb *ToolBox) Execute(name string, args map[string]any) (map[string]any, error) {
	tool, ok := tb.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool.Execute(args)
}
