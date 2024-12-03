package limit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateIsValid(t *testing.T) {
	require.NoError(t, Rate{time.Second, 10}.IsValid())
	require.Error(t, Rate{-time.Second, 10}.IsValid())
	require.Error(t, Rate{0, 10}.IsValid())
	require.Error(t, Rate{time.Second, 0}.IsValid())
}
