package models

import "github.com/google/uuid"

type Chunk struct {
	FileID    uuid.UUID
	StorageID int
	Order     int
}
