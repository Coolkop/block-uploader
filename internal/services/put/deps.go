package put

import (
	"context"

	"github.com/google/uuid"

	"file-storage/internal/models"
)

type client interface {
	Put(ctx context.Context, storageID int, objectName uuid.UUID, data []byte) error
}

type weightsObtainer interface {
	GetLowest(num int) ([]int, error)
}

type repository interface {
	SaveChunks(chunks []models.Chunk) error
}

type chunker interface {
	Chunk(storageIDs []int, data []byte) []DataChunk
}
