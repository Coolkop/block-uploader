package storage

type cache interface {
	Set(newWeights map[int]int64)
}
