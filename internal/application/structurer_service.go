package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"GoApiProject/internal/domain/entity"
	"GoApiProject/internal/ports/output"
)

const maxRetries = 3

// StructurerServiceImpl is the application service that implements StructurerService.
// It orchestrates the LLM call and validates/retries the response.
type StructurerServiceImpl struct {
	llm output.LLMClient
}

// NewStructurerService creates a new StructurerServiceImpl with the given LLM client.
func NewStructurerService(llm output.LLMClient) *StructurerServiceImpl {
	return &StructurerServiceImpl{llm: llm}
}

// Structure sends raw text to the LLM and returns a validated structured JSON response.
// It retries up to 3 times if the LLM returns invalid JSON.
func (s *StructurerServiceImpl) Structure(ctx context.Context, req entity.StructureRequest) (*entity.StructureResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	prompt := buildPrompt(req.RawText)

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("[structurer] attempt %d/%d", attempt, maxRetries)

		rawJSON, err := s.llm.GenerateStructuredJSON(ctx, prompt)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: LLM call failed: %w", attempt, err)
			log.Printf("[structurer] %v", lastErr)
			continue
		}

		response, err := parseAndValidate(rawJSON)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: validation failed: %w", attempt, err)
			log.Printf("[structurer] %v", lastErr)
			continue
		}

		log.Printf("[structurer] success on attempt %d", attempt)
		return response, nil
	}

	return nil, fmt.Errorf("all %d attempts failed, last error: %w", maxRetries, lastErr)
}

// buildPrompt creates the instruction prompt for the LLM.
func buildPrompt(rawText string) string {
	return fmt.Sprintf(`You are a data extraction assistant. Analyze the following unstructured text and extract all relevant structured information from it.

Rules:
1. Return ONLY a valid JSON object, nothing else — no markdown, no explanation, no code fences.
2. Infer the best field names in English (snake_case).
3. Use appropriate types: strings, numbers, arrays, nested objects.
4. Dates should be in "YYYY-MM-DD" format when possible.
5. If the text describes a trip, extract fields like destination, dates, flight numbers, etc.
6. If it's a recipe, extract dish name, ingredients, steps, servings, etc.
7. If it's an event, extract title, date, location, participants, etc.
8. Adapt the schema to whatever content makes sense.

Text to analyze:
"""%s"""`, rawText)
}

// parseAndValidate unmarshals the raw JSON string and validates the structure.
func parseAndValidate(rawJSON string) (*entity.StructureResponse, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(rawJSON), &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	response := &entity.StructureResponse{Data: data}
	if err := response.Validate(); err != nil {
		return nil, err
	}

	return response, nil
}
