package get

import (
	"context"

	"github.com/google/uuid"
)

type processor interface {
	Process(ctx context.Context, fileID uuid.UUID) ([]byte, error)
}
