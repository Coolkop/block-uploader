package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"file-storage/internal/models"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetChunks(fileID uuid.UUID) ([]models.Chunk, error) {
	rows, err := r.db.Queryx(
		`select file_id, storage_id, part_order from chunk where file_id = $1`,
		fileID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("query chunks")
	}

	defer rows.Close()

	var entity Entity
	result := make([]models.Chunk, 0)

	for rows.Next() {
		if err = rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("struct scan for chunk entity: %w", err)
		}

		fileID, err := uuid.Parse(entity.FileID)
		if err != nil {
			return nil, fmt.Errorf("parse entity file id '%s': %w", entity.FileID, err)
		}

		result = append(result, models.Chunk{
			FileID:    fileID,
			StorageID: entity.StorageID,
			Order:     entity.Order,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("scan rows: %w", err)
	}

	if len(result) == 0 {
		return nil, NotFoundErr
	}

	return result, nil
}

func (r *Repository) SaveChunks(chunks []models.Chunk) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	for _, chunk := range chunks {
		err := r.insertChunk(tx, chunk)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("rollback transaction: %w", err)
			}

			return fmt.Errorf("insert chunk: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) insertChunk(tx *sql.Tx, chunk models.Chunk) error {
	_, err := tx.Exec(
		`insert into chunk (file_id, storage_id, part_order) values ($1, $2, $3)`,
		chunk.FileID.String(),
		chunk.StorageID,
		chunk.Order,
	)
	if err != nil {
		return fmt.Errorf("exec insert chunk: %w", err)
	}

	return nil
}
