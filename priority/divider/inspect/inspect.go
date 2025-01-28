package inspect

import (
	"errors"
	"math"

	"github.com/akramarenkov/flow/internal/consts"
	"github.com/akramarenkov/flow/priority/internal/distrib"
	"github.com/akramarenkov/flow/priority/types"

	"github.com/akramarenkov/combin"
	"github.com/akramarenkov/reusable"
	"github.com/akramarenkov/safe"
)

var (
	ErrDividerFailed                   = errors.New("divider failed with an error")
	ErrMonotonicBroken                 = errors.New("monotonic is broken")
	ErrQuantityCalculationFailed       = errors.New("failed to calculate total quantity in distribution")
	ErrQuantityNotFound                = errors.New("failed to find minimum quantity")
	ErrQuantityNotNonFatal             = errors.New("specified quantity is not non-fatal")
	ErrQuantityNotPreserved            = errors.New("total quantity in distribution does not equal to passed quantity")
	ErrQuantityNotSuitable             = errors.New("specified quantity is not suitable")
	ErrReferenceBasisCalculationFailed = errors.New("failed to calculate basis of reference distribution")
	ErrReferenceDistributionUnfilled   = errors.New("reference distribution is not filled")
)

// Checks that the total quantity of data items in the distribution created by the
// divider is equal to the quantity of data handlers passed to the divider.
func IsQuantityPreserved(divider types.Divider, set []Opts) Result {
	for _, opts := range set {
		if result := isQuantityPreserved(divider, opts); result.Conclusion != nil {
			return result
		}
	}

	return Result{}
}

func isQuantityPreserved(divider types.Divider, opts Opts) Result {
	distribution := make(map[uint]uint, len(opts.Priorities))
	priorities := reusable.New[uint](len(opts.Priorities))

	for combination := range combin.Every(opts.Priorities) {
		for quantity := range safe.Iter(0, opts.Quantity) {
			clear(distribution)

			// Protection against modification of slice by divider
			copy(priorities.Get(len(combination)), combination)

			if err := divider(quantity, priorities.Get(0), distribution); err != nil {
				result := Result{
					Conclusion: ErrDividerFailed,
					Err:        err,
					Quantity:   quantity,
					Priorities: combination,
				}

				return result
			}

			distributed, err := distrib.Quantity(combination, distribution)
			if err != nil {
				result := Result{
					Conclusion: ErrQuantityCalculationFailed,
					Err:        err,
					Quantity:   quantity,
					Priorities: combination,
				}

				return result
			}

			if distributed != quantity {
				result := Result{
					Conclusion: ErrQuantityNotPreserved,
					Quantity:   quantity,
					Priorities: combination,
				}

				return result
			}
		}
	}

	return Result{}
}

// Checks that the quantity of data items for each priority in the distribution
// created by the divider is monotonically non-decreasing as the quantity of data
// handlers passed to the divider increases.
func IsMonotonic(divider types.Divider, set []Opts) Result {
	for _, opts := range set {
		if result := isMonotonic(divider, opts); result.Conclusion != nil {
			return result
		}
	}

	return Result{}
}

func isMonotonic(divider types.Divider, opts Opts) Result {
	priorities := reusable.New[uint](len(opts.Priorities))

	for combination := range combin.Every(opts.Priorities) {
		previous := make(map[uint]uint, len(combination))

		for quantity := range safe.Iter(0, opts.Quantity) {
			actual := make(map[uint]uint, len(combination))

			// Protection against modification of slice by divider
			copy(priorities.Get(len(combination)), combination)

			if err := divider(quantity, priorities.Get(0), actual); err != nil {
				result := Result{
					Conclusion: ErrDividerFailed,
					Err:        err,
					Quantity:   quantity,
					Priorities: combination,
				}

				return result
			}

			for _, priority := range combination {
				if actual[priority] < previous[priority] {
					result := Result{
						Conclusion: ErrMonotonicBroken,
						Quantity:   quantity,
						Priorities: combination,
					}

					return result
				}
			}

			previous = actual
		}
	}

	return Result{}
}

// Finds the minimum quantity of data handlers passed to the divider such that there
// are no zero values ​​in the distribution created by the divider.
func FindMinNonFatalQuantity(divider types.Divider, opts Opts) (uint, Result) {
	for quantity := range safe.Inc(1, opts.Quantity) {
		params := Opts{
			Quantity:   quantity,
			Priorities: opts.Priorities,
		}

		if result := IsNonFatalQuantity(divider, params); result.Conclusion != nil {
			continue
		}

		return quantity, Result{}
	}

	result := Result{
		Conclusion: ErrQuantityNotFound,
		Quantity:   opts.Quantity,
		Priorities: opts.Priorities,
	}

	return 0, result
}

// Checks that there are no zero values ​​in the distribution created by the divider for
// the specified quantity of data handlers passed to the divider.
func IsNonFatalQuantity(divider types.Divider, opts Opts) Result {
	distribution := make(map[uint]uint, len(opts.Priorities))
	priorities := reusable.New[uint](len(opts.Priorities))

	for combination := range combin.Every(opts.Priorities) {
		clear(distribution)

		// Protection against modification of slice by divider
		copy(priorities.Get(len(combination)), combination)

		if err := divider(opts.Quantity, priorities.Get(0), distribution); err != nil {
			result := Result{
				Conclusion: ErrDividerFailed,
				Err:        err,
				Quantity:   opts.Quantity,
				Priorities: combination,
			}

			return result
		}

		if !distrib.IsFilled(combination, distribution) {
			result := Result{
				Conclusion: ErrQuantityNotNonFatal,
				Quantity:   opts.Quantity,
				Priorities: combination,
			}

			return result
		}
	}

	return Result{}
}

// Finds the minimum quantity of data handlers passed to the divider such that the
// quantity of data items for each priority in the distribution created by the
// divider will not differ from the real value by the specified relative difference.
//
// Suitable difference is relative and is specified as a percentage.
func FindMinSuitableQuantity(divider types.Divider, opts Opts, suitableDiff float64) (uint, Result) {
	for quantity := range safe.Inc(1, opts.Quantity) {
		params := Opts{
			Quantity:   quantity,
			Priorities: opts.Priorities,
		}

		if result := IsSuitableQuantity(divider, params, suitableDiff); result.Conclusion != nil {
			continue
		}

		return quantity, Result{}
	}

	result := Result{
		Conclusion: ErrQuantityNotFound,
		Quantity:   opts.Quantity,
		Priorities: opts.Priorities,
	}

	return 0, result
}

// Checks that in the distribution created by the divider, for the specified quantity
// of data handlers passed to the divider, the quantity of data items for each
// priority does not differ from the real value by the specified relative difference.
//
// Suitable difference is relative and is specified as a percentage.
func IsSuitableQuantity(divider types.Divider, opts Opts, suitableDiff float64) Result {
	referenceQuantity, referenceRatio, err := referenceBasis(opts)
	if err != nil {
		result := Result{
			Conclusion: ErrReferenceBasisCalculationFailed,
			Err:        err,
			Quantity:   opts.Quantity,
			Priorities: opts.Priorities,
		}

		return result
	}

	priorities := reusable.New[uint](len(opts.Priorities))

	reference := make(map[uint]uint, len(opts.Priorities))
	actual := make(map[uint]uint, len(opts.Priorities))

	for combination := range combin.Every(opts.Priorities) {
		clear(reference)
		clear(actual)

		// Protection against modification of slice by divider
		copy(priorities.Get(len(combination)), combination)

		if err := divider(referenceQuantity, priorities.Get(0), reference); err != nil {
			result := Result{
				Conclusion: ErrDividerFailed,
				Err:        err,
				Quantity:   referenceQuantity,
				Priorities: combination,
			}

			return result
		}

		if !distrib.IsFilled(combination, reference) {
			result := Result{
				Conclusion: ErrReferenceDistributionUnfilled,
				Quantity:   opts.Quantity,
				Priorities: combination,
			}

			return result
		}

		copy(priorities.Get(len(combination)), combination)

		if err := divider(opts.Quantity, priorities.Get(0), actual); err != nil {
			result := Result{
				Conclusion: ErrDividerFailed,
				Err:        err,
				Quantity:   opts.Quantity,
				Priorities: combination,
			}

			return result
		}

		if !isSuitableDiff(combination, reference, actual, referenceRatio, suitableDiff) {
			result := Result{
				Conclusion: ErrQuantityNotSuitable,
				Quantity:   opts.Quantity,
				Priorities: combination,
			}

			return result
		}
	}

	return Result{}
}

func referenceBasis(opts Opts) (uint, uint, error) {
	const factor = 1000

	// Increases the chances of getting correct reference values
	base, err := safe.AddMU(opts.Priorities...)
	if err != nil {
		return 0, 0, err
	}

	ratio, err := safe.Mul(base, factor)
	if err != nil {
		return 0, 0, err
	}

	quantity, err := safe.Mul(opts.Quantity, ratio)
	if err != nil {
		return 0, 0, err
	}

	return quantity, ratio, nil
}

// Suitable difference is relative and is specified as a percentage.
func isSuitableDiff(
	priorities []uint,
	reference map[uint]uint,
	actual map[uint]uint,
	referenceRatio uint,
	suitableDiff float64,
) bool {
	for _, priority := range priorities {
		quantity := actual[priority]
		referenceQuantity := reference[priority]

		diff := 1.0 - float64(referenceRatio)*float64(quantity)/float64(referenceQuantity)
		diff = consts.HundredPercent * math.Abs(diff)

		if diff > suitableDiff {
			return false
		}
	}

	return true
}
