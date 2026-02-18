package output

import "context"

// LLMClient defines the output port for interacting with a Large Language Model.
// Driven adapters (e.g. Gemini client) implement this interface.
type LLMClient interface {
	GenerateStructuredJSON(ctx context.Context, prompt string) (string, error)
}
