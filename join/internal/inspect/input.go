package inspect

import "github.com/akramarenkov/safe"

// Returns a sequence of numbers starting from 1 to value of 'quantity' inclusive
// divided into blocks which should be supplied to the input of the discipline.
func Input(quantity uint, blockSize uint) [][]uint {
	if quantity == 0 {
		return nil
	}

	if blockSize == 0 {
		return nil
	}

	blocksNumber := quantity / blockSize

	if blocksNumber*blockSize != quantity {
		blocksNumber++
	}

	blocks := make([][]uint, 0, blocksNumber)

	for _, base := range safe.IncStep(1, quantity, blockSize) {
		block := make([]uint, 0, blockSize)

		for id := range blockSize {
			item := base + id

			if item > quantity {
				break
			}

			block = append(block, item)
		}

		blocks = append(blocks, block)
	}

	return blocks
}
