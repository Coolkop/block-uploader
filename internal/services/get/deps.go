package get

import (
	"context"

	"github.com/google/uuid"

	"file-storage/internal/models"
)

type client interface {
	Get(ctx context.Context, storageID int, objectName uuid.UUID) ([]byte, error)
}

type repository interface {
	GetChunks(fileID uuid.UUID) ([]models.Chunk, error)
}
