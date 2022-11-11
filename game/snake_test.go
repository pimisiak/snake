package game

import "testing"

func TestHead(t *testing.T) {
	t.Fail()
}

func TestMove(t *testing.T) {
	t.Fail()
}

func TestEat(t *testing.T) {
	snake := &snake{
		body: []coordinate{
			{10, 1},
			{11, 1},
			{12, 1},
		},
		direction: up,
	}

	apple := apple{9, 1}
	snake.eat(apple)

	if !snake.head().equals(apple) {
		t.Errorf("Head is (%d, %d), but should be (%d, %d)", snake.head().x, snake.head().y, apple.x, apple.y)
	}
}

func TestRedirect(t *testing.T) {
	t.Fail()
}
