package storage

import (
	"fmt"
	"sort"
	"sync"
)

type Cache struct {
	sync.RWMutex
	weights []Weight
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Set(newWeights map[int]int64) {
	weights := make([]Weight, 0, len(newWeights))
	for storageID, weight := range newWeights {
		weights = append(weights, Weight{Weight: weight, StorageID: storageID})
	}

	sort.Slice(weights, func(i, j int) bool {
		return weights[i].Weight < weights[j].Weight
	})

	c.Lock()
	defer c.Unlock()

	c.weights = weights
}

// GetLowest получение заданного количества storageID с наименьшими весами
func (c *Cache) GetLowest(num int) ([]int, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.weights) < num {
		return nil, fmt.Errorf("not enough storage weights")
	}

	res := make([]int, 0, num)

	for _, weight := range c.weights[0:num] {
		res = append(res, weight.StorageID)
	}

	return res, nil
}
