package storage_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "file-storage/internal/services/storage"
)

func TestCache(t *testing.T) {
	testCases := []struct {
		name        string
		newWeights  map[int]int64
		nums        int
		expected    []int
		expectedErr error
	}{
		{
			name: "not enough weights",
			newWeights: map[int]int64{
				0: 1,
				1: 2,
			},
			nums:        3,
			expectedErr: fmt.Errorf("not enough storage weights"),
		},
		{
			name: "get lowest N",
			newWeights: map[int]int64{
				0: 10,
				1: 3,
				2: 45,
				3: 9,
				4: 15,
			},
			nums: 3,
			expected: []int{
				1, 3, 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCache()

			c.Set(tc.newWeights)
			actual, err := c.GetLowest(tc.nums)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
