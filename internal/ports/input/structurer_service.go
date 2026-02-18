package input

import (
	"context"

	"GoApiProject/internal/domain/entity"
)

// StructurerService defines the input port for text structuring operations.
// Driving adapters (e.g. HTTP handlers) depend on this interface.
type StructurerService interface {
	Structure(ctx context.Context, req entity.StructureRequest) (*entity.StructureResponse, error)
}
