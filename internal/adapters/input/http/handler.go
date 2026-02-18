package http

import (
	"encoding/json"
	"net/http"

	"GoApiProject/internal/domain/entity"
	"GoApiProject/internal/ports/input"
)

// StructureHandler is the HTTP driving adapter for text structuring operations.
type StructureHandler struct {
	service input.StructurerService
}

// NewStructureHandler creates a new StructureHandler with the given service.
func NewStructureHandler(service input.StructurerService) *StructureHandler {
	return &StructureHandler{service: service}
}

// errorResponse is the standard error response format.
type errorResponse struct {
	Error string `json:"error"`
}

// StructureText handles POST /structure.
//
//	@Summary		Structure unstructured text into JSON
//	@Description	Receives raw unstructured text and uses Google Gemini to extract structured data as JSON. The LLM infers the best schema automatically. Retries up to 3 times on validation failure.
//	@Tags			structurer
//	@Accept			json
//	@Produce		json
//	@Param			request	body		entity.StructureRequest	true	"Raw text to structure"
//	@Success		200		{object}	entity.StructureResponse
//	@Failure		400		{object}	errorResponse	"invalid request body or empty raw_text"
//	@Failure		422		{object}	errorResponse	"LLM failed to produce valid JSON after 3 retries"
//	@Router			/structure [post]
func (h *StructureHandler) StructureText(w http.ResponseWriter, r *http.Request) {
	var req entity.StructureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	result, err := h.service.Structure(r.Context(), req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
