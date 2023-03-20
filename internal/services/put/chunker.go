package put

type Chunker struct{}

func NewChunker() *Chunker {
	return &Chunker{}
}

func (c *Chunker) Chunk(storageIDs []int, data []byte) []DataChunk {
	num := len(storageIDs)

	chunks := make([]DataChunk, 0, len(storageIDs))

	for i := 0; i < num; i++ {
		storageID := storageIDs[i]

		left := i * len(data) / num
		right := (i + 1) * len(data) / num
		if i == num {
			right = len(data)
		}

		chunks = append(chunks, DataChunk{
			Data:      data[left:right],
			StorageID: storageID,
			Order:     i,
		})
	}

	return chunks
}
