package accumulator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	In            []int
	ExpectedState int64
}

var (
	testState int64 = 10

	testsUpdateState = []testCase{
		{
			In:            []int{0, 1, 2},
			ExpectedState: 3,
		},
	}
)

func TestGetState(t *testing.T) {
	a := New()

	a.state = testState

	actualState := a.GetState()

	require.Equal(t, testState, int64(actualState))
}

func TestUpdateState(t *testing.T) {
	a := New()

	for _, test := range testsUpdateState {
		err := a.UpdateState(context.Background(), test.In)
		require.NoError(t, err)

		actualState := a.GetState()

		require.Equal(t, test.ExpectedState, int64(actualState))
	}
}
