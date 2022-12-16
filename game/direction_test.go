package game

import (
	"testing"
)

func TestOpposite(t *testing.T) {
	// given
	testCases := []struct {
		direction direction
		expected  direction
	}{
		{up, down},
		{left, right},
	}

	for _, testCase := range testCases {
		// when && then
		actual := opposite(testCase.direction)
		if testCase.expected != actual {
			t.Errorf("Incorrect direction, is %s but should be %s", testCase.expected, actual)
		}
	}
}
