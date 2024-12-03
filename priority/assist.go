package priority

import (
	"github.com/akramarenkov/flow/priority/internal/distrib"
	"github.com/akramarenkov/flow/priority/types"
)

func divide(
	divider types.Divider,
	quantity uint,
	priorities []uint,
	distribution map[uint]uint,
) error {
	if err := divider(quantity, priorities, distribution); err != nil {
		return err
	}

	distributed, err := distrib.Quantity(priorities, distribution)
	if err != nil {
		return err
	}

	if distributed != quantity {
		return ErrDividerBad
	}

	return nil
}
