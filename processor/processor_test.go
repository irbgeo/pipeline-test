package processor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	In          []int
	ExpectedOut []int
}

var tests = []testCase{
	{
		In:          []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		ExpectedOut: []int{7, 8, 9},
	},
}

func TestFind3Largest(t *testing.T) {
	for _, test := range tests {
		actualOut := find3Largest(test.In)
		require.Equal(t, test.ExpectedOut, actualOut)
	}
}
