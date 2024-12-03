package research

import (
	"testing"
	"time"

	"github.com/akramarenkov/flow/internal/qot"
	"github.com/akramarenkov/flow/priority/internal/measuring"

	chartsopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/require"
)

func TestQuantityPerInterval(t *testing.T) {
	measurements := []measuring.Measure{
		{
			Item:     0,
			Kind:     measuring.KindCompleted,
			Priority: 1,
			Time:     11 * time.Microsecond,
		},
		{
			Item:     0,
			Kind:     measuring.KindProcessed,
			Priority: 1,
			Time:     10 * time.Microsecond,
		},
		{
			Item:     0,
			Kind:     measuring.KindReceived,
			Priority: 1,
			Time:     0,
		},
		{
			Item:     1,
			Kind:     measuring.KindCompleted,
			Priority: 2,
			Time:     25 * time.Microsecond,
		},
		{
			Item:     1,
			Kind:     measuring.KindProcessed,
			Priority: 2,
			Time:     20 * time.Microsecond,
		},
		{
			Item:     1,
			Kind:     measuring.KindReceived,
			Priority: 2,
			Time:     time.Microsecond,
		},
		{
			Item:     2,
			Kind:     measuring.KindProcessed,
			Priority: 3,
			Time:     30 * time.Microsecond,
		},
		{
			Item:     2,
			Kind:     measuring.KindCompleted,
			Priority: 3,
			Time:     33 * time.Microsecond,
		},
		{
			Item:     2,
			Kind:     measuring.KindReceived,
			Priority: 3,
			Time:     0,
		},
	}

	interval := 5 * time.Microsecond

	expected := map[uint][]qot.QoT{
		1: {
			{Quantity: 0, Time: -interval},
			{Quantity: 1, Time: 0},
			{Quantity: 0, Time: interval},
			{Quantity: 2, Time: 2 * interval},
			{Quantity: 0, Time: 3 * interval},
			{Quantity: 0, Time: 4 * interval},
			{Quantity: 0, Time: 5 * interval},
			{Quantity: 0, Time: 6 * interval},
			{Quantity: 0, Time: 7 * interval},
		},
		2: {
			{Quantity: 0, Time: -interval},
			{Quantity: 1, Time: 0},
			{Quantity: 0, Time: interval},
			{Quantity: 0, Time: 2 * interval},
			{Quantity: 0, Time: 3 * interval},
			{Quantity: 1, Time: 4 * interval},
			{Quantity: 1, Time: 5 * interval},
			{Quantity: 0, Time: 6 * interval},
			{Quantity: 0, Time: 7 * interval},
		},
		3: {
			{Quantity: 0, Time: -interval},
			{Quantity: 1, Time: 0},
			{Quantity: 0, Time: interval},
			{Quantity: 0, Time: 2 * interval},
			{Quantity: 0, Time: 3 * interval},
			{Quantity: 0, Time: 4 * interval},
			{Quantity: 0, Time: 5 * interval},
			{Quantity: 2, Time: 6 * interval},
			{Quantity: 0, Time: 7 * interval},
		},
	}

	quantities := QuantityPerInterval(measurements, interval)
	require.Equal(t, expected, quantities)
}

func TestQuantityPerIntervalZeroInput(t *testing.T) {
	quantities := QuantityPerInterval(nil, 5*time.Microsecond)
	require.Equal(t, map[uint][]qot.QoT(nil), quantities)

	quantities = QuantityPerInterval([]measuring.Measure{}, 5*time.Microsecond)
	require.Equal(t, map[uint][]qot.QoT(nil), quantities)
}

func TestQotToLineChart(t *testing.T) {
	interval := 5 * time.Microsecond

	quantities := map[uint][]qot.QoT{
		1: {
			{Quantity: 0, Time: -interval},
			{Quantity: 1, Time: 0},
			{Quantity: 0, Time: interval},
		},
		2: {
			{Quantity: 0, Time: -interval},
			{Quantity: 2, Time: 0},
			{Quantity: 0, Time: interval},
		},
		3: {
			{Quantity: 0, Time: -interval},
			{Quantity: 3, Time: 0},
			{Quantity: 0, Time: interval},
		},
	}

	expectedAxisY := map[uint][]chartsopts.LineData{
		1: {
			{Value: uint(0), Name: "-5µs"},
			{Value: uint(1), Name: "0s"},
			{Value: uint(0), Name: "5µs"},
		},
		2: {
			{Value: uint(0), Name: "-5µs"},
			{Value: uint(2), Name: "0s"},
			{Value: uint(0), Name: "5µs"},
		},
		3: {
			{Value: uint(0), Name: "-5µs"},
			{Value: uint(3), Name: "0s"},
			{Value: uint(0), Name: "5µs"},
		},
	}

	expectedAxisX := []int{-5, 0, 5}

	axisY, axisX := QotToLineChart(quantities, time.Microsecond)
	require.Equal(t, expectedAxisY, axisY)
	require.Equal(t, expectedAxisX, axisX)
}
