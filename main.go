package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// Get the API key from the environment variable.
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\n$> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		response, err := client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(text),
			nil,
		)

		if err != nil {
			return err
		}

		part := response.Candidates[0].Content.Parts[0]
		fmt.Print(part.Text)

		fmt.Print("\n\n$> ")
	}

	return nil
}
