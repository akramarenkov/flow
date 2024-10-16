package inspect

import "github.com/akramarenkov/safe"

type description struct {
	EffectiveJoinSize uint
	EffectiveQuantity uint
	Joins             uint
	RemainderQuantity uint
	UnusedJoinSize    uint
}

// Returns slices expected from the discipline output channel.
func Expected(quantity uint, blockSize uint, joinSize uint) [][]uint {
	return genExpected(1, calcDescription(quantity, blockSize, joinSize))
}

// Returns slices expected from the discipline output channel if there is one delay
// in receiving a data set blocks leading to the formation of a timeout in the
// discipline during accumulating some output slice.
func ExpectedWithTimeout(
	quantity uint,
	pauseAt uint,
	blockSize uint,
	joinSize uint,
) [][]uint {
	descs := calcDescriptionWithTimeout(quantity, pauseAt, blockSize, joinSize)

	blocks := make([][]uint, 0, len(descs))

	begin := uint(1)

	for _, desc := range descs {
		blocks = append(blocks, genExpected(begin, desc)...)

		begin += desc.EffectiveQuantity + desc.RemainderQuantity
	}

	return blocks
}

func calcDescription(quantity uint, blockSize uint, joinSize uint) description {
	if quantity == 0 {
		return description{}
	}

	if blockSize == 0 {
		return description{}
	}

	if joinSize == 0 {
		return description{}
	}

	desc := description{}

	blocksInJoin := joinSize / blockSize

	desc.EffectiveJoinSize = blocksInJoin * blockSize
	desc.UnusedJoinSize = joinSize - desc.EffectiveJoinSize

	if blocksInJoin == 0 {
		desc.EffectiveJoinSize = blockSize
		desc.UnusedJoinSize = 0
	}

	desc.Joins = quantity / desc.EffectiveJoinSize

	desc.EffectiveQuantity = desc.Joins * desc.EffectiveJoinSize
	desc.RemainderQuantity = quantity - desc.EffectiveQuantity

	if desc.Joins == 0 {
		desc.Joins = 1
		desc.EffectiveQuantity = quantity
		desc.RemainderQuantity = 0
	}

	if desc.RemainderQuantity > desc.UnusedJoinSize {
		desc.Joins++
	}

	return desc
}

func calcDescriptionWithTimeout(
	quantity uint,
	pauseAt uint,
	blockSize uint,
	joinSize uint,
) []description {
	if blockSize == 0 {
		return nil
	}

	if pauseAt == 0 {
		return []description{calcDescription(quantity, blockSize, joinSize)}
	}

	if pauseAt > quantity {
		return []description{calcDescription(quantity, blockSize, joinSize)}
	}

	blocksInPauseAt := pauseAt / blockSize
	beforePauseAt := blocksInPauseAt * blockSize

	if beforePauseAt == pauseAt {
		beforePauseAt -= blockSize
	}

	descs := []description{
		calcDescription(beforePauseAt, blockSize, joinSize),
		calcDescription(quantity-beforePauseAt, blockSize, joinSize),
	}

	return descs
}

func genExpected(begin uint, desc description) [][]uint {
	blocks := make([][]uint, 0, desc.Joins)

	effectiveEnd := desc.EffectiveQuantity + begin - 1

	for item := range safe.Inc(begin, effectiveEnd) {
		id := (item - begin) % desc.EffectiveJoinSize

		if id == 0 {
			blocks = append(blocks, make([]uint, 0, desc.EffectiveJoinSize))
		}

		blocks[len(blocks)-1] = append(blocks[len(blocks)-1], item)
	}

	if desc.RemainderQuantity > desc.UnusedJoinSize {
		blocks = append(blocks, make([]uint, 0, desc.RemainderQuantity))

		for base := range safe.Inc(1, desc.RemainderQuantity) {
			item := base + effectiveEnd

			blocks[len(blocks)-1] = append(blocks[len(blocks)-1], item)
		}

		return blocks
	}

	for base := range safe.Inc(1, desc.RemainderQuantity) {
		item := base + effectiveEnd

		blocks[len(blocks)-1] = append(blocks[len(blocks)-1], item)
	}

	return blocks
}
