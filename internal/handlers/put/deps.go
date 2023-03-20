package put

import (
	"context"

	"github.com/google/uuid"
)

type processor interface {
	Process(ctx context.Context, data []byte) (*uuid.UUID, error)
}
