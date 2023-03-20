package put

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"file-storage/internal/models"
)

const (
	numberOfChunks = 5
)

type Processor struct {
	client          client
	weightsObtainer weightsObtainer
	repository      repository
	chunker         chunker
}

func New(client client, weightsObtainer weightsObtainer, repository repository, chunker chunker) *Processor {
	return &Processor{
		client:          client,
		weightsObtainer: weightsObtainer,
		repository:      repository,
		chunker:         chunker,
	}
}

func (p *Processor) Process(ctx context.Context, data []byte) (*uuid.UUID, error) {
	storageIDs, err := p.weightsObtainer.GetLowest(numberOfChunks)
	if err != nil {
		return nil, fmt.Errorf("storages with the lowest weights: %w", err)
	}

	if len(storageIDs) != numberOfChunks {
		return nil, fmt.Errorf("not enough storages")
	}

	fileID := uuid.New()
	chunks := make([]models.Chunk, 0, numberOfChunks)
	wg, ctx := errgroup.WithContext(ctx)

	for _, chunk := range p.chunker.Chunk(storageIDs, data) {
		chunk := chunk

		chunks = append(chunks, models.Chunk{
			FileID:    fileID,
			StorageID: chunk.StorageID,
			Order:     chunk.Order,
		})

		wg.Go(func() error {
			return p.client.Put(ctx, chunk.StorageID, fileID, chunk.Data)
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("upload chunks: %w", err)
	}

	if err := p.repository.SaveChunks(chunks); err != nil {
		// TODO удалить загруженные файлы, чтобы не хранить подвисшие. Можно сделать позже, т.к. БД считаем надежной
		return nil, fmt.Errorf("save chunks: %w", err)
	}

	return &fileID, nil
}
