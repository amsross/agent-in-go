package tools

import (
	"encoding/json"
	"os"

	"google.golang.org/genai"
)

type ListFilesTool struct{}

func (t *ListFilesTool) FunctionDeclaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "list_files",
		Description: "List the contents of of a directory given its path",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"path": {
					Type:        "string",
					Description: "The path to the directory to list. If not provided, the current directory will be listed.",
				},
			},
		},
	}
}

func (t *ListFilesTool) Execute(args map[string]any) (map[string]any, error) {
	path, ok := args["path"].(string)
	if !ok {
		path = "."
	}

	directory_entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var directories []string
	var files []string
	for _, entry := range directory_entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name()+"/")
		} else {
			files = append(files, entry.Name())
		}
	}

	result, err := json.Marshal(append(directories[:], files[:]...))
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"directory_contents": string(result),
	}, nil
}
