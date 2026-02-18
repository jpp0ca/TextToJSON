package gemini

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// Client is the Gemini implementation of the output.LLMClient port.
type Client struct {
	apiKey string
	model  string
}

// NewClient creates a new Gemini Client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  "gemini-3-flash-preview",
	}
}

// GenerateStructuredJSON sends a prompt to Gemini and returns the raw JSON string response.
func (c *Client) GenerateStructuredJSON(ctx context.Context, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  c.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: "You are a data extraction assistant. You analyze unstructured text and return only valid JSON with relevant structured fields. Never include explanations, markdown, or anything other than the JSON object."}},
		},
	}

	result, err := client.Models.GenerateContent(ctx, c.model, genai.Text(prompt), config)
	if err != nil {
		return "", fmt.Errorf("Gemini API error: %w", err)
	}

	text := result.Text()

	// Strip markdown code fences if present (```json ... ```)
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) >= 3 {
			lines = lines[1 : len(lines)-1]
		}
		text = strings.Join(lines, "\n")
	}

	return strings.TrimSpace(text), nil
}
