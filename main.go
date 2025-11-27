package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/amsross/agent-in-go/tools"
	"github.com/charmbracelet/glamour"
	"google.golang.org/genai"
)

const (
	colorRed        = "\033[31m"
	colorGreen      = "\033[32m"
	colorYellow     = "\033[033m"
	colorReset      = "\033[0m"
	geminiAPIKeyEnv = "GEMINI_API_KEY"
	geminiModel     = "gemini-2.5-flash"
	initialPrompt   = "You are a helpful programming assistant that can use basic tools to answer user questions about the local filesystem and about software engineering questions. Your background is that of a senior software engineer with many years of experience in programming NodeJS and Go."
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s%v%s", colorRed, err, colorReset)
	}
}

func newClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv(geminiAPIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("%s environment variable not set", geminiAPIKeyEnv)
	}

	return genai.NewClient(ctx, &genai.ClientConfig{})
}

func setupTools() *tools.ToolBox {
	toolBox := tools.NewToolBox()
	toolBox.AddTool(&tools.TimeTool{})
	toolBox.AddTool(&tools.DateTool{})
	toolBox.AddTool(&tools.ReadFileTool{})
	toolBox.AddTool(&tools.ListFilesTool{})
	return toolBox
}

func handleToolCalls(
	ctx context.Context,
	client *genai.Client,
	toolBox *tools.ToolBox,
	conversation []*genai.Content,
	config *genai.GenerateContentConfig,
) ([]*genai.Content, error) {
	for {
		content := conversation[len(conversation)-1]
		if len(content.Parts) == 0 {
			break
		}

		part := content.Parts[0]
		functionCall := part.FunctionCall

		if functionCall == nil {
			break // Not a function call, break the loop
		}

		functionResponse, err := toolBox.Execute(functionCall.Name, functionCall.Args)
		if err != nil {
			return nil, err
		}

		conversation = append(conversation, &genai.Content{
			Role: "function",
			Parts: []*genai.Part{
				genai.NewPartFromFunctionResponse(
					functionCall.Name,
					functionResponse,
				),
			},
		})

		resp, err := client.Models.GenerateContent(ctx, geminiModel, conversation, config)
		if err != nil {
			return nil, err
		}
		if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
			fmt.Println("No response from the model after tool call.")
			break
		}
		conversation = append(conversation, resp.Candidates[0].Content)
	}

	return conversation, nil
}

func renderModelResponse(conversation []*genai.Content) error {
	if len(conversation) > 0 {
		var stringBuilder strings.Builder
		termRenderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
		if err != nil {
			return err
		}

		for _, part := range conversation[len(conversation)-1].Parts {
			fmt.Fprintf(&stringBuilder, "%s", part.Text)
		}

		output, err := termRenderer.Render(stringBuilder.String())
		if err != nil {
			log.Printf("Error rendering markdown: %v", err)
		} else {
			fmt.Print(output)
		}
	}

	return nil
}

func run() error {
	ctx := context.Background()

	client, err := newClient(ctx)
	if err != nil {
		return err
	}

	toolBox := setupTools()
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: toolBox.FunctionDeclarations()},
		},
	}

	conversation := []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			genai.NewPartFromText(initialPrompt),
		},
	}}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s$> %s", colorGreen, colorYellow)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.EqualFold(text, "") || strings.EqualFold(text, "exit") || strings.EqualFold(text, "quit") {
			break
		}

		conversation = append(conversation, &genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				genai.NewPartFromText(text),
			},
		})

		resp, err := client.Models.GenerateContent(ctx, geminiModel, conversation, config)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
			fmt.Println("No response from the model.")
			continue
		}

		conversation = append(conversation, resp.Candidates[0].Content)

		// Handle tool calls in a loop
		conversation, err = handleToolCalls(ctx, client, toolBox, conversation, config)
		if err != nil {
			return err
		}

		// Render the final response from the model
		renderModelResponse(conversation)

		fmt.Printf("%s$> %s", colorGreen, colorYellow)
	}

	return nil
}
