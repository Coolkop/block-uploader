package storage

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"
)

type WeightLoader struct {
	storagesNumber int
	cache          cache
}

func NewWeightLoader(storagesNumber int, cache cache) *WeightLoader {
	return &WeightLoader{storagesNumber: storagesNumber, cache: cache}
}

// Load заполнение in memory кэша фэйковыми весами хранилищ
func (l *WeightLoader) Load() error {
	newWeights := make(map[int]int64, l.storagesNumber)

	for i := 0; i < l.storagesNumber; i++ {
		// TODO необходимо разработать получение реальных весов хранилищ
		val, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return fmt.Errorf("gen weight for storageID '%d': %w", i, err)
		}

		newWeights[i] = val.Int64()
	}

	l.cache.Set(newWeights)

	return nil
}

func (l *WeightLoader) Watch(ctx context.Context) {
	for {
		select {
		case <-time.After(20 * time.Second):
			err := l.Load()
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
