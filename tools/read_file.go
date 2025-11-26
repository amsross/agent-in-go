package tools

import (
	"os"

	"google.golang.org/genai"
)

type ReadFileTool struct{}

func (t *ReadFileTool) FunctionDeclaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "read_file",
		Description: "Read the contents of a file given its path",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"file_path": {
					Type:        "string",
					Description: "The path to the file to read",
				},
			},
			Required: []string{"file_path"},
		},
	}
}

func (t *ReadFileTool) Execute(args map[string]any) (map[string]any, error) {
	file_path := args["file_path"].(string)
	file_contents, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"file_contents": file_contents,
	}, nil
}
