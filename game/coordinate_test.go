package game

import "testing"

func TestEquals(t *testing.T) {
	testCases := []struct {
		input    coordinate
		expected bool
	}{
		{coordinate{5, 1}, false},
		{coordinate{1, 1}, true},
		{coordinate{0, 44}, false},
		{coordinate{15, 8}, false},
	}

	c := coordinate{1, 1}
	for _, testCase := range testCases {
		if testCase.input.equals(c) != testCase.expected {
			t.Errorf("(%d, %d) == (%d, %d) should be %v", testCase.input.x, testCase.input.y, c.x, c.y, testCase.expected)
		}
	}
}

func TestMove(t *testing.T) {
	testCases := []struct {
		input    coordinate
		dir      direction
		expected coordinate
	}{
		{coordinate{5, 1}, up, coordinate{5, 0}},
		{coordinate{5, 1}, down, coordinate{5, 2}},
		{coordinate{5, 1}, left, coordinate{4, 1}},
		{coordinate{5, 1}, right, coordinate{6, 1}},
	}

	for _, testCase := range testCases {
		moved := testCase.input.move(testCase.dir)
		if !testCase.expected.equals(moved) {
			t.Errorf("(%d, %d) incorrectly %s moved to (%d, %d)", testCase.input.x, testCase.input.y, testCase.dir, moved.x, moved.y)
		}
	}
}
