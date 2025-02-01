package inspect

import (
	"math"
	"testing"

	"github.com/akramarenkov/flow/priority/divider"

	"github.com/akramarenkov/safe"
	"github.com/stretchr/testify/require"
)

func TestIsQuantityPreserved(t *testing.T) {
	result := IsQuantityPreserved(divider.Fair, DefaultSet())
	require.NoError(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.Zero(t, result.Quantity)
	require.Empty(t, result.Priorities)

	result = IsQuantityPreserved(divider.Rate, DefaultSet())
	require.NoError(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.Zero(t, result.Quantity)
	require.Empty(t, result.Priorities)
}

func TestIsQuantityPreservedNegativeConclusion(t *testing.T) {
	overflowing := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		if quantity == 0 {
			return nil
		}

		if len(priorities) == 1 {
			distribution[priorities[0]] = quantity
			return nil
		}

		for _, priority := range priorities {
			distribution[priority] = math.MaxUint
		}

		return nil
	}

	wrong := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		for _, priority := range priorities {
			distribution[priority] = quantity
		}

		return nil
	}

	result := IsQuantityPreserved(dividerFailed, DefaultSet())
	require.Error(t, result.Conclusion)
	require.Error(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)

	result = IsQuantityPreserved(overflowing, DefaultSet())
	require.Error(t, result.Conclusion)
	require.Error(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)

	result = IsQuantityPreserved(wrong, DefaultSet())
	require.Error(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)
}

func TestIsMonotonic(t *testing.T) {
	result := IsMonotonic(divider.Fair, DefaultSet())
	require.NoError(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.Zero(t, result.Quantity)
	require.Empty(t, result.Priorities)

	result = IsMonotonic(divider.Rate, DefaultSet())
	require.NoError(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.Zero(t, result.Quantity)
	require.Empty(t, result.Priorities)
}

func TestIsMonotonicNegativeConclusion(t *testing.T) {
	unmonotonic := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		divisor, err := safe.AddMU(priorities...)
		if err != nil {
			return err
		}

		base := float64(quantity) / float64(divisor)
		remainder := quantity

		for _, priority := range priorities {
			part := uint(base * float64(priority))
			remainder -= part
			distribution[priority] += part
		}

		distribution[priorities[0]] += remainder

		return nil
	}

	result := IsMonotonic(dividerFailed, DefaultSet())
	require.Error(t, result.Conclusion)
	require.Error(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)

	result = IsQuantityPreserved(unmonotonic, DefaultSet())
	require.NoError(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.Zero(t, result.Quantity)
	require.Empty(t, []uint(nil), result.Priorities)

	result = IsMonotonic(unmonotonic, DefaultSet())
	require.Error(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)
}

func TestFindMinNonFatalQuantity(t *testing.T) {
	expected := map[string][]uint{"fair": {10}, "rate": {10}}

	for id, opts := range DefaultSet() {
		quantity, result := FindMinNonFatalQuantity(divider.Fair, opts)
		require.NoError(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.Zero(t, result.Quantity)
		require.Empty(t, result.Priorities)
		require.Equal(t, expected["fair"][id], quantity)
	}

	for id, opts := range DefaultSet() {
		quantity, result := FindMinNonFatalQuantity(divider.Rate, opts)
		require.NoError(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.Zero(t, result.Quantity)
		require.Empty(t, result.Priorities)
		require.Equal(t, expected["rate"][id], quantity)
	}
}

func TestFindMinNonFatalQuantityNegativeConclusion(t *testing.T) {
	for _, opts := range DefaultSet() {
		quantity, result := FindMinNonFatalQuantity(dividerFailed, opts)
		require.Error(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.NotZero(t, result.Quantity)
		require.NotEmpty(t, result.Priorities)
		require.Equal(t, uint(0), quantity)
	}

	fatal := Opts{
		Quantity:   1,
		Priorities: []uint{3, 2, 1},
	}

	quantity, result := FindMinNonFatalQuantity(divider.Fair, fatal)
	require.Error(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)
	require.Equal(t, uint(0), quantity)
}

func TestIsNonFatalQuantity(t *testing.T) {
	for _, opts := range DefaultSet() {
		result := IsNonFatalQuantity(divider.Fair, opts)
		require.NoError(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.Zero(t, result.Quantity)
		require.Empty(t, result.Priorities)
	}

	for _, opts := range DefaultSet() {
		result := IsNonFatalQuantity(divider.Rate, opts)
		require.NoError(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.Zero(t, result.Quantity)
		require.Empty(t, result.Priorities)
	}
}

func TestIsNonFatalQuantityNegativeConclusion(t *testing.T) {
	for _, opts := range DefaultSet() {
		result := IsNonFatalQuantity(dividerFailed, opts)
		require.Error(t, result.Conclusion)
		require.Error(t, result.Err)
		require.NotZero(t, result.Quantity)
		require.NotEmpty(t, result.Priorities)
	}

	fatal := Opts{
		Quantity:   1,
		Priorities: []uint{3, 2, 1},
	}

	result := IsNonFatalQuantity(divider.Fair, fatal)
	require.Error(t, result.Conclusion)
	require.NoError(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)
}

func TestFindMinSuitableQuantity(t *testing.T) {
	expected := map[string][]uint{"fair": {120}, "rate": {644}}

	for id, opts := range DefaultSet() {
		quantity, result := FindMinSuitableQuantity(divider.Fair, opts, 5)
		require.NoError(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.Zero(t, result.Quantity)
		require.Empty(t, result.Priorities)
		require.Equal(t, expected["fair"][id], quantity)
	}

	for id, opts := range DefaultSet() {
		quantity, result := FindMinSuitableQuantity(divider.Rate, opts, 5)
		require.NoError(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.Zero(t, result.Quantity)
		require.Empty(t, result.Priorities)
		require.Equal(t, expected["rate"][id], quantity)
	}
}

func TestFindMinSuitableQuantityNegativeConclusion(t *testing.T) {
	for _, opts := range DefaultSet() {
		quantity, result := FindMinSuitableQuantity(dividerFailed, opts, 5)
		require.Error(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.NotZero(t, result.Quantity)
		require.NotEmpty(t, result.Priorities)
		require.Equal(t, uint(0), quantity)
	}

	for _, opts := range DefaultSet() {
		quantity, result := FindMinSuitableQuantity(divider.Fair, opts, 0.4)
		require.Error(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.NotZero(t, result.Quantity)
		require.NotEmpty(t, result.Priorities)
		require.Equal(t, uint(0), quantity)
	}
}

func TestIsSuitableQuantity(t *testing.T) {
	for _, opts := range DefaultSet() {
		result := IsSuitableQuantity(divider.Fair, opts, 5)
		require.NoError(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.Zero(t, result.Quantity)
		require.Empty(t, result.Priorities)
	}

	for _, opts := range DefaultSet() {
		result := IsSuitableQuantity(divider.Rate, opts, 5)
		require.NoError(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.Zero(t, result.Quantity)
		require.Empty(t, result.Priorities)
	}
}

func TestIsSuitableQuantityNegativeConclusion(t *testing.T) {
	unfilling := func(_ uint, _ []uint, _ map[uint]uint) error {
		return nil
	}

	calls := 0

	double := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		calls++

		if calls == 1 {
			return divider.Fair(quantity, priorities, distribution)
		}

		return dividerFailed(quantity, priorities, distribution)
	}

	for _, opts := range DefaultSet() {
		result := IsSuitableQuantity(dividerFailed, opts, 5)
		require.Error(t, result.Conclusion)
		require.Error(t, result.Err)
		require.NotZero(t, result.Quantity)
		require.NotEmpty(t, result.Priorities)
	}

	for _, opts := range DefaultSet() {
		result := IsSuitableQuantity(unfilling, opts, 5)
		require.Error(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.NotZero(t, result.Quantity)
		require.NotEmpty(t, result.Priorities)
	}

	for _, opts := range DefaultSet() {
		result := IsSuitableQuantity(double, opts, 5)
		require.Error(t, result.Conclusion)
		require.Error(t, result.Err)
		require.NotZero(t, result.Quantity)
		require.NotEmpty(t, result.Priorities)
	}

	for _, opts := range DefaultSet() {
		result := IsSuitableQuantity(divider.Fair, opts, 0.4)
		require.Error(t, result.Conclusion)
		require.NoError(t, result.Err)
		require.NotZero(t, result.Quantity)
		require.NotEmpty(t, result.Priorities)
	}

	opts := Opts{
		Quantity:   1,
		Priorities: []uint{math.MaxUint, 1},
	}

	result := IsSuitableQuantity(divider.Fair, opts, 5)
	require.Error(t, result.Conclusion)
	require.Error(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)

	opts = Opts{
		Quantity:   1,
		Priorities: []uint{math.MaxUint - 1, 1},
	}

	result = IsSuitableQuantity(divider.Fair, opts, 5)
	require.Error(t, result.Conclusion)
	require.Error(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)

	opts = Opts{
		Quantity:   2,
		Priorities: []uint{(math.MaxUint) / 1000 / 2, 1},
	}

	result = IsSuitableQuantity(divider.Fair, opts, 5)
	require.Error(t, result.Conclusion)
	require.Error(t, result.Err)
	require.NotZero(t, result.Quantity)
	require.NotEmpty(t, result.Priorities)
}

func dividerFailed(quantity uint, _ []uint, _ map[uint]uint) error {
	if quantity != 0 {
		return ErrDividerFailed
	}

	return nil
}

func BenchmarkIsQuantityPreserved(b *testing.B) {
	set := DefaultSet()

	var result Result

	for range b.N {
		result = IsQuantityPreserved(divider.Fair, set)
	}

	require.NoError(b, result.Conclusion)
}

func BenchmarkIsMonotonic(b *testing.B) {
	set := DefaultSet()

	var result Result

	for range b.N {
		result = IsMonotonic(divider.Fair, set)
	}

	require.NoError(b, result.Conclusion)
}

func BenchmarkFindMinNonFatalQuantity(b *testing.B) {
	opts := DefaultSet()[0]

	var result Result

	for range b.N {
		_, result = FindMinNonFatalQuantity(divider.Fair, opts)
	}

	require.NoError(b, result.Conclusion)
}

func BenchmarkIsNonFatalQuantity(b *testing.B) {
	opts := DefaultSet()[0]

	var result Result

	for range b.N {
		result = IsNonFatalQuantity(divider.Fair, opts)
	}

	require.NoError(b, result.Conclusion)
}

func BenchmarkFindMinSuitableQuantity(b *testing.B) {
	opts := DefaultSet()[0]

	var result Result

	for range b.N {
		_, result = FindMinSuitableQuantity(divider.Fair, opts, 5)
	}

	require.NoError(b, result.Conclusion)
}

func BenchmarkIsSuitableQuantity(b *testing.B) {
	opts := DefaultSet()[0]

	var result Result

	for range b.N {
		result = IsSuitableQuantity(divider.Fair, opts, 5)
	}

	require.NoError(b, result.Conclusion)
}
