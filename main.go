package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"google.golang.org/genai"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s%v%s", colorRed, err, colorReset)
	}
}

func newClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	return genai.NewClient(ctx, &genai.ClientConfig{})
}

func run() error {
	ctx := context.Background()

	client, err := newClient(ctx)
	if err != nil {
		return err
	}

	termRenderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
	if err != nil {
		return err
	}

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "get_current_time",
					Description: "Get the current time in a given location",
				},
			}},
		},
	}

	var conversation []*genai.Content
	var stringBuilder strings.Builder

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s$> %s", colorGreen, colorReset)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		conversation = append(conversation, &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{{Text: text}},
		})

		resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", conversation, config)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
			fmt.Println("No response from the model.")
			continue
		}

		conversation = append(conversation, resp.Candidates[0].Content)

		// Handle tool calls in a loop
		for {
			if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
				part := resp.Candidates[0].Content.Parts[0]

				if part.FunctionCall == nil {
					break // Not a function call, break the loop
				}

				if part.FunctionCall.Name == "get_current_time" {
					currentTime := getCurrentTime()
					conversation = append(conversation, &genai.Content{
						Role: "model",
						Parts: []*genai.Part{{
							FunctionResponse: &genai.FunctionResponse{
								Name:     "get_current_time",
								Response: map[string]any{"time": currentTime},
							},
						}},
					})

					resp, err = client.Models.GenerateContent(ctx, "gemini-2.5-flash", conversation, config)
					if err != nil {
						return err
					}
					if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
						fmt.Println("No response from the model after tool call.")
						break
					}
					conversation = append(conversation, resp.Candidates[0].Content)
				} else {
					break
				}
			}
		}

		// Render the final response from the model
		if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
			for _, part := range resp.Candidates[0].Content.Parts {
				fmt.Fprintf(&stringBuilder, "%s", part.Text)
			}
		}

		out, err := termRenderer.Render(stringBuilder.String())
		stringBuilder.Reset()
		if err != nil {
			log.Printf("Error rendering markdown: %v", err)
		} else {
			fmt.Print(out)
		}

		fmt.Printf("%s$> %s", colorGreen, colorReset)
	}

	return nil
}
