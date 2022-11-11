package game

import (
	"testing"
)

func TestNewSnake(t *testing.T) {
	// given
	testCases := []struct {
		head     coordinate
		dir      direction
		expected *snake
	}{
		{
			head: coordinate{2, 4},
			dir:  up,
			expected: &snake{
				body:      []coordinate{{2, 4}, {2, 5}, {2, 6}},
				direction: up,
			},
		},
		{
			head: coordinate{2, 4},
			dir:  down,
			expected: &snake{
				body:      []coordinate{{2, 4}, {2, 3}, {2, 2}},
				direction: down,
			},
		},
		{
			head: coordinate{2, 4},
			dir:  left,
			expected: &snake{
				body:      []coordinate{{2, 4}, {3, 4}, {4, 4}},
				direction: left,
			},
		},
		{
			head: coordinate{2, 4},
			dir:  right,
			expected: &snake{
				body:      []coordinate{{2, 4}, {1, 4}, {0, 4}},
				direction: right,
			},
		},
	}

	for _, testCase := range testCases {
		// when
		snake := newSnake(testCase.head, testCase.dir)

		// then
		if snake.body[0] != testCase.expected.body[0] || snake.body[1] != testCase.expected.body[1] || snake.body[2] != testCase.expected.body[2] {
			t.Errorf("Incorrect snake body, is (%d, %d)(%d, %d)(%d, %d), but should be (%d, %d)(%d, %d)(%d, %d)",
				snake.body[0].x, snake.body[0].y,
				snake.body[1].x, snake.body[1].y,
				snake.body[2].x, snake.body[2].y,
				testCase.expected.body[0].x, testCase.expected.body[0].y,
				testCase.expected.body[1].x, testCase.expected.body[1].y,
				testCase.expected.body[2].x, testCase.expected.body[2].y)
		}
		if snake.direction != testCase.expected.direction {
			t.Errorf("Incorrect snake direction, is %s, should be %s", snake.direction, testCase.expected.direction)
		}
	}
}

func TestHead(t *testing.T) {
	// given
	testCases := []struct {
		snake    *snake
		expected coordinate
	}{
		{newSnake(coordinate{10, 10}, up), coordinate{10, 10}},
		{newSnake(coordinate{12, 5}, down), coordinate{12, 5}},
		{newSnake(coordinate{1, 4}, left), coordinate{1, 4}},
		{newSnake(coordinate{13, 10}, right), coordinate{13, 10}},
	}

	for _, testCase := range testCases {
		// when && then
		if !testCase.snake.head().equals(testCase.expected) {
			t.Errorf("Incorrect head, is (%d, %d) but should be (%d, %d)", testCase.snake.head().x, testCase.snake.head().y, testCase.expected.x, testCase.expected.y)
		}
	}
}

func TestMoveSnake(t *testing.T) {
	// given
	testCases := []struct {
		snake *snake
		head  coordinate
		tail  coordinate
	}{
		{newSnake(coordinate{10, 10}, up), coordinate{10, 9}, coordinate{10, 11}},
		{newSnake(coordinate{12, 5}, down), coordinate{12, 6}, coordinate{12, 4}},
		{newSnake(coordinate{1, 4}, left), coordinate{0, 4}, coordinate{2, 4}},
		{newSnake(coordinate{13, 10}, right), coordinate{14, 10}, coordinate{12, 10}},
	}

	for _, testCase := range testCases {
		// when
		testCase.snake.move()

		// then
		if !testCase.snake.head().equals(testCase.head) {
			t.Errorf("Incorrect head, is (%d, %d) should be (%d, %d)", testCase.snake.head().x, testCase.snake.head().y, testCase.head.x, testCase.head.y)
		}
		if !testCase.snake.body[2].equals(testCase.tail) {
			t.Errorf("Incorrect tail, is (%d, %d) should be (%d, %d)", testCase.snake.body[2].x, testCase.snake.body[2].y, testCase.tail.x, testCase.tail.y)
		}
	}
}

func TestEat(t *testing.T) {
	// given
	testCases := []struct {
		snake *snake
		apple apple
	}{
		{newSnake(coordinate{10, 10}, up), coordinate{10, 9}},
		{newSnake(coordinate{12, 5}, down), coordinate{12, 6}},
		{newSnake(coordinate{1, 4}, left), coordinate{0, 4}},
		{newSnake(coordinate{13, 10}, right), coordinate{14, 10}},
	}

	for _, testCase := range testCases {
		// when
		testCase.snake.eat(testCase.apple)

		//then
		if !testCase.snake.head().equals(testCase.apple) {
			t.Errorf("Incorrect head, is (%d, %d) but should be (%d, %d)", testCase.snake.head().x, testCase.snake.head().y, testCase.apple.x, testCase.apple.y)
		}
		if len(testCase.snake.body) != 4 {
			t.Errorf("Incorrect snake's body length, is %d, but should be %d", len(testCase.snake.body), 4)
		}
	}
}

func TestRedirect(t *testing.T) {
	testCases := []struct {
		snake    *snake
		dir      direction
		expected direction
	}{
		{newSnake(coordinate{10, 10}, up), down, up},
		{newSnake(coordinate{12, 5}, down), up, down},
		{newSnake(coordinate{1, 4}, left), right, left},
		{newSnake(coordinate{13, 10}, right), left, right},
		{newSnake(coordinate{10, 10}, up), left, left},
		{newSnake(coordinate{12, 5}, down), right, right},
		{newSnake(coordinate{1, 4}, left), up, up},
		{newSnake(coordinate{13, 10}, right), down, down},
	}

	for _, testCase := range testCases {
		testCase.snake.redirect(testCase.dir)
		if testCase.snake.direction != testCase.expected {
			t.Errorf("Snake moving in incorrect direction, is %s, but should be %s", testCase.snake.direction, testCase.expected)
		}
	}
}
