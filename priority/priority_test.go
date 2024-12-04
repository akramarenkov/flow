package priority

import (
	"testing"
	"time"

	"github.com/akramarenkov/flow/priority/divider"
	"github.com/akramarenkov/flow/priority/internal/measuring"
	"github.com/akramarenkov/flow/priority/internal/research"
	"github.com/akramarenkov/flow/priority/internal/unmanaged"
	"github.com/akramarenkov/flow/priority/types"
	"github.com/akramarenkov/safe"

	"github.com/stretchr/testify/require"
)

func TestOptsAddInput(t *testing.T) {
	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: 6,
	}

	_, err := New(opts)
	require.Error(t, err)

	require.Error(t, opts.AddInput(0, make(chan uint)))
	require.Error(t, opts.AddInput(1, nil))
	require.NoError(t, opts.AddInput(1, make(chan uint)))
	require.Error(t, opts.AddInput(1, make(chan uint)))

	_, err = New(opts)
	require.NoError(t, err)
}

func TestOptsValidation(t *testing.T) {
	opts := Opts[uint]{}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider: divider.Fair,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: 6,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: 6,
		Inputs: map[uint]<-chan uint{
			0: make(chan uint),
		},
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: 6,
		Inputs: map[uint]<-chan uint{
			1: make(chan uint),
			2: nil,
		},
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: 6,
		Inputs: map[uint]<-chan uint{
			1: make(chan uint),
		},
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDisciplineFair(t *testing.T) {
	msr, err := measuring.NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.NoError(t, err)
}

func TestDisciplineRate(t *testing.T) {
	msr, err := measuring.NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.NoError(t, err)
}

func TestDisciplineFairUnbuffered(t *testing.T) {
	msr, err := measuring.NewMeasurer(6, 0)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.NoError(t, err)
}

func TestDisciplineRateUnbuffered(t *testing.T) {
	msr, err := measuring.NewMeasurer(6, 0)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.NoError(t, err)
}

func TestDisciplineError(t *testing.T) {
	msr, err := measuring.NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	wrong := func(_ uint, _ []uint, _ map[uint]uint) error {
		return nil
	}

	opts := Opts[uint]{
		Divider:          wrong,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	_, err = New(opts)
	require.Error(t, err)
}

func TestDisciplineErrorHandlersQuantityTooSmall(t *testing.T) {
	msr, err := measuring.NewMeasurer(2)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	_, err = New(opts)
	require.Error(t, err)
}

func TestDisciplineErrorInPlayOne(t *testing.T) {
	msr, err := measuring.NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	calls := 0

	wrong := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		calls++

		if calls == 1 {
			return divider.Fair(quantity, priorities, distribution)
		}

		return nil
	}

	opts := Opts[uint]{
		Divider:          wrong,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.Error(t, err)
}

func TestDisciplineErrorInPlayTwo(t *testing.T) {
	msr, err := measuring.NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 0)

	calls := 0

	wrong := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		calls++

		if calls == 1 {
			return divider.Fair(quantity, priorities, distribution)
		}

		if len(priorities) == 3 {
			return divider.Fair(quantity, priorities, distribution)
		}

		return nil
	}

	opts := Opts[uint]{
		Divider:          wrong,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.Error(t, err)
}

func TestDisciplineErrorInPlayThree(t *testing.T) {
	msr, err := measuring.NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 0)

	calls := 0

	wrong := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		calls++

		if calls == 1 {
			return divider.Fair(quantity, priorities, distribution)
		}

		if len(priorities) == 2 {
			distribution[priorities[0]] = quantity
			return nil
		}

		return divider.Fair(quantity, priorities, distribution)
	}

	opts := Opts[uint]{
		Divider:          wrong,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.Error(t, err)
}

func TestDisciplineErrorInPlayFour(t *testing.T) {
	msr, err := measuring.NewMeasurer(6)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 0)

	calls := 0

	wrong := func(quantity uint, priorities []uint, distribution map[uint]uint) error {
		calls++

		if calls == 1 {
			return divider.Fair(quantity, priorities, distribution)
		}

		if len(priorities) == 3 {
			return divider.Fair(quantity, priorities, distribution)
		}

		if quantity == msr.HandlersQuantity() {
			return divider.Fair(quantity, priorities, distribution)
		}

		return nil
	}

	opts := Opts[uint]{
		Divider:          wrong,
		HandlersQuantity: msr.HandlersQuantity(),
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	_, err = msr.Play(discipline)
	require.Error(t, err)
}

func TestDisciplineFairOverQuantity(t *testing.T) {
	msr, err := measuring.NewMeasurer(12)
	require.NoError(t, err)

	msr.AddWrite(1, 1000000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 10000)

	opts := Opts[uint]{
		Divider:          divider.Fair,
		HandlersQuantity: msr.HandlersQuantity() / 2,
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	quantities := research.InProcessing(measurements, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, opts.HandlersQuantity)
		}
	}
}

func TestDisciplineRateOverQuantity(t *testing.T) {
	msr, err := measuring.NewMeasurer(12)
	require.NoError(t, err)

	msr.AddWrite(1, 100000)
	msr.AddWrite(2, 100000)
	msr.AddWrite(3, 100000)

	opts := Opts[uint]{
		Divider:          divider.Rate,
		HandlersQuantity: msr.HandlersQuantity() / 2,
		Inputs:           msr.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	measurements, err := msr.Play(discipline)
	require.NoError(t, err)

	quantities := research.InProcessing(measurements, 100*time.Millisecond)

	for priority := range quantities {
		for id := range quantities[priority] {
			require.LessOrEqual(t, quantities[priority][id].Quantity, opts.HandlersQuantity)
		}
	}
}

func BenchmarkDisciplineFair6(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 6)
}

func BenchmarkDisciplineRate6(b *testing.B) {
	benchmarkDiscipline(b, divider.Rate, 6)
}

func BenchmarkUnmanaged6(b *testing.B) {
	benchmarkUnmanaged(b, 6)
}

func BenchmarkDisciplineFair6Unbuffered(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 6, 0)
}

func BenchmarkDisciplineRate6Unbuffered(b *testing.B) {
	benchmarkDiscipline(b, divider.Rate, 6, 0)
}

func BenchmarkUnmanaged6Unbuffered(b *testing.B) {
	benchmarkUnmanaged(b, 6, 0)
}

func BenchmarkDisciplineFair60(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 60)
}

func BenchmarkDisciplineRate60(b *testing.B) {
	benchmarkDiscipline(b, divider.Rate, 60)
}

func BenchmarkUnmanaged60(b *testing.B) {
	benchmarkUnmanaged(b, 60)
}

func BenchmarkDisciplineFair60Unbuffered(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 60, 0)
}

func BenchmarkDisciplineRate60Unbuffered(b *testing.B) {
	benchmarkDiscipline(b, divider.Rate, 60, 0)
}

func BenchmarkUnmanaged60Unbuffered(b *testing.B) {
	benchmarkUnmanaged(b, 60, 0)
}

func BenchmarkDisciplineFair600(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 600)
}

func BenchmarkDisciplineRate600(b *testing.B) {
	benchmarkDiscipline(b, divider.Rate, 600)
}

func BenchmarkUnmanaged600(b *testing.B) {
	benchmarkUnmanaged(b, 600)
}

func BenchmarkDisciplineFair600Unbuffered(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 600, 0)
}

func BenchmarkDisciplineRate600Unbuffered(b *testing.B) {
	benchmarkDiscipline(b, divider.Rate, 600, 0)
}

func BenchmarkUnmanaged600Unbuffered(b *testing.B) {
	benchmarkUnmanaged(b, 600, 0)
}

func BenchmarkDisciplineFair60InputCapacity0(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 60, 0)
}

func BenchmarkDisciplineFair60InputCapacity50(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 60, 30)
}

func BenchmarkDisciplineFair60InputCapacity100(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 60, 60)
}

func BenchmarkDisciplineFair60InputCapacity200(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 60, 120)
}

func BenchmarkDisciplineFair60InputCapacity300(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 60, 180)
}

func BenchmarkDisciplineFair60InputCapacity400(b *testing.B) {
	benchmarkDiscipline(b, divider.Fair, 60, 240)
}

func benchmarkDiscipline(
	b *testing.B,
	divider types.Divider,
	handlersQuantity uint,
	inputCapacity ...uint,
) {
	itemsQuantity, err := safe.IToI[uint](b.N)
	require.NoError(b, err)

	bnch, err := measuring.NewBenchmarker(handlersQuantity, inputCapacity...)
	require.NoError(b, err)

	bnch.AddItems(3, itemsQuantity)
	bnch.AddItems(2, itemsQuantity)
	bnch.AddItems(1, itemsQuantity)

	opts := Opts[uint]{
		Divider:          divider,
		HandlersQuantity: bnch.HandlersQuantity(),
		Inputs:           bnch.Inputs(),
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	b.ResetTimer()

	bnch.Play(discipline)
}

func benchmarkUnmanaged(
	b *testing.B,
	handlersQuantity uint,
	inputCapacity ...uint,
) {
	itemsQuantity, err := safe.IToI[uint](b.N)
	require.NoError(b, err)

	bnch, err := measuring.NewBenchmarker(handlersQuantity, inputCapacity...)
	require.NoError(b, err)

	bnch.AddItems(3, itemsQuantity)
	bnch.AddItems(2, itemsQuantity)
	bnch.AddItems(1, itemsQuantity)

	opts := unmanaged.Opts[uint]{
		HandlersQuantity: bnch.HandlersQuantity(),
		Inputs:           bnch.Inputs(),
	}

	discipline, err := unmanaged.New(opts)
	require.NoError(b, err)

	b.ResetTimer()

	bnch.Play(discipline)
}
