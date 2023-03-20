package get

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	repo "file-storage/internal/repository"
)

type Processor struct {
	client     client
	repository repository
}

func New(client client, repository repository) *Processor {
	return &Processor{
		client:     client,
		repository: repository,
	}
}

// Process возвращает данный в изначальном виде. В случае отсутствия данных возвращает ошибку NotFoundErr
func (p *Processor) Process(ctx context.Context, fileID uuid.UUID) ([]byte, error) {
	chunks, err := p.repository.GetChunks(fileID)
	if err != nil {
		if errors.Is(err, repo.NotFoundErr) {
			return nil, NotFoundErr
		}

		return nil, fmt.Errorf("get chunks meta info: %w", err)
	}

	wg, ctx := errgroup.WithContext(ctx)

	ch := make(chan struct {
		order int
		data  []byte
	})

	for _, chunk := range chunks {
		chunk := chunk
		wg.Go(func() error {
			data, err := p.client.Get(ctx, chunk.StorageID, fileID)
			if err != nil {
				return fmt.Errorf("get chunk: %w", err)
			}

			ch <- struct {
				order int
				data  []byte
			}{order: chunk.Order, data: data}

			return nil
		})
	}

	go func() {
		err = wg.Wait()
		close(ch)
	}()

	var resultSize int64

	chunkDataByOrders := make(map[int][]byte)

	for dataWithOrder := range ch {
		chunkDataByOrders[dataWithOrder.order] = dataWithOrder.data
		resultSize += int64(len(dataWithOrder.data))
	}

	if err != nil {
		return nil, fmt.Errorf("obtain chunks data: %w", err)
	}

	result := make([]byte, 0, resultSize)

	for i := 0; i < len(chunkDataByOrders); i++ {
		result = append(result, chunkDataByOrders[i]...)
	}

	return result, nil
}
