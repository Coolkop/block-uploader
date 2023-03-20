package put_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "file-storage/internal/services/put"
)

func TestChunker_Chunk(t *testing.T) {
	testCases := []struct {
		name       string
		storageIDs []int
		data       []byte
		expected   []DataChunk
	}{
		{
			name:       "data equally divided between storages",
			storageIDs: []int{3, 10, 8},
			data:       []byte("111122223333"),
			expected: []DataChunk{
				{
					StorageID: 3,
					Data:      []byte("1111"),
					Order:     0,
				},
				{
					StorageID: 10,
					Data:      []byte("2222"),
					Order:     1,
				},
				{
					StorageID: 8,
					Data:      []byte("3333"),
					Order:     2,
				},
			},
		},
		{
			name:       "data not equally divided between storages",
			storageIDs: []int{3, 10, 8},
			data:       []byte("11112223"),
			expected: []DataChunk{
				{
					StorageID: 3,
					Data:      []byte("11"),
					Order:     0,
				},
				{
					StorageID: 10,
					Data:      []byte("112"),
					Order:     1,
				},
				{
					StorageID: 8,
					Data:      []byte("223"),
					Order:     2,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewChunker()
			actual := c.Chunk(tc.storageIDs, tc.data)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
